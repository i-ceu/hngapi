package models

import (
	"encoding/json"
	"time"
)

type String struct {
	ID                    string    `json:"id"`
	Value                 string    `json:"value"`
	Length                int       `json:"length"`
	IsPalindrome          bool      `json:"is_palindrome"`
	UniqueChars           int       `json:"unique_characters"`
	WordCount             int       `json:"word_count"`
	CharacterFrequencyMap string    `json:"-"`
	CreatedAt             time.Time `json:"created_at"`
}

func (s String) MarshalJSON() ([]byte, error) {
	type Alias String

	var freqMap map[string]int
	if err := json.Unmarshal([]byte(s.CharacterFrequencyMap), &freqMap); err != nil {
		freqMap = make(map[string]int)
	}

	return json.Marshal(&struct {
		CharacterFrequencyMap map[string]int `json:"character_frequency_map"`
		CreatedAt             string         `json:"created_at"`
		*Alias
	}{
		CharacterFrequencyMap: freqMap,
		CreatedAt:             s.CreatedAt.UTC().Format(time.RFC3339),
		Alias:                 (*Alias)(&s),
	})
}
