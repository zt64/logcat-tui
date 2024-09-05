package main

import "github.com/zt64/logcat-tui/logcat"

type config struct {
	scrollback int
	format     string

	// Profiles are a list of profiles to display
	profiles []profile
}

// Profile is a struct that contains the name of the profile and filters
type profile struct {
	name    string
	buffers []logcat.Buffer
}
