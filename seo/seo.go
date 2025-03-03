package seo

import "fmt"

type SEOMetas struct {
	Title         string
	Robots        string
	Description   string
	Keywords      string
	DCPublisher   string
	DCTitle       string
	DCDescription string
	OGTitle       string
	OGDescription string
	OGImage       string
	TCTitle       string
	TCDescription string
	TCImage       string
	JSONLD        string
}

func CSVHeaderMETAs() []string {
	return []string{
		"title",
		"robots",
		"description",
		"keywords",
		"dc_publisher",
		"dc_title",
		"dc_description",
		"og_title",
		"og_description",
		"og_image",
		"tc_title",
		"tc_description",
		"tc_image",
		"json_ld",
	}
}

func (s *SEOMetas) ToCSVLine() []string {
	return []string{
		s.Title,
		s.Robots,
		s.Description,
		s.Keywords,
		s.DCPublisher,
		s.DCTitle,
		s.DCDescription,
		s.OGTitle,
		s.OGDescription,
		s.OGImage,
		s.TCTitle,
		s.TCDescription,
		s.TCImage,
		s.JSONLD,
	}
}

type SEOLinks struct {
	Links []SEOLink
}

type SEOLink struct {
	URL        string
	StatusCode int
}

func CSVHeaderLinks() []string {
	return []string{
		"status_code",
		"url",
	}
}

func (l *SEOLink) ToCSVLine() []string {
	return []string{
		fmt.Sprint(l.StatusCode),
		l.URL,
	}
}
