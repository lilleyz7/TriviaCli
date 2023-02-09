package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var client *http.Client
var score int
var round = 1

type ApiResponse struct {
	CorrectAnswer    string
	IncorrectAnswers []string
	Question         string
}

func GetQuestion(url string) ApiResponse {
	var response []ApiResponse
	err := GetJSON(url, &response)
	if err != nil {
		panic(err)
	} else {
		return response[0]
	}
}

func GetJSON(url string, target interface{}) error {
	res, err := client.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}

func SetupGame() []string {
	categories := []string{"arts_and_literature", "film_and_tv", "food_and_drink", "general_knowledge", "geography", "history", "music", "science", "society_and_culture", "sports_and_leisure"}
	var randomNum [3]int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		randomNum[i] = rand.Intn(len(categories))
	}
	return categories
}

func GetUserCategory(c []string) string {
	var userSelectedCat string
	fmt.Print("Select one of the following categories by typing the name" + " -> " + c[0] + " - " + c[1] + " - " + c[2] + "\n\n")
	fmt.Println("Your selection: ")
	fmt.Scanln(&userSelectedCat)
	for _, s := range c {
		if userSelectedCat == s {
			return userSelectedCat
		}
	}
	fmt.Println("Please enter a valid name \n\n")
	return GetUserCategory(c)
}

func GetUserDiff(d []string) string {
	var userSelectedDiff string
	fmt.Print("Select one of the following categories by typing the name" + " -> " + d[0] + " - " + d[1] + " - " + d[2] + "\n")
	fmt.Println("Your selection: ")
	fmt.Scanln(&userSelectedDiff)
	for _, s := range d {
		if userSelectedDiff == s {
			return userSelectedDiff
		}
	}
	fmt.Println("Please enter a valid name")
	return GetUserCategory(d)
}

func IncreaseScore(diff string) {
	if diff == "easy" {
		score += 1
	}
	if diff == "medium" {
		score += 2
	} else {
		score += 3
	}
}
func RunGame(c []string, d []string) {
	fmt.Println("Score:", score)
	fmt.Println("Current round: ", round)
	categoryToFetch := GetUserCategory(c)
	diffToFetch := GetUserDiff(d)
	var userResponse int

	url := ("https://the-trivia-api.com/api/questions?categories=" + categoryToFetch + "limit=1&difficulty=" + diffToFetch)
	nextQuestion := GetQuestion(url)
	correct := nextQuestion.CorrectAnswer
	answersArray := append(nextQuestion.IncorrectAnswers, nextQuestion.CorrectAnswer)

	// randomize answers
	rand.Shuffle(len(answersArray), func(i, j int) {
		answersArray[i], answersArray[j] = answersArray[j], answersArray[i]
	})

	fmt.Println(nextQuestion.Question)
	for ind, s := range answersArray {
		fmt.Println(ind, "-", s)
	}
	fmt.Println("Your selection: ")
	fmt.Scanln(&userResponse)
	if answersArray[userResponse] == correct {
		fmt.Println("\n\n Correct!")
		IncreaseScore(diffToFetch)
		return
	} else {
		fmt.Println("Incorrect Response")
		fmt.Println("The correct answer was", correct, "\n\n")
		return
	}

}

func main() {
	client = &http.Client{Timeout: 10 * time.Second}
	difficulties := []string{"easy", "medium", "hard"}
	maxRounds := 5
	cat := SetupGame()
	for round < maxRounds {
		RunGame(cat, difficulties)
		round++
	}
	fmt.Println("Thanks for playing! Your score was", score)
	os.Exit(0)
}
