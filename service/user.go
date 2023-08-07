package service

import (
	"context"
	"crawl/model"
	"crawl/repo"
	"errors"
)

type User struct {
	repo repo.IRepo
}

func NewUser(repo repo.IRepo) *User {
	return &User{
		repo: repo,
	}
}

type IUser interface {
	Login(ctx context.Context, userRequest model.UserRequest) (userResponse model.User, err error)
}

func (s *User) Login(ctx context.Context, userRequest model.UserRequest) (userResponse model.User, err error) {
	user, err := s.repo.GetUserByName(userRequest.Username)
	if err != nil {
		return model.User{}, err
	}
	// Kiểm tra mật khẩu có đúng hay không
	if user.Password != userRequest.Password {
		return model.User{}, errors.New("invalid username or password")
	}

	return user, nil
}
