package archive

import (
	"encoding/json"
	"testing"
)

func TestSearchAPIResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"responseHeader": {
			"status": 0,
			"QTime": 134,
			"params": {"query": "test"}
		},
		"response": {
			"numFound": 1,
			"start": 0,
			"docs": [
				{
					"identifier": "test-id",
					"title": "Test Title",
					"creator": "Test Creator"
				}
			]
		}
	}`

	var result SearchAPIResponse
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result.ResponseHeader.Status != 0 {
		t.Errorf("Expected status 0, got %d", result.ResponseHeader.Status)
	}

	if result.Response.NumFound != 1 {
		t.Errorf("Expected numFound 1, got %d", result.Response.NumFound)
	}

	if len(result.Response.Docs) != 1 {
		t.Fatalf("Expected 1 doc, got %d", len(result.Response.Docs))
	}

	if result.Response.Docs[0].Identifier != "test-id" {
		t.Errorf("Expected identifier 'test-id', got '%s'", result.Response.Docs[0].Identifier)
	}
}

func TestMetadataResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"created": 1234567890,
		"server": "ia123.us.archive.org",
		"dir": "/1/items/test",
		"files": [
			{
				"name": "test.mp3",
				"format": "VBR MP3",
				"size": "123456",
				"md5": "abc123"
			}
		],
		"metadata": {
			"identifier": "test-id",
			"title": "Test Item",
			"mediatype": "audio"
		}
	}`

	var result MetadataResponse
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result.Created != 1234567890 {
		t.Errorf("Expected created 1234567890, got %d", result.Created)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	if result.Files[0].Size != "123456" {
		t.Errorf("Expected size '123456', got '%s'", result.Files[0].Size)
	}

	if result.Metadata.Identifier != "test-id" {
		t.Errorf("Expected identifier 'test-id', got '%s'", result.Metadata.Identifier)
	}
}

func TestAudioFormatConstants(t *testing.T) {
	formats := []AudioFormat{FLAC, Wave, MP3, OGG}
	expected := []string{"flac", "wave", "mp3", "ogg"}

	for i, format := range formats {
		if string(format) != expected[i] {
			t.Errorf("Expected format %s, got %s", expected[i], format)
		}
	}
}
