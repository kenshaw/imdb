package imdb

import (
	"fmt"
	"os"
)

const (
	envKey = "OMDB_APIKEY"
)

func init() {
	// setup default client
	DefaultClient = New(os.Getenv(envKey))
}

// DefaultClient is the default IMDB client.
var DefaultClient *Client

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

// A Rating will hold the related information about rating from different
// sources
type Rating struct {
	Source string
	Value  string
}

// A MovieResult will hold the related information of a single movie.
type MovieResult struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Metascore  string
	ImdbRating string
	ImdbVotes  string
	ImdbID     string
	Type       string
	Ratings    []Rating
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
	Error      string
}

// String satisifies the Stringer interface for MovieResult.
func (mr MovieResult) String() string {
	return fmt.Sprintf("#%s: %s (%s)", mr.ImdbID, mr.Title, mr.Year)
}

// Search searches for movies given the title and optional year using DefaultClient.
func Search(title, year string) (*SearchResponse, error) {
	return DefaultClient.Search(title, year)
}

// MovieByTitle returns a MovieResult given the title and optional year using DefaultClient.
func MovieByTitle(title, year string) (*MovieResult, error) {
	return DefaultClient.MovieByTitle(title, year)
}

// MovieByImdbID performs an API search for a specified movie by the specific
// id (ie, "tt2015381") using DefaultClient.
func MovieByImdbID(id string) (*MovieResult, error) {
	return DefaultClient.MovieByImdbID(id)
}
