package meetup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	baseURL = "https://api.meetup.com"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-Ratelimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
	headerContentType   = "Content-Type"

	mediaType = "application/json"
)

// Client is used to communicate with the Meetup API
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for the meetup api. Have it as a field to allow mock testing
	BaseURL string

	// Authentication can be used to authenticate requests
	Authentication Authenticator

	// Rate limits for the client
	RateLimits *Rate

	sync.Mutex
}

// NewClient returns a new Meetup API client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		client:         httpClient,
		BaseURL:        baseURL,
		Authentication: nil,
		RateLimits:     nil,
	}
}

// Create an API request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, urlStr, buf)
	if err != nil {
		return req, err
	}

	if body != nil {
		req.Header.Set(headerContentType, mediaType)
	}

	if c.Authentication != nil {
		err := c.Authentication.AuthenticateRequest(req)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

// Executes a request to the API and stores the JSON decoded response in v
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	c.Lock()
	c.RateLimits = parseRate(resp.Header)
	c.Unlock()

	err = checkResponse(resp)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

// Rate limits by the Meetup API
type Rate struct {
	// Maximum number of requests that can be made in a window of time
	Limit int

	// Remaining number of requests allowed in the current rate limit window
	Remaining int

	// Number of seconds until the current rate limit window resets
	Reset int
}

// parse rate limits from the Meetup API response header
func parseRate(header http.Header) *Rate {
	var rate Rate

	if limit := header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}

	if remaining := header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}

	if reset := header.Get(headerRateReset); reset != "" {
		rate.Reset, _ = strconv.Atoi(reset)
	}

	return &rate
}

// Error reports a general error from the API
type ErrorResponse struct {
	Response *http.Response
	Errors   []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	message := fmt.Sprintf("%v: %v %v", e.Response.Status, e.Response.Request.Method, e.Response.Request.URL)
	for _, apiErr := range e.Errors {
		message += fmt.Sprintf(" %v,", apiErr.Message)
	}

	return strings.TrimSuffix(message, ",")
}

// Checks whether the API call resulted in an error
func checkResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || body == nil {
		return errorResponse
	}

	json.Unmarshal(body, &errorResponse)

	return errorResponse
}

// params must be a struct / pointer to a struct with possible types: pointers, ints, uints, floats, bool, string
// Encodes params as URL query parameters and returns the resulting url.
// params fields should contain "url" tags, or else would be ignored.
func addQueryParams(rawURL string, params interface{}) (string, error) {
	if params == nil {
		return rawURL, nil
	}

	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return rawURL, nil
		}
		val = val.Elem()
	}

	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return rawURL, fmt.Errorf("expected params to be a struct, was: %v", typ)
	}

	queryParams := make(url.Values)
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		tag := typ.Field(i).Tag.Get("url")
		name := strings.Split(tag, ",")[0]

		omitEmpty := false
		if strings.Contains(strings.TrimPrefix(tag, name), "omitempty") {
			omitEmpty = true
		}

		if tag == "" || name == "" || (omitEmpty && isEmpty(field)) {
			continue
		}

		queryParams.Add(name, fmt.Sprint(field.Interface()))
	}

	rawURL = strings.TrimSuffix(rawURL, "/")
	if len(queryParams) > 0 {
		rawURL = fmt.Sprintf("%v?%v", rawURL, queryParams.Encode())
	}

	return rawURL, nil
}

// Returns whether or not val is the empty value of its type
func isEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Bool:
		return !val.Bool()
	case reflect.String:
		return val.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Array, reflect.Map, reflect.Slice:
		return val.Len() == 0
	}
	return false
}
