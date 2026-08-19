package enums

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "DRAFT"     // 草稿
	ArticleStatusPublished ArticleStatus = "PUBLISHED" // 已发布
	ArticleStatusDeleted   ArticleStatus = "DELETED"   // 已删除
)
