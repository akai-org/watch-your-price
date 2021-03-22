package ceneo

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

// SearchForItem will use the name to search for products on a page and return results
func SearchForItem(name string, firstPageOnly bool) map[string]string {

	results := make(map[string]string)

	url := "https://www.ceneo.pl/;szukaj-" + name

	c := colly.NewCollector(
		colly.AllowedDomains("www.ceneo.pl"),
	)

	q, _ := queue.New(
		4,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	addLimitToCollector(c)

	if !firstPageOnly {
		addNextButtonHandling(c, q)
	}

	handleItemsOnGridView(c, results)

	handleItemsOnListView(c, results)

	//TODO reduce number of results or add some kind of pagination

	q.AddURL(url)
	q.Run(c)

	c.Wait()

	return results
}

// CheckPrice checks the price of the item at the given url
func CheckPrice(url string) string {

	var price string

	c := colly.NewCollector(
		colly.AllowedDomains("www.ceneo.pl"),
	)

	addLimitToCollector(c)

	findPriceTagOnPage(c, &price)

	c.Visit(url)
	return price
}

func addLimitToCollector(collector *colly.Collector) {
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "www.ceneo.pl/*",
		Delay:       3 * time.Second,
		RandomDelay: 1 * time.Second,
	})
}

func addNextButtonHandling(collector *colly.Collector, queue *queue.Queue) {
	collector.OnHTML("a.js_pagination-top-next", func(h *colly.HTMLElement) {
		link := h.Request.AbsoluteURL(h.Attr("href"))
		queue.AddURL(h.Request.AbsoluteURL(link))
	})
}

func handleItemsOnGridView(collector *colly.Collector, results map[string]string) {
	collector.OnHTML("div.grid-row", func(h *colly.HTMLElement) {
		linkTag := h.DOM.Find("a").First()
		if linkTag.HasClass("go-to-shop") {
			return
		}
		relativeLink, _ := linkTag.Attr("href")
		link := h.Request.AbsoluteURL(relativeLink)
		name := linkTag.SiblingsFiltered("div.grid-item__caption").Find("Strong").First().Text()
		results[strings.TrimSpace(name)] = link
	})
}

func handleItemsOnListView(collector *colly.Collector, results map[string]string) {
	collector.OnHTML("strong.cat-prod-row__name", func(h *colly.HTMLElement) {
		linkTag := h.DOM.Find("a").First()
		if linkTag.HasClass("go-to-shop") {
			return
		}
		relativeLink, _ := linkTag.Attr("href")
		link := h.Request.AbsoluteURL(relativeLink)
		results[strings.TrimSpace(linkTag.Text())] = link
	})
}

func findPriceTagOnPage(collector *colly.Collector, price *string) {

	collector.OnHTML("h1.js_product-h1-link", func(h *colly.HTMLElement) {
		h.DOM.ParentsUntil("~").Find("meta").Each(func(_ int, s *goquery.Selection) {
			property, _ := s.Attr("property")
			if strings.EqualFold(property, "product:price:amount") {
				*price, _ = s.Attr("content")
			}
		})
	})
}