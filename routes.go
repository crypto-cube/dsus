package main

import (
	"io"
	"mime/multipart"
	"os"
	"path"

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

func registerRoutes(app *fiber.App, svc *UpdateService) {
	app.Get("/sversion", func(c fiber.Ctx) error {
		return c.SendString(version)
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

		pubKey, err := os.ReadFile(path.Join(svc.certsDir, "/certs/publickey.pub"))
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
