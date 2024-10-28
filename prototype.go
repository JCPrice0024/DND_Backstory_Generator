package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func noFunnyBusiness(str string) string {
	re, err := regexp.Compile(`[^a-zA-Z]`)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return ""
	}

	// Replace non-letter characters with a space
	result := re.ReplaceAllString(str, "")

	return strings.ToLower(result)
}

// function to take in user input for prompt 1:: convert to http once ready
func prompt(invalidAns bool, promptType string, promptMap map[string]struct{}) (string, error) {
	// promptForInput asks the user for input and returns the input as a string.
	if invalidAns {
		fmt.Printf("Invalid %[1]s please try again or type skip to not enter a %[1]s:  ", promptType)
	} else {
		fmt.Printf("Prompt 1: What is your characters %s:  ", promptType)
	}

	scanner := bufio.NewScanner(os.Stdin) // Create a new scanner for reading input

	if scanner.Scan() {
		ans := noFunnyBusiness(scanner.Text())
		_, isInMap := promptMap[ans]
		if isInMap {
			return fmt.Sprintf("Your character's %s is %s", promptType, scanner.Text()), nil // Return the input text
		}
		if ans == "skip" {
			return "", nil
		}
		return prompt(true, promptType, promptMap)

	}

	return "", scanner.Err() // Return an error if scanning fails

}

// function to take in user input for prompt 1
func classPrompt() (string, error) {

	fmt.Println("What is your characters class:  ")
	scanner := bufio.NewScanner(os.Stdin) // Create a new scanner for reading input
	validAns := false
	class := ""
	for !validAns {
		if scanner.Scan() {
			ans := noFunnyBusiness(scanner.Text())
			_, isInMap := dndClassesAndSubclasses[ans]
			if isInMap {
				class = ans
				validAns = true
			}
			if ans == "skip" {
				return "", nil
			}
		} else {
			return "", scanner.Err()
		}

		if !validAns {
			fmt.Println("Invalid class please try again or type skip to not enter a class:  ")
		}

	}

	return prompt(false, "Subclass; if they don't have a subclass yet enter skip", dndClassesAndSubclasses[class])
}

func genPrompt(msg string) (string, error) {
	fmt.Println(msg)
	scanner := bufio.NewScanner(os.Stdin) // Create a new scanner for reading input

	if scanner.Scan() {
		return scanner.Text(), nil // Return the input text
	}

	return "", scanner.Err()
}

func main() {
	/*
		 PROJECT IDEA:
		 create a dnd character backstory generator in which you give the AI info about your character and their personality.
		prompt ideas: Personality, Name, Close Friends, Class, Subclass, Race, etc.
		add selector for serious/edgy or for lighthearted and goofy
		add prompt to filter out topics
		add a refresh button if the backstory isn't for them
		add a complete randomizer button that will completely randomize a character for you

	*/

	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	race, err := prompt(false, "Race", dndRaces)
	if err != nil {
		return
	}

	class, err := classPrompt()
	if err != nil {
		return
	}

	bg, err := prompt(false, "Background", backgrounds)
	if err != nil {
		return
	}

	name, err := genPrompt("What is your character's name")
	if err != nil {
		return
	}

	ideal, err := genPrompt("What is your character's ideal? If they don't have one hit enter")
	if err != nil {
		return
	}

	alignment, err := prompt(false, "Alignment", alignments)
	if err != nil {
		return
	}

	genInfo, err := genPrompt("Do you already have an idea of what you want your character to be? If not hit enter")
	if err != nil {
		return
	}

	ctx := context.Background()

	idealStr := ""
	if ideal != "" {
		idealStr = fmt.Sprintf("Their ideal is to %s", ideal)
	}

	aiRulesStr := fmt.Sprintf("You are making a compelling and unique backstory for a character. %s. %s. %s. %s Remeber that you are only writing a backstory, so don't include any extra fields also make sure not to give the character any items tied to their backstory unless specified in the prompt. Additionally try and make each backstory unique, don't just use the same story for each character", race, class, bg, alignment)
	promptStr := fmt.Sprintf("Can you help me make a backstory for my character %s, %s, here is what I already have: %s", name, idealStr, genInfo)
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, aiRulesStr),
		llms.TextParts(llms.ChatMessageTypeHuman, promptStr),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}), llms.WithTemperature(0.8))
	if err != nil {
		log.Fatal(err)
	}
	_ = completion

}

