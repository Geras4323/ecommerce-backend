package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var ImageErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de las imágenes",

	"NotFound": "Imagen no encontrada",
	"Parse":    "Limite inválido",
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// GET /images
func GetImages(c echo.Context) error {
	rawLimit := c.QueryParam("limit")

	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.SCTMake(ImageErrors["Parse"], err.Error()))
	}

	images := make([]models.Image, 0)

	if err := database.Gorm.Order("RAND()").Limit(limit).Find(&images).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(ImageErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, images)
}

// POST /images/test
func ImagesTest(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["images"]

	fmt.Println(files)

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer src.Close()

		fmt.Println(src)

		// buff := bytes.NewBuffer(nil)
		// io.Copy(buff, src)
		// base64file := toBase64(buff.Bytes())

		// fmt.Println(base64file)
	}

	return nil
}

// POST /images/upload
func UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buff := bytes.NewBuffer(nil)
	io.Copy(buff, src)

	mimetype := http.DetectContentType(buff.Bytes())
	base64file := fmt.Sprintf("data:%s;base64,%s", mimetype, toBase64(buff.Bytes()))
	fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "vouchers")
	fileName := uuid.New().String()

	// filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)
	// fmt.Println(filePath)

	res, err := cloud.Cld.Upload.Upload(cloud.Ctx, base64file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   fileFolder,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.SCTMake(utils.CommonErrors[utils.FileUpload], err.Error()))
	}

	// Save image URL to database
	// image := models.
	// err := database.Gorm

	fmt.Println(res.SecureURL)
	return c.JSON(http.StatusOK, res)
}

func UploadPDF(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buff := bytes.NewBuffer(nil)
	io.Copy(buff, src)

	mimetype := http.DetectContentType(buff.Bytes())
	base64file := fmt.Sprintf("data:%s;base64,%s", mimetype, toBase64(buff.Bytes()))
	fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "payments")
	fileName := uuid.New().String()

	// filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)
	// fmt.Println(filePath)

	res, err := cloud.Cld.Upload.Upload(cloud.Ctx, base64file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   fileFolder,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.SCTMake(utils.CommonErrors[utils.FileUpload], err.Error()))
	}

	return c.JSON(http.StatusOK, res)
}
