package api

import (
	"bytes"
	"heimatcli/src/x/log"
	"strings"

	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

func (api *API) httpPost(
	authtoken string,
	url string,
	queries []Query,
	payload []byte,
) (*http.Response, []*http.Cookie, error) {
	return api.httpRequest(http.MethodPost, authtoken, url, queries, payload)
}

func (api *API) httpGet(
	authtoken string,
	url string,
	queries []Query,
) (*http.Response, []*http.Cookie, error) {
	emptyPayload := make([]byte, 0)
	return api.httpRequest(http.MethodGet, authtoken, url, queries, emptyPayload)
}

//
// USE THE PRECONFIGURED METHODS
//
func (api *API) httpRequest(
	method string,
	authtoken string,
	url string,
	queries []Query,
	payload []byte,
) (*http.Response, []*http.Cookie, error) {

	var body io.Reader

	// Create Request With CookieJar
	jar, err := cookiejar.New(nil)
	emptyCookieJar := make([]*http.Cookie, 0)
	if err != nil {
		log.Error.Printf("Could not create cookie jar: %s\n", err.Error())
	}

	client := http.Client{Jar: jar}

	if len(payload) != 0 {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, emptyCookieJar, err
	}

	// Set Headers
	req.Header.Set("Host", "heimat.sprinteins.com")
	req.Header.Set("authorization", "Bearer "+authtoken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client", "Heimat CLI (hicli)")

	// Add Queries
	q := req.URL.Query()
	for _, query := range queries {
		q.Add(query.key, query.value)
	}
	req.URL.RawQuery = q.Encode()

	// Make Request
	s := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
	s.Start()

	var resp *http.Response
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		resp, err = client.Do(req)
		wg.Done()
	}()
	wg.Wait()
	s.Stop()

	if err != nil {
		return nil, emptyCookieJar, err
	}

	newToken := resp.Header.Get("Authorization")
	if newToken != "" {
		newToken = strings.Replace(newToken, "Bearer ", "", 1)
		api.SetToken(newToken)
	}

	cookies := jar.Cookies(req.URL)
	return resp, cookies, nil
}

func readBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error.Printf("Could not read http resp body: %s\n", err.Error())
	}
	return body
}

// Query _
type Query struct {
	key   string
	value string
}
