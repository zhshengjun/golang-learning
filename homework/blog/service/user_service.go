package service

import (
	"blog/entity"
	"blog/enums"
	apperrors "blog/errors"
	"blog/request"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) UserInfoById(id *int) (request.UserResponse, error) {

	response := request.UserResponse{}
	if id == nil {
		return response, fmt.Errorf("%w: id is require", apperrors.ErrBadRequest)
	}

	var user entity.User

	result := s.db.Model(entity.User{}).Where("id = ?", *id).Find(&user)

	if &user == nil || user.Id == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return response, fmt.Errorf("%w: user not found", apperrors.ErrNotFound)
	}

	if result.Error != nil {
		return response, fmt.Errorf("%w: query user error", result.Error)
	}

	response.UserName = user.UserName
	response.Email = user.Email

	return response, nil
}

func (s *UserService) UserPasswordByName(userName *string) (string, error) {

	if userName == nil {
		return "", fmt.Errorf("%w: user name is require", apperrors.ErrBadRequest)
	}
	var user entity.User

	result := s.db.Model(entity.User{}).Where("user_name = ? and status = 'ACTIVE'", *userName).Find(&user)
	if &user == nil || user.Id == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("%w: userName not found", apperrors.ErrNotFound)
	}
	if result.Error != nil {
		return "", fmt.Errorf("%w: query user errror", result.Error)
	}

	return user.Password, nil
}

func (s *UserService) Register(createRequest *request.UserCreateRequest) error {
	// 先查询用户名是否存在
	var query entity.User
	result := s.db.Model(entity.User{}).Where("user_name = ?", createRequest.UserName).Find(&query)
	if result.Error != nil {
		return fmt.Errorf("query user: %w", result.Error)
	}
	if query.UserName != "" {
		log.Println("user name is exist")
		return fmt.Errorf("%w: user name is exist", apperrors.ErrConflict)
	}

	result = s.db.Model(&entity.User{}).Where("email = ?", createRequest.Email).Find(&query)
	if result.Error != nil {
		return fmt.Errorf("query user: %w", result.Error)
	}

	if query.Email != "" {
		log.Println("email is exist")
		return fmt.Errorf("%w: email is exis", apperrors.ErrConflict)
	}

	var createUser entity.User
	createUser.UserName = createRequest.UserName
	createUser.Email = createRequest.Email
	createUser.Password = HashPassword(createRequest.Password)
	createUser.Status = enums.UserStatusActive
	createUser.Operator = createRequest.UserName

	result = s.db.Model(&entity.User{}).Create(&createUser)
	if result.Error != nil {
		return fmt.Errorf("register user: %w", result.Error)
	}
	return nil
}

func (s *UserService) Updated(updateRequest *request.UserUpdateRequest) error {

	// 先查是否存在
	var user entity.User
	result := s.db.Model(&entity.User{}).Where("id = ?", updateRequest.Id).Find(&user)

	if &user == nil || user.Id == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: user not found", apperrors.ErrNotFound)
	}
	if result.Error != nil {
		return fmt.Errorf("update user: %w", result.Error)
	}

	if updateRequest.Operator != user.UserName {
		return fmt.Errorf("%w: accounts is not permitted", apperrors.ErrForbidden)
	}

	if enums.UserStatusActive != user.Status {
		return fmt.Errorf("%w: user status is not active", apperrors.ErrGone)
	}

	result = s.db.Model(&entity.User{}).
		Where("id = ?", updateRequest.Id).
		UpdateColumns(map[string]any{
			"password": HashPassword(updateRequest.Password),
			"email":    updateRequest.Email,
		})
	if result.Error != nil {
		return fmt.Errorf("update user: %w", result.Error)
	}
	return nil
}

func (s *UserService) Deleted(deletedRequest *request.UserDeletedRequest) error {
	// 先查是否存在
	var user entity.User
	result := s.db.Model(&entity.User{}).Where("id = ?", deletedRequest.Id).Find(&user)

	if &user == nil || user.Id == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: user not found", apperrors.ErrNotFound)
	}

	if result.Error != nil {
		return fmt.Errorf("delete user: %w", result.Error)
	}

	if deletedRequest.Operator != user.UserName {
		return fmt.Errorf("%w: accounts is not permitted", apperrors.ErrForbidden)
	}

	if user.UserName != deletedRequest.UserName {
		return fmt.Errorf("%w: input userName is no match", apperrors.ErrConflict)
	}

	result = s.db.Model(&entity.User{}).
		Where("id = ?", deletedRequest.Id).
		UpdateColumn("status", enums.UserStatusCancelled)
	if result.Error != nil {
		return fmt.Errorf("delete user: %w", result.Error)
	}
	return nil
}
