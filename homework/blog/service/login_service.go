package service

import (
	"blog/request"
	"errors"

	"gorm.io/gorm"
)

type LoginService struct {
	DB          *gorm.DB
	UserService *UserService
}

func NewLoginService(db *gorm.DB, userService *UserService) *LoginService {
	return &LoginService{DB: db, UserService: userService}
}

func (s *LoginService) VerifyPassword(loginRequest *request.LoginRequest) (bool, error) {
	userName := loginRequest.UserName
	password := loginRequest.Password

	passwordHash, err := s.UserService.UserPasswordByName(&userName)
	if err != nil {
		return false, errors.New("登录错误")
	}

	verifyPassword := VerifyPassword(password, passwordHash)

	if !verifyPassword {
		return false, errors.New("密码错误")
	}

	return true, nil
}
