package routes

import (
	"autosalon/controllers"
	"autosalon/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Загружаем все HTML-шаблоны, включая вложенные
	r.LoadHTMLGlob("templates/*.html")

	// ✅ Подключение статики (если будет)
	r.Static("/static", "./static")

	// ✅ Главная страница
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/cars") // или /login по желанию
	})

	// ✅ Публичные HTML-страницы
	r.GET("/register", controllers.ShowRegisterPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/logout", controllers.Logout) // 👈 Новый маршрут
	r.GET("/cars", controllers.ShowCarsPage)

	// ✅ Профиль (только для авторизованных)
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ShowProfilePage)

	// ✅ Админ-панель
	adminPages := r.Group("/admin")
	adminPages.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminPages.GET("/users", controllers.ShowAdminUsersPage)
		adminPages.GET("/cars/add", controllers.ShowCarAddPage)
		adminPages.GET("/cars/edit", controllers.ShowCarEditPage)
		adminPages.GET("/cars", controllers.ShowAdminCarsPage)
	}

	// ✅ REST API
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 🔒 Защищённые API
		secured := api.Group("/")
		secured.Use(middleware.AuthMiddleware())
		{
			// 👤 Профиль
			secured.GET("/profile", controllers.GetProfile)
			secured.PUT("/profile", controllers.UpdateProfile)
			secured.PUT("/profile/password", controllers.ChangePassword)
			secured.DELETE("/profile", controllers.DeleteProfile)
			secured.POST("/profile/avatar", controllers.UploadAvatar)

			// 🚗 Авто
			secured.GET("/cars", controllers.GetCars)
			secured.GET("/cars/:id", controllers.GetCar)

			// 🛠️ Админ-функции
			admin := secured.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("/cars", controllers.CreateCar)
				admin.PUT("/cars/:id", controllers.UpdateCar)
				admin.DELETE("/cars/:id", controllers.DeleteCar)

				admin.GET("/users", controllers.ListUsers)
				admin.PUT("/users/:id/block", controllers.BlockUser)
			}

			// 🔥 Супер-админ
			superadmin := secured.Group("/")
			superadmin.Use(middleware.RoleMiddleware("superadmin"))
			{
				superadmin.DELETE("/users/:id", controllers.DeleteUser)
			}
		}
	}

	return r
}
