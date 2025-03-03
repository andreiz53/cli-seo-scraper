package scraper

import (
	"net/http"
	"slices"
	"strings"
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
	OutputSEOMetas string   `mapstructure:"output_metas" json:"output_metas"`
	OutputSEOLinks string   `mapstructure:"output_links" json:"output_links"`
}

func NewCollector(opts ...colly.CollectorOption) *colly.Collector {
	fakeBrowser := req.DefaultClient().ImpersonateChrome()
	options := slices.Concat([]colly.CollectorOption{colly.UserAgent(fakeBrowser.Headers.Get("user-agent"))}, opts)
	c := colly.NewCollector(options...)
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

func NewScraperConfig(websites []string, outputMETAs string, outputLinks string) *ScraperConfig {
	return &ScraperConfig{
		Websites:       websites,
		OutputSEOMetas: outputMETAs,
		OutputSEOLinks: outputLinks,
	}
}

func (s *Scraper) ScrapeSEO() []seo.SEOMetas {
	var wg sync.WaitGroup
	results := make(chan seo.SEOMetas, len(s.Config.Websites))

	for _, url := range s.Config.Websites {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c := NewCollector()
			settings := seo.SEOMetas{}
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
				results <- settings
			})

			c.Visit(url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var seoSettings []seo.SEOMetas
	for result := range results {
		seoSettings = append(seoSettings, result)
	}

	return seoSettings
}

func (s *Scraper) ScrapeLinks() []seo.SEOLinks {
	brokenLinks := []seo.SEOLinks{}
	var wg sync.WaitGroup
	results := make(chan seo.SEOLinks, len(s.Config.Websites))

	for _, website := range s.Config.Websites {
		wg.Add(1)
		go func(website string) {
			defer wg.Done()
			c := NewCollector(colly.Async(), colly.MaxDepth(2))
			c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
			links := seo.SEOLinks{}

			c.OnHTML("a[href]", func(h *colly.HTMLElement) {
				h.Request.Visit(h.Attr("href"))
			})
			c.OnError(func(r *colly.Response, err error) {
				if strings.HasPrefix(r.Request.URL.String(), "http") {
					links.Links = append(links.Links, seo.SEOLink{
						URL:        r.Request.URL.String(),
						StatusCode: r.StatusCode,
					})
				}
			})
			c.Visit(website)
			c.Wait()
			results <- links
		}(website)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		brokenLinks = append(brokenLinks, result)
	}
	return brokenLinks
}
