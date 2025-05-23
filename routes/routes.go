package routes

import (
	"autosalon/controllers"
	"autosalon/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ Шаблоны и статика
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	// ✅ Публичные HTML-страницы
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/cars")
	})
	r.GET("/register", controllers.ShowRegisterPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/logout", controllers.Logout)
	r.GET("/cars", controllers.ShowCarsPage)

	// ✅ Защищённые HTML-страницы
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ShowProfilePage)

	r.GET("/admin/posts", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.ShowAdminPostsPage)
	r.GET("/admin/users", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":   "Пользователи",
			"Content": "admin_users.html",
		})
	})

	r.GET("/mycars", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":   "Мои объявления",
			"Content": "mycars.html",
		})
	})

	r.GET("/cars/add", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":   "Добавить автомобиль",
			"Content": "car_add.html",
		})
	})

	// ✅ REST API
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

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
			secured.POST("/cars", controllers.CreateCar)

			// 📝 Посты
			secured.POST("/posts", controllers.CreatePost)

			// 🛠️ Админ
			admin := secured.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.PUT("/cars/:id", controllers.UpdateCar)
				admin.DELETE("/cars/:id", controllers.DeleteCar)

				admin.GET("/users", controllers.ListUsers)
				admin.PUT("/users/:id/block", controllers.BlockUser)
				admin.PUT("/users/:id/username", controllers.UpdateUsername)

				admin.GET("/posts", controllers.ListUnapprovedPosts)
				admin.PUT("/posts/:id/approve", controllers.ApprovePost)
				admin.DELETE("/posts/:id", controllers.DeletePost)
			}

			// 🔥 Супер-админ
			superadmin := secured.Group("/")
			superadmin.Use(middleware.RoleMiddleware("superadmin"))
			{
				superadmin.DELETE("/users/:id", controllers.DeleteUser)
				superadmin.PUT("/users/:id/role", controllers.SetUserRole)
				superadmin.PUT("/users/:id/unblock", controllers.UnblockUser) // ✅ Добавлено
			}
		}
	}

	return r
}
