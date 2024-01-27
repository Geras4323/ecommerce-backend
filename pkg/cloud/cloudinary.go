package cloud

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cld *cloudinary.Cloudinary
var Ctx context.Context

func ConnectCloudinary() {
	cld, err := cloudinary.NewFromParams("gerasecommerce", "562628734429182", "kY6vdbsrkOhLFltjp2dzd26kXkk")

	if err != nil {
		log.Fatal("CLOUDINARY: failed to connect")
	}

	Ctx = context.Background()
	Cld = cld
}

var validMimetypes = []string{"image/jpeg", "image/png"}

func chechMimetype(mimetype string) bool {
	for _, v := range validMimetypes {
		if mimetype == v {
			return true
		}
	}
	return false
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func UploadImage(file *multipart.FileHeader, folder string, name string) (*uploader.UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buff := bytes.NewBuffer(nil)
	io.Copy(buff, src)

	mimetype := http.DetectContentType(buff.Bytes())
	if !chechMimetype(mimetype) {
		return nil, errors.New("Invalid image format")
	}

	base64file := fmt.Sprintf("data:%s;base64,%s", mimetype, toBase64(buff.Bytes()))

	res, err := Cld.Upload.Upload(Ctx, base64file, uploader.UploadParams{
		PublicID: name,
		Folder:   folder,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteImage(name string) (*uploader.DestroyResult, error) {
	res, err := Cld.Upload.Destroy(Ctx, uploader.DestroyParams{PublicID: name})
	if err != nil {
		return nil, err
	}

	return res, nil
}
