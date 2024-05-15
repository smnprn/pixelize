package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/smnprn/pixelize/utils"
	"github.com/sunshineplan/imgconv"
)

func HomePage() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose an action:").
				Description("Use 'esc' or 'crtl+c' to exit").
				Options(
					huh.NewOption("Convert image", utils.CONVERSION),
					huh.NewOption("Resize image", utils.RESIZE),
				).
				Key("operation"),
		),
	)

	return form
}

func ConversionPage() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Image name or path").
				Key("oldFileName"),

			huh.NewInput().
				Title("New image name").
				Key("newFileName"),

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
				Key("format"),

			huh.NewConfirm().
				Title("Confirm?").
				Affirmative("Yes").
				Negative("No").
				Key("confirm"),
		),
	)

	return form
}

func ResizePage() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Image name or path").
				Key("fileName"),

			huh.NewSelect[int]().
				Title("How do you want to resize the image?").
				Description("Both options preserve the aspect ratio").
				Options(
					huh.NewOption("Absolute", 0),
					huh.NewOption("Percent", 1),
				).
				Key("resizeMode"),

			huh.NewInput().
				Title("Choose the size").
				Description("Raw number without symbol (e.g. '%' or 'px')").
				Key("newSizeStr"),

			huh.NewConfirm().
				Title("Confirm?").
				Affirmative("Yes").
				Negative("No").
				Key("confirm"),
		),
	)

	return form
}
