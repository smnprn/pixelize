package fileops

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunshineplan/imgconv"
)

func TestResize(t *testing.T) {
	tmp, err := os.MkdirTemp("/tmp", "")
	if err != nil {
		t.Fatalf("\U0001F6AB Failed o create temporary directory %v", err)
	}
	defer os.RemoveAll(tmp)

	testFilePath := filepath.Join(tmp, "testImg.jpg")
	createSampleImage(testFilePath, t)

	err = Resize(testFilePath, 1, 50)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to resize image %v", testFilePath)
	}

	outputFile, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to open output file %v", testFilePath)
	}
	defer outputFile.Close()

	resizedImg, _, err := image.Decode(outputFile)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to decode image %v", testFilePath)
	}

	EXPECTED_WIDTH := 150
	EXPECTED_HEIGHT := 150

	bounds := resizedImg.Bounds()

	if bounds.Dx() != EXPECTED_WIDTH || bounds.Dy() != EXPECTED_HEIGHT {
		t.Errorf(
			"\U0001F6AB Image has wrong dimensions: got %dx%d, want %dx%d",
			bounds.Dx(), bounds.Dy(),
			EXPECTED_WIDTH, EXPECTED_HEIGHT,
		)
	}
}

func TestConvert(t *testing.T) {
	tmp, err := os.MkdirTemp("/tmp", "")
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to create temporary directory %v", err)
	}
	defer os.RemoveAll(tmp)

	inputFilePath := filepath.Join(tmp, "inputTestImg.jpg")
	createSampleImage(inputFilePath, t)

	outputFilePath := filepath.Join(tmp, "outputTestImg.png")
	err = Convert(inputFilePath, outputFilePath, imgconv.PNG)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to convert convert image")
	}
	defer os.RemoveAll(inputFilePath)
	defer os.RemoveAll(outputFilePath)

	inputFormat, err := getFormat(inputFilePath)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to get input file format for %v", inputFilePath)
	}

	outputFormat, err := getFormat(outputFilePath)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to get output file format for %v", outputFilePath)
	}

	if inputFormat == outputFormat {
		t.Errorf("\U0001F6AB Wrong file format: got '.%v', want '.%v'", outputFormat, imgconv.PNG)
	}
}

func createSampleImage(path string, t *testing.T) {
	width, height := 300, 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to create path %v", path)
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		t.Fatalf("\U0001F6AB Failed to encode image %v", path)
	}
}
