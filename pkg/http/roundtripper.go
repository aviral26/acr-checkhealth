package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/aviral26/acr-checkhealth/pkg/io"
	"github.com/opencontainers/go-digest"
	"github.com/rs/zerolog"
)

// HTTP related constants
const (
	HeaderChallenge     = "Www-Authenticate"
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"
)

// Request represents a request made to the registry.
type Request struct {
	Method              string    `json:"method"`
	URL                 *url.URL  `json:"url"`
	HeaderAuthorization string    `json:"authorization"`
	StartedAt           time.Time `json:"startedAt"`
}

// Response respresents a response received from the registry.
type Response struct {
	Code            int           `json:"code"`
	HeaderChallenge string        `json:"Www-Authenticate"`
	HeaderLocation  *url.URL      `json:"redirectLocation"`
	Size            int64         `json:"size"`
	SHA256Sum       digest.Digest `json:"sha256"`

	Body []byte
}

// RoundTripInfo represents information about a network round-trip.
type RoundTripInfo struct {
	Request  `json:"request"`
	Response `json:"response"`
	Elapsed  time.Duration `json:"elapsed"`
}

// RoundTripper provides a means to do an HTTP/HTTPs round trip.
type RoundTripper interface {
	// RoundTrip makes an HTTP request and returns the response with some stats.
	RoundTrip(req *http.Request) (RoundTripInfo, error)
}

// RoundTripperWithContext provides an implementation for RoundTripper.
type RoundTripperWithContext struct {
	Base   http.RoundTripper
	Logger zerolog.Logger
}

// RoundTrip does an HTTP/HTTPs roundtrip and returns the response with some contextual info.
func (r RoundTripperWithContext) RoundTrip(req *http.Request) (RoundTripInfo, error) {
	info := RoundTripInfo{
		Request: Request{
			Method:              req.Method,
			URL:                 req.URL,
			StartedAt:           time.Now(),
			HeaderAuthorization: req.Header.Get(HeaderAuthorization),
		},
	}
	defer func() {
		info.Elapsed = time.Since(info.StartedAt)
		var msg string
		bytes, err := json.Marshal(info)
		if err != nil {
			msg = fmt.Sprintf("marshal_error: %v", err)
		} else {
			msg = string(bytes)
		}
		r.Logger.Trace().Msg(msg)
	}()

	resp, err := r.Base.RoundTrip(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	bodyReader := io.NewReader(resp.Body)
	bodyBytes, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return info, err
	}

	info.Response = Response{
		Code:            resp.StatusCode,
		HeaderChallenge: resp.Header.Get(HeaderChallenge),
		Size:            bodyReader.N(),
		SHA256Sum:       digest.NewDigest(digest.SHA256, bodyReader.SHA256Hash()),
		Body:            bodyBytes,
	}

	locURL, err := resp.Location()
	if err != nil {
		if err != http.ErrNoLocation {
			return info, err
		}
	} else {
		info.Response.HeaderLocation = locURL
	}

	return info, nil
}
