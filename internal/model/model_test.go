package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_NoBreaks(t *testing.T) {
	content := "Slide 1 Hello\n---\nSlide 2 World\n---\nSlide 3 End"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 3)
	assert.Empty(t, m.SlidesWithBreaks)
	assert.Contains(t, m.Slides[0], "Slide 1")
	assert.Contains(t, m.Slides[1], "Slide 2")
	assert.Contains(t, m.Slides[2], "Slide 3")
}

func TestLoad_WithBreaks(t *testing.T) {
	// Note: "# Slide 1" is parsed as metadata frontmatter by meta.Parse if it
	// starts the file. Use content that won't be parsed as YAML metadata.
	content := "Slide 1 content\n---\nSlide 2 intro\n<!-- #break -->\n* Point A\n<!-- #break -->\n* Point B\n---\nSlide 3 content"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	// Slide 1 = 1 page, Slide 2 = 3 pages (1 + 2 breaks), Slide 3 = 1 page = 5 total
	assert.Len(t, m.Slides, 5)
	// Break pages are at indices 2 and 3
	assert.Equal(t, []int{2, 3}, m.SlidesWithBreaks)
}

func TestLoad_BreakCreatesIncrementalContent(t *testing.T) {
	content := "Part A\n<!-- #break -->\nPart B\n<!-- #break -->\nPart C"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 3)
	// First page: only Part A
	assert.Equal(t, "Part A", m.Slides[0])
	// Second page: Part A + Part B (joined with \n since breakDelimiter is removed)
	assert.Equal(t, "Part A\nPart B", m.Slides[1])
	// Third page: all parts
	assert.Equal(t, "Part A\nPart B\nPart C", m.Slides[2])
}

func TestLoad_BreakDoesNotConflictWithSlideDelimiter(t *testing.T) {
	content := "Slide 1 content\n---\nSlide 2 content\n---\nSlide 3 content"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 3)
	assert.Empty(t, m.SlidesWithBreaks)
}

func TestLoad_OnlyBreakDelimiter(t *testing.T) {
	// A single slide with only break delimiters (edge case)
	content := "First\n<!-- #break -->\nSecond\n<!-- #break -->\nThird"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 3)
	assert.Equal(t, []int{1, 2}, m.SlidesWithBreaks)
}

func TestLoad_EmptyBreakSection(t *testing.T) {
	content := "Before\n<!-- #break -->\n\n<!-- #break -->\nAfter"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 3)
	// Second slide should be "Before\n" (the empty section adds just a newline)
	assert.Equal(t, "Before\n", m.Slides[1])
}

func TestLoad_WithMetadataAndBreaks(t *testing.T) {
	content := "---\ntheme: dark\nauthor: Test\n---\n# Slide 1\n<!-- #break -->\n* Point\n---\n# Slide 2"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	// Metadata slide is removed, Slide 1 becomes 2 pages, Slide 2 becomes 1 page = 3
	assert.Len(t, m.Slides, 3)
	assert.Equal(t, "Test", m.Author)
	assert.Contains(t, m.Slides[0], "Slide 1")
}

func TestLoad_InvalidBreakSyntaxIgnored(t *testing.T) {
	// <--break--> is NOT a valid break delimiter
	content := "# Slide 1\n<--break-->\n* Point A\n---\n# Slide 2"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 2)
	assert.Empty(t, m.SlidesWithBreaks)
	// The <--break--> stays as part of the slide content
	assert.Contains(t, m.Slides[0], "<--break-->")
}

func TestLoad_BreakWithoutHashIgnored(t *testing.T) {
	// <!-- break --> (without #) is NOT a valid break delimiter
	content := "Part A\n<!-- break -->\nPart B"
	path := writeTempFile(t, content)

	m := Model{FileName: path}
	err := m.Load()

	require.NoError(t, err)
	assert.Len(t, m.Slides, 1)
	assert.Empty(t, m.SlidesWithBreaks)
}

func TestPaging_WithBreaks(t *testing.T) {
	m := Model{
		Slides:           []string{"a", "ab", "abc", "d", "e"},
		SlidesWithBreaks: []int{1, 2},
		Paging:           "Slide %d / %d",
		Page:             0,
		CurrentSlide:     0,
	}

	// On first page: "Slide 1 / 3" (5 total pages - 2 breaks = 3 logical slides)
	result := m.paging()
	assert.Equal(t, "Slide 1 / 3", result)

	// On a break page: still shows same slide number
	m.Page = 2
	m.CurrentSlide = 0 // still the same logical slide
	result = m.paging()
	assert.Equal(t, "Slide 1 / 3", result)

	// After break group: shows next slide number
	m.Page = 3
	m.CurrentSlide = 1
	result = m.paging()
	assert.Equal(t, "Slide 2 / 3", result)
}

func TestPaging_WithoutBreaks(t *testing.T) {
	m := Model{
		Slides:           []string{"a", "b", "c"},
		SlidesWithBreaks: nil,
		Paging:           "Slide %d / %d",
		Page:             1,
		CurrentSlide:     1,
	}

	result := m.paging()
	assert.Equal(t, "Slide 2 / 3", result)
}

func TestSetPage_ResetsVirtualText(t *testing.T) {
	m := Model{
		Slides:      []string{"a", "b", "c"},
		Page:        0,
		VirtualText: "some output",
	}

	m.SetPage(1)
	assert.Equal(t, 1, m.Page)
	assert.Empty(t, m.VirtualText)
}

func TestSetPage_SamePageNoChange(t *testing.T) {
	m := Model{
		Slides:      []string{"a", "b", "c"},
		Page:        1,
		VirtualText: "keep this",
	}

	m.SetPage(1)
	assert.Equal(t, 1, m.Page)
	assert.Equal(t, "keep this", m.VirtualText) // not reset because page didn't change
}

func TestSetCurrentSlide_ResetsVirtualText(t *testing.T) {
	m := Model{
		Slides:       []string{"a", "b", "c"},
		CurrentSlide: 0,
		VirtualText:  "some output",
	}

	m.SetCurrentSlide(1)
	assert.Equal(t, 1, m.CurrentSlide)
	assert.Empty(t, m.VirtualText)
}

// writeTempFile creates a temporary markdown file for testing and returns its path.
func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}
