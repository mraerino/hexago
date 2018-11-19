package main

import (
	"github.com/ajstarks/svgo"
	"os"
	"math"
)

type Point struct {
	X int
	Y int
}

type Hexagon struct {
	A, B, C, D, E, F Point
	M Point
	R int
}

func (h *Hexagon) calc() {
	// dot above
	h.A.X = h.M.X
	h.A.Y = h.M.Y - h.R
	
	// dot upper right
	h.B.X = h.M.X + int(math.Round(float64(h.R) * math.Cos(math.Pi / 6)))
	h.B.Y = h.M.Y - int(math.Round(float64(h.R) * math.Sin(math.Pi / 6)))
	
	// dot lower right
	h.C.X = h.M.X + int(math.Round(float64(h.R) * math.Cos(math.Pi / 6)))
	h.C.Y = h.M.Y + int(math.Round(float64(h.R) * math.Sin(math.Pi / 6)))
	
	// dot below
	h.D.X = h.M.X
	h.D.Y = h.M.Y + h.R
	
	// dot lower left
	h.E.X = h.M.X - int(math.Round(float64(h.R) * math.Cos(math.Pi / 6)))
	h.E.Y = h.M.Y + int(math.Round(float64(h.R) * math.Sin(math.Pi / 6)))
	
	// dot upper left
	h.F.X = h.M.X - int(math.Round(float64(h.R) * math.Cos(math.Pi / 6)))
	h.F.Y = h.M.Y - int(math.Round(float64(h.R) * math.Sin(math.Pi / 6)))
}

func (h Hexagon) sliceX() []int {
	return []int {
		h.A.X,
		h.B.X,
		h.C.X,
		h.D.X,
		h.E.X,
		h.F.X,
	}
}

func (h Hexagon) sliceY() []int {
	return []int {
		h.A.Y,
		h.B.Y,
		h.C.Y,
		h.D.Y,
		h.E.Y,
		h.F.Y,
	}
}

func (h Hexagon) draw(canvas *svg.SVG, gray int) {
	canvas.Polygon(h.sliceX(), h.sliceY(), canvas.RGB(gray, gray, gray))
}

var grays = []int{
	42, 84, 126, 168,
}

func draw_mupliple_hexagons(canvas *svg.SVG, M Point, rounds int, radius int) {
	first_h := Hexagon{
		M: M,
		R: radius,
	}
	hexagon_width := float64(first_h.R) * 2 * math.Cos(math.Pi / 6)
	first_h.calc()
	first_h.draw(canvas, grays[0])
	for c := 1; c < rounds; c++ {
		count_hex := 6 * c
		h := Hexagon{
			M: Point{
				X: first_h.C.X,
				Y: first_h.C.Y + first_h.R,
			},
			R: first_h.R,
		}
		h.calc()
		h.draw(canvas, grays[0])
		first_h = h
		deg := math.Pi // 180°
		for i := 1; i < count_hex; i++ {
			h = Hexagon{
				M: Point{
					X: h.M.X + int(hexagon_width * math.Cos(deg)),
					Y: h.M.Y + int(hexagon_width * math.Sin(deg)),
				},
				R: h.R,
			}
			h.calc()
			h.draw(canvas, grays[i%4])
			if i % c == 0 {
				deg += math.Pi / 3 // 60°
			}
		}
	}
}

func main() {
	width := 10000
	height := 10000
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	draw_mupliple_hexagons(canvas, Point{X: 1000, Y: 1000}, 10, 50)
	canvas.End()
}
