package ui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/smnprn/pixelize/fileops"
	"github.com/sunshineplan/imgconv"
)

type model struct {
	homeForm       *huh.Form
	conversionForm *huh.Form
	resizeForm     *huh.Form
	currentPage    tea.Model
	width          int
	height         int
	style          *Styles
	errStatus      error
}

func NewModel() model {
	style := DefaultStyle()
	m := model{
		homeForm:       HomePage(),
		conversionForm: ConversionPage(),
		resizeForm:     ResizePage(),
		style:          style,
	}

	m.currentPage = m.homeForm

	return m
}

func (m model) Init() tea.Cmd {
	return m.homeForm.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	currentPage, cmd := m.currentPage.Update(msg)
	if f, ok := currentPage.(*huh.Form); ok {
		m.currentPage = f
		cmds = append(cmds, cmd)
	}

	if m.homeForm.State == huh.StateCompleted {
		operation := m.homeForm.GetInt("operation")

		switch operation {
		case 0:
			m.currentPage = m.conversionForm
			m.currentPage.Init()
		case 1:
			m.currentPage = m.resizeForm
			m.currentPage.Init()
		}
	}

	if m.conversionForm.State == huh.StateCompleted {
		oldFileName := m.conversionForm.GetString("oldFileName")
		newFileName := m.conversionForm.GetString("newFileName")
		format := m.conversionForm.Get("format")
		confirm := m.conversionForm.GetBool("confirm")

		if confirm {
			m.errStatus = fileops.Convert(oldFileName, newFileName, format.(imgconv.Format))
			if m.errStatus != nil {
				m.currentPage = nil
			}
		}

		m.currentPage = nil
	}

	if m.resizeForm.State == huh.StateCompleted {
		fileName := m.resizeForm.GetString("fileName")
		resizeMode := m.resizeForm.GetInt("resizeMode")
		newSizeStr := m.resizeForm.GetString("newSizeStr")
		confirm := m.resizeForm.GetBool("confirm")

		if confirm {
			newSize, err := strconv.ParseFloat(newSizeStr, 64)
			if err != nil {
				m.errStatus = errors.New("could not parse new size")
			}

			m.errStatus = fileops.Resize(fileName, resizeMode, newSize)
			if m.errStatus != nil {
				m.currentPage = nil
			}
		}

		m.currentPage = nil
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var styledForm string
	if m.currentPage != nil {
		styledForm = m.style.Framed.Render(m.currentPage.View())
	} else {
		var success string
		if m.errStatus != nil {
			success = m.style.Failure.Render("failure")
		} else {
			success = m.style.Success.Render("success")
		}

		exit := m.style.Info.Render("You can exit the program using 'esc' or 'ctrl+c'")
		var builder strings.Builder
		fmt.Fprintf(&builder, "Operation status: %s\n", success)
		fmt.Fprintf(&builder, exit)
		styledForm = m.style.Framed.Render(builder.String())
	}

	centeredForm := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		styledForm,
	)

	return centeredForm
}
