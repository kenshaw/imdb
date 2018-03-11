package imdb

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if buf, err := ioutil.ReadFile(".omdbapikey"); err == nil {
		DefaultClient = New(string(bytes.TrimSpace(buf)))
	}
	os.Exit(m.Run())
}

func TestImdbSearchMovies(t *testing.T) {
	res, err := Search("Fight Club", "1999")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(res.Search) < 1 {
		t.Fatalf("expected at least one search result, got: %d", len(res.Search))
	}
	if res.Search[0].Title != "Fight Club" {
		t.Errorf("expected `Fight Club`, got: %q", res.Search[0].Title)
	}
}

func TestImdbGetMovieByTitle(t *testing.T) {
	res, err := MovieByTitle("Fight Club", "1999")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if res.Title != "Fight Club" {
		t.Errorf("expected `Fight Club`, got: %q", res.Title)
	}
}

func TestImdbGetMovieByImdbID(t *testing.T) {
	res, err := MovieByImdbID("tt0137523")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if res.Title != "Fight Club" {
		t.Errorf("expected `Fight Club`, got: %q", res.Title)
	}
}
