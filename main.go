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

//implementAIR

//Ask is the user has an API key setted up in the environment

type ImageFormat int

var imageURL string
var imagePath string
var title string
var isKeySetted bool = false
var finalGroup *huh.Group
var bgColor string = ""

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

	//	prompt := fmt.Sprintf("Image: A Pixel-Art Illustration of a CLI TertMINIMALISTIC Graffiti painted on it. The Graffiti says %s. Background Color: White. Graffiti Outline : %s. Colors: Vibrant Colors. Rules: Only display the Graffiti, DO NOT display anything else, DO NOT display hands!", "TITLE", "thick black")

	finalGroup = huh.NewGroup(

		huh.NewInput().
			Value(&title).
			Title("Title").
			Placeholder("My Title").
			Validate(func(s string) error {
				if s == "" {
					return errors.New("image name cannot be empty")
				}
				return nil
			}),

		huh.NewText().
			Value(&bgColor).
			Placeholder("Skyblue").
			Title("Background Color").
			//	Description("Any specific requirements?").
			CharLimit(400).
			Lines(5),
	)

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("GRAFO").
			Description("Welcome to Grafo the best CLI title Generator in the universe!")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Is your Open AI API key setted in your environment?").
				Value(&isKeySetted).
				Validate(func(b bool) error {

					if !b {
						//os.Exit(1)
						return errors.New("Please set your API key as OPENAI_API_KEY in your environment and the select 'Yes' to continue")
					}
					return nil
				})),

		finalGroup,
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	initImageGen := func() {
		fmt.Println(imagePath)
		prompt := fmt.Sprintf("Image: A Minimalistic Title"+
			"Painting Style: Graffiti "+
			"Effects:  3D effect. "+
			"Title Text %s. "+
			"Paddings: 20px. "+
			"Alignment: Center. "+
			"FontSize: 20px. "+
			"FontWeight: Bold. "+
			// uppercase?
			"Background Color: %s."+
			"Outline : %s."+
			"Colors: The color of the Title is highlighted."+
			"Rules: Only display the Title, DO NOT display anything else!",
			title, bgColor, "thin black")
		imagePath = imgGen.GenerateImg(prompt, imagePath, imageURL)
	}

	_ = spinner.New().Title("Generating your Title...").Accessible(accessible).Action(initImageGen).Run()

	{
		var sb strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&sb,
			"%s\n\n %s  %s  %s.",
			keyword("Your Title has been generated!"),

			lipgloss.NewStyle().Bold(true).Render(imagePath),
			keyword(xstrings.EnglishJoin(order.Image.Dimensions, true)),
			keyword(order.Image.Name),
		)

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
