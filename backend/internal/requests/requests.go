package requests

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/Ryan-Har/site-monitor/backend/internal/models"
	"github.com/prometheus-community/pro-bing"
)

func SendRequest(req models.MakeRequest) (models.MakeRequestResponse, error) {
	switch req.Type {
	case models.RequestTypeICMP:
		return sendICMPRequest(req)
	case models.RequestTypeTCP:
		return sendTCPRequest(req)
	case models.RequestTypeHTTP:
		return sendHTTPRequest(req)
	default:
		return models.MakeRequestResponse{}, errors.New("unsupported request type")
	}
}

func sendICMPRequest(mr models.MakeRequest) (models.MakeRequestResponse, error) {
	runTime := time.Now()
	resp := models.MakeRequestResponse{
		OrigRequest: mr,
		RunTime:     runTime,
	}
	pinger, err := probing.NewPinger(mr.URL)
	if err != nil {
		return resp, fmt.Errorf("error creating pinger for %v: %v", mr.URL, err.Error())
	}
	pinger.Count = 1

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return resp, fmt.Errorf("error pinging %v: %v", mr.URL, err.Error())
	}
	if pinger.Statistics().PacketLoss > 0 {
		return resp, fmt.Errorf("error pinging %v: packet loss", mr.URL)
	}

	resp.Up = true
	resp.ResponseTime = pinger.Statistics().AvgRtt

	return resp, nil
}

func sendHTTPRequest(mr models.MakeRequest) (models.MakeRequestResponse, error) {
	runTime := time.Now()
	resp := models.MakeRequestResponse{
		OrigRequest: mr,
		RunTime:     runTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), mr.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", mr.URL, nil)
	if err != nil {
		return resp, fmt.Errorf("error creating new http request: %v", err.Error())
	}

	client := &http.Client{}

	start := time.Now()
	res, err := client.Do(req)
	end := time.Now()
	if err != nil {
		return resp, fmt.Errorf("error sending http request: %v", err.Error())
	}

	resp.Up = true
	resp.ResponseTime = end.Sub(start)
	resp.HttpResponseCode = res.StatusCode

	return resp, nil
}

func sendTCPRequest(mr models.MakeRequest) (models.MakeRequestResponse, error) {
	runTime := time.Now()
	resp := models.MakeRequestResponse{
		OrigRequest: mr,
		RunTime:     runTime,
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(mr.URL, fmt.Sprint(mr.Port)), mr.Timeout)
	end := time.Now()
	if err != nil {
		return resp, fmt.Errorf("error sending tcp request for %v on port %d: %v", mr.URL, mr.Port, err.Error())
	}
	defer conn.Close()
	resp.Up = true
	resp.ResponseTime = end.Sub(start)

	return resp, nil
}
