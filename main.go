package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func parsePathParams(path string) map[string]string {
	params := map[string]string{}
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Handle paths like /ssl/json/delay/1000 or /ssl/json/statusCode/222
	for i := 0; i < len(segments)-1; i += 2 {
		if i+1 < len(segments) {
			params[segments[i]] = segments[i+1]
		}
	}
	return params
}

type Response struct {
	Message          string `json:"message,omitempty"`
	Server           string `json:"server,omitempty"`
	HTTPS            string `json:"https,omitempty"`
	ContentType      string `json:"headerContentType,omitempty"`
	Accept           string `json:"headerAccept,omitempty"`
	UserAgent        string `json:"headerUserAgent,omitempty"`
	Connection       string `json:"headerConnection,omitempty"`
	HttpHost         string `json:"httpHost,omitempty"`
	ServerAddr       string `json:"httpServerAddr,omitempty"`
	RequestProtocol  string `json:"requestProtocol,omitempty"`
	RequestMethod    string `json:"requestMethod,omitempty"`
	RequestUri       string `json:"requestUri,omitempty"`
	RequestTimestamp int64  `json:"requestTimestamp,omitempty"`
	SslSubject       string `json:"sslSubject,omitempty"`
	SslIssuer        string `json:"sslIssuer,omitempty"`
	SslNotBefore     string `json:"sslNotBefore,omitempty"`
	SslNotAfter      string `json:"sslNotAfter,omitempty"`
}

func sslHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "ssl-server-info; https://github.com/balkin/ssl-server-info")

	if r.TLS == nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	pathParams := parsePathParams(r.URL.Path[len("/ssl/json"):])
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/ssl/json", sslHandler)
	http.HandleFunc("/ssl/json/", sslHandler)

	certFile := os.Getenv("CERT_FILE")
	if certFile == "" {
		certFile = "server.crt"
	}

	keyFile := os.Getenv("KEY_FILE")
	if keyFile == "" {
		keyFile = "server.key"
	}

	flag.StringVar(&certFile, "cert", certFile, "Path to the SSL certificate file")
	flag.StringVar(&keyFile, "key", keyFile, "Path to the SSL key file")
	port := flag.String("port", "443", "Port to listen on (default 443)")
	flag.Parse()

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("Certificate file not found: %s", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("Key file not found: %s", keyFile)
	}

	// Custom TLS configuration to trust all certificates
	tlsConfig := &tls.Config{
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: true,
	}

	server := &http.Server{
		Addr:      ":" + *port,
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting server on :%s with certFile: %s and keyFile: %s", *port, certFile, keyFile)
	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
