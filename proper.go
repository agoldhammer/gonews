package main

import (
	"fmt"
	"regexp"
)

type properNounCounterType map[string]int

// given text of 1 status, find all capped words (excluding caps in URLs)
// return: pointer to a slice containing capped string
func (pnc properNounCounterType) matcher(text *string) *[]string {
	// blank out urls
	reurl := regexp.MustCompile(`\bhttp[s]{0,1}://\S*\b`)
	deurled := reurl.ReplaceAllString(*text, "")
	// find capped words
	re := regexp.MustCompile(`(\b[A-Z]+\S{3,}\b)`)
	matches := re.FindAllString(deurled, 10)
	// fmt.Printf("matches: %v\n", matches)
	return &matches
}

func (pnc properNounCounterType) add(matches *[]string) {
	for _, pnoun := range *matches {
		if n, ok := pnc[pnoun]; ok {
			// fmt.Println("incing", pnoun)
			pnc[pnoun] = n + 1
		} else {
			// fmt.Println("adding", pnoun)
			pnc[pnoun] = 1
		}
	}

}

func (pnc properNounCounterType) print(dth int) {
	for pnoun, count := range pnc {
		if count >= dth {

			fmt.Printf("%s, count: %v\n", pnoun, count)
		}
	}
}
