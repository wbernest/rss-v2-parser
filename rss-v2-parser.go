package rssv2parser

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

type Rss struct {
	XMLName        xml.Name `xml:"rss"`
	Version        string   `xml:"version,attr"`
	Title          string   `xml:"channel>title"`
	Link           string   `xml:"channel>link"`
	Description    string   `xml:"channel>description"`
	Language       string   `xml:"channel>language"`
	Copyright      string   `xml:"channel>copyright"`
	ManagingEditor string   `xml:"channel>managingEditor"`
	WebMaster      string   `xml:"channel>webMaster"`
	PubDate        string   `xml:"channel>pubDate"`
	LastBuildDate  string   `xml:"channel>lastBuildDate"`
	Category       string   `xml:"channel>category"`
	Generator      string   `xml:"channel>generator"`
	Docs           string   `xml:"channel>docs"`
	Image          Image    `xml:"channel>image"`
	Cloud          Cloud    `xml:"channel>cloud"`
	TTL            string   `xml:"channel>ttl"`
	ItemList       []Item   `xml:"channel>item"`
}

type Image struct {
	Url         string `xml:"url"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Width       string `xml:"width"`
	Height      string `xml:"height"`
	Description string `xml:"description"`
}

type Cloud struct {
}

type Item struct {
	Title       string `xml:"title"`
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

// RssParseString will be used to parse strings and will return the Rss object
func RssParseString(s string) (*Rss, error) {
	rss := Rss{}
	if len(s) == 0 {
		return &rss, nil
	}

	decoder := xml.NewDecoder(strings.NewReader(s))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

// RssParseURL will be used to parse a string returned from a url and will return the Rss object
func RssParseURL(url string) (*Rss, string, error) {
	byteValue, err := getContent(url)
	if err != nil {
		return nil, "", err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(byteValue)))
	decoder.CharsetReader = charset.NewReaderLabel
	rss := Rss{}
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, "", err
	}

	return &rss, string(byteValue), nil
}

func getContent(url string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CompareItems - This function will used to compare 2 RSS xml item objects
// and will return a list of differing items
func CompareItems(firstRSS *Rss, secondRSS *Rss) []Item {
	biggerRSS := firstRSS
	smallerRSS := secondRSS
	itemList := []Item{}
	if len(secondRSS.ItemList) > len(firstRSS.ItemList) {
		biggerRSS = secondRSS
		smallerRSS = firstRSS
	} else if len(secondRSS.ItemList) == len(firstRSS.ItemList) {
		return itemList
	}

	for _, item := range smallerRSS.ItemList {
		exists := false
		for _, oldItem := range biggerRSS.ItemList {
			if len(item.Guid) > 0 && oldItem.Guid == item.Guid {
				exists = true
				break
			} else if item.PubDate == oldItem.PubDate && item.Title == oldItem.Title {
				exists = true
				break
			}
		}
		if !exists {
			itemList = append(itemList, item)
		}
	}

	return itemList
}
