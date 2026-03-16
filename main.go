package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "fim-agent/core"
    "fim-agent/monitor"
)

func main() {
    // Command line flags
    configPath := flag.String("config", "", "config file (optional)")
    mode := flag.String("mode", "monitor", "mode: scan|monitor")
    flag.Parse()
    
    fmt.Println("🚀 FIM Agent Starting...")
    
    // Load config
    cfg, err := core.LoadConfig(*configPath)
    if err != nil {
        log.Fatal("Config error:", err)
    }
    
    // Connect to OpenSearch
    sender, err := core.NewSender(&cfg.OpenSearch, cfg.Verbose, cfg.BatchSize)
    if err != nil {
        log.Fatal("OpenSearch connection failed:", err)
    }
    defer sender.Flush()
    
    sender.CreateIndex()
    
    if *mode == "scan" {
        runScanMode(cfg, sender)
    } else {
        runMonitorMode(cfg, sender)
    }
}

func runScanMode(cfg *core.Config, sender *core.Sender) {
    fmt.Println("\n🔍 Running one-time scan...")
    
    scanner := core.NewScanner(cfg)
    detector, _ := core.NewDetector(cfg.StateFile)
    
    result, _ := scanner.Scan()
    changes, _ := detector.Detect(scanner.LastScan())
    
    if len(changes) > 0 {
        sender.Send(changes)
    }
    
    fmt.Printf("\n✅ %s\n", result.Summary())
}

func runMonitorMode(cfg *core.Config, sender *core.Sender) {
    watcher, err := monitor.NewWatcher(cfg, sender)
    if err != nil {
        log.Fatal("Failed to start monitor:", err)
    }
    
    if err := watcher.Start(); err != nil {
        log.Fatal(err)
    }
    
    // Wait for Ctrl+C
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    
    watcher.Stop()
    fmt.Println("\n👋 Goodbye!")
}