package img

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func Load(url string) (image.Image, error) {
	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)

	if err != nil {
		panic(err)
	}

	if err != nil {
		return nil, err
	}

	// Scale down the image
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth := int(float64(width) / 9)
	newHeight := int(float64(height) / 32)

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	// Loop over the height of the new image
	for y := 0; y < newHeight; y++ {
		// Loop over the width of the new image
		for x := 0; x < newWidth; x++ {
			// Set the pixel at the current x, y coordinate of the new image
			// to the pixel at the corresponding x*3, y*3 coordinate of the original image
			// This effectively scales down the image by a factor of 3
			newImg.Set(x, y, img.At(x*9, y*32))
		}
	}

	return newImg, nil
}
