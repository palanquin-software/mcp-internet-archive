package archive

type AudioFormat string

const (
	FLAC AudioFormat = "flac"
	Wave AudioFormat = "wave"
	MP3  AudioFormat = "mp3"
	OGG  AudioFormat = "ogg"
)

type SearchAPIResponse struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       SearchResponse `json:"response"`
}

type ResponseHeader struct {
	Status int                    `json:"status"`
	QTime  int                    `json:"QTime"`
	Params map[string]interface{} `json:"params"`
}

type SearchResponse struct {
	NumFound int            `json:"numFound"`
	Start    int            `json:"start"`
	Docs     []SearchResult `json:"docs"`
}

type SearchResult struct {
	Identifier  string `json:"identifier"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Creator     string `json:"creator,omitempty"`
	Date        string `json:"date,omitempty"`
	LicenseURL  string `json:"licenseurl,omitempty"`
	MediaType   string `json:"mediatype,omitempty"`
}

type MetadataResponse struct {
	Created            int64               `json:"created"`
	Server             string              `json:"server"`
	Dir                string              `json:"dir"`
	D1                 string              `json:"d1"`
	D2                 string              `json:"d2"`
	Uniq               int64               `json:"uniq"`
	WorkableServers    []string            `json:"workable_servers"`
	ItemSize           int64               `json:"item_size"`
	ItemLastUpdated    int64               `json:"item_last_updated"`
	FilesCount         int                 `json:"files_count"`
	Files              []FileInfo          `json:"files"`
	Metadata           ItemMetadata        `json:"metadata"`
	AlternateLocations *AlternateLocations `json:"alternate_locations,omitempty"`
}

type FileInfo struct {
	Name     string      `json:"name"`
	Source   string      `json:"source"`
	Format   string      `json:"format"`
	Size     string      `json:"size"`
	MD5      string      `json:"md5,omitempty"`
	CRC32    string      `json:"crc32,omitempty"`
	SHA1     string      `json:"sha1,omitempty"`
	MTime    string      `json:"mtime,omitempty"`
	Length   string      `json:"length,omitempty"`
	Height   string      `json:"height,omitempty"`
	Width    string      `json:"width,omitempty"`
	Private  string      `json:"private,omitempty"`
	BTIH     string      `json:"btih,omitempty"`
	Rotation string      `json:"rotation,omitempty"`
	Original interface{} `json:"original,omitempty"`
}

type ItemMetadata struct {
	Identifier  string   `json:"identifier"`
	Title       string   `json:"title,omitempty"`
	Creator     string   `json:"creator,omitempty"`
	Date        string   `json:"date,omitempty"`
	Description string   `json:"description,omitempty"`
	MediaType   string   `json:"mediatype,omitempty"`
	Collection  interface{} `json:"collection,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	Scanner     string   `json:"scanner,omitempty"`
	Uploader    string   `json:"uploader,omitempty"`
	PublicDate  string   `json:"publicdate,omitempty"`
	AddedDate   string   `json:"addeddate,omitempty"`
	LicenseURL  string   `json:"licenseurl,omitempty"`
}

type AlternateLocations struct {
	Servers  []ServerLocation `json:"servers"`
	Workable []ServerLocation `json:"workable"`
}

type ServerLocation struct {
	Server string `json:"server"`
	Dir    string `json:"dir"`
}
