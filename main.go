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

func (h Hexagon) draw(canvas *svg.SVG) {
	canvas.Polygon(h.sliceX(), h.sliceY())
}

func main() {
	width := 500
	height := 500
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	h := Hexagon{
		M: Point{X: 250, Y: 250,},
		R: 200,
	}
	h.calc()
	h.draw(canvas)
	canvas.End()
}
