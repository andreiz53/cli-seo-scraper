package scraper

import (
	"net/http"
	"sync"

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
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, url := range s.Config.Websites {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c := NewCollector()
			settings := seo.SEOSettings{}
			c.OnHTML(`head title`, func(h *colly.HTMLElement) {
				settings.Title = h.Text
			})
			c.OnHTML(`head meta[name="robots"]`, func(h *colly.HTMLElement) {
				settings.Robots = h.Attr("content")
			})
			c.OnHTML(`head meta[name="description"]`, func(h *colly.HTMLElement) {
				settings.Description = h.Attr("content")
			})
			c.OnHTML(`head meta[name="keywords"]`, func(h *colly.HTMLElement) {
				settings.Keywords = h.Attr("content")
			})
			c.OnHTML(`head meta[name="dc.publisher"]`, func(h *colly.HTMLElement) {
				settings.DCPublisher = h.Attr("content")
			})
			c.OnHTML(`head meta[name="dc.title"]`, func(h *colly.HTMLElement) {
				settings.DCTitle = h.Attr("content")
			})
			c.OnHTML(`head meta[name="dc.description"]`, func(h *colly.HTMLElement) {
				settings.DCDescription = h.Attr("content")
			})
			c.OnHTML(`head meta[name="og:title"]`, func(h *colly.HTMLElement) {
				settings.OGTitle = h.Attr("content")
			})
			c.OnHTML(`head meta[name="og:description"]`, func(h *colly.HTMLElement) {
				settings.OGDescription = h.Attr("content")
			})
			c.OnHTML(`head meta[name="og:image"]`, func(h *colly.HTMLElement) {
				settings.OGImage = h.Attr("content")
			})
			c.OnHTML(`head meta[name="twitter:title"]`, func(h *colly.HTMLElement) {
				settings.TCTitle = h.Attr("content")
			})
			c.OnHTML(`head meta[name="twitter:description"]`, func(h *colly.HTMLElement) {
				settings.TCDescription = h.Attr("content")
			})
			c.OnHTML(`head meta[name="twitter:image"]`, func(h *colly.HTMLElement) {
				settings.TCImage = h.Attr("content")
			})
			c.OnHTML(`head script[type="application/ld+json"]`, func(h *colly.HTMLElement) {
				settings.JSONLD = h.Text
			})

			c.OnScraped(func(r *colly.Response) {
				mu.Lock()
				seoSettings = append(seoSettings, settings)
				mu.Unlock()
			})
			c.Visit(url)
		}(url)
	}
	wg.Wait()
	return seoSettings
}
