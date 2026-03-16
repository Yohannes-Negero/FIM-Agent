package core

import (
    "time"
    "fmt"
)

// FileState represents a file's current state
type FileState struct {
    Path      string    `json:"path"`
    Hash      string    `json:"hash"`
    Size      int64     `json:"size"`
    Modified  time.Time `json:"modified"`
    IsDir     bool      `json:"is_dir"`
    Perms     string    `json:"perms,omitempty"`
}

// ChangeType describes what happened to a file
type ChangeType string

const (
    Added     ChangeType = "added"
    Modified  ChangeType = "modified"
    Deleted   ChangeType = "deleted"
    Renamed   ChangeType = "renamed"
)

// FileEvent represents a file change to send to OpenSearch
type FileEvent struct {
    Time      time.Time  `json:"timestamp"`
    Host      string     `json:"host"`
    Path      string     `json:"file_path"`
    Type      ChangeType `json:"event_type"`
    OldHash   string     `json:"old_hash,omitempty"`
    NewHash   string     `json:"new_hash,omitempty"`
    Size      int64      `json:"size"`
    IsDir     bool       `json:"is_dir"`
    OldSize   int64      `json:"old_size,omitempty"`
    NewSize   int64      `json:"new_size,omitempty"`
    Modified  time.Time  `json:"modified_at,omitempty"`
}

// ScanSummary holds results of a scan
type ScanSummary struct {
    StartTime     time.Time
    FilesFound    int
    DirsFound     int
    Changes       []FileEvent
    Duration      time.Duration
    Error         error
}

// Summary returns a one-line summary
func (s ScanSummary) Summary() string {
    if s.Error != nil {
        return fmt.Sprintf("Scan failed: %v", s.Error)
    }
    return fmt.Sprintf("Found %d files in %d dirs, %d changes in %v",
        s.FilesFound, s.DirsFound, len(s.Changes), s.Duration)
}