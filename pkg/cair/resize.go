package cair

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

func Resize(img image.Image, newBounds image.Point) image.Image {
	energies := GradientEnergy(img)
	colExclusions := make(map[int]struct{})
	for i := 0; i < img.Bounds().Dx()-newBounds.X; i++ {
		minColEnergy := math.MaxFloat64
		minCol := 0
		for col := 0; col < len(energies[0]); col++ {
			if _, ok := colExclusions[col]; ok {
				continue
			}
			totalColEnergy := 0.0
			for row := 0; row < len(energies); row++ {
				totalColEnergy += energies[row][col]
			}
			if totalColEnergy < minColEnergy {
				minCol = col
				minColEnergy = totalColEnergy
			}
		}
		colExclusions[minCol] = struct{}{}
	}
	rowExclusions := make(map[int]struct{})
	for i := 0; i < img.Bounds().Dy()-newBounds.Y; i++ {
		minRowEnergy := math.MaxFloat64
		minRow := 0
		for row := 0; row < len(energies); row++ {
			if _, ok := rowExclusions[row]; ok {
				continue
			}
			totalRowEnergy := 0.0
			for col := 0; col < len(energies[i]); col++ {
				totalRowEnergy += energies[row][col]
			}
			if totalRowEnergy < minRowEnergy {
				minRow = row
				minRowEnergy = totalRowEnergy
			}
		}
		rowExclusions[minRow] = struct{}{}
	}
	newImg := image.NewRGBA64(image.Rect(0, 0, newBounds.X, newBounds.Y))
	var writePt image.Point
	for i := img.Bounds().Min.X; i < img.Bounds().Max.X; i++ {
		if _, ok := rowExclusions[i]; ok {
			continue
		}
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y; j++ {
			if _, ok := colExclusions[j]; ok {
				continue
			}
			r, g, b, a := img.At(i, j).RGBA()
			fmt.Println(writePt.X, writePt.Y, r, g, b, a)
			newImg.SetRGBA64(writePt.X, writePt.Y, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
			writePt.Y++
		}
		writePt.Y = 0
		writePt.X++
	}
	return newImg
}

func GradientEnergy(img image.Image) [][]float64 {
	var energies [][]float64
	for i := img.Bounds().Min.X; i < img.Bounds().Max.X; i++ {
		var colEnergies []float64
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y; j++ {
			energy := 0.0
			r, g, b, _ := img.At(i, j).RGBA()
			for _, pt := range neighbourPoints(img, i, j) {
				x, y, z, _ := img.At(pt.X, pt.Y).RGBA()
				energy += colorDiff(r, g, b, x, y, z)
			}
			colEnergies = append(colEnergies, energy)
		}
		energies = append(energies, colEnergies)
	}
	return energies
}

func neighbourPoints(img image.Image, i, j int) []image.Point {
	var neighbours []image.Point
	for row := i - 1; row <= i+1; row++ {
		for col := j - 1; col <= j+1; col++ {
			// Add to neighbours if inside image
			if row >= img.Bounds().Min.X && col >= img.Bounds().Min.Y &&
				row < img.Bounds().Max.X && col < img.Bounds().Max.Y {
				neighbours = append(neighbours, image.Pt(row, col))
			}
		}
	}
	return neighbours
}

func colorDiff(r1, g1, b1, r2, g2, b2 uint32) float64 {
	return math.Sqrt(float64((r1-r2)*(r1-r2) + (b1-b2)*(b1-b2) + (g1-g2)*(g1-g2)))
}
