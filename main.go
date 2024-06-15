package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"log"
	"net/url"
	"os"

	"github.com/nbtca/photo-organize/utils"
)

func main() {
	var dir string
	var dest string
	var dry bool
	flag.StringVar(&dir, "dir", ".", "dir for files to organize")
	flag.StringVar(&dest, "dest", "", "dir for output files")
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.Parse()
	log.Printf("dir: %s\n", dir)

	// open and decode image file
	tempDir, _ := os.MkdirTemp(dir, ".temp")
	if dry {
		defer os.RemoveAll(tempDir)
	}
	if dest == "" {
		dest = tempDir
		log.Println("dest is not set, use temp dir: ", dest)
	}
	// create dest if not exist
	os.Mkdir(dest, 0755)

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
				headlessDir := dest + "/" + "headless"
				os.Mkdir(headlessDir, 0755)
				err := utils.CopyFile(filePath, headlessDir+"/"+file.Name())
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Move headless file: %v to dir: %v", file.Name(), headlessDir)
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
			newDir := dest + "/" + currentId
			log.Printf("Create new dir: %v", newDir)
			os.Mkdir(newDir, 0755)
			content := fmt.Sprintf("{ \"id\": \"%v\" , \"status\": \"ready\"}", currentId)
			err := utils.WriteToFile(newDir+"/info.json", content)
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
