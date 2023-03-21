package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/pterm/pterm"
)

var (
	apiKey       = os.Getenv("OPENAI_API_KEY")
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	CharacterMainSet = []string{"Amy the princess", "Julia the queen", "Daniel the king"}

	// Character set and storyline for the fairy tale.
	CharacterSupporterSet = []string{
		"Alice the princess from Spain",
		"Bob the knight",
		"Dennis the hedgehog",
		"Eve the elf",
		"Frank the bird",
		"George the hourse",
		"Linda the kind witch",
		"Tina child from the village",
		"Ursula the evil witch",
		"Victor the dragon",
		"Wendy the fairy",
		"Martina the princess from France and wife of Horst",
		"Horst the king from France and husband of Martina",
		"Angela the bridges troll",
		"Angela the mermaid",
		"Markus the crab",
		"Xavier the wizard",
		"Henry the griffin",
		"Bruno the werewolf",
		"Elsa the ice queen",
		"Anna the ice queen",
		"Penelope the pirate",
	}

	LocationSet = []string{
		"Enchanted forests",
		"Magical kingdoms",
		"Remote villages",
		"Mysterious castles",
		"mountains in snow",
		"Enchanted islands",
		"Dark and mystical caves",
		"Magical schools",
		"Under the sea",
	}

	StoryPlotsSet = []string{
		"A story where the protagonist befriends various animals who help them on their adventure, teaching the importance of kindness and empathy towards all living creatures.",
		"A child discovers a magical garden filled with talking plants and friendly insects, learning about nature and the interconnectedness of all living things.",
		"A young protagonist uses their intelligence and resourcefulness to solve a problem or outwit a more powerful antagonist.",
		"A hero sets out on a journey to retrieve a magical object or sacred artifact that has the power to restore peace or prosperity to their kingdom.",
		"A hero sets out on a journey to rescue a loved one from a villain.",
		"A child's favorite toy goes missing, leading them on a quest to find it, during which they learn about the importance of perseverance and problem-solving.",
		"The protagonist embarks on a journey to find the end of a rainbow, encountering colorful and magical characters that teach them about the beauty of diversity.",
		"A hero sets out on a journey to defeat a villain.",
		"A protagonist sets out on a journey to discover a hidden treasure, learning valuable life lessons along the way",
		"A young protagonist becomes lost or separated from their family or friends and must embark on a journey to find their way back, often with the help of kind strangers or magical creatures.",
		"Stories featuring talking animals that help the protagonist overcome challenges, imparting wisdom or offering companionship along the way.",
		"A protagonist discovers they possess magical powers or undergoes a magical transformation, learning to use their newfound abilities for good.",
		"A young character faces their fears and learns to be brave, often with the support of friends or magical allies.",
		"A young character receives help from a magical creature or being, such as a fairy godmother, who helps them overcome challenges and achieve their dreams.",
		"A story about an unlikely friendship between two characters from different backgrounds or worlds, highlighting the importance of understanding, empathy, and overcoming prejudices.",
		"A magical object or artifact plays a central role in the story, granting the protagonist special abilities or leading them on an adventure.",
		"A wise and experienced mentor helps guide the protagonist through their challenges, imparting wisdom and teaching valuable lessons.",
		"The protagonist takes on the role of a guardian or protector of nature, working with magical creatures to preserve the balance of the natural world and learning the importance of environmental stewardship.",
		"The protagonist is chosen or destined for a special task, such as saving their kingdom or community, and must learn to harness their unique gifts or talents",
	}
)

func main() {
	if apiKey == "" || awsAccessKey == "" || awsSecretKey == "" {
		fmt.Println("Please set the environment variables OPENAI_API_KEY, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY.")
		return
	}

	selectedMainCharaters, selectedSupporterCharaters, location, storyPlot := getFairyTaleOptions()

	pterm.Info.Println("Generating fairy tale...")

	ts := createTimestamp()
	generate(selectedMainCharaters, selectedSupporterCharaters, location, storyPlot, ts)

}

func createTimestamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02_15-04-05")
}

