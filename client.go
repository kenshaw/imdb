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
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/omahaproxy"
	"github.com/kenshaw/diskcache"
	"github.com/kenshaw/httplog"
)

// DefaultTransport is the default http transport.
var DefaultTransport = http.DefaultTransport

// Type values.
const (
	TypeAll     = "al"
	TypeCompany = "co"
	TypeKeyword = "kw"
	TypeName    = "nm"
	TypeTitle   = "tt"
)

// Subtype values.
const (
	SubtypeGame    = "vg"
	SubtypeMovie   = "ft"
	SubtypeSeries  = "tv"
	SubtypeEpisode = "ep"
)

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
		if err = cl.buildUserAgent(ctx); err != nil {
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

// buildUserAgent builds the user agent.
func (cl *Client) buildUserAgent(ctx context.Context) error {
	if cl.UserAgent != "" {
		return nil
	}
	// retrieve latest chrome version
	ver, err := omahaproxy.New(
		omahaproxy.WithTransport(cl.cl.Transport),
	).Latest(ctx, "linux", "stable")
	if err != nil {
		return err
	}
	// build user agent
	cl.UserAgent = fmt.Sprintf("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", ver.Version)
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
	doc.Find(".findResult .result_text").Each(func(i int, s *goquery.Selection) {
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
		typ, subtype := id[:2], id[:2]
		if typ == "tt" {
			// determine subtype
			subtype = SubtypeMovie
			if m := subtypeRE.FindStringSubmatch(s.Text()); m != nil {
				switch m[1] {
				case "TV Mini Series", "TV Series", "TV Short":
					subtype = SubtypeSeries
				case "TV Episode":
					subtype = SubtypeEpisode
				case "Video Game":
					subtype = SubtypeGame
				}
			}
		}
		u.RawQuery = ""
		res = append(res, Result{
			URL:     u.String(),
			ID:      id,
			Title:   a.Text(),
			Type:    typ,
			Subtype: subtype,
			S:       s,
		})
	})
	return res, nil
}

// FindType searches for type and q.
func (cl *Client) FindType(ctx context.Context, typ, q string, params ...string) ([]Result, error) {
	return cl.Find(ctx, q, append(params, "s", typ)...)
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
func (cl *Client) FindTitleSubtype(ctx context.Context, subtype, title string, params ...string) ([]Result, error) {
	return cl.FindTitle(ctx, title, append(params, "ttype", subtype)...)
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
	Type    string
	Subtype string
	S       *goquery.Selection
}

// String satisfies the fmt.Stringer interface.
func (r Result) String() string {
	var year string
	if s := r.Year(); s != "" {
		year = fmt.Sprintf(" (%s)", s)
	}
	return fmt.Sprintf("%s: %q%s %s", r.ID, r.Title, year, r.URL)
}

// Year returns the year from the selection.
func (r Result) Year() string {
	if r.Type == TypeTitle {
		if m := yearRE.FindStringSubmatch(r.S.Text()); m != nil {
			return m[1]
		}
	}
	return ""
}

// yearRE matches a year.
var yearRE = regexp.MustCompile(`\(([0-9]{4})\)`)

// subtypeRE matches subtypes.
var subtypeRE = regexp.MustCompile(`\((TV Mini Series|TV Series|TV Short|TV Episode|Video Game)\)`)

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
