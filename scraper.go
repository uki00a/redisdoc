package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type scraper struct {
	httpClient *http.Client
}

type Scraper interface {
	Scrape(ctx context.Context, url *url.URL) (*Document, error)
}

func NewScraper(httpClient *http.Client) Scraper {
	return &scraper{httpClient}
}

func (s *scraper) Scrape(ctx context.Context, url *url.URL) (*Document, error) {
	resp, err := s.request(ctx, url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return ParseDocument(resp.Body)
}

func (s *scraper) request(ctx context.Context, url *url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	doneCh := make(chan *http.Response)
	errCh := make(chan error)

	go func() {
		resp, err := s.httpClient.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		if resp.StatusCode != 200 {
			errCh <- fmt.Errorf("bad response: %d", resp.StatusCode)
			return
		}
		doneCh <- resp
	}()

	select {
	case resp := <-doneCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled")
	}
}
