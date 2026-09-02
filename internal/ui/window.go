// Package ui builds and runs the fullscreen clock window.
package ui

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clock/internal/assets"
	"clock/internal/config"
	"clock/internal/system"
)

const datePadding = 40

// Run builds and shows the fullscreen clock window. It blocks until
// the window is closed.
func Run() {
	cfg := config.Load()

	a := app.NewWithID("be.jenswbe.fullscreenclock")

	w := a.NewWindow("Clock")
	w.SetFullScreen(true)
	w.SetPadded(false)

	textColor, _ := config.ParseColor(cfg.TextColor)
	backgroundColor, _ := config.ParseColor(cfg.BackgroundColor)

	background := canvas.NewRectangle(backgroundColor)

	timeText := canvas.NewText("", textColor)
	timeText.TextSize = 275
	timeText.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	timeText.FontSource = assets.ResourceChivoMonoExtraBoldTtf

	dateText := canvas.NewText("", textColor)
	dateText.TextSize = 80

	status := newStatusOverview(70)

	// container.NewCenter positions its child based on the child's size at
	// the time the center container itself is laid out. Refreshing a
	// canvas.Text only repaints it in place, so these wrappers must be
	// refreshed too whenever the text content changes, otherwise the text
	// keeps the (stale) position/size from the previous layout pass.
	timeCenter := container.NewCenter(timeText)
	dateCenter := container.NewCenter(dateText)
	statusCenter := container.NewCenter(status.Content)

	update := func() {
		// Check system status
		wifiStatus, wifiOK, wifiErr := system.WifiConnected()
		if wifiErr != nil {
			log.Println("Failed to determine WiFi status:", wifiErr)
		}
		ntpOK := system.NTPSynced()

		if wifiOK && ntpOK {
			// Everything fine => Show time and date
			background.FillColor = backgroundColor
			status.Content.Hide()
			now := time.Now()
			timeText.Text = formatTime(now)
			timeText.Show()
			timeCenter.Refresh()
			dateText.Text = formatDate(now)
			dateText.Show()
			dateCenter.Refresh()
		} else {
			// System issue => Show status overview
			background.FillColor = color.Black
			status.Update(wifiStatus, wifiOK, ntpOK)
			timeText.Hide()
			dateText.Hide()
			status.Content.Show()
			statusCenter.Refresh()
		}
		background.Refresh()
	}

	dateArea := container.NewVBox(
		dateCenter,
		// Keep the date/status area slightly above the bottom edge by
		// reserving a fixed amount of space underneath it.
		newSpacer(datePadding),
	)

	content := container.NewBorder(
		nil,
		dateArea,
		nil,
		nil,
		container.NewStack(timeCenter, statusCenter),
	)

	// Small settings cog in the top-right corner.
	settingsButton := widget.NewButtonWithIcon(
		"",
		theme.SettingsIcon(),
		func() {
			showSettings(
				w,
				&cfg,
				background,
				timeText,
				dateText,
			)
		},
	)
	settingsButton.Importance = widget.LowImportance

	topRight := container.NewHBox(
		layout.NewSpacer(),
		settingsButton,
	)

	hiddenCursor := newHiddenCursorWidget()

	root := container.NewStack(
		background,
		content,
		container.NewBorder(
			topRight,
			nil,
			nil,
			nil,
			nil,
		),
		hiddenCursor,
	)

	w.SetContent(root)

	// All UI mutations, including the initial update, use fyne.Do.
	fyne.Do(update)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fyne.Do(update)
		}
	}()

	w.ShowAndRun()
}

func newSpacer(height float32) *canvas.Rectangle {
	r := canvas.NewRectangle(color.Transparent)
	r.SetMinSize(fyne.NewSize(1, height))
	return r
}
