package navigation

import (
	"strconv"
)

type repeatableFunc func(slide, totalSlides int) int

// State tracks the current buffer, page, and total number of slides
type State struct {
	Buffer string
	// Page is keeping track of each slide. It doesn't matter if it's a break slide or not, it just counts the slides.
	// This is used to render the slide, so we can't skip break slides when rendering, but we can skip them when calculating the current slide number.
	// that's what CurrentSlide is for, it keeps track of the current slide number, excluding break slides, so we can use it to display the correct slide number in the UI.
	Page             int
	CurrentSlide     int
	TotalSlides      int
	SlidesWithBreaks []int
}

// Navigate receives the current State and keyPress, and returns the new State.
func Navigate(state State, keyPress string) State {
	switch keyPress {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		newBuffer := keyPress

		if bufferIsNumeric(state.Buffer) {
			newBuffer = state.Buffer + keyPress
		}

		return State{
			Buffer:           newBuffer,
			Page:             state.Page,
			CurrentSlide:     state.CurrentSlide,
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	case "g":
		switch state.Buffer {
		case "g":
			return State{
				Page:             0,
				CurrentSlide:     0,
				TotalSlides:      state.TotalSlides,
				SlidesWithBreaks: state.SlidesWithBreaks,
			}
		default:
			return State{
				Buffer:           "g",
				Page:             state.Page,
				CurrentSlide:     state.CurrentSlide,
				TotalSlides:      state.TotalSlides,
				SlidesWithBreaks: state.SlidesWithBreaks,
			}
		}
	case "G":
		targetPage := state.TotalSlides - 1
		if bufferIsNumeric(state.Buffer) {
			targetPage = navigateSlide(state.Buffer, state.TotalSlides)
		}

		return State{
			Page:             targetPage,
			CurrentSlide:     calculateCurrentSlideForPage(targetPage, state.SlidesWithBreaks),
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	case " ", "down", "j", "right", "l", "enter", "n", "pgdown":
		nextPage := navigateNext(state)
		return State{
			Page:             nextPage,
			CurrentSlide:     calculateCurrentSlideForPage(nextPage, state.SlidesWithBreaks),
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	case "up", "k", "left", "h", "p", "pgup", "N":
		prevPage := navigatePrevious(state)
		return State{
			Page:             prevPage,
			CurrentSlide:     calculateCurrentSlideForPage(prevPage, state.SlidesWithBreaks),
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	default:
		return State{
			Page:             state.Page,
			CurrentSlide:     state.CurrentSlide,
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	}
}

// calculateCurrentSlideForPage calculates the logical slide number for a given page,
// accounting for break slides. Break slides don't count as separate logical slides.
//
// Example: Pages [0, 1, 2, 3, 4, 5, 6, 7] with SlidesWithBreaks=[4, 5, 6]
//   - Page 3 → 3 breaks before: 0 → CurrentSlide = 3
//   - Page 4 → 1 break ≤ 4      → CurrentSlide = 3
//   - Page 5 → 2 breaks ≤ 5     → CurrentSlide = 3
//   - Page 6 → 3 breaks ≤ 6     → CurrentSlide = 3
//   - Page 7 → 3 breaks ≤ 7     → CurrentSlide = 4
func calculateCurrentSlideForPage(page int, slidesWithBreaks []int) int {
	breaksBefore := 0
	for _, breakPage := range slidesWithBreaks {
		if breakPage <= page {
			breaksBefore++
		}
	}
	return page - breaksBefore
}

func bufferIsNumeric(buffer string) bool {
	_, err := strconv.Atoi(buffer)
	return err == nil
}

func navigateNext(state State) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide < totalSlides-1 {
			return slide + 1
		}

		return totalSlides - 1
	}, state)
}

func navigateSlide(buffer string, totalSlides int) int {
	destinationSlide, _ := strconv.Atoi(buffer)
	destinationSlide--

	if destinationSlide > totalSlides-1 {
		return totalSlides - 1
	}

	if destinationSlide < 0 {
		return 0
	}

	return destinationSlide
}

func navigatePrevious(state State) int {
	return repeatableAction(func(slide, totalSlides int) int {
		if slide > 0 {
			return slide - 1
		}

		return slide
	}, state)
}

func repeatableAction(fn repeatableFunc, state State) int {
	if !bufferIsNumeric(state.Buffer) {
		return fn(state.Page, state.TotalSlides)
	}

	repeat, _ := strconv.Atoi(state.Buffer)
	page := state.Page

	if repeat == 0 {
		// This is how behaviour works in Vim, so following principle of least astonishment.
		return fn(state.Page, state.TotalSlides)
	}

	for range repeat {
		page = fn(page, state.TotalSlides)
	}

	return page
}
