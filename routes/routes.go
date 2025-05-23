package routes

import (
	"autosalon/controllers"
	"autosalon/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// HTML и статика
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// Главная
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/cars")
	})

	// Публичные страницы
	r.GET("/register", controllers.ShowRegisterPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/logout", controllers.Logout)
	r.GET("/cars", controllers.ShowCarsPage)

	// Профиль
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ShowProfilePage)

	// HTML: админ-панель
	adminPages := r.Group("/admin")
	adminPages.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminPages.GET("/users", controllers.ShowAdminUsersPage)
		adminPages.GET("/cars/add", controllers.ShowCarAddPage)
		adminPages.GET("/cars/edit", controllers.ShowCarEditPage)
		adminPages.GET("/cars", controllers.ShowAdminCarsPage)
	}

	// REST API
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// Защищённые API
		secured := api.Group("/")
		secured.Use(middleware.AuthMiddleware())
		{
			// 👤 Профиль
			secured.GET("/profile", controllers.GetProfile)
			secured.PUT("/profile", controllers.UpdateProfile)
			secured.PUT("/profile/password", controllers.ChangePassword)
			secured.DELETE("/profile", controllers.DeleteProfile)
			secured.POST("/profile/avatar", controllers.UploadAvatar)

			// 🚗 Пользователь может просматривать машины
			secured.GET("/cars", controllers.GetCars)
			secured.GET("/cars/:id", controllers.GetCar)

			// 🔒 ADMIN доступ
			admin := secured.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/cars", controllers.CreateCar)
				admin.PUT("/cars/:id/approve", controllers.ApproveCar)
				admin.GET("/users", controllers.ListUsers)
			}

			// 🔥 SUPERADMIN доступ
			superadmin := secured.Group("/")
			superadmin.Use(middleware.RoleMiddleware("superadmin"))
			{
				superadmin.PUT("/cars/:id", controllers.UpdateCar)
				superadmin.DELETE("/cars/:id", controllers.DeleteCar)

				superadmin.DELETE("/users/:id", controllers.DeleteUser)
				superadmin.PUT("/users/:id/role", controllers.ChangeUserRole)
				superadmin.PUT("/users/:id/block", controllers.BlockUser)
			}
		}
	}

	return r
}
