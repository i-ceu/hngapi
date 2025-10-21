package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"profile-api/internal/config"
	"profile-api/internal/models"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddStrings(c *gin.Context) {
	var req struct {
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, gin.H{"error": "Invalid data type for value (must be string)"})
		return
	}
	if strings.TrimSpace(req.Value) == "" {
		c.JSON(400, gin.H{"error": "Invalid request body or missing value field"})
		return
	}

	var exists bool
	config.DB.Model(&models.String{}).
		Select("count(*) > 0").
		Where("value = ?", req.Value).
		Find(&exists)
	if exists {
		c.JSON(409, gin.H{"error": "String already exists in system"})
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(req.Value))
	hash := hex.EncodeToString(hasher.Sum(nil))

	isPalindrome := isPalindrome(req.Value)

	uniqueChars := countUniqueCharacters(req.Value)

	wordCount := len(strings.Fields(req.Value))

	freqMap := getCharacterFrequency(req.Value)
	freqMapJSON, _ := json.Marshal(freqMap)

	storeValue := models.String{
		ID:                    hash,
		Value:                 req.Value,
		Length:                len(req.Value),
		IsPalindrome:          isPalindrome,
		UniqueChars:           uniqueChars,
		WordCount:             wordCount,
		CharacterFrequencyMap: string(freqMapJSON),
		CreatedAt:             time.Now(),
	}

	test := config.DB.Create(&storeValue)

	if test.Error != nil {
		c.JSON(400, test.Error)
	}

	c.JSON(201, gin.H{
		"id":    storeValue.ID,
		"value": storeValue.Value,
		"properties": gin.H{
			"length":                  storeValue.Length,
			"is_palindrome":           storeValue.IsPalindrome,
			"unique_characters":       storeValue.UniqueChars,
			"word_count":              storeValue.WordCount,
			"sha256_hash":             storeValue.ID,
			"character_frequency_map": freqMap,
		},
		"created_at": storeValue.CreatedAt.Format(time.RFC3339),
	})
}

func GetString(c *gin.Context) {
	stringValue := c.Param("string_value")

	var exists models.String
	if err := config.DB.Where("value = ?", stringValue).First(&exists).Error; err != nil {
		c.JSON(404, gin.H{"error": "String does not exist in the system"})
		return
	}

	var freqMap map[string]int
	_ = json.Unmarshal([]byte(exists.CharacterFrequencyMap), &freqMap)

	c.JSON(200, gin.H{
		"id":    exists.ID,
		"value": exists.Value,
		"properties": gin.H{
			"length":                  exists.Length,
			"is_palindrome":           exists.IsPalindrome,
			"unique_characters":       exists.UniqueChars,
			"word_count":              exists.WordCount,
			"sha256_hash":             exists.ID,
			"character_frequency_map": freqMap,
		},
		"created_at": exists.CreatedAt.Format(time.RFC3339),
	})
}

func GetAllStrings(c *gin.Context) {
	query := config.DB.Model(&models.String{})
	filtersApplied := make(map[string]interface{})

	// Apply filters
	if isPalindromeStr := c.Query("is_palindrome"); isPalindromeStr != "" {
		isPalindrome, err := strconv.ParseBool(isPalindromeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid is_palindrome value"})
			return
		}
		query = query.Where("is_palindrome = ?", isPalindrome)
		filtersApplied["is_palindrome"] = isPalindrome
	}

	if minLengthStr := c.Query("min_length"); minLengthStr != "" && reflect.TypeOf(minLengthStr) == reflect.TypeOf(0) {
		minLength, err := strconv.Atoi(minLengthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_length value"})
			return
		}
		query = query.Where("length >= ?", minLength)
		filtersApplied["min_length"] = minLength
	}

	if maxLengthStr := c.Query("max_length"); maxLengthStr != "" && reflect.TypeOf(maxLengthStr) == reflect.TypeOf(0) {
		maxLength, err := strconv.Atoi(maxLengthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_length value"})
			return
		}
		query = query.Where("length <= ?", maxLength)
		filtersApplied["max_length"] = maxLength
	}

	if wordCountStr := c.Query("word_count"); wordCountStr != "" {
		wordCount, err := strconv.Atoi(wordCountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word_count value"})
			return
		}
		query = query.Where("word_count = ?", wordCount)
		filtersApplied["word_count"] = wordCount
	}

	if containsChar := c.Query("contains_character"); containsChar != "" {
		if len(containsChar) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contains_character must be a single character"})
			return
		}
		query = query.Where("value LIKE ?", "%"+containsChar+"%")
		filtersApplied["contains_character"] = containsChar
	}

	var results []models.String
	if err := query.Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch strings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":            results,
		"count":           len(results),
		"filters_applied": filtersApplied,
	})
}

