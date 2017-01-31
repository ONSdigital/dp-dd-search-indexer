package model

// Document contains the base properties for all documents stored in Elastic Search.
type Document struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type Dataset struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url,omitempty"`
	Metadata    *Metadata    `json:"metadata,omitempty"`
	Dimensions  []*Dimension `json:"dimensions,omitempty"`
}

type Metadata struct {
	Description        string         `json:"description,omitempty"`
	Taxonomies         []string       `json:"taxonomies,omitempty"`
	Contact            *Contact       `json:"contact,omitempty"`
	ReleaseDate        string         `json:"releaseDate"`
	NextRelease        string         `json:"nextReleaseDate,omitempty"`
	NationalStatistics bool           `json:"nationalStatistics,omitempty"`
	Publications       []string       `json:"associatedPublications,omitempty"`
	Methodology        []*Methodology `json:"methodology,omitempty"`
	TermsAndConditions string         `json:"termsAndConditions,omitempty"`
}

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Methodology struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Dimension struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"` // Sex
	Type           string             `json:"type"` // Dimension type
	Options        []*DimensionOption `json:"options,omitempty"`
	SelectedOption *DimensionOption   `json:"selectedOption,omitempty"`
}

type DimensionOption struct {
	ID   string `json:"id"`
	Name string `json:"name"` // Male

	Options []*DimensionOption `json:"options,omitempty"`
}
