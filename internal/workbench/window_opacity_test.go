package workbench

import (
	"testing"

	"github.com/vivaui/viva/gfx"
)

// TestBackdropDropsTheShellFrameUnderThePanes locks the paint half of the
// backdrop: while the system material is live and the setting is on, the
// shell must not lay its opaque frame colour over the whole window — the
// panes' translucent grounds would only blend with that colour, and the
// desktop behind the window would never show. The window bar keeps its own
// frame colour; only the full-window rect goes.
func TestBackdropDropsTheShellFrameUnderThePanes(t *testing.T) {
	state, _, _ := tabDragWorld(t)
	size := gfx.Size{Width: 1400, Height: 900}
	full := gfx.Rect{Width: size.Width, Height: size.Height}
	frameRects := func(canvas *gfx.Canvas) int {
		count := 0
		for _, rect := range canvas.Rectangles() {
			if rect.Bounds == full && rect.Color == state.windowFrameColor() {
				count++
			}
		}
		return count
	}
	var off gfx.Canvas
	state.Paint(&off, size)
	if frameRects(&off) == 0 {
		t.Fatal("backdrop off: the shell's full-window frame rect is missing")
	}
	if state.chrome.shellBackdrop {
		t.Error("backdrop off: the shell still thinks the backdrop is showing")
	}
	state.winOpacity = winOpacityMedium
	state.setWindowBackdropActive(true)
	var on gfx.Canvas
	state.Paint(&on, size)
	if count := frameRects(&on); count != 0 {
		t.Errorf("backdrop on: %d full-window frame rects still painted under the panes", count)
	}
	if !state.chrome.shellBackdrop {
		t.Error("backdrop on: the shell does not know the backdrop is showing")
	}
}
