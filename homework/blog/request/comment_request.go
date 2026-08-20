package request

type CommentCreateRequest struct {
	ArticleId int    `json:"articleId" binding:"required,gte=1"`
	AnswerId  int    `json:"answerId"`
	Comment   string `json:"comment" binding:"required"`
	Creator   string `json:"-"`
}

type CommentDeleteRequest struct {
	Id       int    `json:"id" binding:"required,gte=1"`
	Operator string `json:"-"`
}
