package service

import (
	"context"
	"errors"
	"log"
	"user/module/model"
	"user/module/repository"
)

type UserService struct {
	l              *log.Logger
	userRepository repository.UserRepository
}

//sth like a constructor
func NewUserService(l *log.Logger, repository repository.UserRepository) *UserService {
	return &UserService{l, repository}
}
func (u UserService) GetUsers() ([]model.User, error) {

	users, err := u.userRepository.GetUsers()
	if err != nil {
		return nil, errors.New("Cant get users")
	}

	return users, err

}
func (u UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {

	user, err := u.userRepository.GetByUsername(ctx, username)

	if err != nil {
		u.l.Println("Invalid username")
		return nil, err
	}

	return user, nil
}

//GetUserSalt

func (u UserService) GetUserSalt(username string) (string, error) {

	salt, err := u.userRepository.GetUserSalt(username)

	if err != nil {
		return "", err
	}
	return salt, nil
}
func (u UserService) UserExists(username string) error {

	err := u.userRepository.UserExists(username)

	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetUserRole(username string) (string, error) {

	role, err := u.userRepository.GetUserRole(username)

	if err != nil {
		return "", err
	}
	return role, nil
}