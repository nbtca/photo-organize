package utils_test

import (
	"image"
	"log"
	"os"
	"testing"

	"github.com/nbtca/photo-organize/utils"
)

func TestDecodeImage(t *testing.T) {

	t.Run("TestDecodeImage", func(t *testing.T) {
		res, _ := utils.DecodeImage("./examples/qrcode.jpg")
		expected := "https://www.nbtca.space/graduation/download?id=727397df-d0f1-47af-8cad-4af952703d29"
		if res != expected {
			t.Errorf("Expected %s, got %s", expected, res)
		}
	})
}

func TestGetFileExt(t *testing.T) {
	t.Run("TestGetFileExt", func(t *testing.T) {
		res := utils.GetFileExt("test.jpg")
		expected := ".jpg"
		if res != expected {
			t.Errorf("Expected %s, got %s", expected, res)
		}
	})
}

func TestIsPicExt(t *testing.T) {
	t.Run("TestIsPicExt", func(t *testing.T) {
		res := utils.IsPicExt(".jpg")
		expected := true
		if res != expected {
			t.Errorf("Expected %t, got %t", expected, res)
		}
	})
	t.Run("TestIsPicExt", func(t *testing.T) {
		res := utils.IsPicExt(".jpeg")
		expected := true
		if res != expected {
			t.Errorf("Expected %t, got %t", expected, res)
		}
	})
	t.Run("TestIsPicExt", func(t *testing.T) {
		res := utils.IsPicExt(".png")
		expected := true
		if res != expected {
			t.Errorf("Expected %t, got %t", expected, res)
		}
	})
	t.Run("TestIsPicExt", func(t *testing.T) {
		res := utils.IsPicExt(".gif")
		expected := false
		if res != expected {
			t.Errorf("Expected %t, got %t", expected, res)
		}
	})
}

func TestImageToJpeg(t *testing.T) {
	t.Run("TestImageToJpeg", func(t *testing.T) {
		file, err := os.Open("./examples/qrcode.jpg")
		if err != nil {
			t.Error("error at open")
		}
		defer file.Close()

		// decode image
		img, _, err := image.Decode(file)
		if err != nil {
			t.Error("error at decode")
		}

		// create a temp file
		tempFile, err := os.CreateTemp(".", "temp.jpg")
		if err != nil {
			t.Error()
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		var expectedMaxSize int64 = 1000000 * 8
		utils.ImageToJpeg(img, tempFile, expectedMaxSize)
		stat, err := os.Stat(tempFile.Name())
		if err != nil {
			t.Error()
		}
		if stat.Size() > expectedMaxSize {
			t.Errorf("Expected size less than %d, got %d", expectedMaxSize, stat.Size())
		}
	})
}

func TestImageToBase64(t *testing.T) {
	t.Run("TestImageToBase64", func(t *testing.T) {
		file, err := os.Open("./examples/qrcode.jpg")
		if err != nil {
			t.Error("error at open")
		}
		defer file.Close()

		//  to jpeg
		tempFile, err := os.Create("temp.jpg")
		if err != nil {
			t.Error()
		}
		// defer func() {
		// 	tempFile.Close()
		// 	os.Remove(tempFile.Name())
		// }()

		var expectedMaxSize int64 = 1000000 * 8
		prevImg, _, err := image.Decode(file)
		if err != nil {
			t.Error("error at decode")
		}
		utils.ImageToJpeg(prevImg, tempFile, expectedMaxSize)
		tempFile.Close()

		// open temp file
		tempFile, _ = os.Open(tempFile.Name())

		// decode image
		_, _, err = image.Decode(tempFile)
		if err != nil {
			t.Error("error at decode")
		}

		res := utils.ImageToBase64(tempFile.Name())
		// iferr != nil {
		// 	t.Error("Expected not empty, got empty")
		// }
		log.Print(res)
	})
}

func TestDecodeImageCV(t *testing.T) {
	t.Run("TestDecodeImageCV", func(t *testing.T) {
		res, _ := utils.DecodeImageCV("./examples/qrcode_blur.jpg")
		expected := "https://www.nbtca.space/graduation/download?id=a16dcc58-49a5-43d3-aedd-3c41d379ecad"
		if res != expected {
			t.Errorf("Expected %s, got %s", expected, res)
		}
	})
}
