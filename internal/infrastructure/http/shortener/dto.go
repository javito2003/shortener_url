package shortener

type CreateShortenerRequest struct {
	URL string `json:"url" binding:"required,url"`
}
