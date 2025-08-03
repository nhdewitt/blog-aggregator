// Package app contains shared application services and state management.
package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

// fetchFeed retrievves and parses an RSS feed from the given URL.
// It sends an HTTP GET request with a custom User-Agent header, reads the response,
// and unmarshals the XML into an RSSFeed struct.
//
// The function automatically unescapes HTML entities in the feed title, description,
// and all item titles and descriptions to ensure proper display of special characters.
//
// Parameters:
//		- ctx: Context for request cancellation and timeout control
//		- feedURL: The URL of the RSS feed to fetch
//
// Returns the parsed RSSFeed struct or an error if the request fails,
// returns a non-2xx status code, or if XML parsing fails.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create HTTP request with context for cancellation support
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "gator")
	
	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error receiving response: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful HTTP status
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP Status %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	// Parse the XML into RSSFeed struct
	var rss RSSFeed
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling XML: %v", err)
	}

	// Unescape HTML entities in feed content for proper display
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}

	return &rss, nil
}