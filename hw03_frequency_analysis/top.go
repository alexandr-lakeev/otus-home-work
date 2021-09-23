package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`[,.!?:;'"]`)

type item struct {
	word  string
	count int
}

func Top10(text string) []string {
	counter := callculateWordsFrequency(text)
	items := createItemsFromCounter(counter)
	sortItems(items)
	return getTopWordsFromItems(items, 10)
}

func callculateWordsFrequency(text string) map[string]int {
	counter := make(map[string]int)
	splitedText := strings.Fields(text)

	for _, word := range splitedText {
		if word == "-" {
			continue
		}

		word = strings.ToLower(reg.ReplaceAllString(word, ""))
		counter[word]++
	}

	return counter
}

func createItemsFromCounter(counter map[string]int) []item {
	items := make([]item, len(counter))

	i := 0
	for word, count := range counter {
		items[i] = item{word: word, count: count}
		i++
	}

	return items
}

func sortItems(items []item) {
	sort.Slice(items, func(i, j int) bool {
		count1, count2 := items[i].count, items[j].count
		word1, word2 := items[i].word, items[j].word

		return count1 > count2 || (count1 == count2 && word1 < word2)
	})
}

func getTopWordsFromItems(items []item, top int) []string {
	if top > len(items) {
		top = len(items)
	}

	words := make([]string, top)

	for key, item := range items {
		if key == top {
			break
		}
		words[key] = item.word
	}

	return words
}
