package styles_test

import (
  "testing"

  "charm.land/glamour/v2"
  "charm.land/glamour/v2/ansi"
  charmStyles "charm.land/glamour/v2/styles"
  "github.com/stretchr/testify/assert"

  "github.com/maaslalani/slides/styles"
)

func TestSelectTheme(t *testing.T) {
  tests := []struct {
    name    string
    theme   string
    want    ansi.StyleConfig
    wantErr bool
  }{
    {name: "Select dark theme", theme: "dark", want: charmStyles.DarkStyleConfig, wantErr: false},
    {name: "Select light theme", theme: "light", want: charmStyles.LightStyleConfig, wantErr: false},
    {name: "Select ascii theme", theme: "ascii", want: charmStyles.ASCIIStyleConfig, wantErr: false},
    {name: "Select notty theme", theme: "notty", want: charmStyles.NoTTYStyleConfig, wantErr: false},
    {name: "Select theme with error", theme: "notty", want: charmStyles.DarkStyleConfig, wantErr: true},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Execute the theme selection and ensure
      // it returns a non-nil theme
      selectedTheme := styles.SelectTheme(tt.theme)
      assert.NotNil(t, selectedTheme)

      // Initialize renderers to compare output
      gotRenderer, _ := glamour.NewTermRenderer(selectedTheme)
      wantRenderer, _ := glamour.NewTermRenderer(glamour.WithStyles(tt.want))

      // Render a the same string with two different
      // renderers
      gotOutput, _ := gotRenderer.Render(tt.name)
      wantOutput, _ := wantRenderer.Render(tt.name)

      // Inject exception to ensure a style that doesn't match
      // it's associated string
      if tt.wantErr {
        assert.NotEqual(t, wantOutput, gotOutput)
        return
      }

      // Ensure they both match
      assert.Equal(t, wantOutput, gotOutput)
    })
  }
}

func TestSelectTheme_file(t *testing.T) {
  tests := []struct {
    name       string
    theme      string
    fileExists bool
  }{
    {name: "Select custom theme json", theme: "./theme.json", fileExists: true},
    {name: "Use an invalid filepath", theme: "./someinvalidfile.toml", fileExists: false},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Successfully return a theme if a file exists
      assert.NotNil(t, styles.SelectTheme(tt.theme))

      // Successfully return a theme if a file doesn't exist
      if !tt.fileExists {
        assert.NotNil(t, styles.SelectTheme(tt.theme))
      }
    })
  }
}
