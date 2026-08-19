package request

type UserCreateRequest struct {
	UserName string `json:"userName" binding:"required,gte=6"`
	Password string `json:"password" binding:"required,gte=8"`
	Email    string `json:"email" binding:"email"`
}

type UserUpdateRequest struct {
	Id       int64  `json:"id" binding:"required"`
	UserName string `json:"UserName" binding:"required,gte=6"`
	Password string `json:"password" binding:"required,gte=8"`
	Email    string `json:"email" binding:"email"`
}

type UserDeletedRequest struct {
	Id       int64  `json:"id" binding:"required"`
	UserName string `json:"UserName" binding:"required,gte=6"`
}

type UserListRequest struct {
	Keyword     string `json:"keyword"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
}

type UserResponse struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type UserListResponse struct {
	UserName string `json:"UserName"`
	Email    string `json:"email"`
	Operator string `json:"operator"`
	UpdateAt int64  `json:"updateAt"`
}
