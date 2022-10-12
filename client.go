package imdb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kenshaw/diskcache"
	"github.com/kenshaw/httplog"
)

// DefaultTransport is the default http transport.
var DefaultTransport = http.DefaultTransport

// Type is a type.
type Type string

// Type values.
const (
	TypeAll     Type = "al"
	TypeCompany Type = "co"
	TypeKeyword Type = "kw"
	TypeName    Type = "nm"
	TypeTitle   Type = "tt"
)

// String satisfies the fmt.Stringer interface.
func (typ Type) String() string {
	switch typ {
	case TypeAll:
		return "all"
	case TypeCompany:
		return "company"
	case TypeKeyword:
		return "keyword"
	case TypeName:
		return "name"
	case TypeTitle:
		return "title"
	}
	return "Type(" + string(typ) + ")"
}

// Subtype is a subtype.
type Subtype string

// Subtype values.
const (
	SubtypeGame    Subtype = "vg"
	SubtypeMovie   Subtype = "ft"
	SubtypeSeries  Subtype = "tv"
	SubtypeEpisode Subtype = "ep"
)

// String satisfies the fmt.Stringer interface.
func (subtype Subtype) String() string {
	switch subtype {
	case SubtypeGame:
		return "game"
	case SubtypeMovie:
		return "movie"
	case SubtypeSeries:
		return "series"
	case SubtypeEpisode:
		return "episode"
	}
	return "Subtype(" + string(subtype) + ")"
}

// Client is a imdb client.
type Client struct {
	Transport   http.RoundTripper
	UserAgent   string
	AppCacheDir string
	cl          *http.Client
	once        sync.Once
}

// New creates a new imdb client.
func New(opts ...Option) *Client {
	cl := &Client{
		Transport: DefaultTransport,
		UserAgent: "wget",
	}
	for _, o := range opts {
		o(cl)
	}
	return cl
}

// init initializes the client.
func (cl *Client) init(ctx context.Context) error {
	var err error
	cl.once.Do(func() {
		if err = cl.buildClient(ctx); err != nil {
			return
		}
	})
	return err
}

// buildClient builds the http client used for retrievals.
func (cl *Client) buildClient(ctx context.Context) error {
	if cl.cl != nil {
		return nil
	}
	transport := cl.Transport
	if cl.AppCacheDir != "" {
		var err error
		transport, err = diskcache.New(
			diskcache.WithTransport(transport),
			diskcache.WithAppCacheDir(cl.AppCacheDir),
			diskcache.WithTTL(24*time.Hour),
			diskcache.WithHeaderWhitelist("Date", "Set-Cookie", "Content-Type", "Location"),
			diskcache.WithErrorTruncator(),
			diskcache.WithGzipCompression(),
		)
		if err != nil {
			return err
		}
	}
	cl.cl = &http.Client{Transport: transport}
	return nil
}

// buildRequest builds a request for the parameters.
func (cl *Client) buildRequest(method, urlstr string, r io.Reader, params ...string) (*http.Request, error) {
	// check params
	if len(params)%2 != 0 {
		return nil, fmt.Errorf("invalid params length %d", len(params))
	}
	// build query values
	v := make(url.Values, len(params)/2)
	for i := 0; i < len(params); i += 2 {
		v[params[i]] = []string{params[i+1]}
	}
	// build url
	u, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}
	u.RawQuery = v.Encode()
	// create request
	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}
	// add user agent
	if cl.UserAgent != "" {
		req.Header.Set("User-Agent", cl.UserAgent)
	}
	return req, nil
}

