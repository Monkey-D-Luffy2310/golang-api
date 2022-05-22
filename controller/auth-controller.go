package controller

import (
	"golang_api/dto"
	"golang_api/entity"
	"golang_api/helper"
	"golang_api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (a *authController) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authResult := a.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if user, ok := authResult.(entity.User); ok {
		generatedToken := a.jwtService.GenerateToken(strconv.FormatUint(user.ID, 10))
		user.Token = generatedToken
		res := helper.BuildResponse(true, "Login successful", user)
		c.JSON(http.StatusOK, res)
		return
	}
	res := helper.BuildErrorResponse("Please check again your credential", "Invalid credential", helper.EmptyObj{})
	c.AbortWithStatusJSON(http.StatusUnauthorized, res)
}

func (a *authController) Register(c *gin.Context) {
	var registerDTO dto.RegisterDTO
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if !a.authService.IsDuplicateEmail(registerDTO.Email) {
		res := helper.BuildErrorResponse("Failed to process request", "Email already exists", helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusConflict, res)
	} else {
		user := a.authService.CreateUser(registerDTO)
		generatedToken := a.jwtService.GenerateToken(strconv.FormatUint(user.ID, 10))
		user.Token = generatedToken
		res := helper.BuildResponse(true, "Register successful", user)
		c.JSON(http.StatusOK, res)
	}
}
