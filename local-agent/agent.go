package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Rarely used port:
const AgentPort = "42999"

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Local Agent is running!")
	})

	// Dynamic proxy endpoint
	// e.g. /proxy?ip=192.168.200.42&port=8080&proto=http
	// Then the request path & query beyond /proxy will be forwarded to that IP/port
	mux.HandleFunc("/proxy/", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)

		// Extract IP, port, protocol from query params (with defaults)
		q := r.URL.Query()

		targetIP := q.Get("ip")
		if targetIP == "" {
			http.Error(w, "Missing 'ip' query param", http.StatusBadRequest)
			return
		}

		targetPort := q.Get("port")
		if targetPort == "" {
			targetPort = "8080" // default
		}

		proto := q.Get("proto")
		if proto == "" {
			proto = "http"
		}

		// Build the target base URL
		targetURLStr := fmt.Sprintf("%s://%s:%s", proto, targetIP, targetPort)
		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			http.Error(w, "Invalid target URL: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Create a reverse proxy for this target
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Optionally adjust the request path
		// The portion after "/proxy" will be appended to the base
		// e.g. /proxy/register?ip=192.168...
		// means final path is "/register"
		originalPath := r.URL.Path[len("/proxy"):]
		r.URL.Path = originalPath

		// Let the proxy handle the request
		proxy.ServeHTTP(w, r)
	})

	serverAddr := ":" + AgentPort
	fmt.Printf("Local Agent listening on http://localhost%s\n", serverAddr)
	fmt.Printf("Try: http://localhost:%s/ping\n", AgentPort)

	log.Fatal(http.ListenAndServe(serverAddr, mux))
}

// enableCORS adds headers so the browser allows cross-origin requests
func enableCORS(w http.ResponseWriter, r *http.Request) {
	// Adjust or remove as needed
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	}
}
