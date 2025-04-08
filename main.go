package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"main/helpers"
)

type word struct {
	Word string `json:"word"`
}

func main() {

	limiter := helpers.NewRateLimiter()

	guessTheWord(limiter)
	phrase1 := "Elbow Pasta"
	fmt.Printf("the first phrase is %s\n", phrase1)

	// clickbuttons errored out for me - I could not get mac to
	// give permissions to chromedp to open my browser
	// you may need to comment this func out to run the entire
	// program
	clickButtons(limiter)

	GuessNumber(1, 1000000000000000, limiter)
	phrase3 := "Taylor Swift"
	fmt.Printf("the third phrase is %s\n", phrase3)

	Puzzle4(limiter)
	phase4 := "The Seine"
	fmt.Printf("phase 4 is %s", phase4)

	// puzzle(s) 5 and 6 not completed

}

func GuessNumber(min int, max int, limiter *helpers.RateLimiter) {
	low := min
	high := max

	for low <= high {
		guess := (low + high) / 2
		result, adjust := maketheguess(limiter, guess)

		if result {

			break
		} else if adjust == "higher" {
			low = guess + 1
		} else {
			high = guess - 1
		}
	}

}

func maketheguess(limiter *helpers.RateLimiter, num int) (result bool, adjust string) {
	w := struct {
		Guess int `json:"guess"`
	}{
		Guess: num,
	}

	jsonBytes, err := json.Marshal(w)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// fmt.Println("guessing number %d", num)
	url := "https://tht.goethena.com/higher-or-lower?apiKey=323b909a-3d84-47be-8c80-8057fce536cd"
	resp := helpers.NewRequest(limiter, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if resp != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error opening body: %v\n", err)
		}
		something := make(map[string]interface{})
		if err := json.Unmarshal(body, &something); err != nil {
			log.Printf("error unmarshaling: %v\n", err)
		}
		adjust = something["message"].(string)
		if something["success"].(bool) == true {
			return true, adjust
		} else {
			return false, adjust
		}

	}
	return false, adjust
}

func Puzzle4(limiter *helpers.RateLimiter) {
	city := "new-york-city"
	for {
		url := "https://tht.goethena.com/carmen-sandiego/" + city + "?apiKey=323b909a-3d84-47be-8c80-8057fce536cd"
		resp := helpers.NewRequest(limiter, http.MethodGet, url, nil)
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error opening body: %v\n", err)
			}
			something := make(map[string]interface{})
			if err := json.Unmarshal(body, &something); err != nil {
				log.Printf("error unmarshaling: %v\n", err)
			}
			success := something["success"].(bool)
			if success {

				break
			}
			comms := something["communication"].(string)
			a := strings.Split(comms, "+++")
			city = a[1]
			fmt.Println(city)
		}
	}
}

func guessTheWord(limiter *helpers.RateLimiter) {
	words := []string{
		"iceberg",
		"zen",
		"night",
		"art",
		"ocean",
		"cloud",
		"horizon",
		"truth",
		"mystery",
		"harmony",
		"umbrella",
		"unity",
		"quasar",
		"yellow",
		"zenith",
		"tree",
		"youth",
		"orange",
		"novel",
		"piano",
		"orbit",
		"vessel",
		"jewel",
		"yarn",
		"mirror",
		"apple",
		"nature",
		"kindness",
		"yacht",
		"guitar",
		"knowledge",
		"xenon",
		"x-ray",
		"energy",
		"garden",
		"dog",
		"progress",
		"wisdom",
		"zebra",
		"jungle",
		"compass",
		"love",
		"sunshine",
		"globe",
		"kitchen",
		"xylophone",
		"lemon",
		"snowflake",
		"beach",
		"house",
		"adventure",
		"window",
		"travel",
		"elephant",
		"rainbow",
		"light",
		"signal",
		"idea",
		"volcano",
		"justice",
		"echo",
		"kite",
		"flower",
		"banana",
		"solace",
		"wind",
		"opportunity",
		"anchor",
		"lantern",
		"dance",
		"museum",
		"underwater",
		"peace",
		"quest",
		"reflection",
		"bridge",
		"zephyr",
		"telescope",
		"island",
		"notebook",
		"forest",
		"serenity",
		"freedom",
		"whale",
		"cherry",
		"violet",
		"pyramid",
		"courage",
		"river",
		"beauty",
		"desert",
		"happiness",
		"emotion",
		"imagination",
		"quilt",
		"mountain",
		"grace",
		"fire",
		"discovery",
		"road",
	}

	for _, val := range words {

		w := word{
			Word: val,
		}

		jsonBytes, err := json.Marshal(w)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		url := "https://tht.goethena.com/check-word?apiKey=323b909a-3d84-47be-8c80-8057fce536cd"
		resp := helpers.NewRequest(limiter, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("error opening body")
				continue
			}
			something := make(map[string]interface{})
			if err := json.Unmarshal(body, &something); err != nil {
				log.Println("error unmarshaling")
				continue
			}
			if something["success"] == "true" {
				fmt.Printf("%+v/n", something)
			}
		}

	}
}

func clickButtons(limiter *helpers.RateLimiter) {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Define the selectors for the buttons (e.g., CSS selectors or XPath)
	buttonSelectors := []string{}

	for i := 1; i < 200; i++ {
		buttonSelectors = append(buttonSelectors, fmt.Sprintf("#button%d", i))
	}

	// Create a list of tasks to click each button sequentially
	var tasks chromedp.Tasks
	for _, selector := range buttonSelectors {
		tasks = append(
			tasks,
			chromedp.WaitVisible(selector, chromedp.ByID), // Ensure the button is visible
			chromedp.Click(selector, chromedp.ByID),       // Click the button
			chromedp.Sleep(5*time.Second/195),             // Optional: Wait after clicking
		)
	}

	// Run the tasks
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://tht.goethena.com/puzzle-2?apiKey=323b909a-3d84-47be-8c80-8057fce536cd"), // Replace with your target URL
		tasks,
	)
	if err != nil {
		log.Fatal(err)
	}
}
