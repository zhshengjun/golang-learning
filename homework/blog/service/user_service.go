package service

import (
	"blog/entity"
	"blog/request"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) UserInfoById(id *int) (request.UserResponse, error) {

	if id == nil {
		return request.UserResponse{}, nil
	}

	var user entity.User
	response := request.UserResponse{}
	s.DB.Model(entity.User{}).Where("id = ?", *id).Find(&user)

	if &user == nil || user.ID == 0 {
		return response, errors.New("user not found")
	}

	response.UserName = user.UserName
	response.Email = user.Email

	return response, nil
}

func (s *UserService) UserPasswordByName(userName *string) (string, error) {

	if userName == nil {
		return "", nil
	}
	var user entity.User

	s.DB.Model(entity.User{}).Where("user_name = ?", *userName).Find(&user)

	if &user == nil || user.ID == 0 {
		return "", errors.New("user not found")
	}

	return user.Password, nil
}

func (s *UserService) Register(createRequest *request.UserCreateRequest) error {
	// 先查询用户名是否存在
	var query entity.User
	s.DB.Model(entity.User{}).Where("user_name = ?", createRequest.UserName).Find(&query)
	if query.UserName != "" {
		log.Println("user name is exist")
		return errors.New("user name is exist")
	}
	s.DB.Model(entity.User{}).Where("email = ?", createRequest.Email).Find(&query)

	if query.Email != "" {
		log.Println("email is exist")
		return errors.New("email is exist")
	}

	var createUser entity.User
	createUser.UserName = createRequest.UserName
	createUser.Email = createRequest.Email
	createUser.Password = HashPassword(createRequest.Password)
	createUser.Status = true
	createUser.Operator = createRequest.UserName

	err := s.DB.Model(&entity.User{}).Create(&createUser).Error
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s *UserService) Update(updateRequest *request.UserUpdateRequest) error {

	// 先查是否存在
	var user entity.User
	s.DB.Model(&entity.User{}).Where("id = ?", updateRequest.Id).Find(&user)

	if &user == nil || user.ID == 0 {
		return errors.New("user does not exist")
	}

	//user.ID = updateRequest.Id
	user.UserName = updateRequest.UserName
	user.Email = updateRequest.Email
	user.Password = HashPassword(updateRequest.Password)

	s.DB.Model(&user).
		UpdateColumns(map[string]any{
			"user_name": user.UserName,
			"password":  user.Password,
			"email":     user.Email,
		})
	return nil
}

func (s *UserService) Delete(deletedRequest *request.UserDeletedRequest) error {
	// 先查是否存在
	var user entity.User
	s.DB.Model(&entity.User{}).Where("id = ?", deletedRequest.Id).Find(&user)

	if &user == nil || user.ID == 0 {
		return errors.New("user not exist")
	}

	if deletedRequest.Operator != user.UserName {
		return errors.New("deleting other people's accounts is not permitted")
	}

	if user.UserName != deletedRequest.UserName {
		return errors.New("input userName is no match")
	}

	s.DB.Model(&user).
		//Where("id = ?", deletedRequest.Id).
		Delete(&entity.User{})

	return nil
}