func getFairyTaleOptions() (selectedMainCharaters, selectedSupporterCharaters []string, selectedLocation, storyPlot string) {
	selectedMainCharaters, _ = pterm.DefaultInteractiveMultiselect.WithOptions(CharacterMainSet).WithDefaultOptions(CharacterMainSet).WithDefaultText("Select the main characters").Show()

	if len(selectedMainCharaters) == 0 {
		pterm.Error.Println("No main characters selected, a random character will be selected.")
		selectedMainCharaters = append(selectedMainCharaters, CharacterMainSet[rand.Intn(len(CharacterMainSet))])
	}

	pterm.Info.Printfln("Selected main characters: %s", pterm.Green(selectedMainCharaters))

	selectedSupporterCharaters, _ = pterm.DefaultInteractiveMultiselect.WithOptions(CharacterSupporterSet).WithDefaultText("Select the main support characters").Show()

	if len(selectedSupporterCharaters) == 0 {
		pterm.Error.Println("No support characters selected, a random character will be selected.")
		selectedSupporterCharaters = append(selectedSupporterCharaters, CharacterSupporterSet[rand.Intn(len(CharacterSupporterSet))])
	}

	pterm.Info.Printfln("Selected support characters: %s", pterm.Green(selectedSupporterCharaters))

	selectedLocation, _ = pterm.DefaultInteractiveSelect.WithOptions(LocationSet).WithDefaultText("Select the location").Show()

	if selectedLocation == "" {
		pterm.Error.Println("No location selected, a random location will be selected.")
		selectedLocation = LocationSet[rand.Intn(len(LocationSet))]
	}

	pterm.Info.Printfln("Selected location: %s", pterm.Green(selectedLocation))

	storyPlot, _ = pterm.DefaultInteractiveSelect.WithOptions(StoryPlotsSet).WithDefaultText("Select the plot").Show()

	if storyPlot == "" {
		pterm.Error.Println("No plot selected, a random plot will be selected.")
		storyPlot = StoryPlotsSet[rand.Intn(len(StoryPlotsSet))]
	}

	pterm.Info.Printfln("Selected plot: %s", pterm.Green(storyPlot))

	return
}

func generate(mainCharaters, supporterCharaters []string, location, storyPlot, ts string) {
	fairyTale, prompt, err := generateFairyTale(apiKey, mainCharaters, supporterCharaters, location, storyPlot)
	if err != nil {
		fmt.Printf("Error generating fairy tale: %v\n", err)
		return
	}

	f, _ := os.Create(fmt.Sprintf("%s_fairy_tale.txt", ts))
	defer f.Close()

	f.WriteString("Prompt:\n")
	f.WriteString(prompt)
	f.WriteString("\n\nPlot:\n")
	f.WriteString(fairyTale)

	fmt.Println("Generated fairy tale saved!")

	outputFilename := fmt.Sprintf("%s_fairy_tale.mp3", ts)
	err = textToSpeech(fairyTale, outputFilename)
	if err != nil {
		fmt.Printf("Error converting text to speech: %v\n", err)
		return
	}

	fmt.Printf("German audio saved to: %s\n", outputFilename)

	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", outputFilename)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error opening audio file: %v\n", err)
		return
	}
}

func aggregateSlice(input []string) string {
	var output string
	for i, v := range input {
		if i == len(input)-1 {
			output += fmt.Sprintf("%s", v)
		} else {
			output += fmt.Sprintf("%s, ", v)
		}
	}
	return output
}

func generateFairyTale(apiKey string, mainCharaters []string, supporterCharaters []string, location, storyPlot string) (string, string, error) {

	mainCharatersAggregated := aggregateSlice(mainCharaters)
	supporterCharatersAggregated := aggregateSlice(supporterCharaters)

	prompt := fmt.Sprintf(`Create a children fairy tale in German with the following setup.
	
main characters: %s
support characters: %s

the story takes place in %s
the plot of the main characters is %s

The fairy tale should be funny and entertaining for children, it should be about 5-10 pages long.
`, mainCharatersAggregated, supporterCharatersAggregated, location, storyPlot)

	data := map[string]interface{}{
		"model":       "text-davinci-003", // gpt-3.5-turbo or text-davinci-003
		"prompt":      prompt,
		"n":           1,
		"max_tokens":  4097 - len(prompt)/4, // model's maximum context length is 4097 tokens
		"temperature": 0.8,
	}
	// max tokens: https://platform.openai.com/docs/models/gpt-4

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("OpenAI-Organization", "org-K1e47v2URc9rSlglooxQReFe")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result ResponseOpenAi
	if err = json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	// fmt.Println("completions response")
	// fmt.Println(string(body))

	if result.Error.Message != "" {
		return "", "", fmt.Errorf("OpenAI API error: %s", result.Error.Message)
	}

	choices := result.Choices
	firstChoice := choices[0]
	text := firstChoice.Text

	return text, prompt, nil
}

type ResponseOpenAi struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    interface{} `json:"code"`
	} `json:"error"`
}

func textToSpeech(text, outputFilename string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := polly.New(sess, aws.NewConfig().WithRegion("eu-central-1"))

	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		SampleRate:   aws.String("22050"),
		Text:         aws.String(text),
		TextType:     aws.String("text"),
		VoiceId:      aws.String("Daniel"), // German voice
		Engine:       aws.String("neural"),
	}

	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		return err
	}

	defer output.AudioStream.Close()

	buf, err := ioutil.ReadAll(output.AudioStream)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFilename, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}
