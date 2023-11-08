package common

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

func DecodeImageData(imageData string) ([]byte, error) {
	// remove prefix if present
	withoutPrefix := regexp.MustCompile(`^data:image/[^;]+;base64,`).ReplaceAllString(imageData, "")

	imageBytes, err := base64.StdEncoding.DecodeString(withoutPrefix)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func OCRTesseract(imageBytes []byte) (string, error) {

	client := gosseract.NewClient()
	defer client.Close()

	// Set the image data
	err := client.SetImageFromBytes(imageBytes)
	if err != nil {
		fmt.Println("Error setting image:", err)
		return "", err
	}

	// Perform OCR on the image
	recognizedText, err := client.Text()
	if err != nil {
		fmt.Println("Error performing OCR:", err)
		return "", err
	}

	return recognizedText, nil

}
