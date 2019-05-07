package database

type TagModel struct {
	Id int `json:"id"`
	Name string `json:"name"`
	ArticleNumber int `json:"article_number"`
	Description string `json:"description"`
	DeletedAt string `json:"deleted_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
} 