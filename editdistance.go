package main

import (
	"math"
)

func editDistance(word1, word2 string) int {
	cache := make([][]int, len(word1)+1)
	for i := range cache {
		cache[i] = make([]int, len(word2)+1)
		for j := range cache[i] {
			cache[i][j] = math.MaxInt
		}
	}

	for j := 0; j < len(word2)+1; j++ {
		cache[len(word1)][j] = len(word2) - j
	}

	for i := 0; i < len(word1)+1; i++ {
		cache[i][len(word2)] = len(word1) - i
	}

	for i := len(word1) - 1; i >= 0; i-- {
		for j := len(word2) - 1; j >= 0; j-- {
			if word1[i] == word2[j] {
				cache[i][j] = cache[i+1][j+1]
			} else {
				cache[i][j] = 1 + min(cache[i+1][j], cache[i][j+1], cache[i+1][j+1])
			}
		}
	}

	return cache[0][0]
}
