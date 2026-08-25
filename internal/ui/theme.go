package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var lavenderPrimary = color.NRGBA{R: 0x8b, G: 0x6f, B: 0xd9, A: 0xff} // soft violet
var accentPink = color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0xff}

type DarkTheme struct{}

func (t DarkTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x17, G: 0x13, B: 0x21, A: 0xff}
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x1f, G: 0x19, B: 0x2d, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x2b, G: 0x23, B: 0x40, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x26, G: 0x20, B: 0x38, A: 0xff}
	case theme.ColorNameForeground, theme.ColorNameInputBorder:
		return color.NRGBA{R: 0xef, G: 0xe9, B: 0xfa, A: 0xff}
	case theme.ColorNamePrimary:
		return lavenderPrimary
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x4c, G: 0x35, B: 0x75, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xa7, G: 0x8b, B: 0xfa, A: 0x5e}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x33, G: 0x29, B: 0x4a, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x45, G: 0x36, B: 0x66, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x24, G: 0x1d, B: 0x33, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x86, G: 0x78, B: 0xa3, A: 0xff}
	case theme.ColorNameSeparator, theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x3d, G: 0x31, B: 0x58, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x52}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xfb, G: 0x71, B: 0x85, A: 0xff}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xf5, G: 0xc0, B: 0x5f, A: 0xff}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x34, G: 0xd3, B: 0x99, A: 0xff}
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
	if n == theme.SizeNameText {
		return 13.5
	}
	return theme.DefaultTheme().Size(n)
}

type LightTheme struct{}

func (t LightTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0xf7, G: 0xf4, B: 0xfd, A: 0xff}
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0xee, G: 0xe8, B: 0xfa, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0xe6, G: 0xdd, B: 0xf6, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0xea, G: 0xe4, B: 0xf6, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x2c, G: 0x22, B: 0x44, A: 0xff}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x74, G: 0x54, B: 0xc4, A: 0xff}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0xd8, G: 0xca, B: 0xf5, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x8b, G: 0x6f, B: 0xd9, A: 0x50}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0xdf, G: 0xd5, B: 0xf2, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xcd, G: 0xbe, B: 0xec, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0xfc, G: 0xfa, B: 0xff, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x8b, G: 0x80, B: 0xa6, A: 0xff}
	case theme.ColorNameSeparator, theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xd5, G: 0xc9, B: 0xee, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x40, G: 0x30, B: 0x60, A: 0x30}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xdc, G: 0x44, B: 0x5c, A: 0xff}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xb4, G: 0x7d, B: 0x12, A: 0xff}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x14, G: 0x8a, B: 0x62, A: 0xff}
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
	if n == theme.SizeNameText {
		return 13.5
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
