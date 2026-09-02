package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// hiddenCursorWidget is a transparent Fyne widget that causes the
// cursor to be hidden while the pointer is over it.
type hiddenCursorWidget struct {
	widget.BaseWidget
}

func newHiddenCursorWidget() *hiddenCursorWidget {
	w := &hiddenCursorWidget{}
	w.ExtendBaseWidget(w)
	return w
}

func (w *hiddenCursorWidget) Cursor() desktop.Cursor {
	return desktop.HiddenCursor
}

func (w *hiddenCursorWidget) CreateRenderer() fyne.WidgetRenderer {
	return &hiddenCursorRenderer{
		widget: w,
	}
}

type hiddenCursorRenderer struct {
	widget *hiddenCursorWidget
}

func (r *hiddenCursorRenderer) Layout(size fyne.Size) {
	r.widget.Resize(size)
}

func (r *hiddenCursorRenderer) MinSize() fyne.Size {
	return fyne.Size{}
}

func (r *hiddenCursorRenderer) Objects() []fyne.CanvasObject {
	return nil
}

func (r *hiddenCursorRenderer) Refresh() {
}

func (r *hiddenCursorRenderer) Destroy() {
}
