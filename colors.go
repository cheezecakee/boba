package boba

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Color = color.Color

func NewColor(s string) Color {
	return lipgloss.Color(s)
}

func Hex(s string) Color {
	return lipgloss.Color(s)
}

var (
	Black         color.Color = lipgloss.Black
	Red           color.Color = lipgloss.Red
	Green         color.Color = lipgloss.Green
	Yellow        color.Color = lipgloss.Yellow
	Blue          color.Color = lipgloss.Blue
	Magenta       color.Color = lipgloss.Magenta
	Cyan          color.Color = lipgloss.Cyan
	White         color.Color = lipgloss.White
	BrightBlack   color.Color = lipgloss.BrightBlack
	BrightRed     color.Color = lipgloss.BrightRed
	BrightGreen   color.Color = lipgloss.BrightGreen
	BrightYellow  color.Color = lipgloss.BrightYellow
	BrightBlue    color.Color = lipgloss.BrightBlue
	BrightMagenta color.Color = lipgloss.BrightMagenta
	BrightCyan    color.Color = lipgloss.BrightCyan
	BrightWhite   color.Color = lipgloss.BrightWhite
)

// Extended palette
var (
	// Grays
	Gray      color.Color = lipgloss.Color("240")
	DarkGray  color.Color = lipgloss.Color("236")
	LightGray color.Color = lipgloss.Color("250")

	// Your theme colors
	Primary   color.Color = lipgloss.Color("205") // hot pink
	Secondary color.Color = lipgloss.Color("86")  // aqua
	Danger    color.Color = lipgloss.Color("9")   // red
	Warning   color.Color = lipgloss.Color("214") // orange
	Success   color.Color = lipgloss.Color("10")  // green
)