// get retrieves a url.
func (cl *Client) get(ctx context.Context, urlstr string, params ...string) ([]byte, error) {
	// initialize
	if err := cl.init(ctx); err != nil {
		return nil, err
	}
	// build request
	req, err := cl.buildRequest("GET", urlstr, nil, params...)
	if err != nil {
		return nil, err
	}
	// execute
	res, err := cl.cl.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// check status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d != 200", res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}

// Find searches for q.
func (cl *Client) Find(ctx context.Context, q string, params ...string) ([]Result, error) {
	// retrieve
	buf, err := cl.get(ctx, "https://www.imdb.com/find", append(params, "q", q)...)
	if err != nil {
		return nil, err
	}
	// create doc
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	var res []Result
	doc.Find(".find-result-item").Each(func(i int, s *goquery.Selection) {
		// get first a
		a := s.Find("a").First()
		href := a.AttrOr("href", "")
		if href == "" || !strings.HasPrefix(href, "/") {
			return
		}
		// parse url
		u, err := url.Parse("https://www.imdb.com" + href)
		if err != nil {
			return
		}
		// id
		id := path.Base(u.Path)
		typ, subtype := Type(id[:2]), Subtype(id[:2])
		if typ == TypeTitle {
			// determine subtype
			subtype = SubtypeMovie
			switch s.Find("li:nth-child(2)").Text() {
			case "TV Mini Series", "TV Series", "TV Short":
				subtype = SubtypeSeries
			case "TV Episode":
				subtype = SubtypeEpisode
			case "Video Game":
				subtype = SubtypeGame
			}
		}
		u.RawQuery = ""
		var year string
		if typ == TypeTitle {
			year = s.Find("li").First().Text()
		}
		res = append(res, Result{
			URL:     u.String(),
			ID:      id,
			Title:   a.Text(),
			Type:    typ,
			Subtype: subtype,
			Year:    year,
		})
	})
	return res, nil
}

// FindType searches for type and q.
func (cl *Client) FindType(ctx context.Context, typ Type, q string, params ...string) ([]Result, error) {
	return cl.Find(ctx, q, append(params, "s", string(typ))...)
}

// FindCompany searches for a company.
func (cl *Client) FindCompany(ctx context.Context, company string, params ...string) ([]Result, error) {
	return cl.FindType(ctx, TypeCompany, company, params...)
}

// FindKeyword searches for a keyword.
func (cl *Client) FindKeyword(ctx context.Context, keyword string, params ...string) ([]Result, error) {
	return cl.FindType(ctx, TypeKeyword, keyword, params...)
}

// FindName searches for a name.
func (cl *Client) FindName(ctx context.Context, name string, params ...string) ([]Result, error) {
	return cl.FindType(ctx, TypeName, name, params...)
}

// FindTitle searches for a title.
func (cl *Client) FindTitle(ctx context.Context, title string, params ...string) ([]Result, error) {
	return cl.FindType(ctx, TypeTitle, title, params...)
}

// FindTitleSubtype searches for subtype with title.
func (cl *Client) FindTitleSubtype(ctx context.Context, subtype Subtype, title string, params ...string) ([]Result, error) {
	return cl.FindTitle(ctx, title, append(params, "ttype", string(subtype))...)
}

// FindGame searches for a game.
func (cl *Client) FindGame(ctx context.Context, game string, params ...string) ([]Result, error) {
	return cl.FindTitleSubtype(ctx, SubtypeGame, game, params...)
}

// FindMovie searches for a movie.
func (cl *Client) FindMovie(ctx context.Context, movie string, params ...string) ([]Result, error) {
	return cl.FindTitleSubtype(ctx, SubtypeMovie, movie, params...)
}

// FindSeries searches for a series.
func (cl *Client) FindSeries(ctx context.Context, series string, params ...string) ([]Result, error) {
	return cl.FindTitleSubtype(ctx, SubtypeSeries, series, params...)
}

// FindEpisode searches for a episode.
func (cl *Client) FindEpisode(ctx context.Context, episode string, params ...string) ([]Result, error) {
	return cl.FindTitleSubtype(ctx, SubtypeEpisode, episode, params...)
}

// Result is the result of a search.
type Result struct {
	URL     string
	ID      string
	Title   string
	Type    Type
	Subtype Subtype
	Year    string
}

// String satisfies the fmt.Stringer interface.
func (r Result) String() string {
	var year string
	if r.Year != "" {
		year = ", " + r.Year
	}
	return fmt.Sprintf("%s: %q (%s%s) %s", r.ID, r.Title, r.Subtype, year, r.URL)
}

// YearInt returns the year as an int from the selection.
func (r Result) YearInt() int {
	year, _ := strconv.Atoi(cleanRE.ReplaceAllString(r.Year, ""))
	if year > 1800 && year < 2100 {
		return year
	}
	return 0
}

var cleanRE = regexp.MustCompile(`[^0-9]`)

// Option is a imdb client option.
type Option func(*Client)

// WithTransport is a imdb client option to set the http transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.Transport = transport
	}
}

// WithLogf is a imdb client option to set a log handler for http requests and
// responses.
func WithLogf(logf interface{}, opts ...httplog.Option) Option {
	return func(cl *Client) {
		cl.Transport = httplog.NewPrefixedRoundTripLogger(cl.Transport, logf, opts...)
	}
}

// WithAppCacheDir is a imdb client option to set the app cache dir.
func WithAppCacheDir(appCacheDir string) Option {
	return func(cl *Client) {
		cl.AppCacheDir = appCacheDir
	}
}
