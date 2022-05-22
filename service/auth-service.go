package service

import (
	"golang_api/dto"
	"golang_api/entity"
	"golang_api/repository"
	"log"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (a *authService) VerifyCredential(email, password string) interface{} {
	res := a.userRepository.VerifyCredential(email)
	if user, ok := res.(entity.User); ok {
		comparePasswd := comparePassword(user.Password, password)
		if comparePasswd {
			return res
		}
		return false
	}
	return false
}

func (a *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := a.userRepository.InsertUser(userToCreate)
	return res
}

func (a *authService) FindByEmail(email string) entity.User {
	return a.userRepository.FindByEmail(email)
}

func (a *authService) IsDuplicateEmail(email string) bool {
	res := a.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPassword, plainPassword string) bool {
	hashedPasswd := []byte(hashedPassword)
	plainPasswd := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswd, plainPasswd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
