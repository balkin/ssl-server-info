package main

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
