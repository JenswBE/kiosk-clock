package ui

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"clock/internal/config"
	"clock/internal/system"
)

func showSettings(
	w fyne.Window,
	cfg *config.Config,
	background *canvas.Rectangle,
	clockText *canvas.Text,
	dateText *canvas.Text,
) {
	updateColors := func() {
		textColor, _ := config.ParseColor(cfg.TextColor)
		backgroundColor, _ := config.ParseColor(cfg.BackgroundColor)

		clockText.Color = textColor
		dateText.Color = textColor
		background.FillColor = backgroundColor

		clockText.Refresh()
		dateText.Refresh()
		background.Refresh()
	}

	pickColor := func(
		title string,
		description string,
		current string,
		onChanged func(color.Color),
	) {
		currentColor, _ := config.ParseColor(current)

		picker := dialog.NewColorPicker(
			title,
			description,
			func(c color.Color) {
				onChanged(c)
				updateColors()
			},
			w,
		)
		picker.Advanced = true
		picker.SetColor(currentColor)
		picker.Show()
	}

	textColorButton := widget.NewButton(
		cfg.TextColor,
		func() {
			pickColor(
				"Tekstkleur",
				"Kies de kleur voor de klok en datum.",
				cfg.TextColor,
				func(c color.Color) {
					cfg.TextColor = config.ColorHex(c)
					config.Save(*cfg)
				},
			)
		},
	)

	backgroundColorButton := widget.NewButton(
		cfg.BackgroundColor,
		func() {
			pickColor(
				"Achtergrondkleur",
				"Kies de achtergrondkleur.",
				cfg.BackgroundColor,
				func(c color.Color) {
					cfg.BackgroundColor = config.ColorHex(c)
					config.Save(*cfg)
				},
			)
		},
	)

	brightness, err := system.ReadBacklightPercentage()
	if err != nil {
		log.Println("Failed to read backlight brightness:", err)
		brightness = 50
	}

	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Step = 1
	brightnessSlider.Value = float64(brightness)

	brightnessSlider.OnChanged = func(value float64) {
		if err := system.SetBacklightBrightness(uint8(value)); err != nil {
			log.Println("Failed to set backlight brightness:", err)
			return
		}
	}

	rebootButton := widget.NewButton(
		"Herstart",
		func() {
			dialog.NewConfirm(
				"Herstarten",
				"Weet je zeker dat je het toestel wilt herstarten?",
				func(confirmed bool) {
					if !confirmed {
						return
					}

					if err := system.Reboot(); err != nil {
						log.Println("Failed to reboot:", err)
					}
				},
				w,
			).Show()
		},
	)

	shutdownButton := widget.NewButton(
		"Uitschakelen",
		func() {
			dialog.NewConfirm(
				"Uitschakelen",
				"Weet je zeker dat je het toestel wilt uitschakelen?",
				func(confirmed bool) {
					if !confirmed {
						return
					}

					if err := system.Shutdown(); err != nil {
						log.Println("Failed to shut down:", err)
					}
				},
				w,
			).Show()
		},
	)

	closeButton := widget.NewButton("Sluiten", nil)

	content := container.NewVBox(
		widget.NewLabelWithStyle(
			"Instellingen",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),

		widget.NewLabel("Tekstkleur"),
		textColorButton,

		widget.NewLabel("Achtergrondkleur"),
		backgroundColorButton,

		widget.NewLabel("Helderheid"),
		brightnessSlider,

		widget.NewSeparator(),

		rebootButton,
		newSpacer(5),
		shutdownButton,

		newSpacer(5),

		closeButton,
	)

	popup := widget.NewModalPopUp(
		container.NewPadded(content),
		w.Canvas(),
	)

	closeButton.OnTapped = popup.Hide

	popup.Resize(fyne.NewSize(380, 420))
	popup.Show()
}
