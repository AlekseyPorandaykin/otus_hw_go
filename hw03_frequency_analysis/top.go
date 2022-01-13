package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type quantityRepetitionsWord struct {
	word     string
	quantity int
}

func Top10(text string) []string {
	repeatingWords := []quantityRepetitionsWord{}
	result := []string{}

	for word, quantity := range getWordsWithRepetitions(text) {
		repeatingWords = append(repeatingWords, quantityRepetitionsWord{
			word:     word,
			quantity: quantity,
		})
	}

	sort.Slice(repeatingWords, func(indexOne, indexTwo int) bool {
		if repeatingWords[indexOne].quantity == repeatingWords[indexTwo].quantity {
			// Лексикографическая сортировка
			runeOne := []rune(repeatingWords[indexOne].word)
			runeTwo := []rune(repeatingWords[indexTwo].word)
			for index, wordOne := range runeOne {
				if index < len(runeTwo) {
					wordTwo := runeTwo[index]
					// Если первые символы совпадают, то продолжаем дальше проверять
					if wordOne == wordTwo {
						continue
					}

					return wordOne < wordTwo
				}
			}
			return len(runeOne) < len(runeTwo)
		}

		return repeatingWords[indexOne].quantity > repeatingWords[indexTwo].quantity
	})

	for _, repeatingWord := range repeatingWords {
		result = append(result, repeatingWord.word)
	}

	if len(result) > 10 {
		return result[:10]
	}

	return result
}

func getWords(text string) []string {
	result := []string{}
	for _, word := range strings.Fields(text) {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}

func getWordsWithRepetitions(text string) map[string]int {
	mapWords := map[string]int{}
	words := getWords(text)
	for _, word := range words {
		mapWords[word]++
	}

	return mapWords
}
