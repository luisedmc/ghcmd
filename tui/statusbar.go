package tui

import "github.com/knipferrc/teacup/statusbar"

// StatusBar defines a status bar with Github token status and useful commands
func StatusBar(apiKey string, statusText string, status bool) statusbar.Model {
	var sb statusbar.Model

	if apiKey == "" {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: StatusBarForegroundErrorStyle,
				Background: StatusBarBackgroundErrorStyle,
			},
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
		)
	} else {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: StatusBarForegroundSuccessStyle,
				Background: StatusBarBackgroundSuccessStyle,
			},
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
		)
	}

	return sb
}
