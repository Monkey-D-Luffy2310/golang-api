package controller

import (
	"fmt"
	"golang_api/dto"
	"golang_api/entity"
	"golang_api/helper"
	"golang_api/service"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type BookController interface {
	Insert(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	FindAll(c *gin.Context)
	FindById(c *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (b *bookController) Insert(c *gin.Context) {
	bookDTO := dto.BookCreateDTO{}
	if err := c.ShouldBind(&bookDTO); err != nil {
		res := helper.BuildErrorResponse("Failed to bind json", err.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	authHeader := c.GetHeader("Authorization")
	userID, err := strconv.ParseUint(b.getUserIDByToken(authHeader), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	bookDTO.UserID = userID
	createdBook := b.bookService.Insert(bookDTO)
	res := helper.BuildResponse(true, "Insert book success", createdBook)
	c.JSON(http.StatusOK, res)
}

func (b *bookController) Update(c *gin.Context) {
	bookDTO := dto.BookUpdateDTO{}
	if err := c.ShouldBindJSON(&bookDTO); err != nil {
		res := helper.BuildErrorResponse("Failed to bind json", err.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	bookID, err := strconv.ParseUint(c.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	bookDTO.ID = bookID

	book := b.bookService.FindByID(bookID)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "No data with given id", helper.EmptyObj{})
		c.JSON(http.StatusNotFound, res)
		return
	}

	authHeader := c.GetHeader("Authorization")
	userID := b.getUserIDByToken(authHeader)
	if b.bookService.IsAllowedEdit(userID, bookID) {
		parstUserID, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		bookDTO.UserID = parstUserID
		updatedBook := b.bookService.Update(bookDTO)
		res := helper.BuildResponse(true, "Update book success", updatedBook)
		c.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You are not allowed to edit this book", "You are not the owner", helper.EmptyObj{})
		c.JSON(http.StatusUnauthorized, res)
	}
}

func (b *bookController) Delete(c *gin.Context) {
	book := entity.Book{}
	bookID, err := strconv.ParseUint(c.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	book.ID = bookID

	findBook := b.bookService.FindByID(bookID)
	if (findBook == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "No data with given id", helper.EmptyObj{})
		c.JSON(http.StatusNotFound, res)
		return
	}

	authHeader := c.GetHeader("Authorization")
	userID := b.getUserIDByToken(authHeader)
	if b.bookService.IsAllowedEdit(userID, bookID) {
		b.bookService.Delete(book)
		res := helper.BuildResponse(true, "Delete book success", helper.EmptyObj{})
		c.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You are not allowed to delete this book", "You are not the owner", helper.EmptyObj{})
		c.JSON(http.StatusUnauthorized, res)
	}
}

func (b *bookController) FindAll(c *gin.Context) {
	books := b.bookService.FindAll()
	res := helper.BuildResponse(true, "OK!", books)
	c.JSON(http.StatusOK, res)
}

func (b *bookController) FindById(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}
	book := b.bookService.FindByID(bookID)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "No data with given id", helper.EmptyObj{})
		c.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK!", book)
		c.JSON(http.StatusOK, res)
	}
}

func (b *bookController) getUserIDByToken(token string) string {
	authToken, errToken := b.jwtService.ValidateToken(token)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := authToken.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	return userID
}
