package models

type AppStoreReview struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Rating    int    `json:"rating"`
	UpdatedAt string `json:"updatedAt"`
}
