package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/Iconictacoma/daemon/internal/core"
)

const stepCount = 16

var pianoKeyNames = []string{
	"C4", "C#4", "D4", "D#4", "E4", "F4", "F#4", "G4", "G#4", "A4", "A#4", "B4",
}

func makeStepGrid(e *core.Engine, chIdx int, buttons [][]*widget.Button) *fyne.Container {
	grid := container.NewGridWithColumns(stepCount)
	for s := 0; s < stepCount; s++ {
		step := s
		var btn *widget.Button
		btn = widget.NewButton(".", func() {
			active := !e.GetStep(chIdx, step)
			e.SetStep(chIdx, step, active)
			if active {
				btn.SetText("●")
			} else {
				btn.SetText(".")
			}
		})
		buttons[chIdx][step] = btn
		grid.Add(btn)
	}
	return grid
}

func makeChannelColumn(e *core.Engine, chIdx int, buttons [][]*widget.Button, refresh func()) *fyne.Container {
	label := widget.NewLabel(fmt.Sprintf("Chan %d", chIdx+1))
	vol := widget.NewSlider(0, 1)
	vol.Step = 0.01
	vol.Value = float64(e.GetChannelVol(chIdx))
	vol.OnChanged = func(v float64) {
		e.SetChannelVol(chIdx, float32(v))
	}
	mute := widget.NewButton("Mute", func() {
		e.MuteChannel(chIdx, !e.IsChannelMuted(chIdx))
		refresh()
	})
	grid := makeStepGrid(e, chIdx, buttons)
	return container.NewVBox(label, vol, mute, grid)
}

func openPianoRoll(a fyne.App, e *core.Engine) {
	w := a.NewWindow("Piano Roll")
	rows := container.NewVBox()
	for keyIdx, keyName := range pianoKeyNames {
		row := container.NewHBox(widget.NewLabel(keyName))
		for step := 0; step < stepCount; step++ {
			stepIdx := step
			var btn *widget.Button
			btn = widget.NewButton(".", func() {
				active := !e.GetPianoNoteState(keyIdx, stepIdx)
				e.SetPianoNoteState(keyIdx, stepIdx, active)
				if active {
					btn.SetText("●")
				} else {
					btn.SetText(".")
				}
			})
			if e.GetPianoNoteState(keyIdx, stepIdx) {
				btn.SetText("●")
			}
			row.Add(btn)
		}
		rows.Add(row)
	}
	w.SetContent(container.NewVScroll(rows))
	w.Resize(fyne.NewSize(1200, 520))
	w.Show()
}

func pluginDir() string {
	if env := os.Getenv("DAEMON_PLUGIN_DIR"); env != "" {
		return env
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates := []string{
			filepath.Join(exeDir, "plugins"),
			filepath.Join(exeDir, "..", "lib", "daemon", "plugins"),
			filepath.Join(exeDir, "..", "..", "lib", "daemon", "plugins"),
		}
		for _, p := range candidates {
			if info, err := os.Stat(p); err == nil && info.IsDir() {
				return p
			}
		}
	}
	if info, err := os.Stat("/usr/local/lib/daemon/plugins"); err == nil && info.IsDir() {
		return "/usr/local/lib/daemon/plugins"
	}
	if info, err := os.Stat("/usr/lib/daemon/plugins"); err == nil && info.IsDir() {
		return "/usr/lib/daemon/plugins"
	}
	return "plugins"
}

func pluginListText(engine *core.Engine) string {
	infos := engine.PluginInfos()
	if len(infos) == 0 {
		return "Plugins: none"
	}
	lines := make([]string, 0, len(infos))
	for _, info := range infos {
		heading := info.Name
		if heading == "" {
			heading = info.File
		}
		detail := heading
		meta := make([]string, 0, 2)
		if info.Version != "" {
			meta = append(meta, info.Version)
		}
		if info.Author != "" {
			meta = append(meta, "by "+info.Author)
		}
		if len(meta) > 0 {
			detail += " (" + strings.Join(meta, ", ") + ")"
		}
		if info.Description != "" {
			detail += "\n  " + info.Description
		}
		lines = append(lines, detail)
	}
	return strings.Join(lines, "\n")
}

