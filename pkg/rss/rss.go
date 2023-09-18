package rss

import (
	"aggregatenews/pkg/store"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Links       string `xml:"link"`
	Items       []Item `xml:"post"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// ParseFeed read rss-feed and return encoding news
func ParseFeed(url string) ([]store.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var f Feed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}

	data, err := itemToPost(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parsePubTime(pubDate string) (int64, error) {
	pubDate = strings.Replace(pubDate, ",", "", -1)
	t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700 MST", pubDate)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func itemToPost(f Feed) ([]store.Post, error) {
	var data []store.Post
	// Выделение памяти под массив
	data = make([]store.Post, 0, len(f.Channel.Items))
	for _, item := range f.Channel.Items {
		var p store.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = html.EscapeString(p.Content)
		p.Link = item.Link
		t, err := parsePubTime(item.PubDate)
		if err != nil {
			return nil, err
		}
		p.PubTime = t
		data = append(data, p)
	}
	return data, nil
}
