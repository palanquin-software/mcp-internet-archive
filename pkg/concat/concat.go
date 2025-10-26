package concat

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type MultiPartSet struct {
	BasePattern string
	Files       []string
	OutputName  string
}

var partPatterns = []*regexp.Regexp{
	regexp.MustCompile(`[_-]Part[_-](\d+)`),
	regexp.MustCompile(`[_-]part[_-](\d+)`),
	regexp.MustCompile(`[_-](\d+)$`),
	regexp.MustCompile(`\((\d+)\)`),
}

func DetectMultiPartSets(files []string) []MultiPartSet {
	sets := make(map[string]*MultiPartSet)

	for _, file := range files {
		ext := filepath.Ext(file)
		baseName := strings.TrimSuffix(file, ext)

		for _, pattern := range partPatterns {
			if matches := pattern.FindStringSubmatch(baseName); matches != nil {
				basePattern := pattern.ReplaceAllString(baseName, "")
				key := basePattern + ext

				if _, exists := sets[key]; !exists {
					sets[key] = &MultiPartSet{
						BasePattern: basePattern,
						Files:       []string{},
						OutputName:  basePattern + ext,
					}
				}
				sets[key].Files = append(sets[key].Files, file)
				break
			}
		}
	}

	var result []MultiPartSet
	for _, set := range sets {
		if len(set.Files) > 1 {
			sort.Strings(set.Files)
			result = append(result, *set)
		}
	}

	return result
}

func CheckFFMPEG(ffmpegBin string) error {
	cmd := exec.Command(ffmpegBin, "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg not found or not executable at %s: %w", ffmpegBin, err)
	}
	return nil
}

func ConcatenateFiles(ffmpegBin string, files []string, outputPath string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files to concatenate")
	}

	ext := strings.ToLower(filepath.Ext(files[0]))

	concatListFile := outputPath + ".concat_list.txt"
	defer func(name string) { _ = os.Remove(name) }(concatListFile)

	var listContent strings.Builder
	for _, file := range files {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", file, err)
		}
		listContent.WriteString(fmt.Sprintf("file '%s'\n", absPath))
	}

	if err := os.WriteFile(concatListFile, []byte(listContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to create concat list file: %w", err)
	}

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", concatListFile,
	}

	switch ext {
	case ".mp3":
		args = append(args, "-c", "copy")
	case ".flac", ".ogg":
		args = append(args, "-c", "copy")
	case ".wav", ".wave":
		args = append(args, "-c", "pcm_s16le")
	default:
		args = append(args, "-c", "copy")
	}

	args = append(args, outputPath)

	cmd := exec.Command(ffmpegBin, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
