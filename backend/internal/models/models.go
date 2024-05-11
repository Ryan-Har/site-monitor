package models

import (
	"time"
)

// Request defines the information for a request
type MakeRequest struct {
	URL     string
	Type    RequestType   // enum (ICMP, TCP, HTTP)
	Port    int           // optional for TCP requests
	IDs     []int         // slice of IDs for results
	Timeout time.Duration // context timeout in seconds
}

// RequestType defines the supported request types
type RequestType string

const (
	RequestTypeICMP RequestType = "ICMP"
	RequestTypeTCP  RequestType = "TCP"
	RequestTypeHTTP RequestType = "HTTP"
)

type MakeRequestResponse struct {
	OrigRequest      MakeRequest // original request that triggered the response
	Up               bool
	ResponseTime     time.Duration // valid if up
	RunTime          time.Time
	HttpResponseCode int // response code for HTTP requests
}
