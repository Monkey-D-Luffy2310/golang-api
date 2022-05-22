package main

import (
	"golang_api/config"
	"golang_api/controller"
	"golang_api/middleware"
	"golang_api/repository"
	"golang_api/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)

	authService service.AuthService = service.NewAuthService(userRepository)
	userService service.UserService = service.NewUserService(userRepository)
	bookService service.BookService = service.NewBookService(bookRepository)
	jwtService  service.JWTService  = service.NewJWTService()

	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRouters := r.Group("/api/auth")
	{
		authRouters.POST("/login", authController.Login)
		authRouters.POST("/register", authController.Register)
	}

	userRouters := r.Group("/api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRouters.GET("/profile", userController.Profile)
		userRouters.PUT("/profile", userController.Update)
	}

	r.GET("/api/books", bookController.FindAll)
	r.GET("/api/books/:id", bookController.FindById)
	bookRouters := r.Group("/api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRouters.POST("/", bookController.Insert)
		bookRouters.PUT("/:id", bookController.Update)
		bookRouters.DELETE("/:id", bookController.Delete)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
