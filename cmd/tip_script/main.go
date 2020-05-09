package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
)

type response struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Label string `json:"label,omitempty"`
}

type handler interface {
	regex() *regexp.Regexp
	handle(input string) ([]response, error)
}

var handlers = []handler{&cidrHandler{}, &epochHandler{}}

func main() {
	responses, err := handleInput(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(responses)
	if err != nil {
		panic(err)
	}
}

func handleInput(input string) ([]response, error) {
	var responses []response
	for _, h := range handlers {
		if h.regex().MatchString(input) {
			response, err := h.handle(input)
			if err != nil {
				log.Printf("Failed: %v", err)
			} else {
				responses = append(responses, response...)
			}
		}
	}

	return responses, nil
}

func text(text string) response {
	return response{
		Type:  "text",
		Value: text,
	}
}

func link(text string, url string) response {
	return response{
		Type:  "url",
		Value: url,
		Label: text,
	}
}
