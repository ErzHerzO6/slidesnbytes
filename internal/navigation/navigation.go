package navigation

import (
	"slices"
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
		return State{
			Page:             navigateNext(state),
			CurrentSlide:     calculateNextSlide(state),
			TotalSlides:      state.TotalSlides,
			SlidesWithBreaks: state.SlidesWithBreaks,
		}
	case "up", "k", "left", "h", "p", "pgup", "N":
		return State{
			Page:             navigatePrevious(state),
			CurrentSlide:     calculatePrevSlide(state),
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

func calculateNextSlide(state State) int {
	nextPage := navigateNext(state)
	currentSlide := state.CurrentSlide
	// For each page advanced, only increment CurrentSlide if that page is not a break slide
	for p := state.Page + 1; p <= nextPage; p++ {
		if !slices.Contains(state.SlidesWithBreaks, p) {
			currentSlide++
		}
	}
	maxSlide := state.TotalSlides - 1 - len(state.SlidesWithBreaks)
	if currentSlide > maxSlide {
		return maxSlide
	}
	return currentSlide
}

func calculatePrevSlide(state State) int {
	prevPage := navigatePrevious(state)
	// Count how many pages we're going back
	currentSlide := state.CurrentSlide
	for p := state.Page; p > prevPage; p-- {
		if !slices.Contains(state.SlidesWithBreaks, p) {
			currentSlide--
		}
	}
	if currentSlide < 0 {
		return 0
	}
	return currentSlide
}

// calculateCurrentSlideForPage calculates the logical slide number for a given page,
// accounting for break slides. Used when jumping directly to a page (G, gg).
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
