package ui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/smnprn/pixelize/fileops"
	"github.com/smnprn/pixelize/utils"
	"github.com/sunshineplan/imgconv"
)

func HomePage() int {
	var operation int

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose an action:").
				// Description("Lorem ipsum dolor sit amet").
				Options(
					huh.NewOption("Convert image", utils.CONVERSION),
					huh.NewOption("Resize image", utils.RESIZE),
				).
				Value(&operation),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		return -1
	}

	return operation
}

func ConversionPage() {
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
				Title("Confirm?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	if confirm {
		fileops.Convert(oldFileName, newFileName, format)
	}
}

func ResizePage() {
	var fileName string
	var resizeMode int
	var newSizeStr string
	var newSize float64
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Image name or path").
				Value(&fileName),

			huh.NewSelect[int]().
				Title("How do you want to resize the image?").
				Description("Both options preserve the aspect ratio").
				Options(
					huh.NewOption("Absolute", 0),
					huh.NewOption("Percent", 1),
				).
				Value(&resizeMode),

			huh.NewInput().
				Title("Choose the size").
				Description("Raw number without symbol (e.g. '%' or 'px')").
				Value(&newSizeStr),

			huh.NewConfirm().
				Title("Confirm?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	if confirm {
		newSize, err = strconv.ParseFloat(newSizeStr, 64)
		if err != nil {
			log.Fatal(err)
			return
		}

		fileops.Resize(fileName, resizeMode, newSize)
	}
}
