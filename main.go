package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
)

// 8 means 256, accepted values: 6, 7, 8
var scale int = 6
var km int

func main() {

	km = intPow(2, scale)
	downloadImages()
	makeGIF()
}

func makeGIF() {
	// Retrieve the image filenames in the "tmp" folder
	imageFilenames, err := filepath.Glob("tmp/*.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Create a list to store the image frames
	frames := make([]*image.Paletted, 0, len(imageFilenames))
	delays := make([]int, 0, len(imageFilenames))

	// Open and append each image to the frames list
	for _, filename := range imageFilenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		// Resize the image to the desired dimensions
		resizedImg := resize.Resize(577, 400, img, resize.Lanczos3)

		// Create a new paletted image
		palettedImage := image.NewPaletted(resizedImg.Bounds(), palette.Plan9)

		// Quantize the image
		draw.Draw(palettedImage, palettedImage.Bounds(), resizedImg, image.Point{}, draw.Src)

		frames = append(frames, palettedImage)
		delays = append(delays, 30) // Set the delay between frames (in 1/100th of a second)
	}

	gifFilename := fmt.Sprintf("animation_%03d.gif", km)
	err = saveGIF(gifFilename, frames, delays)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("GIF file generated: %s\n", gifFilename)

	// Delete the "tmp" folder and its contents
	err = os.RemoveAll("tmp")
	if err != nil {
		log.Fatal(err)
	}
}

func saveGIF(filename string, frames []*image.Paletted, delays []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	g := gif.GIF{}
	for i, frame := range frames {
		g.Image = append(g.Image, frame)
		g.Delay = append(g.Delay, delays[i])
	}

	err = gif.EncodeAll(file, &g)
	if err != nil {
		return err
	}

	return nil
}

func downloadImages() {
	baseURL := fmt.Sprintf("https://www.weather.gov.hk/wxinfo/radars/rad_%03d_png/2d%03dnradar_", km, km)
	currentTime := time.Now()

	// Create the "tmp" folder if it doesn't exist
	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the 'tmp' folder:", err)
		return
	}

	// Round down the current minute to the nearest multiple of 12
	currentMinute := currentTime.Minute()
	roundedMinute := currentMinute - (currentMinute % 12)
	minuteDiff := currentMinute % 12
	noOfPic := 10
	fmt.Printf("currentMinute: %v, round minute: %v\n", currentMinute, roundedMinute)
	// Iterate over the desired timestamps and download the images
	for i := 0; i < noOfPic; i++ {
		// Calculate the timestamp for the current image
		timestamp := currentTime.Add(-time.Duration(12*(i+1)+minuteDiff)*time.Minute + 8*time.Hour)

		// Format the timestamp to match the image filename
		formattedTimestamp := timestamp.Format("200601021504")

		// Construct the image URL
		imageURL := fmt.Sprintf("%s%s.jpg", baseURL, formattedTimestamp)
		fmt.Printf("the url to be downloaded is: %v\n", imageURL)
		// Set the file path to save the image
		filePath := filepath.Join("tmp", fmt.Sprintf("image_%d.jpg", noOfPic-i-1))

		// Download the image
		err := downloadImage(imageURL, filePath)
		if err != nil {
			fmt.Printf("Error downloading image %d: %v\n", i, err)
		}
	}
}

func downloadImage(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func intPow(base, exponent int) int {
	result := 1

	for i := 0; i < exponent; i++ {
		result *= base
	}

	return result
}
