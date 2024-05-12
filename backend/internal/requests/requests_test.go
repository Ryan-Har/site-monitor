package requests

import (
	"strings"
	"testing"
	"time"
)

func TestSendRequest(t *testing.T) {
	t.Run("ICMP request", func(t *testing.T) {
		req := getPINGmakeRequest(1, 5*time.Second)
		resp, err := sendICMPRequest(req)
		if err != nil {
			t.Errorf("SendRequest failed: %v", err)
		}
		if !resp.Up {
			t.Errorf("Expected Up to be true for ICMP request")
		}
	})

	t.Run("TCP request", func(t *testing.T) {
		req := getTCPmakeRequest(1, 5*time.Second)
		resp, err := sendTCPRequest(req)
		if err != nil {
			t.Errorf("SendRequest failed: %v", err)
		}
		if !resp.Up {
			t.Errorf("Expected Up to be true for TCP request")
		}
	})

	t.Run("HTTP request", func(t *testing.T) {
		req := getHTTPmakeRequest(1, 5*time.Second)
		resp, err := sendHTTPRequest(req)
		if err != nil {
			t.Errorf("SendRequest failed: %v", err)
		}
		if !resp.Up {
			t.Errorf("Expected Up to be true for HTTP request")
		}
		if resp.HttpResponseCode != 200 {
			t.Errorf("Expected HttpResponseCode to be 200, got %d", resp.HttpResponseCode)
		}
	})

	t.Run("Test Sample", func(t *testing.T) {
		reqs := []Requests{
			{
				URL:     "https://example.com",
				RType:   RequestTypeHTTP,
				ID:      1,
				Timeout: 5 * time.Second,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeTCP,
				ID:      2,
				Timeout: 5 * time.Second,
				Port:    80,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeICMP,
				ID:      3,
				Timeout: 5 * time.Second,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeTCP,
				ID:      4,
				Timeout: 5 * time.Nanosecond,
				Port:    80,
			},
			{
				URL:     "::1",
				RType:   RequestTypeICMP,
				ID:      5,
				Timeout: 5 * time.Second,
			},
		}
		resp := Send(reqs)
		for _, item := range resp {
			if item.ID != 4 {
				if item.Err != nil {
					t.Errorf("no error should be produced, got: %v", item.Err.Error())
				}
				if !item.Up {
					t.Errorf("expected up response for id %d", item.ID)
				}
			} else {
				if item.Up {
					t.Errorf("expected to be down due to short timeout")
				}
			}

		}
	})

	t.Run("Error Test Sample", func(t *testing.T) {
		req := []Requests{
			{
				URL:     "https://example.com",
				RType:   "incorrect",
				ID:      1,
				Timeout: 5 * time.Second,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeTCP,
				ID:      2,
				Timeout: 5 * time.Second,
				Port:    81,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeTCP,
				ID:      3,
				Timeout: 5 * time.Nanosecond,
				Port:    80,
			},
			{
				URL:     "https://example.com",
				RType:   RequestTypeHTTP,
				ID:      4,
				Timeout: 5 * time.Microsecond,
			},
			{
				URL:     "example.com",
				RType:   RequestTypeHTTP,
				ID:      5,
				Timeout: 5 * time.Second,
			},
		}

		resp := Send(req)
		respMap := make(map[int]*Response)
		for _, r := range resp {
			respMap[r.ID] = &r
		}

		if respMap[1].Err == nil {
			t.Error("Wrong request type used and no error was produced")
		}
		if respMap[1].Err.Error() != "unsupported request type" {
			t.Error("Wrong request type failed")
		}
		if respMap[1].Up {
			t.Errorf("Expected Up to be false for wrong request")
		}
		if respMap[2].Err == nil {
			t.Error("should have errored due to incorrect port")
		}
		if respMap[2].Up {
			t.Errorf("Expected Up to be false due to incorrect port")
		}
		if respMap[3].Err == nil {
			t.Error("Context timeout should have happened. No error was produced")
		}
		if respMap[3].Up {
			t.Errorf("Expected Up to be false due to context timeout")
		}
		if respMap[4].Err == nil {
			t.Error("Context timeout should have happened. No error was produced")
		}
		if respMap[4].Up {
			t.Errorf("Expected Up to be false due to context timeout")
		}
		if !strings.HasSuffix(respMap[4].Err.Error(), "context deadline exceeded") {
			t.Error("Context timeout check failed", respMap[4].Err.Error())
		}
		if respMap[5].Err == nil {
			t.Error("unsupported protocol scheme should have happened. No error was produced")
		}
		if respMap[5].Up {
			t.Errorf("Expected Up to be false due to unsupported protocol scheme")
		}
		if !strings.HasSuffix(respMap[5].Err.Error(), `unsupported protocol scheme ""`) {
			t.Error("unsupported protocol scheme check failed", respMap[5].Err.Error())
		}

	})
}

func getHTTPmakeRequest(id int, timeout time.Duration) makeRequest {
	return makeRequest{
		url:        "https://example.com",
		rtype:      RequestTypeHTTP,
		maxTimeout: 5 * time.Second,
		idTimeout: map[int]time.Duration{
			id: timeout,
		},
	}
}

func getTCPmakeRequest(id int, timeout time.Duration) makeRequest {
	return makeRequest{
		url:        "example.com",
		rtype:      RequestTypeTCP,
		port:       80,
		maxTimeout: 5 * time.Second,
		idTimeout: map[int]time.Duration{
			id: timeout,
		},
	}
}

func getPINGmakeRequest(id int, timeout time.Duration) makeRequest {
	return makeRequest{
		url:        "example.com",
		rtype:      RequestTypeICMP,
		maxTimeout: 5 * time.Second,
		idTimeout: map[int]time.Duration{
			id: timeout,
		},
	}
}
