// Command quiz is a terminal quiz for German noun articles (der/die/das).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/nimamg/german-flashcards/internal/vocab"
)

func main() {
	path := flag.String("words", "data/seed.csv", "path to the words CSV")
	n := flag.Int("n", 10, "number of questions (0 = all)")
	flag.Parse()

	words, err := vocab.LoadCSV(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	if *n > 0 && *n < len(words) {
		words = words[:*n]
	}

	in := bufio.NewScanner(os.Stdin)
	var missed []vocab.Word

	for i, w := range words {
		prompt := fmt.Sprintf("[%d/%d] ___ %s? ", i+1, len(words), w.Lemma)
		answer, ok := askGender(in, prompt)
		if !ok {
			fmt.Println() // stdin closed (Ctrl-D): stop early
			words = words[:i]
			break
		}

		if answer == w.Gender {
			fmt.Printf("  ✓ %s\n", describe(w))
		} else {
			fmt.Printf("  ✗ it's %s\n", describe(w))
			missed = append(missed, w)
		}
	}

	fmt.Printf("\nScore: %d/%d\n", len(words)-len(missed), len(words))
	if len(missed) > 0 {
		fmt.Println("Review these:")
		for _, w := range missed {
			fmt.Println("  ", describe(w))
		}
	}
}

// askGender prompts until the user enters a valid article. It returns
// ok=false if input ends before a valid answer is given.
func askGender(in *bufio.Scanner, prompt string) (vocab.Gender, bool) {
	for {
		fmt.Print(prompt)
		if !in.Scan() {
			return "", false
		}
		g, err := vocab.ParseGender(in.Text())
		if err == nil {
			return g, true
		}
		fmt.Println("  please type der, die or das")
	}
}

// describe formats a word with its plural and translation,
// e.g. "der Tisch (pl. die Tische) — table".
func describe(w vocab.Word) string {
	s := w.String()
	if w.Plural != "" {
		s += " (pl. die " + w.Plural + ")"
	}
	return s + " — " + w.Translation
}
