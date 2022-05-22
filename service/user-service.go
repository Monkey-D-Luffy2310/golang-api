package service

import (
	"golang_api/dto"
	"golang_api/entity"
	"golang_api/repository"
	"log"

	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (u *userService) Update(user dto.UserUpdateDTO) entity.User {
	userUpdate := entity.User{}
	err := smapping.FillStruct(&userUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	updatedUser := u.userRepository.UpdateUser(userUpdate)
	return updatedUser
}

func (u *userService) Profile(userID string) entity.User {
	return u.userRepository.ProfileUser(userID)
}
