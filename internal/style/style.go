package style

import "github.com/charmbracelet/lipgloss"

var Logo = `
               _               
  ___ ___ _ __| |__   ___ _ __ 
 / __/ _ \ '__| '_ \ / _ \ '__|
| (_|  __/ |  | |_) |  __/ |   
 \___\___|_|  |_.__/ \___|_|   
`

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#26A69A")). // Розовый (Dracula)
	// Background(lipgloss.Color("#232425ff")). // Серый фон
	Padding(0, 1).
	MarginBottom(1)

var NotFoundStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF5722")) // Красный (Dracula)

var SuccessStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#4CAF50")) // Зеленый (Dracula)
	