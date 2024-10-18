package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	r "math/rand/v2"
)

const GridArea = 512
const VectorGridDiv = 128
const VectorGrid = GridArea / VectorGridDiv

func initGrid() [GridArea][GridArea]float32 {
	var grid [GridArea][GridArea]float32
	for i := 0; i < GridArea; i++ {
		for j := 0; j < GridArea; j++ {
			grid[i][j] = 0
		}
	}
	return grid

}

func randVector() (float32, float32) {

	angle := r.Float32() * 2 * math.Pi
	return float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
func Distance(x int, xi int) int {
	return x - xi

}

func CalculateCorners(x int, y int) ([4]int, [4]int) {
	return [4]int{x, x + 1, x, x + 1}, [4]int{y, y, y + 1, y + 1}
}

func CalculateOffset(x int, y int, xi int, yi int) (float32, float32) {
	Xdistance := Distance(x, xi)
	Ydistance := Distance(y, yi)

	return float32(Xdistance) / float32(VectorGridDiv), float32(Ydistance) / float32(VectorGridDiv)

}
func dotProd2D(x float32, y float32, xi float32, yi float32) float32 {
	return x*xi + y*yi

}

// lerp(a,b,t)=a+t⋅(b−a)
func lininterpol(a float32, b float32, t float32) float32 {
	return a + t*(b-a)
}
func fade(w float32) float32 {

	return w * w * w * (w*(w*6-15) + 10)
}
func bilinterpol(prods [4]float32, x int, y int, cornX int, cornY int) float32 {
	dot00 := prods[0]
	dot10 := prods[1]
	dot01 := prods[2]
	dot11 := prods[3]
	//xOffset := (x - x0) / (x1 - x0)
	//yOffset := (y - y0) / (y1 - y0)
	u := float32(x%VectorGridDiv) / float32(VectorGridDiv)
	v := float32(y%VectorGridDiv) / float32(VectorGridDiv)

	R1 := lininterpol(dot00, dot10, fade(u))

	R2 := lininterpol(dot01, dot11, fade(u))
	return lininterpol(R1, R2, fade(v))
}
func Corners(array [GridArea][GridArea]float32) [(GridArea / VectorGridDiv) + 1][1 + (GridArea / VectorGridDiv)][2]float32 {
	const ylen = 1 + (GridArea / VectorGridDiv)
	const xlen = 1 + (GridArea / VectorGridDiv)
	var output [xlen][ylen][2]float32
	for i := 0; i < ylen; i++ {
		for j := 0; j < xlen; j++ {
			output[i][j][0], output[i][j][1] = randVector()
		}
	}
	return output

}

func main() {
	rl.InitWindow(1200, 600, "Perlin Noise")
	defer rl.CloseWindow()
	grid := initGrid()
	corners := Corners(grid)
	for i := 0; i < GridArea; i++ {
		for j := 0; j < GridArea; j++ {
			var prods [4]float32
			XCorners, YCorners := CalculateCorners(j/VectorGridDiv, i/VectorGridDiv)
			for k := 0; k < 4; k++ {
				offsetX, offsetY := CalculateOffset(XCorners[k]*GridArea/VectorGridDiv, YCorners[k]*GridArea/VectorGridDiv, j, i)
				prods[k] = dotProd2D(corners[YCorners[k]][XCorners[k]][0], corners[YCorners[k]][XCorners[k]][1], float32(offsetX), float32(offsetY))
			}
			grid[i][j] = bilinterpol(prods, j, i, XCorners[0]*4, YCorners[0]*4)
		}

	}
	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		for l := 0; l < GridArea; l++ {
			for m := 0; m < GridArea; m++ {
				rl.DrawRectangle(int32(m*600/GridArea), int32(l*600/GridArea), 600/GridArea, 600/GridArea, rl.Color{0, 0, 0, uint8(200 * grid[l][m])})
			}
		}

		rl.EndDrawing()
	}

}
