package navigation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavigation(t *testing.T) {
	tests := []struct {
		keys   string
		target int
	}{
		{target: 0},
		{keys: "l", target: 1},
		{keys: "jjjjjjjjjj", target: 10},
		{keys: "jjjjjjjjjjjjj", target: 10},
		{keys: "G", target: 10},
		{keys: "llgg", target: 0},
		{keys: "2j", target: 2},
		{keys: "0j", target: 1},
		{keys: "-11G", target: 10},
		{keys: "0G", target: 0},
		{keys: "3G", target: 2},
		{keys: "11G", target: 10},
		{keys: "101G", target: 10},
		{keys: "nnN", target: 1},
	}

	for _, tt := range tests {
		t.Run(tt.keys, func(t *testing.T) {
			currentState := State{
				Buffer:      "",
				Page:        0,
				TotalSlides: 11,
			}

			for _, key := range strings.Split(tt.keys, "") {
				currentState = Navigate(currentState, key)
			}

			// Without any SlidesWithBreaks, CurrentSlide should always equal Page
			targetState := State{Page: tt.target, CurrentSlide: tt.target, TotalSlides: 11}
			assert.Equal(t, targetState, currentState)
		})
	}
}

func TestNavigationWithBreaks(t *testing.T) {
	// Simulate the real slides.md structure:
	// Page 0: Slide 0 (Welcome)
	// Page 1: Slide 1 (Everything is markdown)
	// Page 2: Slide 2 (h1 h2 h3)
	// Page 3: Slide 3 (Markdown components - first part)
	// Page 4: Slide 3 break (+ bullet 1)   <- break
	// Page 5: Slide 3 break (+ bullet 2)   <- break
	// Page 6: Slide 3 break (+ numbered)   <- break
	// Page 7: Slide 4 (Tables)
	// Page 8: Slide 5 (Graphs)
	// Page 9: Slide 6 (separator explanation)
	slidesWithBreaks := []int{4, 5, 6}
	totalSlides := 10

	tests := []struct {
		name          string
		keys          string
		startPage     int
		startSlide    int
		expectedPage  int
		expectedSlide int
	}{
		// Forward navigation
		{
			name:          "forward normal slide",
			keys:          "l",
			startPage:     0,
			startSlide:    0,
			expectedPage:  1,
			expectedSlide: 1,
		},
		{
			name:          "forward into first break",
			keys:          "l",
			startPage:     3,
			startSlide:    3,
			expectedPage:  4,
			expectedSlide: 3, // still logical slide 3
		},
		{
			name:          "forward through breaks",
			keys:          "l",
			startPage:     4,
			startSlide:    3,
			expectedPage:  5,
			expectedSlide: 3, // still logical slide 3
		},
		{
			name:          "forward from last break to next real slide",
			keys:          "l",
			startPage:     6,
			startSlide:    3,
			expectedPage:  7,
			expectedSlide: 4, // new logical slide
		},
		// Backward navigation
		{
			name:          "backward from real slide into last break",
			keys:          "h",
			startPage:     7,
			startSlide:    4,
			expectedPage:  6,
			expectedSlide: 3, // back to logical slide 3
		},
		{
			name:          "backward through breaks",
			keys:          "h",
			startPage:     6,
			startSlide:    3,
			expectedPage:  5,
			expectedSlide: 3, // still logical slide 3
		},
		{
			name:          "backward from first break to slide start",
			keys:          "h",
			startPage:     4,
			startSlide:    3,
			expectedPage:  3,
			expectedSlide: 3, // still logical slide 3 (this is its first page)
		},
		{
			name:          "backward from slide start to previous slide",
			keys:          "h",
			startPage:     3,
			startSlide:    3,
			expectedPage:  2,
			expectedSlide: 2, // now on logical slide 2
		},
		// Multiple backward steps through all breaks
		{
			name:          "backward multiple steps through all breaks",
			keys:          "hhhh",
			startPage:     7,
			startSlide:    4,
			expectedPage:  3,
			expectedSlide: 3, // still logical slide 3 (the head of the break group)
		},
		{
			name:          "backward all the way from after breaks",
			keys:          "hhhhh",
			startPage:     7,
			startSlide:    4,
			expectedPage:  2,
			expectedSlide: 2, // past the break group to previous slide
		},
		// Jump navigation
		{
			name:          "jump to end",
			keys:          "G",
			startPage:     0,
			startSlide:    0,
			expectedPage:  9,
			expectedSlide: 6,
		},
		{
			name:          "jump to start",
			keys:          "gg",
			startPage:     7,
			startSlide:    4,
			expectedPage:  0,
			expectedSlide: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentState := State{
				Buffer:           "",
				Page:             tt.startPage,
				CurrentSlide:     tt.startSlide,
				TotalSlides:      totalSlides,
				SlidesWithBreaks: slidesWithBreaks,
			}

			for _, key := range strings.Split(tt.keys, "") {
				currentState = Navigate(currentState, key)
			}

			assert.Equal(t, tt.expectedPage, currentState.Page, "Page mismatch")
			assert.Equal(t, tt.expectedSlide, currentState.CurrentSlide, "CurrentSlide mismatch")
		})
	}
}
