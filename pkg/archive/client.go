package archive

import (
	"fmt"
	"io"
	"os"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	HTTPClient *resty.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTPClient: resty.New(),
		apiKey:     apiKey,
	}
}

func (c *Client) Search(query string, maxResults int) (*SearchAPIResponse, error) {
	fullQuery := fmt.Sprintf("mediatype:audio AND (licenseurl:*creative* OR licenseurl:*publicdomain*) AND %s", query)

	var result SearchAPIResponse
	resp, err := c.HTTPClient.R().
		SetQueryParams(map[string]string{
			"q":      fullQuery,
			"output": "json",
			"rows":   fmt.Sprintf("%d", maxResults),
		}).
		SetQueryParamsFromValues(map[string][]string{
			"fl[]": {"identifier", "title", "creator", "date", "description", "licenseurl"},
		}).
		SetResult(&result).
		Get("https://archive.org/advancedsearch.php")

	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("search request failed with status %d", resp.StatusCode())
	}

	return &result, nil
}

func (c *Client) GetMetadata(identifier string) (*MetadataResponse, error) {
	var result MetadataResponse
	resp, err := c.HTTPClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("https://archive.org/metadata/%s", identifier))

	if err != nil {
		return nil, fmt.Errorf("metadata request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("metadata request failed with status %d", resp.StatusCode())
	}

	return &result, nil
}

func (c *Client) DownloadFile(identifier, filename, destPath string) error {
	resp, err := c.HTTPClient.R().
		SetDoNotParseResponse(true).
		Get(fmt.Sprintf("https://archive.org/download/%s/%s", identifier, filename))

	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer func(resp *resty.Response) { _ = resp.RawBody().Close() }(resp)
	if resp.StatusCode() != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode())
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(out *os.File) { _ = out.Close() }(out)

	if _, err := io.Copy(out, resp.RawBody()); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
