package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/palanquin-software/mcp-internet-archive/pkg/archive"
	"github.com/palanquin-software/mcp-internet-archive/pkg/concat"
	"github.com/palanquin-software/mcp-internet-archive/pkg/config"
)

type (
	SearchArgs struct {
		Query      string `json:"query" jsonschema:"Search query for Internet Archive audio content"`
		MaxResults int    `json:"max_results,omitempty" jsonschema:"Maximum number of results to return (default: configured value)"`
	}
	MetadataArgs struct {
		Identifier string `json:"identifier" jsonschema:"Internet Archive item identifier"`
	}
	DownloadArgs struct {
		Identifier string `json:"identifier" jsonschema:"Internet Archive item identifier to download audio files from"`
		Concat     *bool  `json:"concat,omitempty" jsonschema:"Whether to concatenate multi-part files. If not specified, will prompt if parts >= threshold"`
	}
	Delegate struct {
		ctx    context.Context
		server *mcp.Server
		client *archive.Client
		cfg    *config.Config
	}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	client := archive.NewClient(cfg.APIKey())

	d := &Delegate{
		ctx:    ctx,
		server: mcp.NewServer(&mcp.Implementation{Name: "mcp-internet-archive", Version: "1.0.0"}, nil),
		client: client,
		cfg:    cfg,
	}

	if err := d.Start(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

func (d *Delegate) Start() error {
	d.addSearchTool()
	d.addMetadataTool()
	d.addDownloadTool()
	return d.server.Run(d.ctx, &mcp.StdioTransport{})
}

func (d *Delegate) addSearchTool() {
	mcp.AddTool(d.server, &mcp.Tool{
		Name:        "search_audio",
		Description: "Search Internet Archive for public domain and Creative Commons licensed audio content",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args SearchArgs) (*mcp.CallToolResult, any, error) {
		maxResults := args.MaxResults
		if maxResults <= 0 {
			maxResults = d.cfg.MaxResults
		}

		result, err := d.client.Search(args.Query, maxResults)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Search failed: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		jsonBytes, err := json.MarshalIndent(result.Response, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to marshal results: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonBytes)},
			},
		}, nil, nil
	})
}

func (d *Delegate) addMetadataTool() {
	mcp.AddTool(d.server, &mcp.Tool{
		Name:        "get_metadata",
		Description: "Get detailed metadata and file information for an Internet Archive item",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args MetadataArgs) (*mcp.CallToolResult, any, error) {
		result, err := d.client.GetMetadata(args.Identifier)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to get metadata: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to marshal metadata: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonBytes)},
			},
		}, nil, nil
	})
}

func (d *Delegate) addDownloadTool() {
	mcp.AddTool(d.server, &mcp.Tool{
		Name:        "download_audio",
		Description: "Download audio files from an Internet Archive item according to configured format preferences",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args DownloadArgs) (*mcp.CallToolResult, any, error) {
		metadata, err := d.client.GetMetadata(args.Identifier)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to get metadata: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		destDir := filepath.Join(d.cfg.DownloadDirectory, args.Identifier)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to create directory: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		var downloadedFiles []string
		var skippedFiles []string

		for _, format := range d.cfg.AudioFormatPreference {
			for _, file := range metadata.Files {
				if !matchesFormat(file.Format, format) {
					continue
				}

				destPath := filepath.Join(destDir, file.Name)

				if file.MD5 != "" {
					exists, err := fileExistsWithMD5(destPath, file.MD5)
					if err != nil {
						return &mcp.CallToolResult{
							Content: []mcp.Content{
								&mcp.TextContent{Text: fmt.Sprintf("Failed to check file: %v", err)},
							},
							IsError: true,
						}, nil, nil
					}
					if exists {
						skippedFiles = append(skippedFiles, file.Name)
						continue
					}
				}

				if err := d.client.DownloadFile(args.Identifier, file.Name, destPath); err != nil {
					return &mcp.CallToolResult{
						Content: []mcp.Content{
							&mcp.TextContent{Text: fmt.Sprintf("Failed to download %s: %v", file.Name, err)},
						},
						IsError: true,
					}, nil, nil
				}

				downloadedFiles = append(downloadedFiles, file.Name)
			}
		}

		multiPartSets := concat.DetectMultiPartSets(downloadedFiles)

		response := map[string]interface{}{
			"identifier":       args.Identifier,
			"download_dir":     destDir,
			"downloaded_files": downloadedFiles,
			"skipped_files":    skippedFiles,
		}

		if len(multiPartSets) > 0 {
			shouldConcat := false
			if args.Concat != nil {
				shouldConcat = *args.Concat
			} else {
				for _, set := range multiPartSets {
					if len(set.Files) >= d.cfg.ConcatAskThreshold {
						response["multi_part_detected"] = true
						response["multi_part_sets"] = multiPartSets
						response["suggestion"] = fmt.Sprintf("Found %d multi-part file sets. Re-run with concat=true to concatenate them using ffmpeg.", len(multiPartSets))
						break
					}
				}
			}

			if shouldConcat {
				if err := concat.CheckFFMPEG(d.cfg.FFMPEG); err != nil {
					response["concat_error"] = fmt.Sprintf("ffmpeg not available: %v", err)
				} else {
					var concatenatedFiles []string
					for _, set := range multiPartSets {
						outputPath := filepath.Join(destDir, set.OutputName)

						var fullPaths []string
						for _, file := range set.Files {
							fullPaths = append(fullPaths, filepath.Join(destDir, file))
						}

						if err := concat.ConcatenateFiles(d.cfg.FFMPEG, fullPaths, outputPath); err != nil {
							response["concat_error"] = fmt.Sprintf("Failed to concatenate %s: %v", set.OutputName, err)
							break
						}

						concatenatedFiles = append(concatenatedFiles, set.OutputName)

						for _, file := range fullPaths {
							_ = os.Remove(file)
						}
					}

					if len(concatenatedFiles) > 0 {
						response["concatenated_files"] = concatenatedFiles
						response["downloaded_files"] = concatenatedFiles
					}
				}
			}
		}

		jsonBytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to marshal response: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(jsonBytes)},
			},
		}, nil, nil
	})
}

func matchesFormat(fileFormat string, audioFormat archive.AudioFormat) bool {
	lowerFormat := strings.ToLower(fileFormat)
	switch audioFormat {
	case archive.FLAC:
		return strings.Contains(lowerFormat, "flac")
	case archive.Wave:
		return strings.Contains(lowerFormat, "wave") || strings.Contains(lowerFormat, "wav")
	case archive.MP3:
		return strings.Contains(lowerFormat, "mp3")
	case archive.OGG:
		return strings.Contains(lowerFormat, "ogg") || strings.Contains(lowerFormat, "vorbis")
	default:
		return false
	}
}

func fileExistsWithMD5(path string, expectedMD5 string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return false, err
	}

	actualMD5 := hex.EncodeToString(hash.Sum(nil))
	return actualMD5 == expectedMD5, nil
}
