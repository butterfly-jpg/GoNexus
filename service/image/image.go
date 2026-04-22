package image

import (
	"GoNexus/common/image"
	"io"
	"log"
	"mime/multipart"
)

// RecognizeImage 服务层识别图片方法
func RecognizeImage(file *multipart.FileHeader) (string, error) {
	modelPath := ""
	labelPath := ""
	inputH, inputW := 224, 224
	imageRecognizer, err := image.NewImageRecognizer(modelPath, labelPath, inputH, inputW)
	if err != nil {
		log.Println("NewImageRecognizer failed. err: ", err)
		return "", err
	}
	defer imageRecognizer.Close()
	src, err := file.Open()
	if err != nil {
		log.Println("file.Open failed. err: ", err)
		return "", err
	}
	defer src.Close()
	imageBuffer, err := io.ReadAll(src)
	if err != nil {
		log.Println("io.ReadAll failed. err: ", err)
		return "", err
	}
	return imageRecognizer.PredictFromBuffer(imageBuffer)
}
