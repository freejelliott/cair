package main

import (
	"jelliott.dev/cair/pkg/cair"

	"image"
	"image/jpeg"
	"log"
	"os"
)

var filename = ""

func main() {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("err opening file: %s", err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("err decoding image: %s", err)
	}
	// energies := cair.GradientEnergy(img)
	// fmt.Println(energies)
	// newImg := image.NewGray16(img.Bounds())
	// for i := 0; i < len(energies); i++ {
	// 	for j := 0; j < len(energies[i]); j++ {
	// 		newImg.SetGray16(newImg.Bounds().Min.X+i, newImg.Bounds().Min.Y+j, color.Gray16{Y: uint16(energies[i][j])})
	// 	}
	// }
	newImg := cair.Resize(img, image.Pt(int(float64(img.Bounds().Dx())*0.66), int(float64(img.Bounds().Dy())*0.66)))
	newFile, err := os.Create("")
	if err != nil {
		log.Fatalf("err creating new image: %s", err)
	}
	err = jpeg.Encode(newFile, newImg, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("err encoding jpeg image: %s", err)
	}
}
