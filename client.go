package splunk

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is an interface for interacting with Splunk's REST API
type Client interface {
	URL() string
	NewRequest(method, uri string, body io.Reader) (*Response, error)
	Knowledge() KnowledgeService
}

// ClientConfig is used to configure new API clients
type ClientConfig struct {
	URL                   string
	Username              string
	Password              string
	TLSInsecureSkipVerify bool
}

// ListOptions represents query parametes for pagination and filtering
type ListOptions struct {
	Offset int
	Count  int
}

// Response represents a standard splunk API json response object
// TODO: This is missing a number of standard fields
type Response struct {
	// TODO: This is probably a bad format to deliver as a user interface, but it worked for my experiment
	Entry json.RawMessage `json:"entry"`
}

// client is the default implementation of the Client interface
type client struct {
	client    *http.Client
	rawurl    string
	username  string
	password  string
	service   service
	knowledge KnowledgeService
}

// service is a generic object that gets converted in to various service implementations
type service struct {
	client Client
}

// NewClient constructs a new client object from url and basic auth credentials
func NewClient(config ClientConfig) (Client, error) {
	// raise error on client creation if the url is invalid
	neturl, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	httpClient := http.DefaultClient

	if config.TLSInsecureSkipVerify {
		httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	c := &client{
		client:   httpClient,
		rawurl:   neturl.String(),
		username: config.Username,
		password: config.Password,
	}

	// create a single service object and reuse it for each API service
	c.service.client = c
	c.knowledge = (*knowledgeService)(&c.service)

	return c, nil
}

// URL returns the client's raw URL string
func (c *client) URL() string {
	return c.rawurl
}

// TODO: better request methods for handling standard request types. Get(), Post(), Delete(), etc.

// NewRequest builds an http.Request and sends the Response.Body as an io.ReadCloser
func (c *client) NewRequest(method, uri string, body io.Reader) (*Response, error) {
	// build standard query parameters
	params := url.Values{}
	params.Set("output_mode", "json")

	// build request url
	fullpath := c.URL() + uri + "?" + params.Encode()

	req, err := http.NewRequest(method, fullpath, body)
	if err != nil {
		return nil, err
	}

	// set standard headers
	req.Header.Set("User-Agent", "go-splunk")
	req.SetBasicAuth(c.username, c.password)

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// raise error if request was not successful
	if resp.StatusCode > 399 {
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}

	// marshal body into a new Response object
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := new(Response)
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *client) Knowledge() KnowledgeService {
	return c.knowledge
}
