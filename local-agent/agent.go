package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const AgentPort = "42999"

func main() {
    mux := http.NewServeMux()

    // Health check endpoint
    mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w, r)
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Local Agent is running!")
    })

    // Proxy endpoint
    mux.HandleFunc("/proxy/", handleProxy)

    server := &http.Server{
        Addr:         ":" + AgentPort,
        Handler:      mux,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        // Optional: set a TLSConfig if you want agent to serve HTTPS.
    }

    log.Printf("Local Agent listening on http://localhost:%s\n", AgentPort)
    log.Printf("Try: http://localhost:%s/ping\n", AgentPort)
    log.Fatal(server.ListenAndServe())
}

// handleProxy sets up the ReverseProxy with custom director and modifies the response.
func handleProxy(w http.ResponseWriter, r *http.Request) {
    enableCORS(w, r)

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    q := r.URL.Query()
    targetIP := q.Get("ip")
    if targetIP == "" {
        http.Error(w, "Missing 'ip' parameter", http.StatusBadRequest)
        return
    }

    targetPort := q.Get("port")
    if targetPort == "" {
        targetPort = "8080"
    }

    proto := q.Get("proto")
    if proto == "" {
        proto = "http"
    }

    // Construct base target, e.g.: "https://192.168.100.73:8080"
    baseTarget := fmt.Sprintf("%s://%s:%s", proto, targetIP, targetPort)
    targetURL, err := url.Parse(baseTarget)
    if err != nil {
        http.Error(w, "Invalid target URL: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Create a reverse proxy
    proxy := &httputil.ReverseProxy{
        Director: func(req *http.Request) {
            // Director is called before the request is sent to the backend (device).
            // We'll rewrite the URL’s scheme, host, and path to match the target device.
            req.URL.Scheme = targetURL.Scheme
            req.URL.Host = targetURL.Host

            // The portion after "/proxy" is appended to the base path
            // e.g. if the user requests "/proxy/register?ip=..."
            // final path -> "/register"
            originalPath := strings.TrimPrefix(req.URL.Path, "/proxy")
            if !strings.HasPrefix(originalPath, "/") {
                originalPath = "/" + originalPath
            }
            req.URL.Path = originalPath

            // Remove the local query params "ip, port, proto" so they don't get forwarded
            // to the device if that's undesired. But if you want them forwarded, remove this block.
            filteredQ := req.URL.Query()
            filteredQ.Del("ip")
            filteredQ.Del("port")
            filteredQ.Del("proto")
            req.URL.RawQuery = filteredQ.Encode()

            // Set Host header to the device host
            req.Host = targetURL.Host

            // Log the final request
            log.Printf("[Proxy Request] %s %s", req.Method, req.URL.String())
        },
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
            Proxy: http.ProxyFromEnvironment,
            // Optional timeouts, keep-alives, etc.
            DialContext: (&net.Dialer{
                Timeout:   15 * time.Second,
                KeepAlive: 15 * time.Second,
            }).DialContext,
            ForceAttemptHTTP2:     true,
            MaxIdleConns:          100,
            IdleConnTimeout:       90 * time.Second,
            TLSHandshakeTimeout:   10 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
        },
        ModifyResponse: func(resp *http.Response) error {
            // ModifyResponse is called AFTER the device responds, but BEFORE we send back to client.
            // We can log the status or do custom changes. By default, we do nothing special—
            // the ReverseProxy will forward the exact status, headers, and body to the frontend.
            log.Printf("[Proxy Response] %d %s", resp.StatusCode, resp.Request.URL.String())
            // If you want to read or modify the body, you can do so here.
            // e.g. read all bytes, parse JSON, etc. But typically not needed for a pass-through.
            return nil
        },
        ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
            // If there's a network error or dial error, we handle it here.
            log.Printf("[Proxy Error] %v", err)
            http.Error(rw, "Proxy error: "+err.Error(), http.StatusBadGateway)
        },
    }

    // Now run the proxy
    proxy.ServeHTTP(w, r)
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
