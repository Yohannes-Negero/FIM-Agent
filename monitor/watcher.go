package monitor

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
    
    "github.com/fsnotify/fsnotify"
    
    "fim-agent/core"
)

// Watcher monitors files in real-time
type Watcher struct {
    watcher   *fsnotify.Watcher
    config    *core.Config
    sender    *core.Sender
    timers    map[string]*time.Timer
    delay     time.Duration
    done      chan bool
    watchCnt  int
    mu        sync.Mutex
    host      string
}

// NewWatcher creates a new file watcher
func NewWatcher(cfg *core.Config, sender *core.Sender) (*Watcher, error) {
    w, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }
    
    delay := 5 * time.Second
    if cfg.Debounce != "" {
        if d, err := time.ParseDuration(cfg.Debounce); err == nil {
            delay = d
        }
    }
    
    host, _ := os.Hostname()
    
    return &Watcher{
        watcher: w,
        config:  cfg,
        sender:  sender,
        timers:  make(map[string]*time.Timer),
        delay:   delay,
        done:    make(chan bool),
        host:    host,
    }, nil
}

// Start begins watching directories
func (w *Watcher) Start() error {
    fmt.Println("\n👁️  Starting real-time monitor...")
    
    for _, root := range w.config.WatchRoots {
        abs, _ := filepath.Abs(root)
        if !w.dirExists(abs) {
            fmt.Printf("⚠️  Skipping missing: %s\n", abs)
            continue
        }
        
        cnt, _ := w.walkAndWatch(abs, 0)
        w.watchCnt += cnt
        fmt.Printf("   %s: %d directories\n", abs, cnt)
    }
    
    fmt.Printf("\n✅ Watching %d directories\n", w.watchCnt)
    fmt.Println("   Press Ctrl+C to stop\n")
    
    go w.loop()
    return nil
}

func (w *Watcher) walkAndWatch(root string, depth int) (int, error) {
    cnt := 0
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil || !info.IsDir() {
            return nil
        }
        
        rel, _ := filepath.Rel(root, path)
        curDepth := depth + strings.Count(rel, string(os.PathSeparator))
        
        if w.config.MaxDepth > 0 && curDepth > w.config.MaxDepth {
            return filepath.SkipDir
        }
        
        if w.shouldSkip(path) {
            return filepath.SkipDir
        }
        
        if err := w.watcher.Add(path); err == nil {
            cnt++
        }
        return nil
    })
    return cnt, err
}

func (w *Watcher) shouldSkip(path string) bool {
    low := strings.ToLower(filepath.ToSlash(path))
    for _, ex := range w.config.ExcludeDirs {
        if strings.Contains(low, strings.ToLower(filepath.ToSlash(ex))) {
            return true
        }
    }
    return false
}

func (w *Watcher) loop() {
    for {
        select {
        case e, ok := <-w.watcher.Events:
            if !ok {
                return
            }
            w.handle(e)
            
        case err, ok := <-w.watcher.Errors:
            if !ok {
                return
            }
            fmt.Printf("⚠️ Watcher error: %v\n", err)
            
        case <-w.done:
            return
        }
    }
}

func (w *Watcher) handle(e fsnotify.Event) {
    if w.shouldSkip(e.Name) {
        return
    }
    
    if e.Op&fsnotify.Create == fsnotify.Create {
        if info, err := os.Stat(e.Name); err == nil && info.IsDir() {
            w.watcher.Add(e.Name)
            fmt.Printf("📁 Now watching: %s\n", e.Name)
            return
        }
    }
    
    w.debounce(e)
}

func (w *Watcher) debounce(e fsnotify.Event) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if t, ok := w.timers[e.Name]; ok {
        t.Stop()
    }
    
    w.timers[e.Name] = time.AfterFunc(w.delay, func() {
        w.send(e)
        w.mu.Lock()
        delete(w.timers, e.Name)
        w.mu.Unlock()
    })
}

func (w *Watcher) send(e fsnotify.Event) {
    typ := w.eventType(e)
    if typ == "" {
        return
    }
    
    var size int64
    var mod time.Time
    var isDir bool
    
    if info, err := os.Stat(e.Name); err == nil {
        size = info.Size()
        mod = info.ModTime()
        isDir = info.IsDir()
    }
    
    ev := core.FileEvent{
        Time:     time.Now(),
        Host:     w.host,
        Path:     e.Name,
        Type:     typ,
        Size:     size,
        IsDir:    isDir,
        Modified: mod,
    }
    
    if typ == core.Modified && !isDir {
        if hash, err := hashFile(e.Name); err == nil {
            ev.NewHash = hash
        }
    }
    
    fmt.Printf("%s %s\n", icon(typ), e.Name)
    w.sender.Queue(ev)
}

func (w *Watcher) eventType(e fsnotify.Event) core.ChangeType {
    switch {
    case e.Op&fsnotify.Create == fsnotify.Create:
        return core.Added
    case e.Op&fsnotify.Write == fsnotify.Write:
        return core.Modified
    case e.Op&fsnotify.Remove == fsnotify.Remove:
        return core.Deleted
    case e.Op&fsnotify.Rename == fsnotify.Rename:
        return core.Renamed
    }
    return ""
}

func (w *Watcher) dirExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}

// Stop shuts down the watcher
func (w *Watcher) Stop() error {
    fmt.Println("\n🛑 Stopping...")
    w.done <- true
    w.sender.Flush()
    return w.watcher.Close()
}

func icon(t core.ChangeType) string {
    switch t {
    case core.Added:
        return "➕"
    case core.Modified:
        return "✏️"
    case core.Deleted:
        return "❌"
    case core.Renamed:
        return "📤"
    }
    return "•"
}

func hashFile(path string) (string, error) {
    // Reuse hashFile from scanner, or import it
    // For brevity, assuming it's available
    return "", nil
}