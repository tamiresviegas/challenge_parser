package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Crawler is responsible for crawling the web pages on a specific domain.
type Crawler struct {
	BaseURL   *url.URL
	Subdomain string
	Visited   map[string]bool
}

func NewCrawler(baseURL, subdomain string) (*Crawler, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		BaseURL:   u,
		Subdomain: subdomain,
		Visited:   make(map[string]bool),
	}, nil
}

// Run starts tracing from the starting URL.
func (c *Crawler) Run() error {
	return c.crawl(c.BaseURL)
}

func (c *Crawler) crawl(u *url.URL) error {
	if u.Hostname() != c.BaseURL.Hostname() || strings.HasPrefix(u.Path, c.Subdomain) {
		return nil
	}

	if c.Visited[u.String()] {
		return nil
	}

	c.Visited[u.String()] = true

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to crawl %s: %s", u.String(), resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	c.extractLinks(doc)

	return nil
}

// extract Links extract the links from the HTML document and print them.
func (c *Crawler) extractLinks(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				linkURL, err := c.BaseURL.Parse(attr.Val)
				if err == nil {
					fmt.Println(linkURL.String())
				}
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.extractLinks(child)
	}
}

func main() {
	crawler, err := NewCrawler("https://parserdigital.com/", "parserdigital.com")
	if err != nil {
		fmt.Println("Failed to create crawler:", err)
		return
	}

	err = crawler.Run()
	if err != nil {
		fmt.Println("Crawling failed:", err)
		return
	}
}
