🔍 FIM-Agent

A lightweight, real-time File Integrity Monitoring (FIM) agent that tracks file changes across your system and sends events to OpenSearch for analysis and visualization.

📦 Prerequisites

Go 1.21+ – Download

OpenSearch 2.x – Download or use Docker

OpenSearch Dashboards (optional, for visualization)

⚡ Quick OpenSearch Setup with Docker

docker run -d -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" opensearchproject/opensearch:latest

docker run -d -p 5601:5601 --link <opensearch-container> opensearchproject/opensearch-dashboards:latest

🚀 Installation

Clone the Repository

git clone https://github.com/yourusername/fim-agent.git
cd fim-agent

Install Dependencies

go mod tidy

Build from Source

# Build for your current platform
go build -o fim-agent main.go

# Cross-compile for Linux from Windows
GOOS=linux GOARCH=amd64 go build -o fim-agent-linux main.go

# Cross-compile for Windows from Linux
GOOS=windows GOARCH=amd64 go build -o fim-agent.exe main.go

⚙️ Configuration

Windows → configure config.json

Linux → configure config_kali.json

Configuration Options

Option

Type

Default

Description

watch_roots

array

—

Directories to monitor

exclude_dirs

array

—

Skip paths containing these strings

max_depth

int

10

Maximum subdirectory depth (0 = unlimited)

debounce_time

string

"5s"

Wait time before sending events

batch_size

int

50

Events per batch to OpenSearch

opensearch.url

string

"http://localhost:9200"

OpenSearch endpoint

opensearch.username

string

"admin"

OpenSearch username

opensearch.password

string

"admin"

OpenSearch password

opensearch.index

string

"fim-events"

OpenSearch index name

state_file

string

"fim-state.json"

File to store previous state

verbose

bool

true

Show detailed output

👾 Usage

One-Time Scan Mode

Scan once and exit:

go run main.go -mode scan -config config.json

Continuous Monitoring Mode (Default)

Run continuously until Ctrl+C:

go run main.go -mode monitor -config config.json

Run as Executable

go build -o fim-agent main.go
./fim-agent -mode monitor -config config.json

💡 For Linux, use config_kali.json instead of config.json.

✅ Expected Output

🚀 FIM Agent Starting...
✅ Connected to OpenSearch: opensearch_url

👁️  Starting real-time monitor...
   C:\Users\johnn\Desktop: 1 directories
   C:\Users\johnn\Documents: 28 directories

---# 🔍 FIM-Agent

A lightweight, real-time File Integrity Monitoring (FIM) agent that tracks file changes across your system and sends events to OpenSearch for analysis and visualization.

📦 Prerequisites

Go 1.21+ – Download

OpenSearch 2.x – Download or use Docker

OpenSearch Dashboards (optional, for visualization)

⚡ Quick OpenSearch Setup with Docker

docker run -d -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" opensearchproject/opensearch:latest

docker run -d -p 5601:5601 --link <opensearch-container> opensearchproject/opensearch-dashboards:latest

🚀 Installation

Clone the Repository

git clone https://github.com/yourusername/fim-agent.git
cd fim-agent

Install Dependencies

go mod tidy

Build from Source

# Build for your current platform
go build -o fim-agent main.go

# Cross-compile for Linux from Windows
GOOS=linux GOARCH=amd64 go build -o fim-agent-linux main.go

# Cross-compile for Windows from Linux
GOOS=windows GOARCH=amd64 go build -o fim-agent.exe main.go

⚙️ Configuration

Windows → configure config.json

Linux → configure config_kali.json

Configuration Options

Option

Type

Default

Description

watch_roots

array

—

Directories to monitor

exclude_dirs

array

—

Skip paths containing these strings

max_depth

int

10

Maximum subdirectory depth (0 = unlimited)

debounce_time

string

"5s"

Wait time before sending events

batch_size

int

50

Events per batch to OpenSearch

opensearch.url

string

"http://localhost:9200"

OpenSearch endpoint

opensearch.username

string

"admin"

OpenSearch username

opensearch.password

string

"admin"

OpenSearch password

opensearch.index

string

"fim-events"

OpenSearch index name

state_file

string

"fim-state.json"

File to store previous state

verbose

bool

true

Show detailed output

👾 Usage

One-Time Scan Mode

Scan once and exit:

go run main.go -mode scan -config config.json

Continuous Monitoring Mode (Default)

Run continuously until Ctrl+C:

go run main.go -mode monitor -config config.json

Run as Executable

go build -o fim-agent main.go
./fim-agent -mode monitor -config config.json

💡 For Linux, use config_kali.json instead of config.json.

✅ Expected Output

🚀 FIM Agent Starting...
✅ Connected to OpenSearch: opensearch_url

👁️  Starting real-time monitor...
   C:\Users\johnn\Desktop: 1 directories
   C:\Users\johnn\Documents: 28 directories

---
