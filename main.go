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
	history [][]int
}

func (r *randomColor) isCollision(proposal, circleIndex int, prevCircleNeighbor float64) bool {
	if proposal == -1 {
		return true
	}
	if circleIndex == 0 {
		return false
	}
	if circleIndex == 1 && proposal == r.history[0][0] {
		return true
	}
	currentIndex := len(r.history[circleIndex])
	circleSize := circleIndex * 6
	if currentIndex+1 == circleSize && proposal == r.history[circleIndex][0] {
		return true
	}
	if currentIndex > 0 && proposal == r.history[circleIndex][currentIndex-1] {
		return true
	}
	neighborRight := int(math.Floor(prevCircleNeighbor)) % len(r.history[circleIndex-1])
	if proposal == r.history[circleIndex-1][neighborRight] {
		return true
	}
	neighborLeft := int(math.Ceil(prevCircleNeighbor)) % len(r.history[circleIndex-1])
	if proposal == r.history[circleIndex-1][neighborLeft] {
		return true
	}
	return false
}

func (r *randomColor) get(circleIndex, localIndex int) int {
	if len(r.history) <= circleIndex {
		r.history = append(r.history, []int{})
	}

	ratio := float64(circleIndex) / (float64(circleIndex) - 1)
	prevCircleNeighbor := float64(localIndex) / ratio

	newColorIndex := -1
	for r.isCollision(newColorIndex, circleIndex, prevCircleNeighbor) {
		newIndexBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(grays))))
		newColorIndex = int(newIndexBig.Int64())
		if err != nil {
			panic(err)
		}
	}
	r.history[circleIndex] = append(r.history[circleIndex], newColorIndex)
	return grays[newColorIndex]
}

type HexagonGenerator struct {
	canvas *svg.SVG
	radius int
}

func (gen *HexagonGenerator) drawSingle(pos Point, color int) Hexagon {
	hex := Hexagon{
		M: pos,
		R: gen.radius,
	}
	hex.calc()
	hex.draw(gen.canvas, color)
	return hex
}

func (gen *HexagonGenerator) drawMultipleHexagons(M Point, rounds int) {
	colorGen := randomColor{}

	color := colorGen.get(0, 0)
	firstH := gen.drawSingle(M, color)

	hexagonWidth := float64(firstH.R) * 2 * math.Cos(math.Pi/6)

	for c := 1; c < rounds; c++ {
		hexInCircleCount := 6 * c

		color := colorGen.get(c, 0)
		h := gen.drawSingle(Point{
			X: firstH.C.X,
			Y: firstH.C.Y + firstH.R,
		}, color)
		firstH = h

		deg := math.Pi // 180°
		for i := 1; i < hexInCircleCount; i++ {
			color := colorGen.get(c, i)
			h = gen.drawSingle(Point{
				X: h.M.X + int(hexagonWidth*math.Cos(deg)),
				Y: h.M.Y + int(hexagonWidth*math.Sin(deg)),
			}, color)

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
	gen := HexagonGenerator{canvas, 100}
	gen.drawMultipleHexagons(Point{X: 1000, Y: 1000}, 5)
	canvas.End()
}

func main() {
	http.Handle("/", http.HandlerFunc(graphicHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
