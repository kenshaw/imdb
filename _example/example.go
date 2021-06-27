package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kenshaw/imdb"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	typ := flag.String("t", "all", "type")
	q := flag.String("q", "", "query")
	flag.Parse()
	if err := run(context.Background(), *verbose, *typ, *q); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, verbose bool, typ, q string) error {
	var opts []imdb.Option
	if verbose {
		opts = append(opts, imdb.WithLogf(fmt.Printf))
	}
	opts = append(opts, imdb.WithAppCacheDir("go-imdb-example"))
	cl := imdb.New(opts...)
	// determine what to search for
	var f func(context.Context, string, ...string) ([]imdb.Result, error)
	switch typ {
	case "all":
		f = cl.Find
	case "company":
		f = cl.FindCompany
	case "keyword":
		f = cl.FindKeyword
	case "name":
		f = cl.FindName
	case "title":
		f = cl.FindTitle
	case "movie":
		f = cl.FindMovie
	case "series":
		f = cl.FindSeries
	case "episode":
		f = cl.FindEpisode
	case "game":
		f = cl.FindGame
	default:
		return fmt.Errorf("unknown -t flag: %q", typ)
	}
	// find
	res, err := f(ctx, q)
	if err != nil {
		return err
	}
	for i, r := range res {
		fmt.Printf("%d: %v\n", i, r)
		fmt.Printf("  url: %s\n", r.URL)
	}
	return nil
}
