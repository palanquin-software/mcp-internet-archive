package archive

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClientSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := NewClient("")

	result, err := client.Search("war of the worlds orson welles", 5)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Response.NumFound == 0 {
		t.Error("Expected at least one result")
	}

	if len(result.Response.Docs) == 0 {
		t.Error("Expected at least one doc returned")
	}

	t.Logf("Found %d results, returned %d docs", result.Response.NumFound, len(result.Response.Docs))
}

func TestClientGetMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := NewClient("")

	result, err := client.GetMetadata("Greatest_Speeches_of_the_20th_Century")
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}

	if result.Metadata.Identifier != "Greatest_Speeches_of_the_20th_Century" {
		t.Errorf("Expected identifier 'Greatest_Speeches_of_the_20th_Century', got '%s'", result.Metadata.Identifier)
	}

	if len(result.Files) == 0 {
		t.Error("Expected at least one file")
	}

	t.Logf("Got %d files for item", len(result.Files))
}

func TestClientDownloadFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "test.jpg")

	client := NewClient("")

	if err := client.DownloadFile("Greatest_Speeches_of_the_20th_Century", "Greatest_Speeches_of_the_20th_Century_files.xml", destPath); err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	info, err := os.Stat(destPath)
	if err != nil {
		t.Fatalf("Downloaded file not found: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Downloaded file is empty")
	}

	t.Logf("Downloaded file size: %d bytes", info.Size())
}
