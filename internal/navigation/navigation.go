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
			Buffer:      newBuffer,
			Page:        state.Page,
			TotalSlides: state.TotalSlides,
		}
	case "g":
		switch state.Buffer {
		case "g":
			return State{
				Page:        0,
				TotalSlides: state.TotalSlides,
			}
		default:
			return State{
				Buffer:      "g",
				Page:        state.Page,
				TotalSlides: state.TotalSlides,
			}
		}
	case "G":
		targetSlide := state.TotalSlides - 1
		if bufferIsNumeric(state.Buffer) {
			targetSlide = navigateSlide(state.Buffer, state.TotalSlides)
		}

		return State{
			Page:        targetSlide,
			TotalSlides: state.TotalSlides,
		}
	case " ", "down", "j", "right", "l", "enter", "n", "pgdown":
		return State{
			Page:         navigateNext(state),
			CurrentSlide: calculateNextSlide(state),
			TotalSlides:  state.TotalSlides,
		}
	case "up", "k", "left", "h", "p", "pgup", "N":
		return State{
			Page:         navigatePrevious(state),
			CurrentSlide: calculatePrevSlide(state),
			TotalSlides:  state.TotalSlides,
		}
	default:
		return State{
			Page:         state.Page,
			CurrentSlide: state.CurrentSlide,
			TotalSlides:  state.TotalSlides,
		}
	}
}

func calculatePrevSlide(state State) int {
	currentSlide := state.CurrentSlide
	indexOfBreakSlide := slices.Index(state.SlidesWithBreaks, state.Page+1)
	if indexOfBreakSlide != -1 {
		currentSlide = state.Page - indexOfBreakSlide
		return currentSlide
	}

	if state.CurrentSlide > 0 {
		return state.CurrentSlide - 1
	}

	return currentSlide
}
func calculateNextSlide(state State) int {
	currentSlide := state.CurrentSlide
	indexOfBreakSlide := slices.Index(state.SlidesWithBreaks, state.Page+1)
	if indexOfBreakSlide != -1 {
		currentSlide = state.Page - indexOfBreakSlide
		return currentSlide
	}

	if (state.CurrentSlide + len(state.SlidesWithBreaks)) < state.TotalSlides-1 {
		return state.CurrentSlide + 1
	}

	return currentSlide
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
