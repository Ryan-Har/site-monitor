package requests

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"net"
	"net/http"
	"sync"
	"time"
)

//this package is to be used with a slice of Requests,
//create a slice of Requests and Send(slice...)
//pro-bing ping needs system permissions for sockets
//sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"

type Requests struct {
	URL     string
	RType   RequestType   // enum (ICMP, TCP, HTTP)
	Port    int           // optional for TCP requests
	ID      int           // ID of the monitor
	Timeout time.Duration // context timeout in seconds
}

type Response struct {
	Requests
	CheckResponse
	Err error
}

type CheckResponse struct {
	Up               bool
	ResponseTime     time.Duration // valid if up
	RunTime          time.Time
	HttpResponseCode int // response code for HTTP requests
}

type makeRequest struct {
	url        string
	rtype      RequestType           // enum (ICMP, TCP, HTTP)
	port       int                   // optional for TCP requests
	idTimeout  map[int]time.Duration // map of id's to the timeout they requested - needed to filter later
	maxTimeout time.Duration         // max context timeout in seconds
}

// RequestType defines the supported request types
type RequestType string

const (
	RequestTypeICMP RequestType = "ICMP"
	RequestTypeTCP  RequestType = "TCP"
	RequestTypeHTTP RequestType = "HTTP"
)

type makeRequestResponseWithErr struct {
	makeRequestResponse
	err error
}

type makeRequestResponse struct {
	origRequest makeRequest // original request that triggered the response
	CheckResponse
}

func Send(req ...Requests) []Response {
	organised := organiseMonitorRequests(req...)

	respChan := make(chan makeRequestResponseWithErr, len(organised))
	var wg sync.WaitGroup
	for _, item := range organised {
		wg.Add(1)
		go item.sendRequest(&wg, respChan)
	}
	wg.Wait()
	close(respChan)

	//create map for efficient lookups
	reqMap := make(map[int]*Requests)
	for _, r := range req {
		reqMap[r.ID] = &r
	}

	var sentSlice []Response
	for sent := range respChan {
		for id, timeout := range sent.origRequest.idTimeout {
			responseItem := Response{
				Requests: *reqMap[id],
				Err:      sent.err,
			}
			if sent.err != nil || !sent.Up || timeout < sent.ResponseTime { //errored or responded in a longer time than was requested
				responseItem.CheckResponse = sent.CheckResponse
				responseItem.CheckResponse.Up = false
			} else {
				responseItem.CheckResponse = sent.CheckResponse
			}
			sentSlice = append(sentSlice, responseItem)
		}
	}
	return sentSlice
}

func organiseMonitorRequests(reqs ...Requests) []makeRequest {
	var resp []makeRequest
	for _, req := range reqs {
		foundIndex, err := req.findMatch(resp)
		if err != nil { //create and append
			resp = append(resp, makeRequest{
				url:   req.URL,
				rtype: req.RType,
				port:  req.Port,
				idTimeout: map[int]time.Duration{
					req.ID: req.Timeout,
				},
				maxTimeout: req.Timeout,
			})
		} else { // add id: timeout map to existing record and update max timeout if needed
			resp[foundIndex].idTimeout[req.ID] = req.Timeout
			timeDifference := resp[foundIndex].maxTimeout - req.Timeout
			if timeDifference < 0 {
				resp[foundIndex].maxTimeout = req.Timeout
			}
		}
	}
	return resp
}

// returns the index of the slice where the match exists
func (req *Requests) findMatch(mr []makeRequest) (int, error) {
	for index, item := range mr {
		matchCount := 0
		if item.url == req.URL {
			matchCount++
		}
		if item.rtype == req.RType {
			matchCount++
		}
		if item.port == req.Port {
			matchCount++
		}
		if matchCount == 3 {
			return index, nil
		}
	}
	return -1, errors.New("item not found in provided slice")
}

func (req makeRequest) sendRequest(wg *sync.WaitGroup, rchan chan<- makeRequestResponseWithErr) {
	defer wg.Done()
	switch req.rtype {
	case RequestTypeICMP:
		check, err := sendICMPRequest(req)
		resp := makeRequestResponseWithErr{
			makeRequestResponse: check,
			err:                 err,
		}
		rchan <- resp
	case RequestTypeTCP:
		check, err := sendTCPRequest(req)
		resp := makeRequestResponseWithErr{
			makeRequestResponse: check,
			err:                 err,
		}
		rchan <- resp
	case RequestTypeHTTP:
		check, err := sendHTTPRequest(req)
		resp := makeRequestResponseWithErr{
			makeRequestResponse: check,
			err:                 err,
		}
		rchan <- resp
	default:
		runTime := time.Now()
		mrr := makeRequestResponse{
			origRequest: req,
			CheckResponse: CheckResponse{
				RunTime: runTime,
			},
		}

		resp := makeRequestResponseWithErr{
			makeRequestResponse: mrr,
			err:                 errors.New("unsupported request type"),
		}
		rchan <- resp
	}
}

func sendICMPRequest(mr makeRequest) (makeRequestResponse, error) {
	runTime := time.Now()
	resp := makeRequestResponse{
		origRequest: mr,
		CheckResponse: CheckResponse{
			RunTime: runTime,
		},
	}
	pinger, err := probing.NewPinger(mr.url)
	if err != nil {
		return resp, fmt.Errorf("error creating pinger for %v: %v", mr.url, err.Error())
	}
	pinger.Count = 1

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return resp, fmt.Errorf("error pinging %v: %v", mr.url, err.Error())
	}
	if pinger.Statistics().PacketLoss > 0 {
		return resp, fmt.Errorf("error pinging %v: packet loss", mr.url)
	}

	resp.Up = true
	resp.ResponseTime = pinger.Statistics().AvgRtt

	return resp, nil
}

func sendHTTPRequest(mr makeRequest) (makeRequestResponse, error) {
	runTime := time.Now()
	resp := makeRequestResponse{
		origRequest: mr,
		CheckResponse: CheckResponse{
			RunTime: runTime,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), mr.maxTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", mr.url, nil)
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

func sendTCPRequest(mr makeRequest) (makeRequestResponse, error) {
	runTime := time.Now()
	resp := makeRequestResponse{
		origRequest: mr,
		CheckResponse: CheckResponse{
			RunTime: runTime,
		},
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(mr.url, fmt.Sprint(mr.port)), mr.maxTimeout)
	end := time.Now()
	if err != nil {
		return resp, fmt.Errorf("error sending tcp request for %v on port %d: %v", mr.url, mr.port, err.Error())
	}
	defer conn.Close()
	resp.Up = true
	resp.ResponseTime = end.Sub(start)

	return resp, nil
}
