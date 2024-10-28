package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Character struct {
	Name         string   `json:"name"`
	Class        string   `json:"class"`
	Gender       string   `json:"gender"`
	Subclass     string   `json:"sclass"`
	Multiclasses []string `json:"mclass"`
	Background   string   `json:"bg"`
	Race         string   `json:"race"`
	Alignment    string   `json:"align"`
	Ideal        string   `json:"ideal"`
	GenBInfo     string   `json:"gen"`
	Prev         string   `json:"prev"`
	Additions    string   `json:"add"`
}

const aiRulesStr = `1. Generate unique and creative content based on the given prompt.
2. Focus on originality; avoid rehashing common themes or clichés.
3. Do not add unnecessary backstory elements unless explicitly requested.
4. Ensure that character actions and dialogues are aligned with their established traits and motivations.
5. Keep responses concise and relevant to the prompt without unnecessary elaboration.
6. Encourage imaginative scenarios while respecting the context provided.
7. Avoid filler content; every element should serve a purpose in the narrative.
8. Write the backstory directly without any introductory phrases or framing. Start with the character's name and action, e.g., 'Arin the Brave grew up in a small village...'.`

var ctx context.Context

func main() {
	ctx = context.Background()

	http.HandleFunc("/generateCharacter", generate)
	http.HandleFunc("/adjustCharacter", buildingOff)
	http.HandleFunc("/reroll", reset)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var c Character
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	// c.readyForPrompt()

	if c.Subclass != "" {
		c.Subclass = fmt.Sprintf("- Subclass : %s", c.Subclass)
	}

	mStr := ""
	if len(c.Multiclasses) != 0 {
		mStr = fmt.Sprintf("- Multiclasses: %s", strings.Join(c.Multiclasses, ", "))
	}

	if c.GenBInfo != "" {
		c.GenBInfo = fmt.Sprintf(`Start with this premise - "%s"`, c.GenBInfo)
	}

	promptStr := fmt.Sprintf(`Generate a backstory for a character with the following details:
		- Name: %s
		- Class: %s
		%s
		%s
		- Race: %s
		- Gender: %s
		- Ideal: %s
		- Background: %s
		- Alignment: %s
		
		%s
		
		Ensure the backstory is engaging and flows naturally, without introductory phrases or framing. Begin directly with the character's name and their story.`, c.Name, c.Class, c.Subclass, mStr, c.Race, c.Gender, c.Ideal, c.Background, c.Alignment, c.GenBInfo)
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, aiRulesStr),
		llms.TextParts(llms.ChatMessageTypeHuman, promptStr),
	}
	aiRes, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.5))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Class, promptStr)

	if len(aiRes.Choices) > 0 && aiRes.Choices[0].Content != "" {
		response := map[string]string{
			"message": aiRes.Choices[0].Content,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Something went wrong, please try again later.", http.StatusBadRequest)
		return
	}

}

func buildingOff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var c Character
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if c.Prev == "" || c.Additions == "" {
		http.Error(w, "Missing content aborting", http.StatusBadRequest)
		return
	}

	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	// c.readyForPrompt()

	if c.Subclass != "" {
		c.Subclass = fmt.Sprintf("- Subclass : %s", c.Subclass)
	}

	mStr := ""
	if len(c.Multiclasses) != 0 {
		mStr = fmt.Sprintf("- Multiclasses: %s", strings.Join(c.Multiclasses, ", "))
	}

	if c.GenBInfo != "" {
		c.GenBInfo = fmt.Sprintf(`Start with this premise - "%s"`, c.GenBInfo)
	}

	promptStr := fmt.Sprintf(`Adjust the backstory for the character based on the previous narrative provided below. Expand on their journey, motivations, or relationships while maintaining the established tone and details.

		**Previous Backstory:**
		"%s"
		- Name: %s
		- Class: %s
		%s
		%s
		- Race: %s
		- Gender: %s
		- Ideal: %s
		- Background: %s
		- Alignment: %s

		-UserAdditions: %s
		
		Ensure the adjustment flows naturally from the existing narrative and includes the additions provided by the user.`, c.Prev, c.Name, c.Class, c.Subclass, mStr, c.Race, c.Gender, c.Ideal, c.Background, c.Alignment, c.Additions)
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, aiRulesStr),
		llms.TextParts(llms.ChatMessageTypeHuman, promptStr),
	}
	aiRes, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.5))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Class, promptStr)

	if len(aiRes.Choices) > 0 && aiRes.Choices[0].Content != "" {
		response := map[string]string{
			"message": aiRes.Choices[0].Content,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Something went wrong, please try again later.", http.StatusBadRequest)
		return
	}
}

func reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var c Character
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if c.Prev == "" {
		http.Error(w, "Missing content aborting", http.StatusBadRequest)
		return
	}

	llm, err := ollama.New(ollama.WithModel("llama3.1"))
	if err != nil {
		log.Fatal(err)
	}

	// c.readyForPrompt()

	if c.Subclass != "" {
		c.Subclass = fmt.Sprintf("- Subclass : %s", c.Subclass)
	}

	mStr := ""
	if len(c.Multiclasses) != 0 {
		mStr = fmt.Sprintf("- Multiclasses: %s", strings.Join(c.Multiclasses, ", "))
	}

	if c.GenBInfo != "" {
		c.GenBInfo = fmt.Sprintf(`Start with this premise - "%s"`, c.GenBInfo)
	}

	if c.GenBInfo != "" {
		c.GenBInfo = fmt.Sprintf(`Start with this premise - "%s"`, c.GenBInfo)
	}

	promptStr := fmt.Sprintf(`Give me a new backstory for the character below. Do not use any of the material from the previous backstory mentioned below.

		**Previous Backstory:**
		"%s"
		- Name: %s
		- Class: %s
		%s
		%s
		- Race: %s
		- Gender: %s
		- Ideal: %s
		- Background: %s
		- Alignment: %s

		%s
		
		Ensure the new backstory flows naturally and does not use the same themes or content as the previous one.`, c.Prev, c.Name, c.Class, c.Subclass, mStr, c.Race, c.Gender, c.Ideal, c.Background, c.Alignment, c.GenBInfo)
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, aiRulesStr),
		llms.TextParts(llms.ChatMessageTypeHuman, promptStr),
	}
	aiRes, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.5))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Class, promptStr)

	if len(aiRes.Choices) > 0 && aiRes.Choices[0].Content != "" {
		response := map[string]string{
			"message": aiRes.Choices[0].Content,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Something went wrong, please try again later.", http.StatusBadRequest)
		return
	}
}
