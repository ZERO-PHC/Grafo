package imgGen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"graffiti/img"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

var asciiChars = "@%#*+=-:. "

const density = "Ñ@#W$9876543210?!abc;:+=-,._ "
const max = 765

func GenerateImg(prompt string, imagePath string, imageURL string) string {
	var _ = "1792x1024"
	var _ = "1024x1792"
	//imageURL = "https://oaidalleapiprodscus.blob.core.windows.net/private/org-yMUJFRLsbHiDT7L6B3xxdyO7/user-g8jvMhQIkeeR5EXwNu0cL9nD/img-zbhqqxJYDjVGewk42ibVV9Nx.png?st=2024-02-09T16%3A58%3A02Z&se=2024-02-09T18%3A58%3A02Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2024-02-08T23%3A10%3A37Z&ske=2024-02-09T23%3A10%3A37Z&sks=b&skv=2021-08-06&sig=qQCUhc1Rw8g/08okKbQ5xW0VgI1/9aQaEkRHgx9%2B/2E%3D"

	//fmt.Println(escena_res)
	//fmt.Println(personaje_res)

	data := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": prompt,
		"n":      1, // only
		//	"size":            "256x256",
		"style":           "natural",
		"response_format": "url",
		//"quality":         "hd",
	}

	jsonData, err := marshalData(data)
	if err != nil {
		log.Fatalf("%v", err)
	}

	req, err := createRequest(jsonData)
	if err != nil {
		log.Fatalf("%v", err)
	}

	resp, err := sendRequest(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer resp.Body.Close()

	body, err := readResponseBody(resp)
	if err != nil {
		log.Fatalf("%v", err)
	}

	imageURL, err = unmarshalResponse(body)
	if err != nil {
		log.Fatalf("%v", err)
	}

	imageBytes, err := downloadImage(imageURL)
	if err != nil {
		log.Fatalf("%v", err)
	}

	imagePath, err = saveImage(imageBytes)
	if err != nil {
		log.Fatalf("%v", err)
	}

	//println(imageURL)
	loadImage(imageURL)

	return imagePath
}

func marshalData(data map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error encoding JSON: %v", err)
	}
	return jsonData, nil
}

func createRequest(jsonData []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	return req, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	return resp, nil
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}
	return body, nil
}

func unmarshalResponse(body []byte) (string, error) {
	var result struct {
		Created int `json:"created"`
		Data    []struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("Error unmarshalling response body: %v", err)
	}
	if result.Data == nil || len(result.Data) == 0 {
		return "", fmt.Errorf("No image data found in the response")
	}
	return result.Data[0].URL, nil
}

func downloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer resp.Body.Close()
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading image data: %v", err)
	}
	return imageBytes, nil
}

func saveImage(imageBytes []byte) (string, error) {
	uid := time.Now().Format("20060102150405")
	imagePath := fmt.Sprintf("gens/image_%s.png", uid)
	err := os.WriteFile(imagePath, imageBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("Error saving image to file: %v", err)
	}
	return imagePath, nil
}

// This function will map a brightness value to a color attribute
func getBrightnessColorAttribute(brightness uint32) color.Attribute {
	// Map the brightness to a range of 0-5, for example
	colorRange := brightness * 5 / max

	// Map this range to a specific color
	switch colorRange {
	case 0:
		return color.FgHiBlack
	case 1:
		return color.FgRed
	case 2:
		return color.FgGreen
	case 3:
		return color.FgYellow
	case 4:
		return color.FgBlue
	default:
		return color.FgHiWhite
	}
}

func loadImage(url string) {
	println("Loading image...")
	//staticUrl := "https://cdn.discordapp.com/attachments/1167981493872758844/1205246473109766174/image_20240208113238.png?ex=65d7ac3f&is=65c5373f&hm=4edc50967cc620cb5cb5bb58b43bd69963a2550ebefe3256025d49a7ddf7e634&"

	//img, err := img.Load(url)
	img, err := img.Load(url)

	if err != nil {
		panic(err)
	}

	f, err := os.Create("./result.txt")

	if err != nil {
		panic(err)
	}

	err = os.Truncate("./result.txt", 0)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for i := 0; i < img.Bounds().Max.Y-1; i++ {
		for j := 0; j < img.Bounds().Max.X-1; j++ {
			r8, g8, b8, _ := img.At(j, i).RGBA()

			r := r8 >> 8
			g := g8 >> 8
			b := b8 >> 8

			sum := r + g + b
			percent := int((sum * 100) / max)

			floored := (len(density) * percent) / 100

			if floored >= len(density) {
				floored = len(density) - 1
			}
			char := density[floored]

			// Calculate a color attribute based on the brightness
			//colorAttr := color.Attribute(int(float64(sum) * 165 / 765))
			colorAttr := getBrightnessColorAttribute(sum)

			// Create a new color with the calculated attribute
			c := color.New(colorAttr)

			// Print the character with the color
			c.Printf("%s", string(char))
			f.Write([]byte(string(char)))

		}
		f.Write([]byte("\n"))
		fmt.Printf("\n")
	}
}

func imageToASCII(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening image file: %v", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image file: %v", err)
		return
	}

	bounds := img.Bounds()
	ditherMatrix := [4][4]int{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 10 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 10 {
			var sum uint32
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if x+j < bounds.Max.X && y+i < bounds.Max.Y {
						r, g, b, _ := img.At(x+j, y+i).RGBA()
						brightness := (299*r + 587*g + 114*b) / (1000 * 65535)
						if brightness > uint32(ditherMatrix[i][j]*16) {
							sum += 255
						}
					}
				}
			}
			fmt.Print(string(toASCII(sum / 16)))
		}
		fmt.Println()
	}
}

func toASCII(brightness uint32) rune {
	index := brightness * uint32(len(density)) / 256
	return rune(density[index])
}
