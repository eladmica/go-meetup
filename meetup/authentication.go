package meetup

import "net/http"

// used to authenticate requests
type Authenticator interface {
	AuthenticateRequest(*http.Request) error
}

// KeyAuth is used to Authenticate requests with the user's key
// Meetup API docs: https://www.meetup.com/meetup_api/auth/#keys
type KeyAuth struct {
	Key string
}

func NewKeyAuth(key string) Authenticator {
	return &KeyAuth{Key: key}
}

func (auth *KeyAuth) AuthenticateRequest(req *http.Request) error {
	params := req.URL.Query()
	params.Set("key", auth.Key)
	req.URL.RawQuery = params.Encode()
	return nil
}
