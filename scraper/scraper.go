package scraper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"

	"cli-seo-scraper/seo"
)

type Scraper struct {
	Collector *colly.Collector
}

func NewScraper() *Scraper {
	fakeBrowser := req.DefaultClient().ImpersonateChrome()
	c := colly.NewCollector(
		colly.UserAgent(fakeBrowser.Headers.Get("user-agent")),
	)
	c.SetClient(&http.Client{
		Transport: fakeBrowser.Transport,
	})

	return &Scraper{
		Collector: c,
	}
}

func (s *Scraper) WithSEOSettings(url string) *seo.SEOSettings {
	settings := seo.SEOSettings{}
	s.Collector.OnHTML(`head title`, func(h *colly.HTMLElement) {
		settings.Title = h.Text
	})
	s.Collector.OnHTML(`head meta[name="robots"]`, func(h *colly.HTMLElement) {
		settings.Robots = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="description"]`, func(h *colly.HTMLElement) {
		settings.Description = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="keywords"]`, func(h *colly.HTMLElement) {
		settings.Keywords = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="dc.publisher"]`, func(h *colly.HTMLElement) {
		settings.DCPublisher = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="dc.title"]`, func(h *colly.HTMLElement) {
		settings.DCTitle = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="dc.description"]`, func(h *colly.HTMLElement) {
		settings.DCDescription = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="og:title"]`, func(h *colly.HTMLElement) {
		settings.OGTitle = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="og:description"]`, func(h *colly.HTMLElement) {
		settings.OGDescription = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="og:image"]`, func(h *colly.HTMLElement) {
		settings.OGImage = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="twitter:title"]`, func(h *colly.HTMLElement) {
		settings.TCTitle = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="twitter:description"]`, func(h *colly.HTMLElement) {
		settings.TCDescription = h.Attr("content")
	})
	s.Collector.OnHTML(`head meta[name="twitter:image"]`, func(h *colly.HTMLElement) {
		settings.TCImage = h.Attr("content")
	})
	s.Collector.OnHTML(`head script[type="application/ld+json"]`, func(h *colly.HTMLElement) {
		settings.JSONLD = h.Text
	})
	s.Collector.Visit(url)
	return &settings
}
