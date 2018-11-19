package main

import (
	"github.com/ajstarks/svgo"
	"os"
	"math"
)

func hexagon(canvas *svg.SVG, x int, y int, r int)  {
	dots_x := []int{}
	dots_y := []int{}
	
	// dot above
	dots_x = append(dots_x, x)
	dots_y = append(dots_y, y - r)
	
	// dot upper right
	dots_x = append(dots_x, x + int(math.Round(float64(r) * math.Cos(math.Pi / 6))))
	dots_y = append(dots_y, y - int(math.Round(float64(r) * math.Sin(math.Pi / 6))))
	
	// dot lower right
	dots_x = append(dots_x, x + int(math.Round(float64(r) * math.Cos(math.Pi / 6))))
	dots_y = append(dots_y, y + int(math.Round(float64(r) * math.Sin(math.Pi / 6))))
	
	// dot below
	dots_x = append(dots_x, x)
	dots_y = append(dots_y, y + r)
	
	// dot lower left
	dots_x = append(dots_x, x - int(math.Round(float64(r) * math.Cos(math.Pi / 6))))
	dots_y = append(dots_y, y + int(math.Round(float64(r) * math.Sin(math.Pi / 6))))
	
	// dot upper left
	dots_x = append(dots_x, x - int(math.Round(float64(r) * math.Cos(math.Pi / 6))))
	dots_y = append(dots_y, y - int(math.Round(float64(r) * math.Sin(math.Pi / 6))))
	
	// actually draw
	canvas.Polygon(dots_x, dots_y)
}

func main() {
	width := 500
	height := 500
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)
	hexagon(canvas, 250, 250, 200)
	canvas.End()
}
