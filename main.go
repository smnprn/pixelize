package main

import (
	"log"

	"github.com/charmbracelet/huh"
	"github.com/sunshineplan/imgconv"
)

func main() {
	var oldFileName string
	var newFileName string
	var format imgconv.Format
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Image name or path").
				Value(&oldFileName),

			huh.NewInput().
				Title("New image name").
				Value(&newFileName),

			huh.NewSelect[imgconv.Format]().
				Title("New image format").
				Options(
					huh.NewOption("JPEG", imgconv.JPEG),
					huh.NewOption("PNG", imgconv.PNG),
					huh.NewOption("GIF", imgconv.GIF),
					huh.NewOption("TIFF", imgconv.TIFF),
					huh.NewOption("BMP", imgconv.BMP),
					huh.NewOption("PDF", imgconv.PDF),
				).
				Value(&format),

			huh.NewConfirm().
				Title("Are you sure?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if confirm {
		convert(oldFileName, newFileName, format)
	}

}

func convert(oldFileName string, newFileName string, format imgconv.Format) {
	src, err := imgconv.Open(oldFileName)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	err = imgconv.Save(newFileName, src, &imgconv.FormatOption{Format: format})
	if err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
}
