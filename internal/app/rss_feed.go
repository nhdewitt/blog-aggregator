// Package app contains shared application services and state management.
package app

type RSSFeed struct {
	Channel struct {
		Title		string	`xml:"title"`
		Link		string	`xml:"link"`
		Description	string	`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}