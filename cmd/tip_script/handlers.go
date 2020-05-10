package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
)

var _ handler = &cidrHandler{}

type cidrHandler struct {
}

func (h *cidrHandler) regex() *regexp.Regexp {
	return regexp.MustCompile("^(\\d{1,3}(\\.\\d{1,3}(\\.\\d{1,3}(\\.\\d{1,3})?)?)?)/(\\d{1,3})$")
}

func (h *cidrHandler) handle(input string) ([]response, error) {
	groups := h.regex().FindStringSubmatch(input)
	missingOctets := 0
	for _, group := range groups {
		if group == "" {
			missingOctets++
		}
	}
	input = fmt.Sprintf("%s%s/%s", groups[1], strings.Repeat(".0", missingOctets), groups[5])

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

var _ handler = &epochHandler{}

type epochHandler struct {
}

func (_ *epochHandler) regex() *regexp.Regexp {
	return regexp.MustCompile("^\\d+$")
}

func (_ *epochHandler) handle(input string) ([]response, error) {
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

var _ handler = &jiraHandler{}

type jiraHandler struct {
	url string
}

func (j *jiraHandler) regex() *regexp.Regexp {
	return regexp.MustCompile("^[A-Z]{3,4}-\\d+$")
}

func (j *jiraHandler) handle(input string) ([]response, error) {
	base, err := url.Parse(j.url)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	base.Path += input

	return []response{link(fmt.Sprintf("%s JIRA", base.Hostname()), base.String())}, nil
}
