package main

import (
	"errors"
	"fmt"
	"grafo/imgGen"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

//TODO auto breakliner

var asciiTitle = ":!!!!!!!!6!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!7!\n" +
	":!!!!!!!!5!!!!!!!!!!!!!!!@#W#######6!!!!!!a###########W!!!!ÃÃ###Ã##W!!!!!!!!1Ã3###84671!!!!!!!0@59W56!!!!!!!8!\n" +
	":!!!!!!!!6!!!!!!!!!!Ã###100000000000##3ÃÃ@#10000000000000W#ÃÃ##100000#W!!ÃÃÃÃ##000000000000#####000000000#W?!!$!\n" +
	":!!!!!!!!6!!!!!!!ÃÃ##100000000000000?##ÃÃ#10000000000000002###00000000##ÃÃÃÃ##0000000000000##00000000000000Ã#!W!\n" +
	":!!!!!!!!4!!!!$ÃÃ##10000000@####000@##ÃÃÃ#900000##Ã##000000##0000600000##ÃÃÃ#1000Ã##WÃÃÃÃ##00000######900000Ã#!\n" +
	":!!!!!!!!3!!!ÃÃÃ@#1000000##ÃÃÃÃÃÃÃ#####WÃ##00000Ã##000000?##0000Ã##0000?##ÃÃ#00000000#ÃÃ#10000##ÃÃÃÃÃÃ#W00000#W!\n" +
	":!!!!!!!!4a!ÃÃÃÃ##000000Ã#0$ÃÃ##20000001##W00000000000####$0000#600000005#@#W000000###Ã#10000Ã#a99ÃÃÃÃ#W00000WÃ!\n" +
	":!!!!!!!!3a$ÃÃÃÃ##0000006##Ã##00000000000Ã#0000000000000Ã##10000000000000###1000Ã#ÃÃÃÃÃ#000000##9ÃÃÃ$#100000##@!\n" +
	":!!!!!!!!3?$WÃÃÃÃ##0000000Ã####0W#6000000Ã##00000Ã#000000000#####ÃÃÃ##0000Ã#0000##ÃÃÃÃÃ##00000000110000000#?!@!\n" +
	":!!!!!!!!1a$$ÃÃÃÃÃ##00000000000000000000###000000####8000000Ã#ÃÃÃÃÃÃ#W0000Ã####W9ÃÃÃÃÃÃÃ#00000000000000##!!!!@!\n" +
	":!!!!!!!!2!!$$9ÃÃÃÃÃ##600000000000000###ÃÃÃ#W00000Ã#ÃÃÃÃ#####?!$$ÃÃÃÃÃ#W000##!!!!$$ÃÃÃÃÃÃÃ##W##W#####0!!!!!!!@!\n" +
	":!!!!!!!!1!!!a$$$ÃÃÃÃÃÃÃ###########Ã1$$ÃÃÃÃÃÃ####WÃ$$ÃÃÃÃ1!!!!!!9$9ÃÃÃÃÃÃÃ!!!!!!!!!9$$$ÃÃÃÃÃÃÃÃÃÃ!!!!!!!!!!!!!\n" +
	":!!!!!!!!1a!!!!!?$$$$ÃÃÃÃÃÃÃÃÃÃÃ$!!!!!?$$$ÃÃÃÃ#!!!!!!!!!!!!!!!!!!!2$$9!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Ã!\n" +
	":!!!!!!!!0!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Ã!\n"

var imageURL string
var imagePath string
var title string
var isKeySetted bool = false
var finalGroup *huh.Group
var bgColor string = ""

func main() {

	finalGroup = huh.NewGroup(

		huh.NewInput().
			Value(&title).
			Title("Title").
			Placeholder("My Title").
			Validate(func(s string) error {
				if s == "" {
					return errors.New("Title name cannot be empty")
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
			"%s\n\n %s .",
			keyword("Your Title has been generated!"),

			lipgloss.NewStyle().Bold(true).Render(imagePath),
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
