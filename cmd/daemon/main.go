package main

import (
    "fmt"
    "strconv"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
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
    for key := range pianoKeyNames {
        row := container.NewHBox(widget.NewLabel(pianoKeyNames[key]))
        for step := 0; step < stepCount; step++ {
            keyIdx := key
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
            if e.GetPianoNoteState(key, step) {
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

func refreshPatternGrid(e *core.Engine, buttons [][]*widget.Button, patternSelect *widget.Select, tempoEntry *widget.Entry, statusLabel *widget.Label) {
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
    statusLabel.SetText(fmt.Sprintf("Pattern %s • BPM %d", e.CurrentPatternName(), e.Tempo()))
}

func main() {
    a := app.New()
    w := a.NewWindow("Daemon")
    engine := core.NewEngine()
    engine.LoadPlugins("plugins")
    patternButtons := make([][]*widget.Button, engine.ChannelCount())
    for i := range patternButtons {
        patternButtons[i] = make([]*widget.Button, stepCount)
    }
    statusLabel := widget.NewLabel("")
    tempoEntry := widget.NewEntry()
    tempoEntry.SetText(strconv.Itoa(engine.Tempo()))
    var playBtn *widget.Button
    recordBtn := widget.NewButton("Record", func() {})
    patternSelect := widget.NewSelect(engine.PatternNames(), func(s string) {
        engine.SetPatternByName(s)
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
        refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
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
        }
    }
    newPatternBtn := widget.NewButton("New Pattern", func() {
        engine.AddPattern(fmt.Sprintf("Pattern %d", engine.PatternCount()+1))
        refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
    })
    saveBtn := widget.NewButton("Save .dmon", func() {
        if err := engine.SaveProject("project.dmon"); err != nil {
            statusLabel.SetText(err.Error())
            return
        }
        statusLabel.SetText("Saved project.dmon")
    })
    loadBtn := widget.NewButton("Load .dmon", func() {
        if err := engine.LoadProject("project.dmon"); err != nil {
            statusLabel.SetText(err.Error())
            return
        }
        refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
        statusLabel.SetText("Loaded project.dmon")
    })
    channelColumns := container.NewHBox()
    for ch := 0; ch < engine.ChannelCount(); ch++ {
        channelColumns.Add(makeChannelColumn(engine, ch, patternButtons, func() {
            refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
        }))
    }
    pianoBtn := widget.NewButton("Piano Roll", func() {
        openPianoRoll(a, engine)
    })
    toolbar := container.NewHBox(playBtn, stopBtn, recordBtn, tempoEntry, patternSelect)
    actions := container.NewHBox(saveBtn, loadBtn, newPatternBtn, pianoBtn)
    left := container.NewBorder(toolbar, actions, nil, nil, container.NewVScroll(channelColumns))
    right := container.NewVBox(widget.NewLabel("Status"), statusLabel)
    content := container.NewHSplit(left, right)
    w.SetContent(content)
    w.Resize(fyne.NewSize(1400, 900))
    refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
    go func() {
        ticker := time.NewTicker(300 * time.Millisecond)
        for range ticker.C {
            refreshPatternGrid(engine, patternButtons, patternSelect, tempoEntry, statusLabel)
        }
    }()
    w.ShowAndRun()
}

