package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type CeneoScraper struct {
	Domain       string
	DomainGlob   string
	Delay        time.Duration
	RandomDelay  time.Duration
	QueueStorage int
	QueueThreads int
	BaseUrl      string
}

func NewCeneoScraper() *CeneoScraper {
	return &CeneoScraper{
		Domain:      "www.ceneo.pl",
		DomainGlob:  "www.ceneo.pl/*",
		Delay:       3 * time.Second,
		RandomDelay: 1 * time.Second,
	}
}

func (cs *CeneoScraper) Search(phrase string, page int) (SearchResult, error) {

	url := cs.createSearchUrl(phrase, page)

	result, err := cs.search(url, phrase, page)
	if err != nil {
		logrus.WithError(err).Error("can't process search request")
		return SearchResult{}, err
	}

	return result, nil
}

func (cs *CeneoScraper) CheckPrice(url string) (CheckResult, error) {

	result, err := cs.check(url)
	if err != nil {
		logrus.WithError(err).Error("can't process check request")
		return CheckResult{}, err
	}

	return result, nil
}

func (cs *CeneoScraper) addLimitToCollector(collector *colly.Collector) error {
	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  cs.DomainGlob,
		Delay:       cs.Delay,
		RandomDelay: cs.RandomDelay,
	})
	return err
}

func (cs *CeneoScraper) createQueue() (*queue.Queue, error) {
	q, err := queue.New(
		cs.QueueThreads,
		&queue.InMemoryQueueStorage{MaxSize: cs.QueueStorage},
	)

	if err != nil {
		logrus.WithError(err).Error("can not create ceneo queue")
		return nil, err
	}

	return q, nil
}

func (cs *CeneoScraper) check(url string) (CheckResult, error) {
	var price string

	c := colly.NewCollector(
		colly.AllowedDomains(cs.Domain),
	)

	err := cs.addLimitToCollector(c)
	if err != nil {
		logrus.WithError(err).Error("error while limiting collector")
		return CheckResult{}, err
	}

	findPriceTagOnPage(c, &price)

	err = c.Visit(url)
	if err != nil {
		logrus.WithError(err).Error("error while running collector")
		return CheckResult{}, err
	}
	return CheckResult{Price: price}, nil
}

func findPriceTagOnPage(collector *colly.Collector, price *string) {

	collector.OnHTML("html", func(h *colly.HTMLElement) {
		h.DOM.Find("meta").Each(func(_ int, s *goquery.Selection) {
			property, _ := s.Attr("property")
			if strings.EqualFold(property, "product:price:amount") {
				*price, _ = s.Attr("content")
			}
		})
	})
}
