package meetup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	baseURL = "https://api.meetup.com"
)

// Client is used to communicate with the Meetup API
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// TODO: Add authentication

	// TODO: Add rate limit
}

// NewClient returns a new Meetup API client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		client: httpClient,
	}
}

// Create an API request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// TODO: request body
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return req, err
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

// Encodes params struct as URL query parameters and returns the resulting url.
// params fields should contain "url" tags, or else would be ignored.
// TODO: Support maps / slices / interfaces, params as map?
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
