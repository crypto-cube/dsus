package main

import (
	"crypto/subtle"
	"log"
	"os"

	"github.com/crypto-cube/dsus/services"
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

	devicesPrefix := os.Getenv("DSUS_DEVICES_PREFIX")

	certsDir := "."
	filesDir := "."
	wgDir := "./wg"
	if isDebug != "true" {
		if authUser == "" || authPass == "" {
			log.Fatal("DSUS_USER and DSUS_PASS must be set")
		}
		if devicesPrefix == "" {
			log.Fatal("DSUS_DEVICES_PREFIX must be set")
		}
		certsDir = "/etc/dsus"
		filesDir = "/var/lib/dsus/files"
		wgDir = "/var/lib/dsus/wg"
		os.MkdirAll(filesDir, 0755)
		os.MkdirAll(wgDir, 0755)
	}

	log.Println("Serving content from " + filesDir)

	svc := services.NewUpdateService(certsDir, filesDir)
	provSvc := services.NewProvisionService(devicesPrefix, wgDir)
	registerRoutes(app, svc, provSvc)

	if filesDir == "." {
		app.Use("/", static.New("./files"))
	} else {
		app.Use("/", static.New(filesDir))
	}

	log.Fatal(app.Listen(port))
}
