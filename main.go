package main

import (
	"flag"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

const (
	canvasWidth  = 1080
	canvasHeight = 1920

	numRows float64 = 28
	numCols float64 = 15

	maxRotation  = 360
	minLineWidth = 1.0
	maxLineWidth = 3.5
	minAlpha     = 76
	maxAlpha     = 255
)

func main() {
	rand.Seed(time.Now().UnixNano())

	context := gg.NewContext(canvasWidth, canvasHeight)

	context.SetRGB(1, 1, 1)
	context.Clear()

	color := flag.Bool("color", false, "color the squares orange")
	padding := flag.Float64("padding", 0, "padding for top and left (px)")

	flag.Parse()

	sideLength := float64(canvasWidth-2**padding) / numCols
	maxJitter := 0.9 * sideLength

	for row := float64(0); row < numRows; row++ {
		// this parameter induces change in all aspects: line width, color, jitter, angle
		damper := math.Pow(row/numRows, 1.5) // linear growth induces chaos too quickly

		lineWidth := minLineWidth + damper*(maxLineWidth-minLineWidth)
		context.SetLineWidth(lineWidth)

		if *color {
			alpha := int(math.Round(minAlpha + damper*(maxAlpha-minAlpha)))
			context.SetRGBA255(229, 122, 50, alpha) // orange
		}

		y := *padding + row*sideLength

		for col := float64(0); col < numCols; col++ {
			x := *padding + col*sideLength

			context.Push()

			angle := ((rand.Float64() - 0.5) * maxRotation * math.Pi / 180.0) * damper
			jitterX := (rand.Float64() - 0.5) * maxJitter * damper
			jitterY := (rand.Float64() - 0.5) * maxJitter * damper

			context.RotateAbout(angle, x+sideLength/2, y+sideLength/2)
			context.Translate(jitterX, jitterY)
			context.DrawRectangle(x, y, sideLength, sideLength) // top-left corner

			context.Pop()
		}

		if *color {
			context.FillPreserve() // fill the path, but preserve it after filling it
		}

		context.SetRGB(0, 0, 0)
		context.Stroke() // stroke the path
	}

	context.SavePNG("out.png")
}
