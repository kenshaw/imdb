// Package imdb is a simple imdb client.
package imdb

import (
	"context"
)

// Find searches for q.
func Find(ctx context.Context, q string, opts ...Option) ([]Result, error) {
	return New(opts...).Find(ctx, q)
}

// FindType searches for type and q.
func FindType(ctx context.Context, typ, q string, opts ...Option) ([]Result, error) {
	return New(opts...).FindType(ctx, typ, q)
}

// FindCompany searches for a company.
func FindCompany(ctx context.Context, company string, opts ...Option) ([]Result, error) {
	return New(opts...).FindCompany(ctx, company)
}

// FindKeyword searches for a keyword.
func FindKeyword(ctx context.Context, keyword string, opts ...Option) ([]Result, error) {
	return New(opts...).FindKeyword(ctx, keyword)
}

// FindName searches for a name.
func FindName(ctx context.Context, name string, opts ...Option) ([]Result, error) {
	return New(opts...).FindName(ctx, name)
}

// FindTitle searches for a title.
func FindTitle(ctx context.Context, title string, opts ...Option) ([]Result, error) {
	return New(opts...).FindTitle(ctx, title)
}

// FindTitleSubtype searches for subtype with title.
func FindTitleSubtype(ctx context.Context, subtype, title string, opts ...Option) ([]Result, error) {
	return New(opts...).FindTitleSubtype(ctx, subtype, title)
}

// FindGame searches for a game.
func FindGame(ctx context.Context, game string, opts ...Option) ([]Result, error) {
	return New(opts...).FindGame(ctx, game)
}

// FindMovie searches for a movie.
func FindMovie(ctx context.Context, movie string, opts ...Option) ([]Result, error) {
	return New(opts...).FindMovie(ctx, movie)
}

// FindSeries searches for a series.
func FindSeries(ctx context.Context, series string, opts ...Option) ([]Result, error) {
	return New(opts...).FindSeries(ctx, series)
}

// FindEpisode searches for a episode.
func FindEpisode(ctx context.Context, episode string, opts ...Option) ([]Result, error) {
	return New(opts...).FindEpisode(ctx, episode)
}
