package main

import (
	"fmt"
	"time"
)

type Embed struct {
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Timestamp   JSONTime  `json:"timestamp"`
	Color       int       `json:"color"`
	Footer      footer    `json:"footer"`
	Image       image     `json:"image"`
	Thumbnail   thumbnail `json:"thumbnail"`
	Video       video     `json:"video"`
	Provider    provider  `json:"provider"`
	Author      author    `json:"author"`
	Fields      fields    `json:"fields"`
}

type footer struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

type image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type thumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type video struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type provider struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type author struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

type fields struct {
	Field []field
}

type field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// Marshaler interface for time
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// JSONTime is custom type
type JSONTime time.Time

// MarshalJSON interface for time
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
	return []byte(stamp), nil
}

func embedMessage() Embed {
	message := Embed{
		Title:       "Title",
		Description: "Description",
		URL:         "https://boardgamegeek.com/boardgame/13",
	}

	return message
}
