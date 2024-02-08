package main

import (
	"errors"
	"fmt"
	"graffiti/imgGen"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
)

// preview in ASCII:

// TODO getOsEnv
// const OPENAI_API_KEY = "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"

//implementAIR

//Ask is the user has an API key setted up in the environment

type ImageFormat int

var imageURL string
var imagePath string
var apiKey string

const (
	JPEG ImageFormat = iota + 1
	PNG
	GIF
)

func (f ImageFormat) String() string {
	switch f {
	case JPEG:
		return "JPEG "
	case PNG:
		return "PNG "
	case GIF:
		return "GIF "
	default:
		return ""
	}
}

type Image struct {
	Format     ImageFormat
	Dimensions []string
	Name       string
}

type Order struct {
	Image        Image
	Instructions string
}

func main() {
	var image Image
	var order = Order{Image: image}

	prompt := fmt.Sprintf("Image: A Digital Illustration of the final result of a  Graffiti that says %s. Background Color: White. Graffiti Outline : %s. Colors: Vibrant Colors. Rules: Only display the Graffiti, DO NOT display anything else, DO NOT display hands!", "GRAFFITI", "thick black")
	imagePath = imgGen.GenerateImg(prompt, imagePath, imageURL)

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Image Generator").
			Description("Welcome to _Image Generator™_.\n\nHow may we assist you?")),

		// Choose an image format.
		//huh.NewGroup(
		//huh.NewSelect[ImageFormat]().
		//	Options(
		//		huh.NewOption("JPEG", JPEG).Selected(true),
		//		huh.NewOption("PNG", PNG),
		//		huh.NewOption("GIF", GIF),
		//	).
		//	Title("Choose your image format").
		//	Description("We support JPEG, PNG, and GIF formats.").
		//	Value(&order.Image.Format),

		//	huh.NewMultiSelect[string]().
		//		Title("Dimensions").
		//		Description("Choose the width and height.").
		//		Options(
		//			huh.NewOption("800x600", "800x600").Selected(true),
		//			huh.NewOption("1024x768", "1024x768"),
		//			huh.NewOption("1280x720", "1280x720"),
		//			huh.NewOption("1920x1080", "1920x1080"),
		//		).
		//Validate(func(t []string) error {
		//	if len(t) != 1 {
		//		return fmt.Errorf("please select width and height")
		//		}
		//	return nil
		//	}).
		//	Value(&order.Image.Dimensions).
		//	Filterable(true).
		//	Limit(2),
		//),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().Value(&apiKey).
				Title("API Key").
				Placeholder("Your API key").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("API key cannot be empty")
					}
					return nil
				}).
				Description("Your API key is required to generate the image."),

			huh.NewInput().
				Value(&order.Image.Name).
				Title("What's the image name?").
				Placeholder("my_image").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("image name cannot be empty")
					}
					return nil
				}).
				Description("For your reference."),

			huh.NewText().
				Value(&order.Instructions).
				Placeholder("Please use vibrant colors").
				Title("Special Instructions").
				Description("Any specific requirements?").
				CharLimit(400).
				Lines(5),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// sprint f var called prompt

	initImageGen := func() {
		fmt.Println(imagePath)
		//imageToASCII(imagePath)

	}

	_ = spinner.New().Title("Generating your image...").Accessible(accessible).Action(initImageGen).Run()

	// Print order summary.
	{
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&sb,
			"%s\n\nOne %s image of dimensions %s named %s.",
			lipgloss.NewStyle().Bold(true).Render(imagePath),

			keyword(order.Image.Format.String()),
			keyword(xstrings.EnglishJoin(order.Image.Dimensions, true)),
			keyword(order.Image.Name),
		)

		//print image url

		//printapi key
		//fmt.Fprintf(&sb, "\n\nAPI Key: %s", keyword(apiKey))

		fmt.Fprintf(&sb, "\n\nThanks for the interaction!")

		fmt.Println(
			lipgloss.NewStyle().
				Width(100).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}

func scaleDown(c uint32, max uint32) uint8 {
	return uint8((c*255 + max/2) / max)
}
