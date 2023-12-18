package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kominkamen/rootds-admin/db"
	_ "github.com/kominkamen/rootds-admin/docs"
	"github.com/kominkamen/rootds-admin/routes"
)

func serveStatic(app *fiber.App) {
	app.Static("/", "./views")
}

// @title ROOTDS API
// @version 1.0
// @description Golang GoFiber swagger auto generate step by step by swaggo
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey X-Access-Token
// @in header
// @name  X-Access-Token
func main() {
	dbPool, _ := db.ConnectToDB()
	defer dbPool.Close()

	app := fiber.New()
	serveStatic(app)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization",
	}))

	routes.SetupUserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
