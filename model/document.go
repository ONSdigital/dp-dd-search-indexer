package model

// Document contains the base properties for all documents stored in Elastic Search.
type Document struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type Area struct {
	Title  string `json:"title"`   // labels.label[1]
	Type   string `json:"type"`    // areaType.codename
	TypeId string `json:"type_id"` // areaType.codename
}

type Dataset struct {
	ID                  string                        `json:"id"`
	Title               string                        `json:"title"`
	URL                 string                        `json:"url,omitempty"`
	Metadata            *Metadata                     `json:"metadata,omitempty"`
	Dimensions          []*Dimension                  `json:"dimensions,omitempty"`
	GeographicHierarchy []*GeographicHierarchySummary `json:"geographic_hierarchy,omitempty"`
}

type GeographicHierarchySummary struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	AreaTypes []*AreaType `json:"area_types"`
}

type AreaType struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Level int    `json:"level"`
}

type Metadata struct {
	Description        string         `json:"description,omitempty"`
	Taxonomies         []string       `json:"taxonomies,omitempty"`
	Contact            *Contact       `json:"contact,omitempty"`
	ReleaseDate        string         `json:"release_date"`
	NextRelease        string         `json:"next_release_date,omitempty"`
	NationalStatistics bool           `json:"is_national_statistic,omitempty"`
	Publications       []string       `json:"associated_publications,omitempty"`
	Methodology        []*Methodology `json:"methodology,omitempty"`
	TermsAndConditions string         `json:"terms_and_conditions,omitempty"`
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
	SelectedOption *DimensionOption   `json:"selected_option,omitempty"`
}

type DimensionOption struct {
	ID   string `json:"id"`
	Name string `json:"name"` // Male

	Options []*DimensionOption `json:"options,omitempty"`
}
