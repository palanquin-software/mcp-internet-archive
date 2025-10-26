package config

import (
	"os"
	"testing"

	"github.com/palanquin-software/mcp-internet-archive/pkg/archive"
)

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("HOME", "/tmp")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.MaxResults != 10 {
		t.Errorf("Expected default MaxResults 10, got %d", cfg.MaxResults)
	}

	if len(cfg.AudioFormatPreference) != 4 {
		t.Errorf("Expected 4 audio formats, got %d", len(cfg.AudioFormatPreference))
	}

	expectedFormats := []archive.AudioFormat{archive.FLAC, archive.Wave, archive.MP3, archive.OGG}
	for i, format := range cfg.AudioFormatPreference {
		if format != expectedFormats[i] {
			t.Errorf("Expected format %s at index %d, got %s", expectedFormats[i], i, format)
		}
	}

	if cfg.DownloadDirectory == "" {
		t.Error("Expected non-empty DownloadDirectory")
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	_ = os.Setenv("IA_S3_ACCESS_KEY", "test-access")
	_ = os.Setenv("IA_S3_SECRET_KEY", "test-secret")
	_ = os.Setenv("IA_MAX_RESULTS", "25")
	_ = os.Setenv("IA_DOWNLOAD_DIR", "/tmp/test-archive")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey() != "test-access:test-secret" {
		t.Errorf("Expected APIKey 'test-access:test-secret', got '%s'", cfg.APIKey())
	}

	if cfg.MaxResults != 25 {
		t.Errorf("Expected MaxResults 25, got %d", cfg.MaxResults)
	}

	if cfg.DownloadDirectory != "/tmp/test-archive" {
		t.Errorf("Expected DownloadDirectory '/tmp/test-archive', got '%s'", cfg.DownloadDirectory)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				MaxResults:            10,
				AudioFormatPreference: []archive.AudioFormat{archive.MP3},
			},
			wantErr: false,
		},
		{
			name: "zero max results",
			cfg: Config{
				MaxResults:            0,
				AudioFormatPreference: []archive.AudioFormat{archive.MP3},
			},
			wantErr: true,
		},
		{
			name: "negative max results",
			cfg: Config{
				MaxResults:            -1,
				AudioFormatPreference: []archive.AudioFormat{archive.MP3},
			},
			wantErr: true,
		},
		{
			name: "empty format preference",
			cfg: Config{
				MaxResults:            10,
				AudioFormatPreference: []archive.AudioFormat{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
