# FIM-Agent 🔍

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/dl/)
[![OpenSearch](https://img.shields.io/badge/OpenSearch-2.x-005EB8?style=flat&logo=opensearch)](https://opensearch.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A lightweight, real-time File Integrity Monitoring (FIM) agent that tracks file changes across your system and sends events to OpenSearch for analysis and visualization.

## ✨ Features

- 🔄 **Real-time monitoring** - Instant detection of file changes
- 📊 **OpenSearch integration** - Centralized logging and visualization
- ⚡ **Lightweight** - Minimal system resource usage
- 🔧 **Configurable** - Flexible monitoring options
- 🖥️ **Cross-platform** - Works on Windows, Linux, and macOS

## 📋 Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [OpenSearch 2.x](https://opensearch.org/downloads.html) (or use Docker)
- [OpenSearch Dashboards](https://opensearch.org/downloads.html) (optional, for visualization)

## 🐳 Quick OpenSearch Setup with Docker

```bash
# Start OpenSearch
docker run -d -p 9200:9200 -p 9600:9600 \
  -e "discovery.type=single-node" \
  --name opensearch \
  opensearchproject/opensearch:latest

# Start OpenSearch Dashboards
docker run -d -p 5601:5601 \
  --link opensearch \
  opensearchproject/opensearch-dashboards:latest
🚀 Installation
Clone the Repository
bash
git clone https://github.com/yourusername/fim-agent.git
cd fim-agent
Install Dependencies
bash
go mod tidy
Build from Source
bash
# Build for your current platform
go build -o fim-agent main.go

# Cross-compile for Linux from Windows
GOOS=linux GOARCH=amd64 go build -o fim-agent-linux main.go

# Cross-compile for Windows from Linux
GOOS=windows GOARCH=amd64 go build -o fim-agent.exe main.go
⚙️ Configuration
Create a config.json file (Windows) or config_kali.json (Linux) with the following structure:

json
{
  "watch_roots": [
    "C:\\Users\\username\\Documents",
    "C:\\Users\\username\\Desktop"
  ],
  "exclude_dirs": [
    "\\.git",
    "node_modules",
    "temp"
  ],
  "max_depth": 10,
  "debounce_time": "5s",
  "batch_size": 50,
  "opensearch": {
    "url": "http://localhost:9200",
    "username": "admin",
    "password": "admin",
    "index": "fim-events"
  },
  "state_file": "fim-state.json",
  "verbose": true
}
Configuration Options
Option	Type	Default	Description
watch_roots	array	required	Directories to monitor
exclude_dirs	array	[]	Skip paths containing these strings
max_depth	int	10	Maximum subdirectory depth (0 = unlimited)
debounce_time	string	"5s"	Wait time before sending events
batch_size	int	50	Events per batch to OpenSearch
opensearch.url	string	"http://localhost:9200"	OpenSearch endpoint
opensearch.username	string	"admin"	OpenSearch username
opensearch.password	string	"admin"	OpenSearch password
opensearch.index	string	"fim-events"	OpenSearch index name
state_file	string	"fim-state.json"	File to store previous state
verbose	bool	true	Show detailed output
🎮 Usage
One-Time Scan Mode
bash
# Windows
go run main.go -mode scan -config config.json

# Linux
go run main.go -mode scan -config config_kali.json
Continuous Monitoring Mode (Default)
bash
# Windows
go run main.go -mode monitor -config config.json

# Linux
go run main.go -mode monitor -config config_kali.json
Run as Executable
bash
# Build and run
go build -o fim-agent main.go

# Windows
./fim-agent -mode monitor -config config.json

# Linux
./fim-agent -mode monitor -config config_kali.json
📊 Expected Output
text
🚀 FIM Agent Starting...
✅ Connected to OpenSearch: http://localhost:9200

👁️  Starting real-time monitor...
   C:\Users\johnn\Desktop: 1 directories
   C:\Users\johnn\Documents: 28 directories

✅ Watching 29 directories
   Press Ctrl+C to stop

✏️ C:\Users\johnn\Desktop\notes.txt
📤 Sent 1 events to OpenSearch
