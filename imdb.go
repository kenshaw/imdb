// Package imdb is a simple imdb client.
package imdb

import (
	"context"
)

// FindCompany searches for a company.
func FindCompany(company string, opts ...Option) ([]Result, error) {
	return New(opts...).FindCompany(context.Background(), company)
}

// FindKeyword searches for a keyword.
func FindKeyword(keyword string, opts ...Option) ([]Result, error) {
	return New(opts...).FindKeyword(context.Background(), keyword)
}

// FindName searches for a name.
func FindName(name string, opts ...Option) ([]Result, error) {
	return New(opts...).FindName(context.Background(), name)
}

// FindTitle searches for a title.
func FindTitle(title string, opts ...Option) ([]Result, error) {
	return New(opts...).FindTitle(context.Background(), title)
}

// FindGame searches for a game.
func FindGame(game string, opts ...Option) ([]Result, error) {
	return New(opts...).FindGame(context.Background(), game)
}

// FindMovie searches for a movie.
func FindMovie(movie string, opts ...Option) ([]Result, error) {
	return New(opts...).FindMovie(context.Background(), movie)
}

// FindSeries searches for a series.
func FindSeries(series string, opts ...Option) ([]Result, error) {
	return New(opts...).FindSeries(context.Background(), series)
}

// FindEpisode searches for a episode.
func FindEpisode(episode string, opts ...Option) ([]Result, error) {
	return New(opts...).FindEpisode(context.Background(), episode)
}