func saveProjectDialog(w fyne.Window, engine *core.Engine, statusLabel *widget.Label, refresh func()) {
	dlg := dialog.NewFileSave(func(writeCloser fyne.URIWriteCloser, err error) {
		if err != nil || writeCloser == nil {
			return
		}
		path := writeCloser.URI().Path()
		writeCloser.Close()
		if !strings.HasSuffix(strings.ToLower(path), ".dmon") {
			path += ".dmon"
		}
		if err := engine.SaveProject(path); err != nil {
			statusLabel.SetText(err.Error())
			return
		}
		engine.AddRecentProject(path)
		_ = engine.SaveWorkspace()
		statusLabel.SetText("Saved " + path)
		if refresh != nil {
			refresh()
		}
	}, w)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".dmon"}))
	dlg.Show()
}

func loadProjectDialog(w fyne.Window, engine *core.Engine, statusLabel *widget.Label, refresh func()) {
	dlg := dialog.NewFileOpen(func(readCloser fyne.URIReadCloser, err error) {
		if err != nil || readCloser == nil {
			return
		}
		path := readCloser.URI().Path()
		readCloser.Close()
		if err := engine.LoadProject(path); err != nil {
			statusLabel.SetText(err.Error())
			return
		}
		engine.AddRecentProject(path)
		_ = engine.SaveWorkspace()
		statusLabel.SetText("Loaded " + path)
		if refresh != nil {
			refresh()
		}
	}, w)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".dmon"}))
	dlg.Show()
}

func refreshPatternGrid(e *core.Engine, buttons [][]*widget.Button, patternSelect *widget.Select, tempoEntry *widget.Entry, statusLabel *widget.Label, pluginLabel *widget.Label) {
	options := e.PatternNames()
	patternSelect.Options = options
	patternSelect.Refresh()
	patternSelect.SetSelected(e.CurrentPatternName())
	for ch := 0; ch < len(buttons); ch++ {
		for step := 0; step < stepCount; step++ {
			if buttons[ch][step] == nil {
				continue
			}
			if e.GetStep(ch, step) {
				buttons[ch][step].SetText("●")
			} else {
				buttons[ch][step].SetText(".")
			}
		}
	}
	tempoEntry.SetText(strconv.Itoa(e.Tempo()))
	statusLabel.SetText(fmt.Sprintf("%s • Pattern %s • BPM %d", e.ProjectName(), e.CurrentPatternName(), e.Tempo()))
	pluginLabel.SetText(pluginListText(e))
}

