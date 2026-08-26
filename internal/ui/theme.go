package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var lavenderPrimary = color.NRGBA{R: 0x8b, G: 0x6f, B: 0xd9, A: 0xff}
var accentPink = color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0xff}

type DarkTheme struct{}

func (t DarkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x14, G: 0x10, B: 0x1e, A: 0xff}
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x1c, G: 0x16, B: 0x2a, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x28, G: 0x20, B: 0x3c, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x22, G: 0x1c, B: 0x34, A: 0xff}
	case theme.ColorNameForeground, theme.ColorNameInputBorder:
		return color.NRGBA{R: 0xed, G: 0xe6, B: 0xf8, A: 0xff}
	case theme.ColorNamePrimary:
		return lavenderPrimary
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x45, G: 0x30, B: 0x6e, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xa0, G: 0x85, B: 0xf5, A: 0x55}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x30, G: 0x26, B: 0x48, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x40, G: 0x32, B: 0x62, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x20, G: 0x1a, B: 0x30, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x80, G: 0x72, B: 0xa0, A: 0xff}
	case theme.ColorNameSeparator, theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x38, G: 0x2e, B: 0x54, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x48}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xfb, G: 0x71, B: 0x85, A: 0xff}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xf5, G: 0xc0, B: 0x5f, A: 0xff}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x34, G: 0xd3, B: 0x99, A: 0xff}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 0x18, G: 0x13, B: 0x24, A: 0xff}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (DarkTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (DarkTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (t DarkTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBarSmall:
		return 6
	}
	return theme.DefaultTheme().Size(n)
}

type LightTheme struct{}

func (t LightTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0xf5, G: 0xf1, B: 0xfc, A: 0xff}
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0xec, G: 0xe5, B: 0xf8, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0xe2, G: 0xd9, B: 0xf3, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0xe8, G: 0xe1, B: 0xf4, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x2a, G: 0x20, B: 0x42, A: 0xff}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x70, G: 0x50, B: 0xc0, A: 0xff}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0xd4, G: 0xc5, B: 0xf2, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x88, G: 0x6c, B: 0xd6, A: 0x48}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0xdc, G: 0xd2, B: 0xf0, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xc8, G: 0xba, B: 0xea, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0xfb, G: 0xf8, B: 0xff, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x88, G: 0x7c, B: 0xa4, A: 0xff}
	case theme.ColorNameSeparator, theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xd2, G: 0xc6, B: 0xeb, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x40, G: 0x30, B: 0x60, A: 0x28}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xdc, G: 0x44, B: 0x5c, A: 0xff}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xb4, G: 0x7d, B: 0x12, A: 0xff}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x14, G: 0x8a, B: 0x62, A: 0xff}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 0xf0, G: 0xea, B: 0xf8, A: 0xff}
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (LightTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (LightTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (t LightTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBarSmall:
		return 6
	}
	return theme.DefaultTheme().Size(n)
}

func ApplyTheme(a fyne.App, light bool) {
	if light {
		a.Settings().SetTheme(LightTheme{})
	} else {
		a.Settings().SetTheme(DarkTheme{})
	}
}
