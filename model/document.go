package model

// Document contains the base properties for all documents stored in Elastic Search.
type Document struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}
