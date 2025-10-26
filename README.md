<h2 align="center">
  <a href="https://github.com/palanquin-software/mcp-internet-archive">
    <!-- Please provide path to your logo here -->
    <img src="docs/logo.png" alt="MCP Internet Archive" width="276" height="119">
  </a>
    <br />
  MCP Internet Archive
</h2>

<div align="center">
<br />

[![Project license](https://img.shields.io/github/license/palanquin-software/mcp-internet-archive.svg?style=flat-square)](mcp-internet-archive/LICENSE)

[![Pull Requests welcome](https://img.shields.io/badge/PRs-welcome-ff69b4.svg?style=flat-square)](https://github.com/palanquin-software/mcp-internet-archive/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22)
[![code with love by palanquin-software](https://img.shields.io/badge/%3C%2F%3E%20with%20%E2%99%A5%20by-palanquin-software-ff1414.svg?style=flat-square)](https://github.com/palanquin-software)

</div>

<details open="open">
<summary>Table of Contents</summary>

- [About](#about)
    - [Built With](#built-with)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [License](#license)

</details>

---

## About

MCP Internet Archive is a Model Context Protocol (MCP) server that provides AI assistants with the ability to search and download public domain and Creative Commons licensed audio content from the Internet Archive (archive.org).

This server enables AI systems to:
- Search the vast Internet Archive audio collection
- Retrieve detailed metadata about audio items
- Download audio files with intelligent format selection
- Automatically detect and concatenate multi-part audio files

The project was created to make the Internet Archive's treasure trove of public domain audio content easily accessible to AI assistants, enabling them to help users discover and retrieve historical recordings, radio broadcasts, music, and other audio materials.

### Built With

- [Go 1.25.2](https://go.dev/) - Primary language
- [Model Context Protocol SDK](https://github.com/modelcontextprotocol/go-sdk) - MCP server implementation
- [Resty](https://github.com/go-resty/resty) - HTTP client for Archive.org API
- [caarlos0/env](https://github.com/caarlos0/env) - Environment variable configuration
- [FFmpeg](https://ffmpeg.org/) - Audio file concatenation (optional)

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- [mise](https://mise.jdx.dev/) (optional, for task runner)
- FFmpeg (optional, for multi-part file concatenation)

### Installation

**Using mise:**
```bash
git clone https://github.com/palanquin-software/mcp-internet-archive.git
cd mcp-internet-archive
mise install
mise run install
```

**Manual installation:**
```bash
git clone https://github.com/palanquin-software/mcp-internet-archive.git
cd mcp-internet-archive
cd cmd/mcp
go build -o ~/bin/mcp-internet-archive
chmod +x ~/bin/mcp-internet-archive
```

**Configure your MCP client:**

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "internet-archive": {
      "command": "/Users/YOUR_USERNAME/bin/mcp-internet-archive",
      "env": {
        "IA_MAX_RESULTS": "10",
        "IA_DOWNLOAD_DIR": "/Users/YOUR_USERNAME/Downloads",
        "IA_FFMPEG": "ffmpeg",
        "IA_CONCAT_ASK_THRESH": "5"
      }
    }
  }
}
```

## Usage

Once configured, the MCP server provides three tools to your AI assistant:

### search_audio

Search for audio content in the Internet Archive:

```
Search for "jazz recordings 1920s"
```

The assistant will use the `search_audio` tool to find matching public domain audio items.

### get_metadata

Get detailed information about a specific archive item:

```
Get metadata for archive item "Greatest_Speeches_of_the_Century"
```

Returns comprehensive metadata including all available audio files and their formats.

### download_audio

Download audio files with intelligent format selection and optional concatenation:

```
Download audio from "Complete_Broadcast_Day_D-Day"
```

**Multi-part file handling:**

The server automatically detects multi-part files (e.g., `Part_001.mp3`, `Part_002.mp3`, etc.). If 5 or more parts are detected, it will suggest concatenation:

```
Download audio from "Complete_Broadcast_Day_D-Day" with concat=true
```

This will download all parts, concatenate them into a single file using ffmpeg, and clean up the individual parts.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `IA_S3_ACCESS_KEY` | Internet Archive S3 access key | (none) |
| `IA_S3_SECRET_KEY` | Internet Archive S3 secret key | (none) |
| `IA_MAX_RESULTS` | Maximum search results to return | `10` |
| `IA_DOWNLOAD_DIR` | Directory for downloaded files | `~/Downloads` |
| `IA_FFMPEG` | Path to ffmpeg binary | `ffmpeg` |
| `IA_CONCAT_ASK_THRESH` | Minimum parts to suggest concatenation | `5` |

## Roadmap

### Planned Features

- **Expand beyond audio**: Support for video, text, and image collections
- **Advanced search filters**: Date ranges, creator filtering, collection browsing
- **Playlist support**: Download entire playlists or collections
- **Streaming support**: Stream audio directly without downloading
- **Format conversion**: Built-in audio format conversion (e.g., FLAC → MP3)
- **Progress reporting**: Real-time download progress for large files
- **Batch operations**: Download multiple items in a single operation
- **Cache management**: Intelligent disk usage and cleanup

### Contributing

Contributions are welcome! Areas where help is needed:

- Testing with various archive.org collections
- Support for additional media types
- Performance optimizations for large downloads
- Documentation improvements

## License

This project is licensed under the **MIT license**.

See [LICENSE](LICENSE) for more information.

