package imdb

import (
	"testing"
)

func TestFindName(t *testing.T) {
	t.Parallel()
	res, err := FindName("brad pitt", WithLogf(t.Logf), WithAppCacheDir("imdb-test"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(res) < 1 {
		t.Fatalf("expected at least one search result, got: %d", len(res))
	}
	if exp := "Brad Pitt"; res[0].Title != exp {
		t.Errorf("expected %q, got: %q", exp, res[0].Title)
	}
	if exp, got := "", res[0].Year(); got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "nm", res[0].Subtype; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "nm0000093", res[0].ID; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
}

func TestFindTitle(t *testing.T) {
	t.Parallel()
	res, err := FindTitle("bobs burgers", WithLogf(t.Logf), WithAppCacheDir("imdb-test"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(res) < 1 {
		t.Fatalf("expected at least one search result, got: %d", len(res))
	}
	if exp := "Bob's Burgers"; res[0].Title != exp {
		t.Errorf("expected %q, got: %q", exp, res[0].Title)
	}
	if exp, got := "2011", res[0].Year(); got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "tv", res[0].Subtype; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "tt1561755", res[0].ID; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
}

func TestFindMovie(t *testing.T) {
	res, err := FindMovie("fight club", WithLogf(t.Logf), WithAppCacheDir("imdb-test"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(res) < 1 {
		t.Fatalf("expected at least one search result, got: %d", len(res))
	}
	if exp := "Fight Club"; res[0].Title != exp {
		t.Errorf("expected %q, got: %q", exp, res[0].Title)
	}
	if exp, got := "1999", res[0].Year(); got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "ft", res[0].Subtype; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
	if exp, got := "tt0137523", res[0].ID; got != exp {
		t.Errorf("expected %q, got: %q", exp, got)
	}
}
