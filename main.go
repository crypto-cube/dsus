package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/gofiber/fiber"
)

const (
	port = 8443
)

var isDebug = "true"

func verifySignature(pubKey []byte, file []byte, signature []byte) bool {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		log.Fatal("Invalid public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	hash := sha256.Sum256(file)
	return (rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash[:], signature) == nil)

}

func main() {
	app := fiber.New()
	app.Settings.BodyLimit = 200 * 1024 * 1024 // 200 MB
	var ext = "."
	var fpath = "."
	var err error
	if isDebug != "true" {
		ext = "/etc/dsus"
		fpath, err = os.UserConfigDir()
		if err != nil {
			log.Panic("Config folder not found!")
		}
		fpath = path.Join(fpath, "dsus/files")
		os.MkdirAll(fpath, 755)
	}
	cer, err := tls.LoadX509KeyPair(ext+"/certs/server.crt", ext+"/certs/server.key")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Serving content from " + fpath)
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	app.Post("/upload", func(c *fiber.Ctx) {
		log.Println("fpath = " + fpath)
		e1, err := c.FormFile("executable")
		if err == nil {
			s1, err := c.FormFile("signature")
			if err == nil {
				exe := make([]byte, e1.Size)
				exefile, err := e1.Open()
				if err == nil {
					exefile.Read(exe)
					sig := make([]byte, s1.Size)
					signature, err := s1.Open()
					if err == nil {
						signature.Read(sig)
						pubKey, err := ioutil.ReadFile(ext + "/certs/publickey.pub")
						if err == nil {
							if verifySignature(pubKey, exe, sig) {
								log.Println("Serving files from:" + fpath)
								log.Println("Storing ... " + path.Join(fpath, "/latest"))
								c.SaveFile(e1, path.Join(fpath, "/latest"))
								log.Println("Storing ... " + path.Join(fpath, "/signature"))
								c.SaveFile(s1, path.Join(fpath, "/signature"))

								hash := sha256.Sum256(exe)
								log.Println("Storing ... " + path.Join(fpath, "/version"))
								fis, err := os.Create(path.Join(fpath, "/version"))
								if err == nil {
									fis.WriteString(fmt.Sprintf("%x", hash))
									c.Send("OK")
									return
								}
							}
						}
					}
				}
			}
		}

		c.Status(422).Send("ERROR: Required fields not set or signature invalid")
	})

	if fpath == "." {
		app.Static("/", "./files")
	} else {
		app.Static("/", fpath)
	}
	app.Listen(port, config)

}
