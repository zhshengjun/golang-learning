package service

import (
	"blog/entity"
	"blog/request"
	"errors"
	"fmt"
	"log"
	"time"

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

func (s *UserService) Register(createRequest request.UserCreateRequest) error {
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
	createUser.Password = createRequest.Password
	createUser.CreateAt = time.Now()
	createUser.UpdateAt = time.Now()

	err := s.DB.Model(&entity.User{}).Create(&createUser).Error
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s *UserService) Update(updateRequest request.UserUpdateRequest) error {

	// 先查是否存在
	var user entity.User
	s.DB.Model(&entity.User{}).Where("id = ?", updateRequest.Id).Find(&user)

	if &user == nil || user.ID == 0 {
		return errors.New("user id is exist")
	}

	//user.ID = updateRequest.Id
	user.UserName = updateRequest.UserName
	user.Email = updateRequest.Email
	user.Password = updateRequest.Password

	s.DB.Model(&user).
		UpdateColumns(map[string]any{
			"user_name": updateRequest.UserName,
			"password":  updateRequest.Password,
			"email":     updateRequest.Email,
		})
	return nil
}

func (s *UserService) Delete(deletedRequest request.UserDeletedRequest) error {
	// 先查是否存在
	var user entity.User
	s.DB.Model(&entity.User{}).Where("id = ?", deletedRequest.Id).Find(&user)

	if &user == nil || user.ID == 0 {
		return errors.New("user id is exist")
	}

	if user.UserName != deletedRequest.UserName {
		return errors.New("input userName is no match")
	}

	s.DB.Model(&user).
		//Where("id = ?", deletedRequest.Id).
		Delete(&entity.User{})

	return nil
}
