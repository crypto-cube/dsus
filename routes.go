package main

import (
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/crypto-cube/dsus/services"
	"github.com/gofiber/fiber/v3"
)

func readFormFile(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func registerRoutes(app *fiber.App, svc *services.UpdateService, provSvc *services.ProvisionService) {
	app.Get("/sversion", func(c fiber.Ctx) error {
		return c.SendString(version)
	})

	app.Post("/provision", func(c fiber.Ctx) error {
		configBytes, deviceName, token, err := provSvc.Provision()
		if err != nil {
			return c.Status(500).SendString("provision failed")
		}
		c.Set("X-Device-Name", deviceName)
		c.Set("X-Auth-Token", token)
		c.Set("Content-Type", "text/plain")
		return c.Send(configBytes)
	})

	app.Post("/upload", func(c fiber.Ctx) error {
		fail := func() error {
			return c.Status(422).SendString("ERROR: Required fields not set or signature invalid")
		}

		e1, err := c.FormFile("executable")
		if err != nil {
			return fail()
		}
		s1, err := c.FormFile("signature")
		if err != nil {
			return fail()
		}

		exe, err := readFormFile(e1)
		if err != nil {
			return fail()
		}
		sig, err := readFormFile(s1)
		if err != nil {
			return fail()
		}

		pubKey, err := os.ReadFile(path.Join(svc.CertsDir, "/certs/publickey.pub"))
		if err != nil {
			return fail()
		}
		if !svc.VerifySignature(pubKey, exe, sig) {
			return fail()
		}
		if err := svc.StoreUpdate(c, e1, s1, exe); err != nil {
			return fail()
		}
		return c.SendString("OK")
	})
}