var (
	dndRaces = map[string]struct{}{
		"human":      {},
		"elf":        {},
		"dwarf":      {},
		"halfling":   {},
		"orc":        {},
		"dragonborn": {},
		"halfelf":    {},
		"halforc":    {},
		"tiefling":   {},
		"gnome":      {},
		"tabaxi":     {},
		"aarakocra":  {},
		"genasi":     {},
		"goliath":    {},
		"kenku":      {},
		"triton":     {},
		"firbolg":    {},
		"yuanti":     {},
		"warforged":  {},
	}
	backgrounds = map[string]struct{}{
		"acolyte":           {},
		"charlatan":         {},
		"criminal":          {},
		"entertainer":       {},
		"folkhero":          {},
		"guildartisan":      {},
		"herald":            {},
		"knight":            {},
		"mercenary":         {},
		"noble":             {},
		"outlander":         {},
		"pirate":            {},
		"sage":              {},
		"sailor":            {},
		"soldier":           {},
		"urchin":            {},
		"urbanbountyhunter": {},
		"citywatch":         {},
		"clancrafter":       {},
		"factionagent":      {},
		"rebel":             {},
		"wildling":          {},
	}
	alignments = map[string]struct{}{
		"lawfulgood":     {},
		"neutralgood":    {},
		"chaoticgood":    {},
		"lawfulneutral":  {},
		"trueneutral":    {},
		"chaoticneutral": {},
		"lawfulevil":     {},
		"neutralevil":    {},
		"chaoticevil":    {},
	}
	fighterSubclasses = map[string]struct{}{
		"battlemaster":   {},
		"cavalier":       {},
		"champion":       {},
		"eldritchknight": {},
		"psiwarrior":     {},
	}

	rogueSubclasses = map[string]struct{}{
		"arcanetrickster": {},
		"assassin":        {},
		"inquisitive":     {},
		"mastermind":      {},
		"phantom":         {},
		"scout":           {},
		"swashbuckler":    {},
	}

	wizardSubclasses = map[string]struct{}{
		"abjuration":    {},
		"conjuration":   {},
		"divination":    {},
		"enchantment":   {},
		"evocation":     {},
		"illusion":      {},
		"necromancy":    {},
		"transmutation": {},
	}

	clericSubclasses = map[string]struct{}{ // The best class
		"knowledge": {},
		"life":      {},
		"light":     {},
		"nature":    {},
		"tempest":   {},
		"trickery":  {},
		"war":       {},
	}

	bardSubclasses = map[string]struct{}{
		"collegeofcreation":  {},
		"collegeofeloquence": {},
		"collegeoflore":      {},
		"collegeofswords":    {},
		"collegeofvalor":     {},
		"collegeofwhispers":  {},
	}

	rangerSubclasses = map[string]struct{}{
		"beastmaster":   {},
		"gloomstalker":  {},
		"horizonwalker": {},
		"monsterslayer": {},
		"swarmkeeper":   {},
	}

	sorcererSubclasses = map[string]struct{}{
		"aberrantmind":      {},
		"clockworksoul":     {},
		"draconicbloodline": {},
		"divinesoul":        {},
		"wildmagic":         {},
	}

	warlockSubclasses = map[string]struct{}{
		"archfey":     {},
		"celestial":   {},
		"fiend":       {},
		"greatoldone": {},
		"hexblade":    {},
	}

	paladinSubclasses = map[string]struct{}{
		"oathoftheancients": {},
		"oathofconquest":    {},
		"oathofdevotion":    {},
		"oathofthecrown":    {},
		"oathofvengeance":   {},
	}

	monkSubclasses = map[string]struct{}{
		"ancienttradition": {},
		"drunkenmaster":    {},
		"fourelements":     {},
		"kensei":           {},
		"longdeath":        {},
		"mercy":            {},
		"shadow":           {},
		"sunsoul":          {},
	}

	barbarianSubclasses = map[string]struct{}{
		"ancestralprotectors": {},
		"battlerager":         {},
		"berserker":           {},
		"stormherald":         {},
		"totemwarrior":        {},
		"wildmagic":           {},
	}

	druidSubclasses = map[string]struct{}{
		"circleoftheland":     {},
		"circleofthemoon":     {},
		"circleoftheshepherd": {},
		"circleofthespores":   {},
		"circleofthestars":    {},
		"circleofthewildfire": {},
	}

	artificerSubclasses = map[string]struct{}{
		"alchemist":     {},
		"artillerist":   {},
		"battlesmith":   {},
		"clockworksoul": {},
		"godoftheforge": {},
		"wildmagic":     {},
	}

	dndClassesAndSubclasses = map[string]map[string]struct{}{
		"barbarian": barbarianSubclasses,
		"bard":      bardSubclasses,
		"cleric":    clericSubclasses,
		"druid":     druidSubclasses,
		"fighter":   fighterSubclasses,
		"monk":      monkSubclasses,
		"paladin":   paladinSubclasses,
		"ranger":    rangerSubclasses,
		"rogue":     rogueSubclasses,
		"sorcerer":  sorcererSubclasses,
		"warlock":   warlockSubclasses,
		"wizard":    wizardSubclasses,
		"artificer": artificerSubclasses,
	}
)
