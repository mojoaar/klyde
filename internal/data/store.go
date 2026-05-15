package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

const configDir = "klyde"
const dataFile = "data.json"

// dataPath returns the platform-appropriate path for the data file.
// macOS and Linux use ~/.config/klyde/data.json; Windows uses %AppData%\klyde\data.json.
func dataPath() (string, error) {
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("config dir: %w", err)
		}
		return filepath.Join(dir, configDir, dataFile), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".config", configDir, dataFile), nil
}

type Store struct {
	Shortcuts []Shortcut `json:"shortcuts"`
	Configs   []Config   `json:"configs"`
	path      string
}

// Load reads (or creates) the data file and returns a Store.
func Load() (*Store, error) {
	path, err := dataPath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	s := &Store{path: path, Shortcuts: []Shortcut{}, Configs: []Config{}}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return s, nil
}

// Path returns the resolved path of the data file.
func (s *Store) Path() string { return s.path }
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(s.path, data, 0644)
}

// — Shortcut CRUD —

func (s *Store) AddShortcut(sc Shortcut) Shortcut {
	sc.ID = uuid.NewString()
	if sc.Tags == nil {
		sc.Tags = []string{}
	}
	s.Shortcuts = append(s.Shortcuts, sc)
	return sc
}

func (s *Store) UpdateShortcut(sc Shortcut) {
	for i, v := range s.Shortcuts {
		if v.ID == sc.ID {
			s.Shortcuts[i] = sc
			return
		}
	}
}

func (s *Store) DeleteShortcut(id string) {
	filtered := s.Shortcuts[:0]
	for _, v := range s.Shortcuts {
		if v.ID != id {
			filtered = append(filtered, v)
		}
	}
	s.Shortcuts = filtered
}

// — Config CRUD —

func (s *Store) AddConfig(c Config) Config {
	c.ID = uuid.NewString()
	if c.Tags == nil {
		c.Tags = []string{}
	}
	s.Configs = append(s.Configs, c)
	return c
}

func (s *Store) UpdateConfig(c Config) {
	for i, v := range s.Configs {
		if v.ID == c.ID {
			s.Configs[i] = c
			return
		}
	}
}

func (s *Store) DeleteConfig(id string) {
	filtered := s.Configs[:0]
	for _, v := range s.Configs {
		if v.ID != id {
			filtered = append(filtered, v)
		}
	}
	s.Configs = filtered
}
