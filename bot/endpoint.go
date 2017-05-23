package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ungerik/go-rss"
)

var (
	client       = &http.Client{}
	NoAuthCookie = http.Cookie{
		Name:  "session_id",
		Value: "noauth",
	}
)

type Endpoint struct {
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

type EndpointItem struct {
	ID          string
	Link        string
	Title       string
	Categories  []string
	Description string
}

func (e *Endpoint) readRSS() ([]EndpointItem, error) {
	var itemList []EndpointItem
	feed, err := rss.Read(e.URL)
	if err != nil {
		return itemList, err
	}
	if len(feed.Item) == 0 {
		return itemList, errors.New("Empty item list")
	}
	for _, i := range feed.Item {
		itemList = append(itemList, EndpointItem{
			ID:          i.GUID,
			Link:        i.Link,
			Title:       i.Title,
			Categories:  i.Category,
			Description: i.Description,
		})
	}
	return itemList, nil
}

func (e *Endpoint) getCleverTitleResponse() (CleverTitleResponse, error) {
	var cleverResponse CleverTitleResponse
	req, err := http.NewRequest("GET", e.URL, nil)
	if err != nil {
		return cleverResponse, err
	}
	req.AddCookie(&NoAuthCookie)
	req.Header.Set("X-Requested-With", "xmlhttprequest")
	resp, err := client.Do(req)
	if err != nil {
		return cleverResponse, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return cleverResponse, fmt.Errorf("Incorrect response HTTP code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cleverResponse, err
	}
	err = json.Unmarshal(body, &cleverResponse)
	if err != nil {
		return cleverResponse, err
	}
	return cleverResponse, nil
}

func (e *Endpoint) readCleverTitle(cleverTitleType string) ([]EndpointItem, error) {
	var itemList []EndpointItem
	resp, err := e.getCleverTitleResponse()
	if err != nil {
		return itemList, err
	}
	return resp.ParseItems(cleverTitleType)
}

func (e *Endpoint) Read() ([]EndpointItem, error) {
	switch strings.ToLower(e.Type) {
	case "clever_title_series":
		return e.readCleverTitle(e.Type)
	default:
		return e.readRSS()
	}
}
