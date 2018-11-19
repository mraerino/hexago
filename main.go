package main

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	"net/http"

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

func (r *randomColor) isCollision(proposal, circleIndex int) bool {
	globalIndex := len(r.history)
	if proposal == -1 {
		return true
	}
	if globalIndex == 0 {
		return false
	}
	if circleIndex == 1 && proposal == r.history[0] {
		return true
	}
	circleSize := circleIndex * 6
	indexToCheck := globalIndex - (circleSize - 1)
	if indexToCheck > 0 && proposal == r.history[indexToCheck] {
		return true
	}
	if proposal == r.history[globalIndex-1] {
		return true
	}
	return false
}

func (r *randomColor) get(circleIndex int) int {
	newColorIndex := -1
	for r.isCollision(newColorIndex, circleIndex) {
		newIndexBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(grays))))
		newColorIndex = int(newIndexBig.Int64())
		if err != nil {
			panic(err)
		}
	}
	r.history = append(r.history, newColorIndex)
	return grays[newColorIndex]
}

func drawMultipleHexagons(canvas *svg.SVG, M Point, rounds int, radius int) {
	colorGen := randomColor{}
	firstH := Hexagon{
		M: M,
		R: radius,
	}
	hexagonWidth := float64(firstH.R) * 2 * math.Cos(math.Pi/6)
	firstH.calc()
	color := colorGen.get(0)
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
		color := colorGen.get(c)
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

			color := colorGen.get(c)
			h.draw(canvas, color)
			if i%c == 0 {
				deg += math.Pi / 3 // 60°
			}
		}
	}
}

func graphicHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	width := 10000
	height := 10000
	canvas := svg.New(w)
	canvas.Start(width, height)
	drawMultipleHexagons(canvas, Point{X: 1000, Y: 1000}, 10, 50)
	canvas.End()
}

func main() {
	http.Handle("/", http.HandlerFunc(graphicHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
