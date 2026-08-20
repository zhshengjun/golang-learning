package request

type ArticleCreateRequest struct {
	Title   string `json:"title" binding:"required,lte=500"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"`
}

type ArticleUpdateRequest struct {
	Id       int    `json:"id" binding:"required,gte=1"`
	Title    string `json:"title" binding:"required,lte=500"`
	Content  string `json:"content" binding:"required"`
	Operator string `json:"operator"`
}

type ArticlePublishedRequest struct {
	Id       int    `json:"id" binding:"required,gte=1"`
	Operator string `json:"operator"`
}

type ArticleDeleteRequest struct {
	Id       int    `json:"id" binding:"required,gte=1"`
	Operator string `json:"operator"`
}

type ArticleListRequest struct {
	Author      string `json:"author"`
	CurrentPage int    `json:"currentPage" form:"currentPage"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
}
