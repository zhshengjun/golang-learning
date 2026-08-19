package service

import (
	apperrors "blog/errors"
	"blog/request"
	"fmt"

	"gorm.io/gorm"
)

type LoginService struct {
	db          *gorm.DB
	userService *UserService
}

func NewLoginService(db *gorm.DB, userService *UserService) *LoginService {
	return &LoginService{db: db, userService: userService}
}

func (s *LoginService) VerifyPassword(loginRequest *request.LoginRequest) (bool, error) {
	userName := loginRequest.UserName
	password := loginRequest.Password

	passwordHash, err := s.userService.UserPasswordByName(&userName)
	if err != nil {
		return false, fmt.Errorf("%w: login error", apperrors.ErrBadRequest)
	}

	verifyPassword := VerifyPassword(password, passwordHash)

	if !verifyPassword {
		return false, fmt.Errorf("%w: password wrong", apperrors.ErrBadRequest)
	}

	return true, nil
}
