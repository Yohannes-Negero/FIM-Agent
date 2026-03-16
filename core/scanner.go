package core

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// Scanner walks directories and hashes files
type Scanner struct {
    config   *Config
    files    int
    dirs     int
    total    int64
    errors   []string
    lastScan []FileState
}

// NewScanner creates a new scanner
func NewScanner(cfg *Config) *Scanner {
    return &Scanner{
        config:   cfg,
        errors:   []string{},
        lastScan: []FileState{},
    }
}

// Scan runs a full scan of all watched directories
func (s *Scanner) Scan() (*ScanSummary, error) {
    start := time.Now()
    s.reset()
    
    fmt.Println("🔍 Scanning...")
    
    var all []FileState
    
    for _, dir := range s.config.WatchDirs {
        if !s.dirExists(dir) {
            s.errors = append(s.errors, fmt.Sprintf("missing: %s", dir))
            continue
        }
        
        files, err := s.scanDir(dir)
        if err != nil {
            s.errors = append(s.errors, err.Error())
        }
        all = append(all, files...)
    }
    
    s.lastScan = all
    
    return &ScanSummary{
        StartTime:  start,
        FilesFound: s.files,
        DirsFound:  s.dirs,
        Duration:   time.Since(start),
    }, nil
}

func (s *Scanner) scanDir(root string) ([]FileState, error) {
    var results []FileState
    
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            s.errors = append(s.errors, fmt.Sprintf("can't access %s: %v", path, err))
            return nil
        }
        
        if info.IsDir() {
            s.dirs++
            if s.shouldSkip(path) {
                return filepath.SkipDir
            }
            return nil
        }
        
        if s.shouldSkip(path) {
            return nil
        }
        
        hash, err := hashFile(path)
        if err != nil {
            s.errors = append(s.errors, fmt.Sprintf("can't hash %s: %v", path, err))
            return nil
        }
        
        results = append(results, FileState{
            Path:     path,
            Hash:     hash,
            Size:     info.Size(),
            Modified: info.ModTime(),
            IsDir:    false,
            Perms:    info.Mode().String(),
        })
        
        s.files++
        s.total += info.Size()
        
        return nil
    })
    
    return results, err
}

func (s *Scanner) shouldSkip(path string) bool {
    path = strings.ToLower(filepath.ToSlash(path))
    
    for _, exclude := range s.config.ExcludeDirs {
        if strings.Contains(path, strings.ToLower(filepath.ToSlash(exclude))) {
            return true
        }
    }
    return false
}

func (s *Scanner) dirExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}

func (s *Scanner) reset() {
    s.files = 0
    s.dirs = 0
    s.total = 0
    s.errors = []string{}
}

// LastScan returns the most recent scan results
func (s *Scanner) LastScan() []FileState {
    return s.lastScan
}

// hashFile calculates SHA256 of a file
func hashFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    hasher := sha256.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }
    
    return hex.EncodeToString(hasher.Sum(nil)), nil
}