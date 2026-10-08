package vocab

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadCSV reads words from a CSV file with the header
// lemma,gender,plural,translation.
func LoadCSV(path string) ([]Word, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseCSV(f)
}

// ParseCSV reads words from any io.Reader (a file, a string, an HTTP body...).
// The first row is treated as a header and skipped.
func ParseCSV(r io.Reader) ([]Word, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 4

	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading csv: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("csv has no words")
	}

	words := make([]Word, 0, len(records)-1)
	for i, rec := range records[1:] {
		gender, err := ParseGender(rec[1])
		if err != nil {
			// i+2: one for the header, one because lines are 1-based.
			return nil, fmt.Errorf("line %d: %w", i+2, err)
		}
		words = append(words, Word{
			Lemma:       strings.TrimSpace(rec[0]),
			Gender:      gender,
			Plural:      strings.TrimSpace(rec[2]),
			Translation: strings.TrimSpace(rec[3]),
		})
	}
	return words, nil
}
