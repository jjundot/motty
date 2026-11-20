package workbench

import (
	"github.com/vivaui/viva/app"
)

// A translucent window, in the manner of Xshell's transparency slider: the
// window takes the system's acrylic backdrop material and the terminal
// grounds are painted at the chosen strength over it, so the desktop shows
// through the terminals while the window's own chrome — its bar, sidebar,
// pages and dialogs — stays opaque and readable. The framework has no
// whole-window alpha of its own; the backdrop material is the platform's one
// translucent surface, and it exists on Windows (11 22H2+ in full, 10 1809+
// through the accent policy). Where the system reports no material — another
// platform, high contrast, transparency effects off — the setting changes
// nothing on screen: the grounds keep their full strength, because
// winBackdropActive never turns on.

// windowOpacity is how much of the window's ground stays: the alpha of the
// terminal surfaces over the backdrop.
type windowOpacity uint8

const (
	winOpaque windowOpacity = iota
	winOpacityLight
	winOpacityMedium
	winOpacityStrong
)

// factor is the ground's alpha at this level: what fraction of the terminal's
// own colour covers the backdrop.
func (opacity windowOpacity) factor() float32 {
	switch opacity {
	case winOpacityLight:
		return 0.85
	case winOpacityMedium:
		return 0.70
	case winOpacityStrong:
		return 0.55
	}
	return 1
}

func (opacity windowOpacity) saved() string {
	switch opacity {
	case winOpacityLight:
		return "light"
	case winOpacityMedium:
		return "medium"
	case winOpacityStrong:
		return "strong"
	}
	return ""
}

func windowOpacityFromString(value string) windowOpacity {
	switch value {
	case "light":
		return winOpacityLight
	case "medium":
		return winOpacityMedium
	case "strong":
		return winOpacityStrong
	}
	return winOpaque
}

// windowBackdropStyle is the material the window asks for, read by the
// framework after each paint: acrylic at the chosen setting, none at opaque.
// The callback is installed even when the setting is opaque, because a
// non-nil callback is what opts the window into composition at creation;
// without it a later change of the setting could not take effect until the
// window was rebuilt.
func (state *State) windowBackdropStyle() app.BackdropStyle {
	if state.winOpacity == winOpaque {
		return app.BackdropStyle{Kind: app.BackdropNone}
	}
	return app.BackdropStyle{Kind: app.BackdropAcrylic, Dark: state.pal != nil && state.pal.dark}
}

// setWindowBackdropActive records whether the system is really compositing a
// material behind the window. Only then do the terminal grounds go
// translucent; the repaint it asks for is the one that shows it.
func (state *State) setWindowBackdropActive(active bool) {
	if state.winBackdropActive == active {
		return
	}
	state.winBackdropActive = active
	state.Wake()
}

// windowOpacityChoices are the levels Settings offers, weakest first after
// opaque.
var windowOpacityChoices = [...]windowOpacity{winOpaque, winOpacityLight, winOpacityMedium, winOpacityStrong}

func (state *State) windowOpacityLabel(opacity windowOpacity) string {
	m := state.messages.Tools
	switch opacity {
	case winOpacityLight:
		return m.WindowOpacityLight
	case winOpacityMedium:
		return m.WindowOpacityMedium
	case winOpacityStrong:
		return m.WindowOpacityStrong
	}
	return m.WindowOpacityOff
}

// paintWindowOpacitySettings declares the transparency choice of the terminal
// settings page, in its own group under the background picture's: the two are
// cousins, one putting a picture behind the terminals and the other the
// desktop itself.
func (state *State) paintWindowOpacitySettings(page *settingsPage) {
	m := state.messages.Tools
	opacity := state.winOpacity
	page.group(m.WindowOpacityGroup, func() {
		settingsChoiceRow(page, m.WindowOpacityRow, m.WindowOpacityDescription, &opacity, windowOpacityChoices[:], state.windowOpacityLabel,
			func() {
				state.winOpacity = opacity
				state.save()
			}, nil)
	})
}
