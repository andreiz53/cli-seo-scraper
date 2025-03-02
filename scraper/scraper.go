package scraper

import (
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/imroc/req/v3"

	"cli-seo-scraper/seo"
)

type Scraper struct {
	Collector *colly.Collector
	Config    ScraperConfig
}

type ScraperConfig struct {
	Websites       []string `mapstructure:"websites" json:"websites"`
	OutputFilename string   `mapstructure:"output_filename" json:"output_filename"`
}

func NewCollector() *colly.Collector {
	fakeBrowser := req.DefaultClient().ImpersonateChrome()
	c := colly.NewCollector(
		colly.UserAgent(fakeBrowser.Headers.Get("user-agent")),
	)
	c.SetClient(&http.Client{
		Transport: fakeBrowser.Transport,
	})

	return c
}

func NewScraper(coll *colly.Collector, cfg ScraperConfig) *Scraper {
	return &Scraper{
		Collector: coll,
		Config:    cfg,
	}
}

func NewScraperConfig(websites []string, output string) *ScraperConfig {
	return &ScraperConfig{
		Websites:       websites,
		OutputFilename: output,
	}
}

func (s *Scraper) ScrapeSEO() []seo.SEOSettings {
	var seoSettings []seo.SEOSettings
	for _, url := range s.Config.Websites {
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

		s.Collector.OnScraped(func(r *colly.Response) {
			seoSettings = append(seoSettings, settings)
		})
		s.Collector.Visit(url)

	}
	return seoSettings
}
