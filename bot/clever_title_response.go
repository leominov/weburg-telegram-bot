package bot

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/xmlpath.v2"
)

const linkPrefix = "http://weburg.net"

var (
	xpathSerialsContent         = xmlpath.MustCompile(`//div[@class="wb-serials-content"]`)
	xpathSerialsTitle           = xmlpath.MustCompile(`.//h3/a`)
	xpathSerialsLink            = xmlpath.MustCompile(`.//h3/a/@href`)
	xpathSerialsTitleOriginal   = xmlpath.MustCompile(`.//div[@class="wb-serials-title"]`)
	xpathSerialsDescription     = xmlpath.MustCompile(`.//div[@class="wb-serials-last-series"]/a`)
	xpathSerialsCategories      = xmlpath.MustCompile(`.//ul[@class="wb-serials-tags"]/li`)
	xpathSerialsCategoriesTitle = xmlpath.MustCompile(`.//a`)
)

type CleverTitleResponse struct {
	Items         string `json:"items"`
	NextPage      bool   `json:"next_page"`
	LastElementID int    `json:"last_element_id"`
}

func (c *CleverTitleResponse) processSeriesNode(iter *xmlpath.Iter) ([]EndpointItem, error) {
	var itemList []EndpointItem
	categories := []string{}
	title, ok := xpathSerialsTitle.String(iter.Node())
	if !ok {
		return itemList, errors.New("Cant get title")
	}
	title = strings.TrimSpace(title)
	titleOriginal, ok := xpathSerialsTitleOriginal.String(iter.Node())
	if !ok {
		return itemList, errors.New("Cant get original title")
	}
	titleOriginal = strings.TrimSpace(titleOriginal)
	link, ok := xpathSerialsLink.String(iter.Node())
	if !ok {
		return itemList, errors.New("Cant get link")
	}
	link = strings.TrimSpace(link)
	description, ok := xpathSerialsDescription.String(iter.Node())
	if !ok {
		return itemList, errors.New("Cant get description")
	}
	description = strings.TrimSpace(description)
	categoriesIter := xpathSerialsCategories.Iter(iter.Node())
	for categoriesIter.Next() {
		category, ok := xpathSerialsCategoriesTitle.String(categoriesIter.Node())
		if !ok {
			continue
		}
		categories = append(categories, strings.TrimSpace(category))
	}
	itemList = append(itemList, EndpointItem{
		ID:          linkPrefix + link,
		Link:        linkPrefix + link,
		Title:       fmt.Sprintf("%s / %s", title, titleOriginal),
		Description: description,
		Categories:  categories,
	})
	return itemList, nil
}

func (c *CleverTitleResponse) ParseItems(cleverTitleType string) ([]EndpointItem, error) {
	var itemList []EndpointItem
	if len(c.Items) == 0 {
		return itemList, errors.New("Empty item list")
	}
	root, err := xmlpath.ParseHTML(strings.NewReader(c.Items))
	if err != nil {
		return itemList, err
	}
	iter := xpathSerialsContent.Iter(root)
	for iter.Next() {
		switch cleverTitleType {
		case "clever_title_series":
			items, err := c.processSeriesNode(iter)
			if err != nil {
				return itemList, err
			}
			for _, item := range items {
				itemList = append(itemList, item)
			}
		}
	}
	return itemList, nil
}
