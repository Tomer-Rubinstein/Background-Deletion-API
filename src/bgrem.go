package main

import (
	"fmt"
	"encoding/base64"
	"strings"
	"image"
	"image/png"
	"image/color"
	"os"
	"math"
)

// go run bgrem.go utils.go

func main() {
	/* DEBUG: opening a file containing base64-encoded png image data */
	encodedFile, err := os.Open("../tests/base64-text.txt")
	isError(err)
	buffer := make([]byte, 334476)
	n, err := encodedFile.Read(buffer)
	isError(err)
	
	/* converting base64 to PNG */
	data, err := base64.StdEncoding.DecodeString(string(buffer[:n]))
	isError(err)
	myReader := strings.NewReader(string(data))
	img, err := png.Decode(myReader)
	isError(err)

	result := RembgPNG(img)
	outputFile, err := os.Create("outputyesssss.png")
	isError(err)
	defer outputFile.Close()
	err := png.Encode(outputFile, result)
	isError(err)
}


/*
	upperBound := []float64{255,255,180}
	lowerBound := []float64{40,75,0}
*/

func RembgPNG(img image.Image) *image.NRGBA {
	newImg := make([][3]float64, img.Bounds().Size().X*img.Bounds().Size().Y)
	k := 0

	// ranging lower & upper bounds
	upperBound := []float64{200,130,100}
	lowerBound := []float64{0,75,0}

	// creating the HSV version of the img
	for row:=0; row < img.Bounds().Size().Y; row++ {
		for col:=0; col < img.Bounds().Size().X; col++ {
			point := img.At(row, col)
			H,S,V := RGBtoHSV(point)
			if (lowerBound[0] < H && H < upperBound[0]) && (lowerBound[1] < S && S < upperBound[1]) && (lowerBound[2] < V && V < upperBound[2]) {
				newImg[k][0],newImg[k][1],newImg[k][2] = H,S,V
			}else{
				newImg[k][0],newImg[k][1],newImg[k][2] = 0,0,0
			}
			k=k+1
		}
	}

	// img (bitwise and) newImg
	k = 0
	rgbImg := image.NewNRGBA(img.Bounds())
	transparent := color.RGBA{0,0,0,0}
	for row:=0; row < img.Bounds().Size().Y; row++ {
		for col:=0; col < img.Bounds().Size().X; col++ {
			rgbImg.Set(row, col, transparent)
			if (newImg[k][0] == 0 && newImg[k][1] == 0 && newImg[k][2] == 0) {
				rgbImg.Set(row, col, img.At(row, col))
			}
			k=k+1
		}
	}

	fmt.Println(rgbImg)
	return rgbImg
}


func RGBtoHSV(point color.Color) (float64, float64, float64){
 	/* converting to HSV */
	_r,_g,_b,_ := point.RGBA()
	r,g,b := float64(_r)/257,float64(_g)/257,float64(_b)/257
	M := math.Max(r, math.Max(g, b))
	m := math.Min(r, math.Min(g, b))
	var V float64 = M/255
	var S float64 = 0
	if M > 0 {
		S = 1 - m/M
	}
	var H float64 = -1
	if g >= b {
		H = (180/math.Pi)*math.Acos((r - 0.5*g - 0.5*b)/(math.Sqrt(math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2) - r*g - r*b - g*b)))
	} else {
		H = 360 - (180/math.Pi)*math.Acos((r - 0.5*g - 0.5*b)/(math.Sqrt(math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2) - r*g - r*b - g*b)))
	}
	//fmt.Println(H, S, V)
	return H, S*100, V*100
}
