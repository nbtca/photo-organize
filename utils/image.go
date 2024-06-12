package utils

import (
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"

	"errors"
	"image/jpeg"
	"io"
	"math"

	"golang.org/x/image/draw"
)

func DecodeImageCV(filePath string) (string, error) {
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	mats := make([]gocv.Mat, 0)
	path := "../model"
	wq := contrib.NewWeChatQRCode(path+"/detect.prototxt", path+"/detect.caffemodel",
		path+"/sr.prototxt", path+"/sr.caffemodel")
	got := wq.DetectAndDecode(img, &mats)
	justString := strings.Join(got, " ")
	if justString == "" {
		return "", errors.New("no QR code found")
	}
	return justString, nil
}

func DecodeImage(filePath string) (string, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}

func GetFileExt(filename string) string {
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}
	return ext
}

func IsPicExt(ext string) bool {
	exts := []string{".jpg", ".jpeg", ".png"}
	for _, e := range exts {
		if e == ext {
			return true
		}
	}
	return false
}

// If one of parameters below 1 - resize with proper aspect ratio.
//
// If image has an abnormal aspect ratio, it will be reduced to within 500 pixels.
func ResizeImage(img image.Image, width, height int) *image.RGBA {
	bounds := img.Bounds()

	if width == 0 && height == 0 {
		return nil
	}

	if width == 0 {
		width = bounds.Dx() * height / bounds.Dy()
	}
	if height == 0 {
		height = bounds.Dy() * width / bounds.Dx()
	}

	if width > 500 || height > 500 {
		scaleFactor := float64(500) / math.Max(float64(width), float64(height))
		width = int(float64(width) * scaleFactor)
		height = int(float64(height) * scaleFactor)
	}

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

	return newImg
}

// maxSize: max size in bytes. Dont uses compression if maxSize <= 0.
func ImageToJpeg(in image.Image, out *os.File, maxSize int64) error {
	quality := 100

	if err := jpeg.Encode(out, in, &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	if maxSize < 1 {
		return nil
	}

	stat, err := out.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()

	for size > maxSize && quality > 0 {
		if err = out.Truncate(0); err != nil {
			return err
		}
		if _, err = out.Seek(0, io.SeekStart); err != nil {
			return err
		}
		quality -= 10
		if err = jpeg.Encode(out, in, &jpeg.Options{Quality: quality}); err != nil {
			return err
		}
		stat, err = out.Stat()
		if err != nil {
			return err
		}
		size = stat.Size()
	}

	if quality <= 0 {
		return errors.New("can't resize image (quality below zero)")
	}

	return err
}

func ImageToBase64(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var base64Encoding string
	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)
	return base64Encoding
}
