package core

import (
    "bytes"
    "context"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/opensearch-project/opensearch-go/v2"
    "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

// Sender sends events to OpenSearch
type Sender struct {
    client  *opensearch.Client
    index   string
    buffer  []FileEvent
    maxSize int
    verbose bool
}

// NewSender creates a new OpenSearch sender
func NewSender(cfg *OpenSearchConfig, verbose bool, batchSize int) (*Sender, error) {
    client, err := opensearch.NewClient(opensearch.Config{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        Addresses: []string{cfg.URL},
        Username:  cfg.Username,
        Password:  cfg.Password,
    })
    if err != nil {
        return nil, err
    }
    
    // Test connection
    if _, err := client.Info(); err != nil {
        return nil, fmt.Errorf("can't connect: %v", err)
    }
    
    fmt.Printf("✅ Connected to OpenSearch: %s\n", cfg.URL)
    
    return &Sender{
        client:  client,
        index:   cfg.Index,
        buffer:  []FileEvent{},
        maxSize: batchSize,
        verbose: verbose,
    }, nil
}

// Send sends events immediately
func (s *Sender) Send(events []FileEvent) error {
    if len(events) == 0 {
        return nil
    }
    
    var buf bytes.Buffer
    for _, e := range events {
        meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" } }%s`, s.index, "\n"))
        buf.Write(meta)
        
        data, _ := json.Marshal(e)
        buf.Write(data)
        buf.Write([]byte("\n"))
    }
    
    req := opensearchapi.BulkRequest{Body: &buf}
    resp, err := req.Do(context.Background(), s.client)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.IsError() {
        return fmt.Errorf("bulk error: %s", resp.String())
    }
    
    if s.verbose {
        fmt.Printf("📤 Sent %d events\n", len(events))
    }
    return nil
}

// Queue adds an event to buffer, sends when full
func (s *Sender) Queue(event FileEvent) error {
    s.buffer = append(s.buffer, event)
    
    if len(s.buffer) >= s.maxSize {
        return s.Flush()
    }
    return nil
}

// Flush sends all buffered events
func (s *Sender) Flush() error {
    if len(s.buffer) == 0 {
        return nil
    }
    err := s.Send(s.buffer)
    s.buffer = []FileEvent{}
    return err
}

// CreateIndex sets up the index with proper mappings
func (s *Sender) CreateIndex() error {
    exists := opensearchapi.IndicesExistsRequest{Index: []string{s.index}}
    resp, err := exists.Do(context.Background(), s.client)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
        return nil // already exists
    }
    
    mapping := `{
        "mappings": {
            "properties": {
                "timestamp": { "type": "date" },
                "host": { "type": "keyword" },
                "file_path": { "type": "keyword" },
                "event_type": { "type": "keyword" }
            }
        }
    }`
    
    create := opensearchapi.IndicesCreateRequest{
        Index: s.index,
        Body:  bytes.NewReader([]byte(mapping)),
    }
    
    resp, err = create.Do(context.Background(), s.client)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    fmt.Printf("✅ Created index: %s\n", s.index)
    return nil
}