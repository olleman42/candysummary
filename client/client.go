package client

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/olleman42/candysummary/summarizer"
)

// Client is the container type for querying the entry source
type Client struct {
	uri string
}

// New returns a new instance of the HTTP client
func New(uri string) *Client {
	return &Client{uri}
}

// GetEntries uses this implementation
func (c Client) GetEntries() ([]summarizer.HistoryEntry, error) {
	resp, err := http.Get(c.uri)
	// resp, err := http.Get("https://candystore.zimpler.net/")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// typically check for proper status code, but omitted in this case

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	entries := []summarizer.HistoryEntry{} // potential optimization for large datasets (which for some reason get delivered using HTML)
	parseError := false

	sel := doc.Find("#top\\.customers tbody tr")

	if sel.Size() == 0 {
		return nil, errors.New("Could not find matching pattern on source page")
	}

	sel.Each(func(i int, s *goquery.Selection) {
		if parseError {
			return
		}
		entry := summarizer.HistoryEntry{}
		s.Children().Each(func(i int, s *goquery.Selection) {

			switch i {
			case 0:
				entry.Name = s.Text()
			case 1:
				entry.Candy = s.Text()
			case 2:
				entry.Eaten, err = strconv.Atoi(s.Text())
				if err != nil {
					parseError = true
				}
			}
		})
		entries = append(entries, entry)
	})

	if parseError {
		return nil, errors.New("parsing failed")
	}

	return entries, nil
}
