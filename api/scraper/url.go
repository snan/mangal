package scraper

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"time"
)

type URL struct {
	// Address of this URL
	Address string
	// Info about this URL (e.g. page title)
	Info string
	// Source parent of this URL
	Source *Source
}

// Get http response from url
func (u URL) Get() (*http.Response, error) {
	// Set default http timeout to 5 seconds
	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(u.Address)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New(fmt.Sprintf("Cannot get '%s'. Status %d - %s", u.Address, resp.StatusCode, resp.Status))
	}

	return resp, nil
}

// getDocument Gets goquery Document from URL
func (u URL) document() (*goquery.Document, error) {
	resp, err := u.Get()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Unexpected error while closing http body")
		}
	}(resp.Body)

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}
