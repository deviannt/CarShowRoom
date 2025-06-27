package routes

import (
	"autosalon/controllers"
	"autosalon/middleware"
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*.html
var templatesFS embed.FS

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 🧩 HTML шаблоны и статика
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "./static")

	// 🔓 Публичные HTML страницы
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/cars")
	})
	r.GET("/register", controllers.ShowRegisterPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.GET("/logout", controllers.Logout)
	r.GET("/cars", controllers.ShowCarsPage)

	// 🔒 Авторизованные страницы
	r.GET("/profile", middleware.AuthMiddleware(), controllers.ShowProfilePage)
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
	r.GET("/support", middleware.AuthMiddleware(), controllers.ShowSupportPage)

	// 🔒 Админ-панели (HTML)
	r.GET("/admin/users", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title":   "Пользователи",
			"Content": "admin_users.html",
		})
	})
	r.GET("/admin/posts", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.ShowAdminPostsPage)
	r.GET("/admin/cars", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), controllers.ShowAdminCarsPage)

	r.GET("/cars/:id", controllers.CarDetailPage)

	// ✅ REST API
	api := r.Group("/api")
	{
		// 🔓 Публичный API
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
			secured.POST("/cars", controllers.CreateCar)

			// 📝 Посты
			secured.POST("/posts", controllers.CreatePost)

			// 💬 Поддержка
			secured.GET("/support", controllers.GetSupportMessages)
			secured.POST("/support", controllers.SendSupportMessage)

			// 🛠️ Админ
			admin := secured.Group("/")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/admin/cars/pending", controllers.ListPendingCars)
				admin.PUT("/admin/cars/:id/approve", controllers.ApproveCar)

				admin.PUT("/cars/:id", controllers.UpdateCar)
				admin.DELETE("/cars/:id", controllers.DeleteCar)

				admin.GET("/users", controllers.ListUsers)
				admin.PUT("/users/:id/block", controllers.BlockUser)
				admin.PUT("/users/:id/username", controllers.UpdateUsername)

				admin.GET("/posts", controllers.ListUnapprovedPosts)
				admin.PUT("/posts/:id", controllers.UpdatePost)
				admin.PUT("/posts/:id/approve", controllers.ApprovePost)
				admin.DELETE("/posts/:id", controllers.DeletePost)
			}

			// 🔥 Супер-админ
			superadmin := secured.Group("/")
			superadmin.Use(middleware.RoleMiddleware("superadmin"))
			{
				superadmin.DELETE("/users/:id", controllers.DeleteUser)
				superadmin.PUT("/users/:id/role", controllers.SetUserRole)
				superadmin.PUT("/users/:id/unblock", controllers.UnblockUser)
			}
		}
	}

	return r
}
