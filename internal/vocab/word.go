// Package vocab defines the core vocabulary types and how to load them.
package vocab

import (
	"fmt"
	"strings"
)

// Gender is the grammatical gender of a German noun, represented by its
// definite article.
type Gender string

const (
	Masculine Gender = "der"
	Feminine  Gender = "die"
	Neuter    Gender = "das"
)

// ParseGender converts user or file input ("der", " Die ") into a Gender.
func ParseGender(s string) (Gender, error) {
	switch g := Gender(strings.ToLower(strings.TrimSpace(s))); g {
	case Masculine, Feminine, Neuter:
		return g, nil
	}
	return "", fmt.Errorf("invalid gender %q (want der, die or das)", s)
}

// Word is a single German noun with its article, plural and translation.
type Word struct {
	Lemma       string // dictionary form, e.g. "Tisch"
	Gender      Gender
	Plural      string // empty if the noun has no (common) plural
	Translation string
}

// String implements fmt.Stringer, so fmt.Println(w) prints "der Tisch".
func (w Word) String() string {
	return string(w.Gender) + " " + w.Lemma
}
