package ui

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type chartPt struct {
	X float64
	Y float64
}

type chartSeries struct {
	Label  string
	Color  color.NRGBA
	Points []chartPt
	Labels []string
}

func (s *chartSeries) labelAt(i int) string {
	if s == nil || i < 0 || i >= len(s.Labels) {
		return ""
	}
	return s.Labels[i]
}

func niceCeil(v float64) float64 {
	if v <= 0 {
		return 1
	}
	p := 1.0
	for p*10 <= v {
		p *= 10
	}
	m := v / p
	switch {
	case m <= 1:
		return p
	case m <= 2:
		return 2 * p
	case m <= 5:
		return 5 * p
	default:
		return 10 * p
	}
}

var (
	chartBG     = color.NRGBA{}
	gridDark    = color.NRGBA{R: 0x46, G: 0x38, B: 0x66, A: 0x70}
	textLight   = color.NRGBA{R: 0xb9, G: 0xab, B: 0xd8, A: 0xff}
	gridLight   = color.NRGBA{R: 0xc9, G: 0xbb, B: 0xe6, A: 0x90}
	textDarkCol = color.NRGBA{R: 0x5a, G: 0x4b, B: 0x80, A: 0xff}
)

func drawText(img *image.RGBA, x, y int, s string, col color.NRGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}

func renderChart(series []chartSeries, w, h int, dark bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := chartBG
	draw.Draw(img, img.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	grid, txt := gridDark, textLight
	if dark == false {
		grid, txt = gridLight, textDarkCol
	}

	max := 1.0
	n := 0
	for _, s := range series {
		if len(s.Points) > n {
			n = len(s.Points)
		}
		for _, p := range s.Points {
			if p.Y > max {
				max = p.Y
			}
		}
	}
	top := niceCeil(max * 1.06)
	if n < 2 {
		n = 2
	}

	const pl, pr, pt, pb = 56, 14, 14, 26
	iw := w - pl - pr
	ih := h - pt - pb
	xAt := func(i int) int { return pl + i*(iw)/(n-1) }
	yAt := func(v float64) int { return pt + ih - int(v/top*float64(ih)) }

	for f := 0.0; f <= 1.001; f += 0.25 {
		y := pt + ih - int(f*float64(ih))
		lineColor := grid
		if f == 0 {
			lineColor = txt
		}
		for x := pl; x < w-pr; x++ {
			img.Set(x, y, lineColor)
		}
		val := top * f
		lbl := fmtFloatShort(val)
		drawText(img, pl-6-len(lbl)*7, y+4, lbl, txt)
	}

	if n > 1 && len(series) > 0 && len(series[0].Points) > 0 {
		first := &series[0]
		ticks := []int{0, (n - 1) / 2, n - 1}
		for _, ti := range ticks {
			lbl := first.labelAt(ti)
			anchor := xAt(ti)
			if ti == n-1 {
				anchor -= len(lbl) * 7
			} else if ti != 0 {
				anchor -= len(lbl) * 7 / 2
			}
			drawText(img, anchor, h-8, lbl, txt)
		}
	}

	for si, s := range series {
		if len(s.Points) < 2 {
			continue
		}
		col := s.Color
		fill := color.NRGBA{R: col.R, G: col.G, B: col.B, A: 0x30}
		prevY := yAt(s.Points[0].Y)
		for i := 1; i < len(s.Points); i++ {
			x0, x1 := xAt(i-1), xAt(i)
			y0, y1 := prevY, yAt(s.Points[i].Y)
			stepX := 1
			if x1 < x0 {
				stepX = -1
			}
			for x := x0; x != x1; x += stepX {
				t := float64(x-x0) / float64(x1-x0)
				y := y0 + int(t*float64(y1-y0))
				yy := y
				for yy <= pt+ih {
					img.Set(x, yy, fill)
					yy++
				}
				img.Set(x, y, col)
			}
			prevY = y1
		}
		_ = si
	}
	return img
}

func fmtFloatShort(v float64) string {
	switch {
	case v >= 1e9:
		return trimF(v/1e9) + "B"
	case v >= 1e6:
		return trimF(v/1e6) + "M"
	case v >= 1e3:
		return trimF(v/1e3) + "k"
	default:
		return trimF(v)
	}
}

func trimF(v float64) string {
	s := ""
	intPart := int(v)
	frac := int((v - float64(intPart)) * 10)
	s = itoaStr(intPart)
	if frac > 0 {
		s += "." + itoaStr(frac)
	}
	return s
}

func itoaStr(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [12]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}
