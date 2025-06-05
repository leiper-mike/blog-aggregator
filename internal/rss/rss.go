package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200{
		return nil, fmt.Errorf("error: %v", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	feed.decodeText()
	return &feed, nil
}


func (feed *RSSFeed) decodeText() {
	channel := &feed.Channel
	channel.Title = html.UnescapeString(channel.Title)
	channel.Description = html.UnescapeString(channel.Description)
	for i := range channel.Item{
		channel.Item[i].decodeText()
	}
}
func (item *RSSItem) decodeText(){
	item.Description  = html.UnescapeString(item.Description)
	item.Title = html.UnescapeString(item.Title)
}
func (item RSSItem) Format() string{
	return fmt.Sprintf("Post Title: %v \nLink: %v \nPublished: %v \nDescription: %v", item.Title, item.Link, item.PubDate, item.Description)
}
func (feed RSSFeed) Format() string{
	return fmt.Sprintf("Feed Title: %v \nLink: %v \nDescription: %v", feed.Channel.Title, feed.Channel.Link, feed.Channel.Description)
}