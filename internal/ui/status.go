package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

var (
	statusOKColor      = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 255} // Green
	statusFailColor    = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 255} // Red
	statusWaitingColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 255} // Yellow
	statusLabelColor   = color.White
)

// statusOverview shows the WiFi and NTP sync status, replacing the
// date when the clock is not fully online.
type statusOverview struct {
	wifiLabel *canvas.Text
	wifiValue *canvas.Text
	ntpLabel  *canvas.Text
	ntpValue  *canvas.Text
	Content   fyne.CanvasObject
}

func newStatusOverview(textSize float32) *statusOverview {
	wifiLabel := canvas.NewText("WiFi: ", statusLabelColor)
	wifiLabel.TextSize = textSize
	wifiLabel.Alignment = fyne.TextAlignLeading

	wifiValue := canvas.NewText("", statusFailColor)
	wifiValue.TextSize = textSize
	wifiValue.Alignment = fyne.TextAlignLeading

	ntpLabel := canvas.NewText("Tijd (NTP): ", statusLabelColor)
	ntpLabel.TextSize = textSize
	ntpLabel.Alignment = fyne.TextAlignLeading

	ntpValue := canvas.NewText("", statusFailColor)
	ntpValue.TextSize = textSize
	ntpValue.Alignment = fyne.TextAlignLeading

	// WiFi row: label and value side by side, left-aligned within centered container
	wifiRow := container.NewHBox(wifiLabel, wifiValue)

	// NTP row: label and value side by side, left-aligned within centered container
	ntpRow := container.NewHBox(ntpLabel, ntpValue)

	return &statusOverview{
		wifiLabel: wifiLabel,
		wifiValue: wifiValue,
		ntpLabel:  ntpLabel,
		ntpValue:  ntpValue,
		Content:   container.NewVBox(wifiRow, ntpRow),
	}
}

// Update refreshes the displayed status text and colors.
func (s *statusOverview) Update(wifiStatus string, wifiOK, ntpOK bool) {
	// WiFi status
	s.wifiValue.Color = statusFailColor
	switch wifiStatus {
	case "connected":
		s.wifiValue.Text = "Verbonden"
		s.wifiValue.Color = statusOKColor
	case "disconnected":
		s.wifiValue.Text = "Niet verbonden"
	default:
		s.wifiValue.Text = wifiStatus
	}
	s.wifiLabel.Refresh()
	s.wifiValue.Refresh()

	// NTP status
	if ntpOK {
		s.ntpValue.Text = "Gesynchroniseerd"
		s.ntpValue.Color = statusOKColor
	} else {
		s.ntpValue.Text = "Wachten op synchronisatie"
		s.ntpValue.Color = statusWaitingColor
	}
	s.ntpLabel.Refresh()
	s.ntpValue.Refresh()
}
