package database

var (
	QueryUserByUsername = "SELECT `id`, `username`, `password`, `nickname` FROM users WHERE username = ?"
	QueryUserByUserId   = "SELECT `id`, `username`, `nickname` FROM users WHERE id = ?"

	QueryTagByName = "SELECT `id`, `name`, `article_number`, `description`, `created_at`, 'updated_at' from tags where name = ? AND deleted_at is null"
	QueryTagById   = "SELECT `id`, `name`, `article_number`, `description`, `created_at`, 'updated_at' from tags where id = ? AND deleted_at is null"
	InsertTag      = "INSERT INTO tags(`name`, `description`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?)"
	// DeleteTag      = "DELETE from tags WHERE id = ?"
	SoftDeleteTagById  = "UPDATE tags SET deleted_at = ? WHERE id = ?"

	CreateNewArticle = "INSERT INTO articles(user_id, title, content, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)"

	//CreateArticleTagRelation = "INSERT INTO article_tag(article_id, tag_id) VALUES(?, ?)"
	//QueryArticleById = "select `id`, `user_id`, `title`, "
)
