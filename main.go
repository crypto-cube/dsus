package main

import (
	"crypto/subtle"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const (
	port = ":8080"
)

var isDebug = "true"
var version = "dev"

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024, // 200 MB
	})

	authUser := os.Getenv("DSUS_USER")
	authPass := os.Getenv("DSUS_PASS")
	if authUser != "" && authPass != "" {
		app.Use(basicauth.New(basicauth.Config{
			Authorizer: func(user, pass string, _ fiber.Ctx) bool {
				return subtle.ConstantTimeCompare([]byte(user), []byte(authUser)) == 1 &&
					subtle.ConstantTimeCompare([]byte(pass), []byte(authPass)) == 1
			},
		}))
	}

	certsDir := "."
	filesDir := "."
	if isDebug != "true" {
		certsDir = "/etc/dsus"
		filesDir = "/var/lib/dsus/files"
		os.MkdirAll(filesDir, 0755)
	}

	log.Println("Serving content from " + filesDir)

	svc := NewUpdateService(certsDir, filesDir)
	registerRoutes(app, svc)

	if filesDir == "." {
		app.Use("/", static.New("./files"))
	} else {
		app.Use("/", static.New(filesDir))
	}

	log.Fatal(app.Listen(port))
}
