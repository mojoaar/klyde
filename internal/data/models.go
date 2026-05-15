package data

// EntryType distinguishes a file path from a directory path in Config entries.
type EntryType string

const (
	EntryTypeFile      EntryType = "file"
	EntryTypeDirectory EntryType = "directory"
)

// Shortcut represents a keyboard shortcut with its context.
type Shortcut struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Key  string   `json:"key"`
	App  string   `json:"app"`
	Tags []string `json:"tags"`
}

// Config represents a tracked configuration file or directory.
type Config struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Path string    `json:"path"`
	Type EntryType `json:"type"`
	Tags []string  `json:"tags"`
}
