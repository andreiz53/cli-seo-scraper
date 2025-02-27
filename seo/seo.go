package seo

type SEOSettings struct {
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

func CSVHeader() []string {
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

func (s *SEOSettings) ToCSVLine() []string {
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
