package link

type Link struct {
	ID         string `json:"id,omitempty" redis:"id"`
	URL        string `json:"url,omitempty" redis:"url"`
	ShortCode  string `json:"short_code,omitempty" redis:"short_code"`
	ClickCount int    `json:"click_count,omitempty" redis:"click_count"`
}
