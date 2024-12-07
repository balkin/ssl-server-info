package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// Root redirect handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/balkin/ssl-server-info", http.StatusFound)
	})
	http.HandleFunc("/ssl/json", SslJsonHandler)
	http.HandleFunc("/ssl/json/", SslJsonHandler)

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
