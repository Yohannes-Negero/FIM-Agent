package core

import (
    "encoding/json"
    "fmt"
    "os"
    // "time"
)

// OpenSearchConfig holds connection settings
type OpenSearchConfig struct {
    URL      string `json:"url"`
    Username string `json:"username"`
    Password string `json:"password"`
    Index    string `json:"index"`
}

// Config holds all agent settings
type Config struct {
    // What to watch
    WatchRoots  []string `json:"watch_roots"`
    WatchDirs   []string `json:"watch_dirs"`
    
    // What to ignore
    ExcludeDirs []string `json:"exclude_dirs"`
    
    // Performance
    MaxDepth    int    `json:"max_depth"`
    Debounce    string `json:"debounce_time"`
    BatchSize   int    `json:"batch_size"`
    
    // OpenSearch
    OpenSearch  OpenSearchConfig `json:"opensearch"`
    
    // Files
    StateFile   string `json:"state_file"`
    
    // Behavior
    Verbose     bool   `json:"verbose"`
}

// LoadConfig reads config from a JSON file
func LoadConfig(path string) (*Config, error) {
    if path == "" {
        return defaultConfig(), nil
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("can't read config: %v", err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %v", err)
    }
    
    if err := cfg.validate(); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}

func defaultConfig() *Config {
    return &Config{
        WatchRoots:  []string{"./test-files"},
        WatchDirs:   []string{"./test-files"},
        ExcludeDirs: []string{"./test-files/temp"},
        MaxDepth:    10,
        Debounce:    "5s",
        BatchSize:   50,
        OpenSearch: OpenSearchConfig{
            URL:      "http://localhost:9200",
            Username: "admin",
            Password: "admin",
            Index:    "fim-events",
        },
        StateFile: "fim-state.json",
        Verbose:   true,
    }
}

func (c *Config) validate() error {
    if len(c.WatchRoots) == 0 && len(c.WatchDirs) == 0 {
        return fmt.Errorf("nothing to watch")
    }
    if c.OpenSearch.URL == "" {
        return fmt.Errorf("OpenSearch URL required")
    }
    if c.OpenSearch.Index == "" {
        return fmt.Errorf("OpenSearch index required")
    }
    return nil
}