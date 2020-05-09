package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
)

type response struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Label string `json:"label,omitempty"`
}

var cidrRegex = regexp.MustCompile("^\\d{1,3}(\\.\\d{1,3}(\\.\\d{1,3}(\\.\\d{1,3})?)?)?/\\d{1,3}$")
var epochNumber = regexp.MustCompile("^\\d+$")

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
	if cidrRegex.MatchString(input) {
		response, err := handleCidr(input)
		if err != nil {
			log.Printf("Failed to parse CIDR: %v", err)
		} else {
			return response, nil
		}
	}

	if epochNumber.MatchString(input) {
		response, err := handleEpoch(input)
		if err != nil {
			log.Printf("Failed to parse epoch number: %v", err)
		} else {
			return response, nil
		}
	}

	return []response{}, nil
}

func handleCidr(input string) ([]response, error) {
	_, c, err := net.ParseCIDR(input)
	if err != nil {
		return nil, err
	}

	prefix, _ := c.Mask.Size()
	start, end := cidr.AddressRange(c)

	responses := []response{text(fmt.Sprintf("%s - %s", start, end))}

	if subnet, exceeded := cidr.PreviousSubnet(c, prefix); !exceeded {
		responses = append(responses, text(fmt.Sprintf("Previous: %s", subnet)))
	}
	if subnet, exceeded := cidr.NextSubnet(c, prefix); !exceeded {
		responses = append(responses, text(fmt.Sprintf("Next: %s", subnet)))
	}

	return responses, nil
}

func handleEpoch(input string) ([]response, error) {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return nil, err
	}

	var t time.Time
	switch {
	case i < 100000000000: // seconds
		t = time.Unix(i, 0)
	case i < 100000000000000: // millis
		t = time.Unix(0, i*1000000)
	default: // nanos
		t = time.Unix(0, i)
	}

	return []response{text(t.UTC().String())}, nil
}

func text(text string) response {
	return response{
		Type:  "text",
		Value: text,
	}
}
