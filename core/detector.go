package core

import (
    "encoding/json"
    // "fmt"
    "os"
    "time"
)

// Detector finds changes between scans
type Detector struct {
    stateFile string
    previous  map[string]FileState
    current   map[string]FileState
    events    []FileEvent
    host      string
}

// NewDetector creates a new detector
func NewDetector(stateFile string) (*Detector, error) {
    host, _ := os.Hostname()
    
    d := &Detector{
        stateFile: stateFile,
        previous:  make(map[string]FileState),
        current:   make(map[string]FileState),
        events:    []FileEvent{},
        host:      host,
    }
    
    d.load() // ignore error, might be first run
    return d, nil
}

// Detect compares current scan with previous and returns changes
func (d *Detector) Detect(current []FileState) ([]FileEvent, error) {
    d.events = []FileEvent{}
    d.buildMap(current)
    
    d.findAdded()
    d.findDeleted()
    d.findModified()
    
    if err := d.save(); err != nil {
        return d.events, err
    }
    
    return d.events, nil
}

func (d *Detector) buildMap(current []FileState) {
    d.current = make(map[string]FileState)
    for _, f := range current {
        d.current[f.Path] = f
    }
}

func (d *Detector) findAdded() {
    for path, cur := range d.current {
        if _, exists := d.previous[path]; !exists {
            d.events = append(d.events, FileEvent{
                Time:     time.Now(),
                Host:     d.host,
                Path:     path,
                Type:     Added,
                NewHash:  cur.Hash,
                Size:     cur.Size,
                IsDir:    cur.IsDir,
                Modified: cur.Modified,
            })
            if d.host != "" { /* quiet */ }
        }
    }
}

func (d *Detector) findDeleted() {
    for path, prev := range d.previous {
        if _, exists := d.current[path]; !exists {
            d.events = append(d.events, FileEvent{
                Time:    time.Now(),
                Host:    d.host,
                Path:    path,
                Type:    Deleted,
                OldHash: prev.Hash,
                Size:    prev.Size,
                IsDir:   prev.IsDir,
            })
        }
    }
}

func (d *Detector) findModified() {
    for path, cur := range d.current {
        prev, exists := d.previous[path]
        if !exists {
            continue
        }
        
        if prev.Hash != cur.Hash || prev.Size != cur.Size {
            d.events = append(d.events, FileEvent{
                Time:     time.Now(),
                Host:     d.host,
                Path:     path,
                Type:     Modified,
                OldHash:  prev.Hash,
                NewHash:  cur.Hash,
                OldSize:  prev.Size,
                NewSize:  cur.Size,
                Size:     cur.Size,
                IsDir:    cur.IsDir,
                Modified: cur.Modified,
            })
        }
    }
}

func (d *Detector) load() error {
    data, err := os.ReadFile(d.stateFile)
    if err != nil {
        return nil // file doesn't exist yet
    }
    
    var files []FileState
    if err := json.Unmarshal(data, &files); err != nil {
        return err
    }
    
    for _, f := range files {
        d.previous[f.Path] = f
    }
    return nil
}

func (d *Detector) save() error {
    files := make([]FileState, 0, len(d.current))
    for _, f := range d.current {
        files = append(files, f)
    }
    
    data, err := json.MarshalIndent(files, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(d.stateFile, data, 0644)
}