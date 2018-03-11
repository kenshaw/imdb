package imdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// DefaultURL is the default request URL for API requests.
	DefaultURL = "http://www.omdbapi.com/"
)

// Option is a client option.
type Option func(*Client)

// Client is a imdb client.
type Client struct {
	cl     *http.Client
	url    string
	apikey string
}

// New creates a new client.
func New(apikey string, opts ...Option) *Client {
	return &Client{
		apikey: apikey,
		url:    DefaultURL,
	}
}

// Do wraps sending a request.
func (c *Client) Do(params url.Values) (*http.Response, error) {
	if c.apikey == "" {
		return nil, errors.New("must provide apikey")
	}
	params.Set("apikey", c.apikey)

	// setup URL
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimSuffix(u.Path, "/") + "/"
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	cl := c.cl
	if cl == nil {
		cl = http.DefaultClient
	}

	// do
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		return nil, fmt.Errorf("received status code %d", res.StatusCode)
	}

	return res, nil
}

// Search searches for movies given the title and optional year.
func (c *Client) Search(title, year string) (*SearchResponse, error) {
	res, err := c.Do(url.Values{
		"s": []string{title},
		"y": []string{year},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r := new(SearchResponse)
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// MovieByTitle returns a MovieResult given the title and optional year.
func (c *Client) MovieByTitle(title, year string) (*MovieResult, error) {
	res, err := c.Do(url.Values{
		"t":        []string{title},
		"y":        []string{year},
		"plot":     []string{"plot"},
		"tomatoes": []string{"true"},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r := &MovieResult{}
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// MovieByImdbID performs an API search for a specified movie by the specific
// id (ie, "tt2015381").
func (c *Client) MovieByImdbID(id string) (*MovieResult, error) {
	res, err := c.Do(url.Values{
		"i":        []string{id},
		"plot":     []string{"plot"},
		"tomatoes": []string{"true"},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r := &MovieResult{}
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}
