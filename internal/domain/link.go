package link

type Link struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
}
