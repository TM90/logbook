package main

import "testing"

func TestEditDistance(t *testing.T) {
	var tests = []struct {
		name     string
		s1       string
		s2       string
		distance int
	}{
		{"Two empty strings should result in a distance of 0", "", "", 0},
		{"abc string and empty string should result in a distance of 3", "abc", "", 3},
		{"empty string and abc string should result in a distance of 3", "", "abc", 3},
		{"abc string and abc string should result in a distance of 0", "abc", "abc", 0},
		{"abc string and abcd string should result in a distance of 1", "abc", "abcd", 1},
		{"abcd string and abc string should result in a distance of 1", "abcd", "abc", 1},
		{"ac string and abc string should result in a distance of 1", "ac", "abc", 1},
		{"abc string and ac string should result in a distance of 1", "abc", "ac", 1},
		{"ab string and abc string should result in a distance of 0", "ab", "abc", 1},
		{"abc string and ab string should result in a distance of 0", "abc", "ab", 1},
		{"a string and abc string should result in a distance of 2", "a", "abc", 2},
		{"abc string and a string should result in a distance of 2", "abc", "a", 2},
		{"horse string and ros string should result in a distance of 0", "horse", "ros", 3},
		{"ros string and horse string should result in a distance of 0", "ros", "horse", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := editDistance(tt.s1, tt.s2)
			if tt.distance != distance {
				t.Errorf("got %d expected %d!", distance, tt.distance)
			}
		})
	}
}
