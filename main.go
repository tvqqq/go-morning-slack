package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type quote struct {
	Text   string `json:"q"`
	Author string `json:"a"`
}

type slackDataSend struct {
	Text string `json:"text"`
}

func main() {
	jsonQuote := getQuote()
	sendSlack(jsonQuote)
}

func getQuote() quote {
	quoteUrl := "https://zenquotes.io/api/today"

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, reqErr := http.NewRequest(http.MethodGet, quoteUrl, nil)
	if reqErr != nil {
		log.Fatal(reqErr)
	}

	res, resErr := httpClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var quotes []quote
	jsonErr := json.Unmarshal(body, &quotes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return quotes[0]
}

func sendSlack(jsonQuote quote) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackUrl := os.Getenv("SLACK_WEBHOOK_URL")
	dataSend := slackDataSend{
		Text: fmt.Sprintf("%s - %s", jsonQuote.Text, jsonQuote.Author),
	}

	jsonDataSend, jsonErr := json.Marshal(dataSend)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	byteDataSend := bytes.NewBuffer(jsonDataSend)

	req, resErr := http.NewRequest(http.MethodPost, slackUrl, byteDataSend)
	if resErr != nil {
		log.Fatal(resErr)
	}

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	_, resErr = httpClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
}
