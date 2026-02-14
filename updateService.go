package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/gofiber/fiber/v3"
)

type UpdateService struct {
	certsDir string
	filesDir string
}

func NewUpdateService(certsDir, filesDir string) *UpdateService {
	return &UpdateService{certsDir: certsDir, filesDir: filesDir}
}

func (s *UpdateService) VerifySignature(pubKey, file, signature []byte) bool {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		log.Fatal("Invalid public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	hash := sha256.Sum256(file)
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash[:], signature) == nil
}

func (s *UpdateService) StoreUpdate(c fiber.Ctx, executable, signature *multipart.FileHeader, exeBytes []byte) error {
	log.Println("Serving files from:" + s.filesDir)
	log.Println("Storing ... " + path.Join(s.filesDir, "/latest"))
	c.SaveFile(executable, path.Join(s.filesDir, "/latest"))
	log.Println("Storing ... " + path.Join(s.filesDir, "/signature"))
	c.SaveFile(signature, path.Join(s.filesDir, "/signature"))

	hash := sha256.Sum256(exeBytes)
	log.Println("Storing ... " + path.Join(s.filesDir, "/version"))
	fis, err := os.Create(path.Join(s.filesDir, "/version"))
	if err != nil {
		return err
	}
	defer fis.Close()
	fis.WriteString(fmt.Sprintf("%x", hash))
	return nil
}
