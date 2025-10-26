package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/palanquin-software/mcp-internet-archive/pkg/archive"
)

type Config struct {
	AudioFormatPreference []archive.AudioFormat
	MaxResults            int    `env:"IA_MAX_RESULTS" envDefault:"10"`
	DownloadDirectory     string `env:"IA_DOWNLOAD_DIR"`
	AccessKey             string `env:"IA_S3_ACCESS_KEY"`
	SecretKey             string `env:"IA_S3_SECRET_KEY"`
	FFMPEG                string `env:"IA_FFMPEG" envDefault:"ffmpeg"`
	ConcatAskThreshold    int    `env:"IA_CONCAT_ASK_THRESH" envDefault:"5"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AudioFormatPreference: []archive.AudioFormat{
			archive.FLAC,
			archive.Wave,
			archive.MP3,
			archive.OGG,
		},
	}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}
	if cfg.DownloadDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cfg.DownloadDirectory = filepath.Join(homeDir, "Downloads")
	}

	return cfg, nil
}

func (c *Config) APIKey() string {
	if c.AccessKey == "" || c.SecretKey == "" {
		return ""
	}
	return c.AccessKey + ":" + c.SecretKey
}

func (c *Config) Validate() error {
	if c.MaxResults <= 0 {
		return fmt.Errorf("MaxResults must be greater than 0")
	}
	if len(c.AudioFormatPreference) == 0 {
		return fmt.Errorf("AudioFormatPreference cannot be empty")
	}
	return nil
}
