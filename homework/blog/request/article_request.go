package request

type ArticleCreateRequest struct {
	Title   string `json:"title" binding:"required,lte=500"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"`
}

type ArticleUpdateRequest struct {
	ID       int64  `json:"id" binding:"required,gte=1"`
	Title    string `json:"title" binding:"required,lte=500"`
	Content  string `json:"content" binding:"required"`
	Operator string `json:"operator"`
}

type ArticleDeleteRequest struct {
	ID       int64  `json:"id" binding:"required,gte=1"`
	Title    string `json:"title" binding:"required,lte=500"`
	Operator string `json:"operator"`
}
