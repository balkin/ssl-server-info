package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// parsePathParams now handles both path parameters and query parameters
func parsePathParams(path string, query url.Values) map[string]string {
	params := map[string]string{}
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Handle path parameters like /ssl/json/delay/1000 or /ssl/json/statusCode/222
	for i := 0; i < len(segments)-1; i += 2 {
		if i+1 < len(segments) {
			params[segments[i]] = segments[i+1]
		}
	}

	// Handle query parameters like ?delay=1000 or ?statusCode=222
	for key, values := range query {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params
}

func SslJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "ssl-server-info; https://github.com/balkin/ssl-server-info")
	w.Header().Add("Content-Type", "application/json")

	if r.TLS == nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse path parameters and query parameters
	pathParams := parsePathParams(r.URL.Path[len("/ssl/json"):], r.URL.Query())

	response := Response{
		Server:           "github.com/balkin/ssl-server-info",
		HTTPS:            "on",
		ContentType:      r.Header.Get("Content-Type"),
		Accept:           r.Header.Get("Accept"),
		UserAgent:        r.Header.Get("User-Agent"),
		HttpHost:         r.Host,
		Connection:       r.Header.Get("Connection"),
		ServerAddr:       r.RemoteAddr,
		RequestProtocol:  r.Proto,
		RequestMethod:    r.Method,
		RequestUri:       r.RequestURI,
		RequestTimestamp: time.Now().Unix(),
	}

	// Handle delay
	if delay, ok := pathParams["delay"]; ok {
		if delayValue, err := strconv.Atoi(delay); err == nil {
			time.Sleep(time.Duration(delayValue) * time.Millisecond)
		}
	}

	// Handle statusCode
	statusCode := 200
	if code, ok := pathParams["statusCode"]; ok {
		if statusCodeValue, err := strconv.Atoi(code); err == nil {
			statusCode = statusCodeValue
		}
	}

	if len(r.TLS.PeerCertificates) > 0 {
		cert := r.TLS.PeerCertificates[0]
		response.SslSubject = cert.Subject.CommonName
		response.SslIssuer = cert.Issuer.CommonName
		response.SslNotBefore = cert.NotBefore.String()
		response.SslNotAfter = cert.NotAfter.String()
		w.WriteHeader(statusCode)
	} else {
		response.Message = "No mTLS certificate provided"
		log.Println("No peer certificates, returning 403")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
