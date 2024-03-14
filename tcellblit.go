package tcellblit

import (
	"image"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gdamore/tcell/v2"
	"github.com/nfnt/resize"
	"golang.org/x/exp/constraints"
)

func Draw(s tcell.Screen, img image.Image, rx int, ry int, fill bool) {
	// Get terminal size and cursor width/height ratio
	termWidth, termHeight := rx, ry

	whratio := defaultRatio

	bounds := img.Bounds()
	imgW, imgH := bounds.Dx(), bounds.Dy()

	imgScale := scale(imgW, imgH, termWidth, termHeight, whratio)
	if fill {
		imgScaleX := scale(imgW, imgH, termWidth, 10000000, whratio)
		imgScaleY := scale(imgW, imgH, 10000000, termHeight, whratio)
		imgScale = imgScaleX
		if imgScaleY < imgScaleX {
			imgScale = imgScaleY
		}
	}

	// Resize canvas to fit scaled image
	width, height := int(float64(imgW)/imgScale), int(float64(imgH)/(imgScale*whratio))

	m := resize.Resize(uint(width), uint(height)*2, img, resize.MitchellNetravali)

	s.Clear()
	s.Sync()
	render(s, m)
	time.Sleep(10 * time.Millisecond)
	if w2, h2 := s.Size(); w2 != termWidth || h2 != termHeight || !verify(s, m) {
		Draw(s, img, rx, ry, fill)
		return
	}
	s.Show() // heavier than s.Show() but it's the only way to make it work
	s.Sync()
}

func verify(s tcell.Screen, img image.Image) bool {
	w, h := s.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			colorUp := tcell.FromImageColor(img.At(x, y*2))
			colorDown := tcell.FromImageColor(img.At(x, y*2+1))
			_, _, style, _ := s.GetContent(x, y)
			f, b, _ := style.Decompose()
			if f != colorDown || b != colorUp {
				return false
			}
		}
	}
	return true
}

func render(s tcell.Screen, img image.Image) {
	w, h := s.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			colorUp := tcell.FromImageColor(img.At(x, y*2))
			colorDown := tcell.FromImageColor(img.At(x, y*2+1))
			s.SetContent(x, y, '▄', nil, tcell.Style{}.Foreground(colorDown).Background(colorUp)) // use tcell
		}
	}
}

const defaultRatio float64 = 7.0 / 3.0 // The termiqnal's default cursor width/height ratio

// scales calculates the image scale to fit within the terminal width/height
func scale(imgW, imgH, termW, termH int, whratio float64) float64 {
	hr := float64(imgH) / (float64(termH) * whratio)
	wr := float64(imgW) / float64(termW)
	return max(hr, wr, 1)
}

// max returns the maximum value
func max[T constraints.Ordered](values ...T) T {
	var m T
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}
