package websites

import (
	"github.com/sirupsen/logrus"
	"search_module/scraper"
	"time"
)

type ceneoCheck struct {
	queueStorage int
	queueThreads int
	domain       string
	domainGlob   string
	baseUrl      string
	delay        time.Duration
	randomDelay  time.Duration
}

func newCeneoCheck() *ceneoCheck {
	return &ceneoCheck{
		queueStorage: 100,
		queueThreads: 4,
		domain:       "www.ceneo.pl",
		domainGlob:   "www.ceneo.pl/*",
		delay:        3 * time.Second,
		randomDelay:  1 * time.Second,
	}
}

func (cc *ceneoCheck) GetResults(url string) (scraper.CheckResult, error) {

	result, err := cc.check(url)
	if err != nil {
		logrus.WithError(err).Error("can't process check request")
		return scraper.CheckResult{}, err
	}

	return result, nil
}

func (cc *ceneoCheck) check(url string) (scraper.CheckResult, error) {
	result := scraper.CheckResult{
		Price: "10", //TODO
	}

	return result, nil
}