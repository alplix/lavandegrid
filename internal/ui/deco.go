package ui

import (
	"image"
	"image/color"
	"math"
)

func drawPetal(img *image.RGBA, cx, cy, rx, ry, angle float64, col color.NRGBA) {
	dr := img.Bounds()
	sin, cos := math.Sin(angle), math.Cos(angle)
	for y := int(cy - ry - 1); y <= int(cy+ry+1); y++ {
		for x := int(cx-rx-1); x <= int(cx+rx+1); x++ {
			if x < dr.Min.X || x >= dr.Max.X || y < dr.Min.Y || y >= dr.Max.Y {
				continue
			}
			dx, dy := float64(x)-cx, float64(y)-cy
			lx := dx*cos + dy*sin
			ly := -dx*sin + dy*cos
			nx := lx / rx
			ny := ly / ry
			d2 := nx*nx + ny*ny
			if d2 > 1 {
				continue
			}
			alpha := float64(col.A) * math.Max(0, 1-d2*0.7)
			if alpha < 1 {
				continue
			}
			px := x + y*img.Stride
			img.Pix[px+0] = col.R
			img.Pix[px+1] = col.G
			img.Pix[px+2] = col.B
			img.Pix[px+3] = uint8(alpha)
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r float64, col color.NRGBA) {
	dr := img.Bounds()
	for y := int(cy - r - 1); y <= int(cy+r+1); y++ {
		for x := int(cx-r-1); x <= int(cx+r+1); x++ {
			if x < dr.Min.X || x >= dr.Max.X || y < dr.Min.Y || y >= dr.Max.Y {
				continue
			}
			dx, dy := float64(x)-cx, float64(y)-cy
			d := math.Sqrt(dx*dx + dy*dy)
			if d > r {
				continue
			}
			alpha := float64(col.A) * math.Max(0, 1-d/(r+0.5))
			if alpha < 1 {
				continue
			}
			px := x + y*img.Stride
			img.Pix[px+0] = col.R
			img.Pix[px+1] = col.G
			img.Pix[px+2] = col.B
			img.Pix[px+3] = uint8(alpha)
		}
	}
}

func renderFlower(w, h int, dark bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if dark {
		for i := range img.Pix {
			img.Pix[i] = 0
		}
	} else {
		for i := range img.Pix {
			img.Pix[i] = 0xff
		}
	}

	petalCol := color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0x90}
	petalCol2 := color.NRGBA{R: 0xa7, G: 0x8b, B: 0xfa, A: 0x70}
	centerCol := color.NRGBA{R: 0xf5, G: 0xc0, B: 0x5f, A: 0xd0}

	flowers := []struct {
		x, y, s float64
		petals  int
		rot     float64
	}{
		{float64(w) * 0.88, float64(h) * 0.12, 28, 5, 0},
		{float64(w) * 0.92, float64(h) * 0.35, 18, 5, 0.4},
		{float64(w) * 0.06, float64(h) * 0.88, 22, 5, 0.8},
		{float64(w) * 0.95, float64(h) * 0.72, 14, 5, 1.2},
		{float64(w) * 0.03, float64(h) * 0.55, 10, 5, 0.2},
		{float64(w) * 0.78, float64(h) * 0.92, 16, 5, 1.8},
	}

	for _, f := range flowers {
		for i := 0; i < f.petals; i++ {
			a := f.rot + float64(i)*2*math.Pi/float64(f.petals)
			px := f.x + math.Cos(a)*f.s*0.45
			py := f.y + math.Sin(a)*f.s*0.45
			if i%2 == 0 {
				drawPetal(img, px, py, f.s*0.38, f.s*0.18, a, petalCol)
			} else {
				drawPetal(img, px, py, f.s*0.35, f.s*0.15, a, petalCol2)
			}
		}
		drawCircle(img, f.x, f.y, f.s*0.15, centerCol)
	}

	for i := 0; i < 12; i++ {
		x := float64(w) * (0.1 + 0.8*float64(i)/11)
		y := float64(h)*0.95 + math.Sin(float64(i)*0.7)*float64(h)*0.03
		drawPetal(img, x, y, 6, 3, float64(i)*0.5, color.NRGBA{R: 0xc4, G: 0x95, B: 0xf0, A: 0x35})
	}

	return img
}

func renderGradient(w, h int, c1, c2 color.NRGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(float64(c1.R)*(1-t) + float64(c2.R)*t)
		g := uint8(float64(c1.G)*(1-t) + float64(c2.G)*t)
		b := uint8(float64(c1.B)*(1-t) + float64(c2.B)*t)
		for x := 0; x < w; x++ {
			px := (y*w + x) * 4
			img.Pix[px+0] = r
			img.Pix[px+1] = g
			img.Pix[px+2] = b
			img.Pix[px+3] = 0xff
		}
	}
	return img
}

func renderCardGradient(w, h int, base color.NRGBA, light bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	top := color.NRGBA{R: base.R + 12, G: base.G + 10, B: base.B + 14, A: 0xff}
	if light {
		top = color.NRGBA{R: base.R - 8, G: base.G - 8, B: base.B - 5, A: 0xff}
	}
	for y := 0; y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(float64(top.R)*(1-t) + float64(base.R)*t)
		g := uint8(float64(top.G)*(1-t) + float64(base.G)*t)
		b := uint8(float64(top.B)*(1-t) + float64(base.B)*t)
		for x := 0; x < w; x++ {
			px := (y*w + x) * 4
			img.Pix[px+0] = r
			img.Pix[px+1] = g
			img.Pix[px+2] = b
			img.Pix[px+3] = 0xff
		}
	}
	return img
}

func renderScatter(w, h int, count int, base color.NRGBA, dark bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if dark {
		for i := range img.Pix {
			img.Pix[i] = 0
		}
	} else {
		for i := range img.Pix {
			img.Pix[i] = 0xff
		}
	}
	type dot struct{ x, y, r float64 }
	dots := make([]dot, count)
	for i := range dots {
		seed := float64(i*7919 + 31)
		dots[i] = dot{
			x: math.Mod(seed*1.3, float64(w)),
			y: math.Mod(seed*2.7, float64(h)),
			r: 2 + math.Mod(seed, 5),
		}
	}
	for _, d := range dots {
		col := base
		col.A = 40 + uint8(math.Mod(d.x+d.y, 60))
		drawCircle(img, d.x, d.y, d.r, col)
	}
	return img
}

func mixImage(base, overlay *image.RGBA) {
	bd := base.Bounds()
	od := overlay.Bounds()
	for y := 0; y < min(bd.Dy(), od.Dy()); y++ {
		for x := 0; x < min(bd.Dx(), od.Dx()); x++ {
			bp := (y*base.Stride + x*4)
			op := (y*overlay.Stride + x*4)
			srcA := float64(overlay.Pix[op+3]) / 255
			if srcA < 0.01 {
				continue
			}
			dstA := float64(base.Pix[bp+3]) / 255
			outA := srcA + dstA*(1-srcA)
			if outA < 0.01 {
				continue
			}
			for c := 0; c < 3; c++ {
				src := float64(overlay.Pix[op+c])
				dst := float64(base.Pix[bp+c])
				base.Pix[bp+c] = uint8((src*srcA + dst*dstA*(1-srcA)) / outA)
			}
			base.Pix[bp+3] = uint8(outA * 255)
		}
	}
}
