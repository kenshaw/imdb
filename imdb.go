/*
Copyright 2014 Kaissersoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package imdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseURL = "http://www.omdbapi.com/?"
)

// A SearchResult represents a single API search result.
type SearchResult struct {
	Title  string
	Year   string
	ImdbID string
	Type   string
}

// String satisfies the Stringer interface for SearchResult.
func (sr SearchResult) String() string {
	return fmt.Sprintf("#%s: %s (%s) Type: %s", sr.ImdbID, sr.Title, sr.Year, sr.Type)
}

// A SearchResponse is the surrounding container holding multiple SearchResults.
type SearchResponse struct {
	Search   []SearchResult
	Response string
	Error    string
}

// Search searches for movies given the title and optional year.
func Search(title string, year string) (*SearchResponse, error) {
	resp, err := doRequest(url.Values{
		"s": []string{title},
		"y": []string{year},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// A MovieResult will hold the related information of a single movie.
type MovieResult struct {
	Title             string
	Year              string
	Rated             string
	Released          string
	Runtime           string
	Genre             string
	Director          string
	Writer            string
	Actors            string
	Plot              string
	Language          string
	Country           string
	Awards            string
	Poster            string
	Metascore         string
	ImdbRating        string
	ImdbVotes         string
	ImdbID            string
	Type              string
	TomatoMeter       string
	TomatoImage       string
	TomatoRating      string
	TomatoReviews     string
	TomatoFresh       string
	TomatoRotten      string
	TomatoConsensus   string
	TomatoUserMeter   string
	TomatoUserRating  string
	TomatoUserReviews string
	DVD               string
	BoxOffice         string
	Production        string
	Website           string
	Response          string
	Error             string
}

// String satisifies the Stringer interface for MovieResult.
func (mr MovieResult) String() string {
	return fmt.Sprintf("#%s: %s (%s)", mr.ImdbID, mr.Title, mr.Year)
}

// MovieByTitle returns a MovieResult given the title and optional year.
func MovieByTitle(title string, year string) (*MovieResult, error) {
	resp, err := doRequest(url.Values{
		"t":        []string{title},
		"y":        []string{year},
		"plot":     []string{"plot"},
		"tomatoes": []string{"true"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &MovieResult{}
	err = json.NewDecoder(resp.Body).Decode(r)
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
func MovieByImdbID(id string) (*MovieResult, error) {
	resp, err := doRequest(url.Values{
		"i":        []string{id},
		"plot":     []string{"plot"},
		"tomatoes": []string{"true"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &MovieResult{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	if r.Response == "False" {
		return r, errors.New(r.Error)
	}

	return r, nil
}

// doRequest handles actual request to the API.
func doRequest(params url.Values) (resp *http.Response, err error) {
	var u *url.URL
	u, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u.Path += "/"
	u.RawQuery = params.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d received from imdb", res.StatusCode)
	}
	return res, nil
}
