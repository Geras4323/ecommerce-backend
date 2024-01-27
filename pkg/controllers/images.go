package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// POST /images/upload
func UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	fmt.Println(file.Filename)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buff := bytes.NewBuffer(nil)
	io.Copy(buff, src)

	mimetype := http.DetectContentType(buff.Bytes())
	base64file := fmt.Sprintf("data:%s;base64,%s", mimetype, toBase64(buff.Bytes()))
	fileFolder := "products"
	fileName := uuid.New().String()
	filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)

	fmt.Println(filePath)

	res, err := cloud.Cld.Upload.Upload(cloud.Ctx, base64file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   fileFolder,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Save image URL to database
	// image := models.
	// err := database.Gorm

	fmt.Println(res.SecureURL)
	return c.JSON(http.StatusOK, res)
}

// res, err := cloud.Cld.Upload.Destroy(cloud.Ctx, uploader.DestroyParams{PublicID: "products/testDesdeGo2"})
