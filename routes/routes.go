package routes

import (
	"autosalon/controllers"
	"autosalon/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Загружаем все HTML-шаблоны, включая вложенные
	r.LoadHTMLGlob("templates/*.html")

	// ✅ Подключение статики (изображения, CSS и т.д.)
	r.Static("/static", "./static")

	// ✅ Главная страница
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/cars")
	})

	// ✅ Публичные HTML-страницы
	r.GET("/register", controllers.ShowRegisterPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/logout", controllers.Logout)
	r.GET("/cars", controllers.ShowCarsPage)

	// ✅ Страница профиля (только для авторизованных)
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ShowProfilePage)

	// ✅ Страница модерации постов
	r.GET("/admin/posts", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.ShowAdminPostsPage)

	r.GET("/mycars", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout", gin.H{
			"Title":   "Мои объявления",
			"Content": "mycars.html",
		})
	})

	// ✅ REST API
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 🔒 Защищённые маршруты
		secured := api.Group("/")
		secured.Use(middleware.AuthMiddleware())
		{
			// 👤 Профиль
			secured.GET("/profile", controllers.GetProfile)
			secured.PUT("/profile", controllers.UpdateProfile)
			secured.PUT("/profile/password", controllers.ChangePassword)
			secured.DELETE("/profile", controllers.DeleteProfile)
			secured.POST("/profile/avatar", controllers.UploadAvatar)

			// 🚗 Автомобили
			secured.GET("/cars", controllers.GetCars)
			secured.GET("/cars/:id", controllers.GetCar)
			secured.GET("/mycars", controllers.GetMyCars)

			// 📝 Посты
			secured.POST("/posts", controllers.CreatePost)

			// 🛠️ Админ-функции
			admin := secured.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				// 🚗 Автомобили
				admin.POST("/cars", controllers.CreateCar)
				admin.PUT("/cars/:id", controllers.UpdateCar)
				admin.DELETE("/cars/:id", controllers.DeleteCar)

				// 👥 Пользователи
				admin.GET("/users", controllers.ListUsers)
				admin.PUT("/users/:id/block", controllers.BlockUser)

				// 📝 Посты (модерация)
				admin.GET("/posts", controllers.ListUnapprovedPosts)
				admin.PUT("/posts/:id/approve", controllers.ApprovePost)
				admin.DELETE("/posts/:id", controllers.DeletePost)
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
