package main

import (
	"crypto/rand"
	"math"
	"math/big"
	"os"

	"github.com/ajstarks/svgo"
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
	h.B.X = h.M.X + int(math.Round(float64(h.R)*math.Cos(math.Pi/6)))
	h.B.Y = h.M.Y - int(math.Round(float64(h.R)*math.Sin(math.Pi/6)))

	// dot lower right
	h.C.X = h.M.X + int(math.Round(float64(h.R)*math.Cos(math.Pi/6)))
	h.C.Y = h.M.Y + int(math.Round(float64(h.R)*math.Sin(math.Pi/6)))

	// dot below
	h.D.X = h.M.X
	h.D.Y = h.M.Y + h.R

	// dot lower left
	h.E.X = h.M.X - int(math.Round(float64(h.R)*math.Cos(math.Pi/6)))
	h.E.Y = h.M.Y + int(math.Round(float64(h.R)*math.Sin(math.Pi/6)))

	// dot upper left
	h.F.X = h.M.X - int(math.Round(float64(h.R)*math.Cos(math.Pi/6)))
	h.F.Y = h.M.Y - int(math.Round(float64(h.R)*math.Sin(math.Pi/6)))
}

func (h Hexagon) sliceX() []int {
	return []int{
		h.A.X,
		h.B.X,
		h.C.X,
		h.D.X,
		h.E.X,
		h.F.X,
	}
}

func (h Hexagon) sliceY() []int {
	return []int{
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

type randomColor struct {
	history []int
}

func (r *randomColor) get() int {
	newIndex := r.history[len(r.history)-1]
	for newIndex == r.history[len(r.history)-1] {
		newIndexBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(grays))))
		newIndex = int(newIndexBig.Int64())
		if err != nil {
			panic(err)
		}
	}
	r.history = append(r.history, newIndex)
	return grays[newIndex]
}

func drawMultipleHexagons(canvas *svg.SVG, M Point, rounds int, radius int) {
	colorGen := randomColor{
		history: []int{0},
	}
	firstH := Hexagon{
		M: M,
		R: radius,
	}
	hexagonWidth := float64(firstH.R) * 2 * math.Cos(math.Pi/6)
	firstH.calc()
	color := colorGen.get()
	firstH.draw(canvas, color)
	for c := 1; c < rounds; c++ {
		hexInCircleCount := 6 * c
		h := Hexagon{
			M: Point{
				X: firstH.C.X,
				Y: firstH.C.Y + firstH.R,
			},
			R: firstH.R,
		}
		h.calc()
		color := colorGen.get()
		h.draw(canvas, color)
		firstH = h
		deg := math.Pi // 180°
		for i := 1; i < hexInCircleCount; i++ {
			h = Hexagon{
				M: Point{
					X: h.M.X + int(hexagonWidth*math.Cos(deg)),
					Y: h.M.Y + int(hexagonWidth*math.Sin(deg)),
				},
				R: h.R,
			}
			h.calc()
			color := colorGen.get()
			h.draw(canvas, color)
			if i%c == 0 {
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
	drawMultipleHexagons(canvas, Point{X: 1000, Y: 1000}, 10, 50)
	canvas.End()
}
