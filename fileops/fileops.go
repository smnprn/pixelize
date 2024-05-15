package fileops

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/sunshineplan/imgconv"
)

var backupFile = "/tmp/backup"

func Convert(oldFileName string, newFileName string, format imgconv.Format) error {
	src, err := imgconv.Open(oldFileName)
	if err != nil {
		return errors.New(fmt.Sprint("failed to open ", oldFileName))
	}

	err = imgconv.Save(newFileName, src, &imgconv.FormatOption{Format: format})
	if err != nil {
		return errors.New(fmt.Sprint("failed to save ", newFileName))
	}

	return nil
}

func Resize(fileName string, resizeMode int, value float64) error {
	src, err := imgconv.Open(fileName)
	if err != nil {
		return errors.New(fmt.Sprint("failed to open " + fileName))
	}

	imgFormat, err := imgconv.FormatFromExtension(getFormat(fileName))
	if err != nil {
		return fmt.Errorf("could not find image format")
	}

	err = backupImg(src, &imgconv.FormatOption{Format: imgFormat})
	if err != nil {
		return errors.New("could not backup file")
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
		restoreBackup(fileName, &imgconv.FormatOption{Format: imgFormat})
		return fmt.Errorf("failed to save image, invalid size")
	}

	return nil
}

func backupImg(img image.Image, format *imgconv.FormatOption) error {
	err := imgconv.Save(backupFile, img, format)
	if err != nil {
		return err
	}

	return nil
}

func restoreBackup(fileName string, format *imgconv.FormatOption) error {
	defer os.Remove(backupFile)

	src, err := imgconv.Open(backupFile)
	if err != nil {
		return errors.New("failed to open backup")
	}

	err = imgconv.Save(fileName, src, format)
	if err != nil {
		return errors.New("failed to restore backup")
	}

	return nil
}

func getFormat(fileName string) string {
	format := strings.Split(fileName, ".")[1]
	return format
}
