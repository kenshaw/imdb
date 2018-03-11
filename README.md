# About imdb

`imdb` provides a simple wrapper around the [OMDb API][omdb-api].

## Installing

Install in the usual Go fashion:

```sh
$ go get -u github.com/kenshaw/imdb
```

## Using

`imdb` can be used similarly to the following:

```go
import (
    /* ... */
    "github.com/kenshaw/imdb"
)

cl := imdb.New("my-api-key")
res, err := cl.Search("Fight Club", "")
if err != nil { /* ... */ }
log.Printf(">>> results: %+v", res)
```

Please see the [GoDoc][godoc] listing for the full API.

[omdb-api]: http://www.omdbapi.com/
[godoc]: https://godoc.org/github.com/kenshaw/imdb
