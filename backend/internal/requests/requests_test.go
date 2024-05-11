package requests

import (
	"github.com/Ryan-Har/site-monitor/backend/internal/models"
	"strings"
	"testing"
	"time"
)

func TestSendRequest(t *testing.T) {
	t.Run("ICMP request", func(t *testing.T) {
		req := models.MakeRequest{
			URL:  "example.com",
			Type: models.RequestTypeICMP,
		}
		resp, err := SendRequest(req)
		if err != nil {
			t.Errorf("SendRequest failed: %v", err)
		}
		if !resp.Up {
			t.Errorf("Expected Up to be true for ICMP request")
		}
	})

	t.Run("TCP request", func(t *testing.T) {
		req := models.MakeRequest{
			URL:     "example.com",
			Type:    models.RequestTypeTCP,
			Port:    80,
			Timeout: 5 * time.Second,
		}
		resp, err := SendRequest(req)
		if err != nil {
			t.Errorf("SendRequest failed: %v", err)
		}
		if !resp.Up {
			t.Errorf("Expected Up to be true for TCP request")
		}
	})

	t.Run("HTTP request", func(t *testing.T) {
		req := models.MakeRequest{
			URL:     "https://example.com",
			Type:    models.RequestTypeHTTP,
			Timeout: 5 * time.Second,
		}
		resp, err := SendRequest(req)
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

	t.Run("Context timeout", func(t *testing.T) {
		req := models.MakeRequest{
			URL:     "https://example.com",
			Type:    models.RequestTypeHTTP,
			Timeout: 5 * time.Nanosecond,
		}
		resp, err := SendRequest(req)
		if err == nil {
			t.Error("Context timeout should have happened. No error was produced")
		}
		if !strings.HasSuffix(err.Error(), "context deadline exceeded") {
			t.Error("Context timeout check failed", err.Error())
		}
		if resp.Up {
			t.Errorf("Expected Up to be false for context timeout")
		}
	})

	t.Run("Wrong Request Type", func(t *testing.T) {
		req := models.MakeRequest{
			URL:     "https://example.com",
			Type:    "NOTRIGHT",
			Timeout: 5 * time.Nanosecond,
		}
		resp, err := SendRequest(req)
		if err == nil {
			t.Error("Wrong request type used and no error was produced")
		}
		if err.Error() != "unsupported request type" {
			t.Error("Wrong request type failed")
		}
		if resp.Up {
			t.Errorf("Expected Up to be false for wrong request")
		}
	})
}
