package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
)

type scraper struct {
	httpClient *http.Client
}

type commandDescription struct {
	Name     string
	Args     []string
	Metadata []string
}

type Scraper interface {
	Scrape(ctx context.Context, url *url.URL) (*commandDescription, error)
}

func NewScraper(httpClient *http.Client) Scraper {
	return &scraper{httpClient}
}

func (s *scraper) Scrape(ctx context.Context, url *url.URL) (*commandDescription, error) {
	resp, err := s.request(ctx, url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return parseDocument(resp.Body)
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

func parseDocument(body io.Reader) (*commandDescription, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	cmdSelection := doc.Find(".command")
	name := cmdSelection.Find(".name").Text()
	args := []string{}
	cmdSelection.Find(".arg").Each(func(_ int, selection *goquery.Selection) {
		args = append(args, selection.Text())
	})

	articleMainSelection := doc.Find(".article-main")
	metadata := []string{}
	articleMainSelection.Find(".metadata").Each(func(_ int, selection *goquery.Selection) {
		metadata = append(metadata, selection.Text())
	})
	cmdDescription := &commandDescription{
		Name:     name,
		Args:     args,
		Metadata: metadata,
	}
	return cmdDescription, nil
}
