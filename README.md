# FIM-Agent
A lightweight, real-time File Integrity Monitoring (FIM) agent that tracks file changes across your system and sends events to OpenSearch for analysis and visualization.

📦 Prerequisites
Go 1.21+ - Download
OpenSearch 2.x - Download or use Docker

OpenSearch Dashboards (optional, for visualization)

# Quick OpenSearch Setup with Docker
bash
docker run -d -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" opensearchproject/opensearch:latest
docker run -d -p 5601:5601 --link <opensearch-container> opensearchproject/opensearch-dashboards:latest


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
-for windows configure config.json file by making appropriate adjusments
-for linux configure config_kali.json file

# Configuration Options
Option	Type	Default	Description
watch_roots:	Directories to monitor
exclude_dirs: Skip paths containing these strings
max_depth	int	10	: Maximum subdirectory depth (0 = unlimited)
debounce_time	string	"5s"	: Wait time before sending events
batch_size	int	50	: Events per batch to OpenSearch
opensearch.url	string	"http://localhost:9200"	: OpenSearch endpoint
opensearch.username	string	"admin"	: OpenSearch username
opensearch.password	string	"admin"	: OpenSearch password
opensearch.index	string	"fim-events"	: OpenSearch index name
state_file	string	"fim-state.json"	: File to store previous state
verbose	bool	true	: Show detailed output


🎮 Usage
# One-Time Scan Mode
- Scan once and exit
       go run main.go -mode scan -config config.json
  
# Continuous Monitoring Mode (Default)
-Run continuously until Ctrl+C
           go run main.go -mode monitor -config config.json
           
# Run as Executable
- Build and run
go build -o fim-agent main.go
./fim-agent -mode monitor -config config.json

# for linux config.json = config_kali.json

# Expected Output:
🚀 FIM Agent Starting...
✅ Connected to OpenSearch: opensearch_url

👁️  Starting real-time monitor...
   C:\Users\johnn\Desktop: 1 directories
   C:\Users\johnn\Documents: 28 directories

✅ Watching 29 directories
   Press Ctrl+C to stop

✏️ C:\Users\jo\Desktop\notes.txt
📤 Sent 1 events to OpenSearch
