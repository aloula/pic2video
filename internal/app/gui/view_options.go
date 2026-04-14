package gui

import (
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type OptionsView struct {
	Profile         *widget.Select
	ImageEffect     *widget.Select
	ImageDur        *widget.Slider
	ImageDurValue   *widget.Label
	Transition      *widget.Slider
	TransitionValue *widget.Label
	FPS30           *widget.Button
	FPS60           *widget.Button
	FPSSelector     fyne.CanvasObject
	QualityLow      *widget.Button
	QualityMedium   *widget.Button
	QualityHigh     *widget.Button
	QualitySelector fyne.CanvasObject
	ExifOverlay     *widget.Check
	ExifFontSize    *widget.Entry
	OrderMode       *widget.Select
	OrderFile       *widget.SelectEntry
	AudioSource     *widget.Select
	Encoder         *widget.Select
	Overwrite       *widget.Check
	DebugExif       *widget.Check
	Container       *fyne.Container
}

func NewOptionsView(cfg GuiRunConfiguration) *OptionsView {
	profile := widget.NewSelect([]string{"fhd", "uhd"}, nil)
	profile.SetSelected(cfg.Profile)
	imageEffect := widget.NewSelect([]string{"static", "kenburns-low", "kenburns-medium", "kenburns-high"}, nil)
	imageEffect.SetSelected(cfg.ImageEffect)
	imageDur := widget.NewSlider(4, 60)
	imageDur.Step = 1
	if cfg.ImageDuration < 4 {
		cfg.ImageDuration = 5
	}
	if cfg.ImageDuration > 60 {
		cfg.ImageDuration = 60
	}
	imageDur.SetValue(cfg.ImageDuration)
	imageDurValue := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	imageDurValue.SetText("Current: " + strconv.Itoa(int(math.Round(imageDur.Value))) + "s")
	imageDur.OnChanged = func(v float64) {
		imageDurValue.SetText("Current: " + strconv.Itoa(int(math.Round(v))) + "s")
	}
	transition := widget.NewSlider(1, 5)
	transition.Step = 1
	if cfg.Transition < 1 {
		cfg.Transition = 1
	}
	if cfg.Transition > 5 {
		cfg.Transition = 5
	}
	transition.SetValue(cfg.Transition)
	transitionValue := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	transitionValue.SetText("Current: " + strconv.Itoa(int(math.Round(transition.Value))) + "s")
	transition.OnChanged = func(v float64) {
		transitionValue.SetText("Current: " + strconv.Itoa(int(math.Round(v))) + "s")
	}
	fps30 := widget.NewButton("30", nil)
	fps60 := widget.NewButton("60", nil)
	setFPSSelection := func(v int) {
		fps30.Importance = widget.MediumImportance
		fps60.Importance = widget.MediumImportance
		if v == 30 {
			fps30.Importance = widget.HighImportance
		} else {
			fps60.Importance = widget.HighImportance
		}
		fps30.Refresh()
		fps60.Refresh()
	}
	fps30.OnTapped = func() { setFPSSelection(30) }
	fps60.OnTapped = func() { setFPSSelection(60) }
	if cfg.FPS == 30 {
		setFPSSelection(30)
	} else {
		// Keep 60 FPS as the GUI default.
		setFPSSelection(60)
	}
	fpsSelector := container.NewHBox(fps30, fps60)

	qualityLow := widget.NewButton("low", nil)
	qualityMedium := widget.NewButton("medium", nil)
	qualityHigh := widget.NewButton("high", nil)
	setQualitySelection := func(q string) {
		qualityLow.Importance = widget.MediumImportance
		qualityMedium.Importance = widget.MediumImportance
		qualityHigh.Importance = widget.MediumImportance
		switch q {
		case "low":
			qualityLow.Importance = widget.HighImportance
		case "medium":
			qualityMedium.Importance = widget.HighImportance
		default:
			qualityHigh.Importance = widget.HighImportance
		}
		qualityLow.Refresh()
		qualityMedium.Refresh()
		qualityHigh.Refresh()
	}
	qualityLow.OnTapped = func() { setQualitySelection("low") }
	qualityMedium.OnTapped = func() { setQualitySelection("medium") }
	qualityHigh.OnTapped = func() { setQualitySelection("high") }
	selectedQuality := strings.TrimSpace(cfg.Quality)
	if selectedQuality == "" {
		selectedQuality = "high"
	}
	setQualitySelection(selectedQuality)
	qualitySelector := container.NewHBox(qualityLow, qualityMedium, qualityHigh)

	exifOverlay := widget.NewCheck("EXIF overlay", nil)
	exifOverlay.SetChecked(cfg.ExifOverlay)
	exifFont := widget.NewEntry()
	exifFont.SetText(strconv.Itoa(cfg.ExifFontSize))
	orderMode := widget.NewSelect([]string{"name", "time", "exif", "explicit"}, nil)
	selectedOrderMode := cfg.OrderMode
	if selectedOrderMode == "" {
		selectedOrderMode = "name"
	}
	orderMode.SetSelected(selectedOrderMode)
	orderFile := widget.NewSelectEntry([]string{"order.txt", "manifest.txt", "sequence.txt", "list.txt"})
	orderFile.SetPlaceHolder("Choose or type file path")
	orderFile.SetText(cfg.OrderFile)
	setOrderFileEnabled := func(mode string) {
		if strings.TrimSpace(mode) == "explicit" {
			orderFile.Enable()
			return
		}
		orderFile.Disable()
	}
	orderMode.OnChanged = setOrderFileEnabled
	setOrderFileEnabled(selectedOrderMode)
	audioSource := widget.NewSelect([]string{"mp3", "video", "mix"}, nil)
	selectedAudioSource := strings.TrimSpace(cfg.AudioSource)
	if selectedAudioSource == "" {
		selectedAudioSource = "mp3"
	}
	audioSource.SetSelected(selectedAudioSource)
	encoder := widget.NewSelect([]string{"auto", "nvenc", "cpu"}, nil)
	encoder.SetSelected(cfg.Encoder)
	overwrite := widget.NewCheck("Overwrite", nil)
	overwrite.SetChecked(cfg.Overwrite)
	debugExif := widget.NewCheck("Debug EXIF", nil)

	c := container.NewVBox(
		widget.NewLabel("Profile"), profile,
		widget.NewLabel("Image effect"), imageEffect,
		widget.NewLabel("Image duration (s)"), container.NewBorder(nil, nil, nil, imageDurValue, imageDur),
		widget.NewLabel("Transition (s)"), container.NewBorder(nil, nil, nil, transitionValue, transition),
		widget.NewLabel("FPS"), fpsSelector,
		widget.NewLabel("Quality"), qualitySelector,
		exifOverlay,
		widget.NewLabel("EXIF font size"), exifFont,
		widget.NewLabel("Order mode"), orderMode,
		widget.NewLabel("Order file"), orderFile,
		widget.NewLabel("Audio source"), audioSource,
		widget.NewLabel("Encoder"), encoder,
		overwrite,
		debugExif,
	)
	return &OptionsView{
		Profile: profile, ImageEffect: imageEffect, ImageDur: imageDur, ImageDurValue: imageDurValue,
		Transition: transition, TransitionValue: transitionValue,
		FPS30: fps30, FPS60: fps60, FPSSelector: fpsSelector,
		QualityLow: qualityLow, QualityMedium: qualityMedium, QualityHigh: qualityHigh, QualitySelector: qualitySelector,
		ExifOverlay:  exifOverlay,
		ExifFontSize: exifFont, OrderMode: orderMode, OrderFile: orderFile,
		AudioSource: audioSource, Encoder: encoder, Overwrite: overwrite, DebugExif: debugExif, Container: c,
	}
}
