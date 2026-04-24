package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
	"github.com/mvaldes14/task-manager-mcp/internal/tools"
)

func main() {
	baseURL := os.Getenv("DOIT_BASE_URL")
	if baseURL == "" {
		log.Fatal("DOIT_BASE_URL is not set")
	}

	apiKey := os.Getenv("DOIT_API_KEY")
	if apiKey == "" {
		log.Fatal("DOIT_API_KEY is not set")
	}

	addr := os.Getenv("MCP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	c := client.New(baseURL, apiKey)

	s := server.NewMCPServer("doit-mcp", "0.1.0")

	tools.RegisterTaskTools(s, c)
	tools.RegisterProjectTools(s, c)
	tools.RegisterNLPTools(s, c)
	tools.RegisterSubtaskTools(s, c)
	tools.RegisterAITools(s, c)

	h := server.NewStreamableHTTPServer(s)

	mux := http.NewServeMux()
	mux.Handle("/mcp", logMiddleware(h))

	srv := &http.Server{Addr: addr, Handler: mux}
	log.Printf("MCP HTTP server listening on %s/mcp", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.code = code
	sw.ResponseWriter.WriteHeader(code)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, "server-to-client streaming not supported", http.StatusMethodNotAllowed)
			log.Printf("GET %s status=405 (SSE not supported)", r.URL.Path)
			return
		}

		start := time.Now()

		method := ""
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
			var rpc struct {
				Method string `json:"method"`
			}
			if json.Unmarshal(body, &rpc) == nil {
				method = rpc.Method
			}
		}

		rw := &statusWriter{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("%s %s rpc=%s status=%d dur=%s", r.Method, r.URL.Path, method, rw.code, time.Since(start).Round(time.Millisecond))
	})
}
