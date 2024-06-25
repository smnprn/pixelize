package ui

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/smnprn/pixelize/fileops"
	"github.com/sunshineplan/imgconv"
)

func terminateHomePage(m *Model) {
	homeForm := m.pages[m.currentPageIndex].(*huh.Form)
	operation := homeForm.GetInt("operation")

	switch operation {
	case 0:
		m.currentPageIndex = 1
	case 1:
		m.currentPageIndex = 2
	}
}

func terminateConversionPage(m *Model) {
	conversionForm := m.pages[m.currentPageIndex].(*huh.Form)

	oldFileName := conversionForm.GetString("oldFileName")
	newFileName := conversionForm.GetString("newFileName")
	format := conversionForm.Get("format")
	confirm := conversionForm.GetBool("confirm")

	if confirm {
		m.errStatus = fileops.Convert(oldFileName, newFileName, format.(imgconv.Format))
		m.completed = true
	} else {
		resetPages(m)
	}
}

func terminateResizePage(m *Model) {
	resizeForm := m.pages[m.currentPageIndex].(*huh.Form)

	fileName := resizeForm.GetString("fileName")
	resizeMode := resizeForm.GetInt("resizeMode")
	newSizeStr := resizeForm.GetString("newSizeStr")
	confirm := resizeForm.GetBool("confirm")

	if confirm {
		newSize, err := strconv.ParseFloat(newSizeStr, 64)
		if err != nil {
			m.errStatus = errors.New("could not parse new size")
		}

		m.errStatus = fileops.Resize(fileName, resizeMode, newSize)
		m.completed = true
	} else {
		resetPages(m)
	}
}

func resetPages(m *Model) {
	m.pages[0] = HomePage()
	m.pages[1] = ConversionPage()
	m.pages[2] = ResizePage()

	m.currentPageIndex = 0
	m.pages[m.currentPageIndex].(*huh.Form).State = huh.StateNormal
	m.completed = false
	m.errStatus = nil
}
