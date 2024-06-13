package model

import (
	"log"
	"path/filepath"
)

type ModelPaths struct {
	DetectCaffeModel string
	DetectPrototxt   string
	SrCaffeModel     string
	SrPrototxt       string
}

func NewModelPaths() *ModelPaths {
	paths := &ModelPaths{
		DetectCaffeModel: "model/detect.caffemodel",
		DetectPrototxt:   "model/detect.prototxt",
		SrCaffeModel:     "model/sr.caffemodel",
		SrPrototxt:       "model/sr.prototxt",
	}

	// Convert to absolute paths
	var err error
	paths.DetectCaffeModel, err = filepath.Abs(paths.DetectCaffeModel)
	if err != nil {
		log.Fatal(err)
	}
	paths.DetectPrototxt, err = filepath.Abs(paths.DetectPrototxt)
	if err != nil {
		log.Fatal(err)
	}
	paths.SrCaffeModel, err = filepath.Abs(paths.SrCaffeModel)
	if err != nil {
		log.Fatal(err)
	}
	paths.SrPrototxt, err = filepath.Abs(paths.SrPrototxt)
	if err != nil {
		log.Fatal(err)
	}

	return paths
}
