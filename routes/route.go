package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/kominkamen/rootds-admin/controllers"
	_ "github.com/kominkamen/rootds-admin/docs"
)

func SetupUserRoutes(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/api/v1/login", controllers.HandlerAdminLogin)
	/*app.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/api/v1/login" {
			return auth.ValidateAuthToken(c)
		}
		return c.Next()
	})*/

	app.Get("/api/v1/country", controllers.HandlerGetAllCountry)
	app.Post("/api/v1/country", controllers.HandlerCreateCountry)
	app.Get("/api/v1/country/:id", controllers.HandlerGetOneCountry)
	app.Put("/api/v1/country/:id", controllers.HandleUpdateCountry)
	app.Delete("/api/v1/country/:id", controllers.HandleDeleteCountry)

	app.Get("/api/v1/role", controllers.HandlerGetAllRole)
	app.Post("/api/v1/role", controllers.HandlerCreateRole)
	app.Get("/api/v1/role/:id", controllers.HandlerGetOneRole)
	app.Put("/api/v1/role/:id", controllers.HandleUpdateRole)
	app.Delete("/api/v1/role/:id", controllers.HandleDeleteRole)

	app.Get("/api/v1/product", controllers.HandlerGetAllProduct)
	app.Post("/api/v1/product", controllers.HandlerCreateProduct)
	app.Get("/api/v1/product/:id", controllers.HandlerGetOneProduct)
	app.Put("/api/v1/product/:id", controllers.HandleUpdateProduct)
	app.Delete("/api/v1/product/:id", controllers.HandleDeleteProduct)

	app.Post("/api/v1/user", controllers.HandlerCreateUser)
	app.Get("/api/v1/user/:id", controllers.HandlerGetOneUser)
	app.Get("/api/v1/user", controllers.HandlerGetAllUser)
	app.Put("/api/v1/user/:id", controllers.HandleUpdateUser)
	app.Delete("/api/v1/user/:id", controllers.HandleDeleteUser)
}
