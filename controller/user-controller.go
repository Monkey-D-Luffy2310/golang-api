package controller

import (
	"fmt"
	"golang_api/dto"
	"golang_api/helper"
	"golang_api/service"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Update(c *gin.Context)
	Profile(c *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (u *userController) Update(c *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errBind := c.ShouldBind(&userUpdateDTO)
	if errBind != nil {
		res := helper.BuildErrorResponse("Failed to process request", errBind.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := c.GetHeader("Authorization")
	token, errToken := u.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	user := u.userService.Update(userUpdateDTO)
	res := helper.BuildResponse(true, "Update user success", user)
	c.JSON(http.StatusOK, res)
}

func (u *userController) Profile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token, err := u.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := u.userService.Profile(id)
	res := helper.BuildResponse(true, "OK!", user)
	c.JSON(http.StatusOK, res)
}
