package gui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	appversion "github.com/loula/pic2video/internal/app/version"
)

func Run() error {
	baseCfg := DefaultConfiguration()
	state := NewRunStateMachine()
	logs := NewLogStore(500)
	runner := NewRunner("")

	a := app.NewWithID("io.github.loula.pic2video")
	w := a.NewWindow(fmt.Sprintf("pic2video %s", appversion.Short()))
	w.Resize(fyne.NewSize(700, 480))

	form := NewFormView(baseCfg)
	options := NewOptionsView(baseCfg)
	status := NewStatusView()
	logView := NewLogView()
	logController := NewLogController(logs, logView)

	outputPreview := widget.NewLabel(OutputPreviewText(baseCfg))
	outputPreview.Wrapping = fyne.TextTruncate
	outputPreviewAdvanced := widget.NewLabel(OutputPreviewText(baseCfg))
	outputPreviewAdvanced.Wrapping = fyne.TextTruncate
	defaultOutputFolderForInput := func(input string) string {
		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			trimmed = "."
		}
		return filepath.Join(trimmed, "output")
	}
	currentDefaultOutput := func() string {
		return defaultOutputFolderForInput(form.InputEntry.Text)
	}
	isSamePath := func(a, b string) bool {
		if strings.TrimSpace(a) == "" && strings.TrimSpace(b) == "" {
			return true
		}
		return filepath.Clean(a) == filepath.Clean(b)
	}
	outputAuto := isSamePath(form.OutputEntry.Text, currentDefaultOutput())
	refreshPreview := func() {
		cfg := CollectConfiguration(baseCfg, form, options)
		outputPreview.SetText(OutputPreviewText(cfg))
		outputPreviewAdvanced.SetText(OutputPreviewText(cfg))
	}
	form.InputEntry.OnChanged = func(_ string) {
		if outputAuto {
			form.OutputEntry.SetText(currentDefaultOutput())
			form.OutputEntry.CursorColumn = len([]rune(form.OutputEntry.Text))
			form.OutputEntry.Refresh()
		}
		refreshPreview()
	}
	form.OutputEntry.OnChanged = func(_ string) {
		outputAuto = isSamePath(form.OutputEntry.Text, currentDefaultOutput())
		refreshPreview()
	}
	options.OrderMode.OnChanged = func(_ string) { refreshPreview() }
	options.OrderFile.OnChanged = func(_ string) { refreshPreview() }
	options.Profile.OnChanged = func(_ string) { refreshPreview() }
	if options.FPS30 != nil {
		prev := options.FPS30.OnTapped
		options.FPS30.OnTapped = func() {
			if prev != nil {
				prev()
			}
			refreshPreview()
		}
	}
	if options.FPS60 != nil {
		prev := options.FPS60.OnTapped
		options.FPS60.OnTapped = func() {
			if prev != nil {
				prev()
			}
			refreshPreview()
		}
	}
	if options.QualityLow != nil {
		prev := options.QualityLow.OnTapped
		options.QualityLow.OnTapped = func() {
			if prev != nil {
				prev()
			}
			refreshPreview()
		}
	}
	if options.QualityMedium != nil {
		prev := options.QualityMedium.OnTapped
		options.QualityMedium.OnTapped = func() {
			if prev != nil {
				prev()
			}
			refreshPreview()
		}
	}
	if options.QualityHigh != nil {
		prev := options.QualityHigh.OnTapped
		options.QualityHigh.OnTapped = func() {
			if prev != nil {
				prev()
			}
			refreshPreview()
		}
	}

	var startBtnRender *widget.Button
	var startBtnAdvanced *widget.Button
	setStartButtonsEnabled := func(enabled bool) {
		if startBtnRender != nil {
			if enabled {
				startBtnRender.Enable()
			} else {
				startBtnRender.Disable()
			}
		}
		if startBtnAdvanced != nil {
			if enabled {
				startBtnAdvanced.Enable()
			} else {
				startBtnAdvanced.Disable()
			}
		}
	}
	startRender := func(localCfg GuiRunConfiguration) {
		status.SetError("")
		setStartButtonsEnabled(false)

		go func(localCfg GuiRunConfiguration) {
			err := StartRun(context.Background(), runner, state, localCfg, func(s RunStatus) {
				ApplyRunnerStatus(state, status, s)
				if s == RunStatusFinished {
					dialog.ShowInformation("Render Finished", "Processing completed successfully.\n\nOutput: "+localCfg.OutputPath, w)
				}
			}, logController.Append)
			if err != nil {
				logController.Append("system", err.Error())
				status.SetStatus(RunStatusFailed)
				status.StopAnimation()
				dialog.ShowError(fmt.Errorf("render failed: %w", err), w)
			}
			setStartButtonsEnabled(true)
		}(localCfg)
	}
	runFromUI := func() {
		cfg := CollectConfiguration(baseCfg, form, options)
		cfg.LaunchDirectory = ResolveLaunchDirectory()

		res := ValidatePreflight(cfg)
		if !res.OK {
			errMsg := strings.Join(res.Messages, "\n")
			ShowValidationError(status.ErrorLabel, errMsg)
			dialog.ShowError(fmt.Errorf(errMsg), w)
			return
		}

		needsConfirm, msg, err := shouldConfirmShortAudio(cfg)
		if err != nil {
			logController.Append("system", "audio-length check skipped: "+err.Error())
			startRender(cfg)
			return
		}
		if needsConfirm {
			dialog.ShowConfirm("Audio Shorter Than Video", msg, func(ok bool) {
				if ok {
					startRender(cfg)
				}
			}, w)
			return
		}

		startRender(cfg)
	}
	startBtnRender = widget.NewButton("▶  Start Render", runFromUI)
	startBtnAdvanced = widget.NewButton("▶  Start Render", runFromUI)

	// Folder browse helpers
	inputBrowse := widget.NewButton("📁", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			form.InputEntry.SetText(uri.Path())
			refreshPreview()
		}, w)
	})
	outputBrowse := widget.NewButton("📁", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			form.OutputEntry.SetText(uri.Path())
			form.OutputEntry.CursorColumn = len([]rune(form.OutputEntry.Text))
			form.OutputEntry.Refresh()
			refreshPreview()
		}, w)
	})
	inputRow := container.NewBorder(nil, nil, nil, inputBrowse, form.InputEntry)
	outputRow := container.NewBorder(nil, nil, nil, outputBrowse, form.OutputEntry)

	// Render tab: main settings + start button + output preview
	imageDurRow := container.NewBorder(nil, nil, nil, options.ImageDurValue, options.ImageDur)
	transitionRow := container.NewBorder(nil, nil, nil, options.TransitionValue, options.Transition)
	renderForm := widget.NewForm(
		widget.NewFormItem("Input folder", inputRow),
		widget.NewFormItem("Output folder", outputRow),
		widget.NewFormItem("Profile", options.Profile),
		widget.NewFormItem("Image effect", options.ImageEffect),
		widget.NewFormItem("Image duration (s)", imageDurRow),
		widget.NewFormItem("Transition (s)", transitionRow),
		widget.NewFormItem("FPS", options.FPSSelector),
		widget.NewFormItem("Quality", options.QualitySelector),
	)
	renderActions := container.NewVBox(
		widget.NewSeparator(),
		startBtnRender,
		outputPreview,
	)
	renderTab := container.NewBorder(nil, renderActions, nil, nil, renderForm)

	// Advanced tab: encoder, order, exif options
	advancedForm := widget.NewForm(
		widget.NewFormItem("Order mode", options.OrderMode),
		widget.NewFormItem("Order file", options.OrderFile),
		widget.NewFormItem("Audio source", options.AudioSource),
		widget.NewFormItem("Encoder", options.Encoder),
		widget.NewFormItem("EXIF font size", options.ExifFontSize),
	)
	advancedContent := container.NewVBox(
		advancedForm,
		options.ExifOverlay,
		options.Overwrite,
		options.DebugExif,
	)
	advancedActions := container.NewVBox(
		widget.NewSeparator(),
		startBtnAdvanced,
		outputPreviewAdvanced,
	)
	advancedTab := container.NewBorder(nil, advancedActions, nil, nil, advancedContent)

	tabs := container.NewAppTabs(
		container.NewTabItem("Render", renderTab),
		container.NewTabItem("Advanced", advancedTab),
	)

	// Status bar + compact log pinned to the bottom
	bottomPanel := container.NewVBox(
		widget.NewSeparator(),
		status.Container,
		widget.NewSeparator(),
		logView.Container,
	)

	w.SetContent(container.NewBorder(nil, bottomPanel, nil, nil, tabs))
	w.ShowAndRun()
	return nil
}
