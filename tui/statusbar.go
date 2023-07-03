package tui

import "github.com/knipferrc/teacup/statusbar"

// StatusBar defines a status bar with Github token status and useful commands
func StatusBar(apiKey string) statusbar.Model {
	var sb statusbar.Model

	if apiKey == "" {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: StatusBarForegroundErrorStyle,
				Background: StatusBarBackgroundErrorStyle,
			},
			DefaultSBColors,
			DefaultSBColors,
			DefaultSBColors,
		)
	} else {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: StatusBarForegroundSuccessStyle,
				Background: StatusBarBackgroundSuccessStyle,
			},
			DefaultSBColors,
			DefaultSBColors,
			DefaultSBColors,
		)
	}

	return sb
}