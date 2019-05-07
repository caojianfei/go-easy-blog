package database

type ArticleModel struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status int `json:"status"`
	Views int `json:"views"`
	CommentNumber int `json:"comment_number"`
	DeletedAt string `json:"deleted_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
