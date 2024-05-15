package main

import (
	"github.com/smnprn/pixelize/ui"
)

func main() {
	operation := ui.HomePage()

	if operation == 0 {
		ui.ConversionPage()
	} else {
		ui.ResizePage()
	}
}