func main() {
	headless := os.Getenv("DAEMON_HEADLESS") == "1"
	var a fyne.App
	if headless {
		a = test.NewApp()
	} else {
		a = app.New()
	}
	w := a.NewWindow("Daemon")
	engine := core.NewEngine()
	_ = engine.LoadWorkspace()
	engine.LoadPlugins(pluginDir())
	patternButtons := make([][]*widget.Button, engine.ChannelCount())
	for i := range patternButtons {
		patternButtons[i] = make([]*widget.Button, stepCount)
	}
	statusLabel := widget.NewLabel("")
	pluginLabel := widget.NewLabel(pluginListText(engine))
	projectNameEntry := widget.NewEntry()
	projectNameEntry.SetText(engine.ProjectName())
	var playBtn, recordBtn, newPatternBtn *widget.Button
	var patternSelect, recentSelect *widget.Select
	var refreshUI func()
	projectNameEntry.OnChanged = func(text string) {
		engine.SetProjectName(text)
		refreshUI()
	}
	tempoEntry := widget.NewEntry()
	tempoEntry.SetText(strconv.Itoa(engine.Tempo()))
	refreshUI = func() {
		refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel, pluginLabel)
		if recentSelect != nil {
			recentSelect.Options = engine.RecentProjects()
			recentSelect.Refresh()
		}
	}
	patternSelect = widget.NewSelect(engine.PatternNames(), func(s string) {
		engine.SetPatternByName(s)
		refreshUI()
	})
	patternSelect.SetSelected(engine.CurrentPatternName())
	playBtn = widget.NewButton("Play", func() {
		if engine.IsPlaying() {
			engine.Stop()
			playBtn.SetText("Play")
		} else {
			engine.Play()
			playBtn.SetText("Pause")
		}
	})
	stopBtn := widget.NewButton("Stop", func() {
		engine.Stop()
		playBtn.SetText("Play")
		refreshUI()
	})
	recordBtn = widget.NewButton("Record", func() {
		if engine.IsRecording() {
			engine.StopRecording()
			recordBtn.SetText("Record")
			statusLabel.SetText("Recording stopped")
		} else {
			err := engine.StartRecording("session.wav")
			if err != nil {
				statusLabel.SetText(err.Error())
				return
			}
			recordBtn.SetText("Recording")
			statusLabel.SetText("Recording started")
		}
	})
	tempoEntry.OnChanged = func(t string) {
		if v, err := strconv.Atoi(t); err == nil {
			engine.SetTempo(v)
			refreshUI()
		}
	}
	newPatternBtn = widget.NewButton("New Pattern", func() {
		engine.AddPattern(fmt.Sprintf("Pattern %d", engine.PatternCount()+1))
		refreshUI()
	})
	saveBtn := widget.NewButton("Save .dmon", func() {
		saveProjectDialog(w, engine, statusLabel, func() {
			refreshUI()
		})
	})
	loadBtn := widget.NewButton("Load .dmon", func() {
		loadProjectDialog(w, engine, statusLabel, func() {
			refreshUI()
			projectNameEntry.SetText(engine.ProjectName())
		})
	})
	templateSelect := widget.NewSelect([]string{"Blank", "Starter Kit", "Beat Machine"}, func(s string) {})
	templateSelect.SetSelected("Blank")
	newProjectBtn := widget.NewButton("New Project", func() {
		engine.NewProject(projectNameEntry.Text, templateSelect.Selected)
		refreshUI()
	})
	recentSelect = widget.NewSelect(engine.RecentProjects(), func(path string) {
		if path == "" {
			return
		}
		if err := engine.LoadProject(path); err != nil {
			statusLabel.SetText(err.Error())
			return
		}
		engine.AddRecentProject(path)
		_ = engine.SaveWorkspace()
		refreshUI()
		projectNameEntry.SetText(engine.ProjectName())
	})
	channelColumns := container.NewHBox()
	for ch := 0; ch < engine.ChannelCount(); ch++ {
		channelColumns.Add(makeChannelColumn(engine, ch, patternButtons, func() {
			refreshUI()
		}))
	}
	pianoBtn := widget.NewButton("Piano Roll", func() {
		openPianoRoll(a, engine)
	})
	toolbar := container.NewHBox(playBtn, stopBtn, recordBtn, projectNameEntry, templateSelect, patternSelect)
	actions := container.NewHBox(saveBtn, loadBtn, newPatternBtn, newProjectBtn, pianoBtn)
	left := container.NewBorder(toolbar, actions, nil, nil, container.NewVScroll(channelColumns))
	right := container.NewVBox(widget.NewLabel("Status"), statusLabel, widget.NewLabel("Plugin Metadata"), pluginLabel, widget.NewLabel("Recent Projects"), recentSelect)
	content := container.NewHSplit(left, right)
	w.SetContent(content)
	w.Resize(fyne.NewSize(1400, 900))
	refreshUI()
	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		for range ticker.C {
			refreshUI()
		}
	}()

	if headless {
		log.Println("Daemon running in headless mode. Press Ctrl+C to exit.")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		return
	}

	w.ShowAndRun()
}
