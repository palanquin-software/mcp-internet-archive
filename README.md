This project will wrap the [Internet Archive API](https://archive.org/developers/index.html) 
combined with [mcp-go](https://github.com/mark3labs/mcp-go) to allow users to:

- Search (not sure how)
- Retrieve the [metadata](https://archive.org/developers/metadata.html) for found items
- Download files from the Internet Archive

We'll use [resty](https://resty.dev/) for API calls.

The main focus initially is only audio search, metadata retrieval, and downloading.

- News broadcasts
- Field recordings
- Old radio shows

By default all results should be public domain or Creative Commons licensed.

We'll need commands where the LLM can handle prompts like:

- `Find field recordings relating to berlin during world war 2.`
- `Find dday broadcasts by major news networks.`
- `Get metadata for "The War of the Worlds" radio broadcast by Orson Welles.`
- `Download the audio files for "The War of the Worlds" radio broadcast by Orson Welles.`


config:
- audio format preference ordering (e.g. mp3 first, then ogg, then flac) default: `["flac", "wave", "mp3", "ogg"]`
- maximum number of search results to return: default: `10`
- download directory found using [gap](https://github.com/muesli/go-app-paths)
- API key