package main

import (
	"math/rand"
	"strings"

	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"
)

type fairyTaleOptions struct {
	mainCharaters      []string
	supporterCharaters []string
	location           string
	storyPlot          string
}

// getFairyTaleOptions returns the options for the fairy tale from user input.
func getFairyTaleOptions(randomWithNoUserInput bool) fairyTaleOptions {
	if randomWithNoUserInput {

		mainCharacters := []string{CharacterMainSet[rand.Intn(len(CharacterMainSet))]}
		if slices.IndexFunc(mainCharacters, func(c string) bool { return strings.EqualFold(c, "Amy") }) == -1 {
			CharacterSupporterSet = append(CharacterSupporterSet, "Amy the princess")
		}

		return fairyTaleOptions{
			mainCharaters:      mainCharacters,
			supporterCharaters: []string{CharacterSupporterSet[rand.Intn(len(CharacterSupporterSet))]},
			location:           LocationSet[rand.Intn(len(LocationSet))],
			storyPlot:          StoryPlotSet[rand.Intn(len(StoryPlotSet))],
		}
	}

	selectedMainCharaters, _ := pterm.DefaultInteractiveMultiselect.WithOptions(CharacterMainSet).WithDefaultText("Select the main characters").Show()

	if len(selectedMainCharaters) == 0 {
		pterm.Warning.Println("No main characters selected, a random character will be selected.")
		selectedMainCharaters = append(selectedMainCharaters, CharacterMainSet[rand.Intn(len(CharacterMainSet))])
		if slices.IndexFunc(selectedMainCharaters, func(c string) bool { return strings.EqualFold(c, "Amy") }) == -1 {
			CharacterSupporterSet = append(CharacterSupporterSet, "Amy the princess")
		}
	}

	pterm.Info.Printfln("Selected main characters: %s", pterm.Green(selectedMainCharaters))

	selectedSupporterCharaters, _ := pterm.DefaultInteractiveMultiselect.WithOptions(CharacterSupporterSet).WithDefaultText("Select the main support characters").Show()

	if len(selectedSupporterCharaters) == 0 {
		pterm.Warning.Println("No support characters selected, a random character will be selected.")
		selectedSupporterCharaters = append(selectedSupporterCharaters, CharacterSupporterSet[rand.Intn(len(CharacterSupporterSet))])
	}

	pterm.Info.Printfln("Selected support characters: %s", pterm.Green(selectedSupporterCharaters))

	selectedLocation, _ := pterm.DefaultInteractiveSelect.WithOptions(LocationSet).WithDefaultText("Select the location").Show()

	if selectedLocation == "" {
		pterm.Warning.Println("No location selected, a random location will be selected.")
		selectedLocation = LocationSet[rand.Intn(len(LocationSet))]
	}

	pterm.Info.Printfln("Selected location: %s", pterm.Green(selectedLocation))

	storyPlot, _ := pterm.DefaultInteractiveSelect.WithOptions(StoryPlotSet).WithDefaultText("Select the plot").Show()

	if storyPlot == "" {
		pterm.Warning.Println("No plot selected, a random plot will be selected.")
		storyPlot = StoryPlotSet[rand.Intn(len(StoryPlotSet))]
	}

	pterm.Info.Printfln("Selected plot: %s", pterm.Green(storyPlot))

	return fairyTaleOptions{
		mainCharaters:      selectedMainCharaters,
		supporterCharaters: selectedSupporterCharaters,
		location:           selectedLocation,
		storyPlot:          storyPlot,
	}
}
