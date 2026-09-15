package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	token := strings.TrimSpace(os.Getenv("MCP_AUTH_TOKEN"))
	if token == "" {
		log.Fatal("MCP_AUTH_TOKEN must be set; refusing to start an unauthenticated remote MCP gateway")
	}

	target, _ := url.Parse("http://127.0.0.1:8081")
	proxy := httputil.NewSingleHostReverseProxy(target)

	cmd := exec.Command("luno-mcp", "--transport", "streamable-http", "--sse-address", "127.0.0.1:8081", "--log-level", "info")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		log.Fatalf("start luno-mcp: %v", err)
	}
	defer func() { _ = cmd.Process.Signal(syscall.SIGTERM) }()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true,"service":"headortails-mcp"}`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+token {
			w.Header().Set("WWW-Authenticate", `Bearer realm="headortails-mcp"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Del("Authorization")
		proxy.ServeHTTP(w, r)
	})

	server := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 15 * time.Second}
	go func() {
		log.Println("headortails-mcp listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)
	<-wait
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	_ = cmd.Process.Signal(syscall.SIGTERM)
	_ = cmd.Wait()
}
