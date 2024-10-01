package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type IPDetails struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (used by proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, we need the first one
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// If X-Forwarded-For is not present, use the remote address
	return r.RemoteAddr
}

func visitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// Get the real IP (from X-Forwarded-For or RemoteAddr)
	ip := getIP(r)

	// Get the user agent
	userAgent := r.Header.Get("User-Agent")

	ipDetails := IPDetails{
		IP:        ip,
		UserAgent: userAgent,
	}

	// Log the IP and User-Agent
	fmt.Printf("IP: %s, User-Agent: %s\n", ip, userAgent)

	// Send the IP details as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ipDetails)
}

func main() {
	http.HandleFunc("/", visitHandler)
	port := 3333
	fmt.Printf("Server running on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
