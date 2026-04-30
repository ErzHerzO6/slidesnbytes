package server

import (
  "fmt"

  tea "charm.land/bubbletea/v2"
  "github.com/charmbracelet/ssh"
  "charm.land/wish/v2"
  bm "charm.land/wish/v2/bubbletea"
  "github.com/charmbracelet/colorprofile"
)

func slidesMiddleware(srv *Server) wish.Middleware {
  newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
    p := tea.NewProgram(m, opts...)
    return p
  }
  teaHandler := func(s ssh.Session) *tea.Program {
    _, _, active := s.Pty()
    if !active {
      fmt.Println("no active terminal, skipping")
      err := s.Exit(1)
      if err != nil {
        fmt.Println("Error exiting session")
      }
      return nil
    }
    return newProg(srv.presentation, tea.WithInput(s), tea.WithOutput(s), tea.WithColorProfile(colorprofile.ANSI256))
  }
  return bm.MiddlewareWithProgramHandler(teaHandler)
}
