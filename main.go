package main

import (
	"flag"
	_ "image/jpeg"
	"log"
	"net/url"
	"os"

	"github.com/nbtca/photo-organize/utils"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", ".", "dir for files to organize")
	flag.Parse()

	log.Printf("dir: %s\n", dir)

	// open and decode image file
	tempDir, _ := os.MkdirTemp(dir, ".temp")
	log.Println("Temp dir: ", tempDir)

	// list all files under the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var currentId string
	var currentStart int
	for i, file := range files {
		if file.IsDir() {
			log.Printf("Ignore because is a dir: %v", file.Name())
			continue
		}
		ext := utils.GetFileExt(file.Name())
		if !utils.IsPicExt(ext) {
			log.Printf("Ignore because ext not allowed: %v", file.Name())
			continue
		}

		// open file
		filePath := dir + "/" + file.Name()
		res, err := utils.DecodeImageCV(filePath)
		if err != nil {
			log.Printf("No QR code found in file: %v", filePath)
			if currentId == "" {
				os.Mkdir(tempDir+"/"+"headless", 0755)
				err := utils.CopyFile(filePath, tempDir+"/"+"headless"+"/"+file.Name())
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Move headless file: %v to dir: %v", file.Name(), tempDir+"/"+"headless")
			}
			continue
		}
		log.Printf("QR code found in file: %v, content: %v", filePath, res)
		// get id from res's query param
		// res: http://localhost:8080/?id=123
		u, err := url.Parse(res)
		if err != nil {
			log.Fatal("error at parse url", err)
		}
		id := u.Query().Get("id")

		if currentId == "" {
			currentId = id
			currentStart = i
		} else if currentId == id {
			newDir := tempDir + "/" + currentId
			log.Printf("Create new dir: %v", newDir)
			err := os.Mkdir(newDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
			for j := currentStart; j < i; j++ {
				// copy file
				err = utils.CopyFile(dir+"/"+files[j].Name(), newDir+"/"+files[j].Name())
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Move file: %v to dir: %v", files[j].Name(), newDir)
			}
			currentId = ""
			currentStart = -1
		} else {
			// throw error
			log.Fatalf("QR code not match, currentId: %v, newId: %v", currentId, res)
		}
	}

	if currentId != "" {
		// throw error
		log.Fatalf("QR code not match, currentId: %v", currentId)
	}
}