func FilterByNaturalLanguage(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
		return
	}

	filters, err := parseNaturalLanguageQuery(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build database query
	dbQuery := config.DB.Model(&models.String{})
	for key, value := range filters {
		switch key {
		case "is_palindrome":
			dbQuery = dbQuery.Where("is_palindrome = ?", value)
		case "word_count":
			dbQuery = dbQuery.Where("word_count = ?", value)
		case "min_length":
			dbQuery = dbQuery.Where("length >= ?", value)
		case "max_length":
			dbQuery = dbQuery.Where("length <= ?", value)
		case "contains_character":
			dbQuery = dbQuery.Where("value LIKE ?", "%"+value.(string)+"%")
		}
	}

	var results []models.String
	if err := dbQuery.Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch strings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"count": len(results),
		"interpreted_query": gin.H{
			"original":       query,
			"parsed_filters": filters,
		},
	})
}

func DeleteString(c *gin.Context) {
	stringValue := c.Param("string_value")

	result := config.DB.Where("value = ?", stringValue).Delete(&models.String{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete string"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "String does not exist in the system"})
		return
	}

	c.Status(http.StatusNoContent)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

func countUniqueCharacters(s string) int {
	charSet := make(map[rune]bool)
	for _, char := range s {
		charSet[char] = true
	}
	return len(charSet)
}

func getCharacterFrequency(s string) map[string]int {
	freqMap := make(map[string]int)
	for _, char := range s {
		freqMap[string(char)]++
	}
	return freqMap
}

func parseNaturalLanguageQuery(query string) (map[string]interface{}, error) {
	filters := make(map[string]interface{})
	lowerQuery := strings.ToLower(query)

	// Check for palindrome
	if strings.Contains(lowerQuery, "palindrom") {
		filters["is_palindrome"] = true
	}

	// Check for single word
	if strings.Contains(lowerQuery, "single word") {
		filters["word_count"] = 1
	}

	// Check for word count patterns
	if strings.Contains(lowerQuery, "two word") || strings.Contains(lowerQuery, "2 word") {
		filters["word_count"] = 2
	}
	if strings.Contains(lowerQuery, "three word") || strings.Contains(lowerQuery, "3 word") {
		filters["word_count"] = 3
	}

	// Check for length patterns
	if strings.Contains(lowerQuery, "longer than") {
		words := strings.Fields(lowerQuery)
		for i, word := range words {
			if word == "than" && i+1 < len(words) {
				if num, err := strconv.Atoi(words[i+1]); err == nil {
					filters["min_length"] = num + 1
				}
			}
		}
	}

	if strings.Contains(lowerQuery, "shorter than") {
		words := strings.Fields(lowerQuery)
		for i, word := range words {
			if word == "than" && i+1 < len(words) {
				if num, err := strconv.Atoi(words[i+1]); err == nil {
					filters["max_length"] = num - 1
				}
			}
		}
	}

	// Check for contains character
	if strings.Contains(lowerQuery, "containing the letter") || strings.Contains(lowerQuery, "contain the letter") {
		words := strings.Fields(lowerQuery)
		for i, word := range words {
			if word == "letter" && i+1 < len(words) {
				char := strings.TrimSpace(words[i+1])
				if len(char) > 0 {
					filters["contains_character"] = string(char[0])
				}
			}
		}
	}

	// Check for first vowel
	if strings.Contains(lowerQuery, "first vowel") {
		filters["contains_character"] = "a"
	}

	if len(filters) == 0 {
		return nil, fmt.Errorf("unable to parse natural language query")
	}

	return filters, nil
}
