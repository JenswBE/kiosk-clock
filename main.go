package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	configFileName = ".fullscreen-clock.json"
	datePadding    = 40
	backlightDir   = "/sys/class/backlight"
)

var dutchWeekdays = [...]string{
	"zondag",
	"maandag",
	"dinsdag",
	"woensdag",
	"donderdag",
	"vrijdag",
	"zaterdag",
}

var dutchMonths = [...]string{
	"januari",
	"februari",
	"maart",
	"april",
	"mei",
	"juni",
	"juli",
	"augustus",
	"september",
	"oktober",
	"november",
	"december",
}

type Config struct {
	TextColor       string `json:"text_color"`
	BackgroundColor string `json:"background_color"`
}

var defaultConfig = Config{
	TextColor:       "#FFFFFF",
	BackgroundColor: "#000000",
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return configFileName
	}

	return filepath.Join(home, configFileName)
}

func loadConfig() Config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return defaultConfig
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultConfig
	}

	if _, err := parseColor(config.TextColor); err != nil {
		config.TextColor = defaultConfig.TextColor
	}

	if _, err := parseColor(config.BackgroundColor); err != nil {
		config.BackgroundColor = defaultConfig.BackgroundColor
	}

	return config
}

func saveConfig(config Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(configPath(), data, 0o600)
}

func parseColor(value string) (color.Color, error) {
	value = strings.TrimPrefix(value, "#")

	if len(value) != 6 {
		return nil, fmt.Errorf("invalid color: %q", value)
	}

	r, err := strconv.ParseUint(value[0:2], 16, 8)
	if err != nil {
		return nil, err
	}

	g, err := strconv.ParseUint(value[2:4], 16, 8)
	if err != nil {
		return nil, err
	}

	b, err := strconv.ParseUint(value[4:6], 16, 8)
	if err != nil {
		return nil, err
	}

	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}, nil
}

func colorHex(c color.Color) string {
	r, g, b, _ := c.RGBA()

	return fmt.Sprintf(
		"#%02X%02X%02X",
		uint8(r>>8),
		uint8(g>>8),
		uint8(b>>8),
	)
}

func ntpSynced() bool {
	cmd := exec.Command(
		"timedatectl",
		"show",
		"-p", "NTPSynchronized",
		"--value",
	)

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "yes"
}

func dutchDate(t time.Time) string {
	return dutchWeekdays[t.Weekday()] + " " +
		t.Format("2") + " " +
		dutchMonths[t.Month()-1] + " " +
		t.Format("2006")
}

func showSettings(
	w fyne.Window,
	config *Config,
	background *canvas.Rectangle,
	clock *canvas.Text,
	date *canvas.Text,
) {
	updateColors := func() {
		textColor, _ := parseColor(config.TextColor)
		backgroundColor, _ := parseColor(config.BackgroundColor)

		clock.Color = textColor
		date.Color = textColor
		background.FillColor = backgroundColor

		clock.Refresh()
		date.Refresh()
		background.Refresh()
	}

	pickColor := func(
		title string,
		description string,
		current string,
		onChanged func(color.Color),
	) {
		currentColor, _ := parseColor(current)

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
		config.TextColor,
		func() {
			pickColor(
				"Tekstkleur",
				"Kies de kleur voor de klok en datum.",
				config.TextColor,
				func(c color.Color) {
					config.TextColor = colorHex(c)
					saveConfig(*config)
				},
			)
		},
	)

	backgroundColorButton := widget.NewButton(
		config.BackgroundColor,
		func() {
			pickColor(
				"Achtergrondkleur",
				"Kies de achtergrondkleur.",
				config.BackgroundColor,
				func(c color.Color) {
					config.BackgroundColor = colorHex(c)
					saveConfig(*config)
				},
			)
		},
	)

	brightness, err := readBacklightPercentage()
	if err != nil {
		log.Println("Failed to read backlight brightness:", err)
		brightness = 50
	}
	brightnessSlider := widget.NewSlider(0, 100)
	brightnessSlider.Step = 1
	brightnessSlider.Value = float64(brightness)

	brightnessSlider.OnChanged = func(value float64) {
		if err := setBacklightBrightness(uint8(value)); err != nil {
			log.Println("Failed to set backlight brightness:", err)
			return
		}
	}

	settingsObjects := []fyne.CanvasObject{
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

		layout.NewSpacer(),
	}

	content := container.NewVBox(settingsObjects...)

	closeButton := widget.NewButton("Sluiten", nil)

	settingsContent := container.NewBorder(
		nil,
		closeButton,
		nil,
		nil,
		container.NewPadded(content),
	)

	popup := widget.NewModalPopUp(
		settingsContent,
		w.Canvas(),
	)

	closeButton.OnTapped = popup.Hide

	popup.Resize(fyne.NewSize(380, 320))
	popup.Show()
}

func main() {
	config := loadConfig()

	a := app.NewWithID("be.example.fullscreenclock")

	w := a.NewWindow("Clock")
	w.SetFullScreen(true)
	w.SetPadded(false)

	textColor, _ := parseColor(config.TextColor)
	backgroundColor, _ := parseColor(config.BackgroundColor)

	background := canvas.NewRectangle(backgroundColor)

	clock := canvas.NewText("", textColor)
	clock.TextSize = 275
	clock.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	clock.FontSource = resourceChivoMonoExtraBoldTtf

	date := canvas.NewText("", textColor)
	date.TextSize = 80

	update := func() {
		now := time.Now()

		if ntpSynced() {
			// Blink the colon every second.
			//
			// Both ":" and " " occupy one character in Chivo Mono,
			// so the clock never changes width.
			separator := ":"
			if now.Second()%2 == 1 {
				separator = " "
			}

			clock.Text = now.Format("15") + separator + now.Format("04")
		} else {
			clock.Text = "Tijd synchroniseren ..."
		}

		date.Text = dutchDate(now)

		clock.Refresh()
		date.Refresh()
	}

	// Keep the date slightly above the bottom edge by reserving
	// a fixed amount of space underneath it.
	datePaddingObject := canvas.NewRectangle(color.Transparent)
	datePaddingObject.SetMinSize(
		fyne.NewSize(1, datePadding),
	)

	dateArea := container.NewBorder(
		nil,
		datePaddingObject,
		nil,
		nil,
		container.NewCenter(date),
	)

	content := container.NewBorder(
		nil,
		dateArea,
		nil,
		nil,
		container.NewCenter(clock),
	)

	// Small settings cog in the top-right corner.
	settingsButton := widget.NewButtonWithIcon(
		"",
		theme.SettingsIcon(),
		func() {
			showSettings(
				w,
				&config,
				background,
				clock,
				date,
			)
		},
	)
	settingsButton.Importance = widget.LowImportance

	topRight := container.NewHBox(
		layout.NewSpacer(),
		settingsButton,
	)

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
