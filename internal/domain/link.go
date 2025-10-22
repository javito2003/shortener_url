package domain

import "time"

const ExpiresAtMinutes = time.Minute * 60 * 24 * 7 // 7 dias

type Link struct {
	ID         string
	URL        string
	ShortCode  string
	ClickCount int
	ExpiresAt  *time.Time
	UserID     string
}

func (l *Link) SetExpireTime() {
	expiresAt := time.Now().Add(ExpiresAtMinutes)
	l.ExpiresAt = &expiresAt
}

func (l *Link) IsExpired() bool {
	return l.ExpiresAt != nil && time.Now().After(*l.ExpiresAt)
}
