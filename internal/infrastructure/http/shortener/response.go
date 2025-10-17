package shortener

import "github.com/javito2003/shortener_url/internal/domain"

type linkResponse struct {
	ID         string `json:"id"`
	ShortCode  string `json:"short_code"`
	URL        string `json:"url"`
	ClickCount int    `json:"click_count"`
}

func toLinkResponse(link *domain.Link) *linkResponse {
	return &linkResponse{
		ID:         link.ID,
		ShortCode:  link.ShortCode,
		URL:        link.URL,
		ClickCount: link.ClickCount,
	}
}
