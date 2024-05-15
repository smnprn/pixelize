package fileops

import (
	"image"
	"log"
	"strings"

	"github.com/sunshineplan/imgconv"
)

func Convert(oldFileName string, newFileName string, format imgconv.Format) {
	src, err := imgconv.Open(oldFileName)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	err = imgconv.Save(newFileName, src, &imgconv.FormatOption{Format: format})
	if err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
}

func Resize(fileName string, resizeMode int, value float64) {
	src, err := imgconv.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	imgFormat, err := imgconv.FormatFromExtension(getFormat(fileName))
	if err != nil {
		log.Fatalf("failed to get image format: %v", err)
	}

	var resizedImg image.Image
	switch resizeMode {
	case 0:
		resizedImg = imgconv.Resize(src, &imgconv.ResizeOption{Width: int(value)})
	case 1:
		resizedImg = imgconv.Resize(src, &imgconv.ResizeOption{Percent: value})
	}

	err = imgconv.Save(fileName, resizedImg, &imgconv.FormatOption{Format: imgFormat})
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func getFormat(fileName string) string {
	format := strings.Split(fileName, ".")[1]
	return format
}
