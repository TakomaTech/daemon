package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "strconv"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"

    "github.com/Iconictacoma/daemon/internal/core"
)

type Project struct {
    Name  string `json:"name"`
    Tempo int    `json:"tempo"`
}

func makeStepGrid(e *core.Engine, chIdx int) *fyne.Container {
    grid := container.NewGridWithColumns(16)
    for s := 0; s < 16; s++ {
        idx := s
        btn := widget.NewButton(".", func() {
            e.ToggleStep(chIdx, idx)
        })
        grid.Add(btn)
    }
    return grid
}

func makeChannelColumn(e *core.Engine, chIdx int) *fyne.Container {
    lbl := widget.NewLabel("Chan")
    vol := widget.NewSlider(0, 1)
    vol.Value = 0.8
    vol.OnChanged = func(v float64) { e.SetChannelVol(chIdx, float32(v)) }
    grid := makeStepGrid(e, chIdx)
    col := container.NewVBox(lbl, vol, grid)
    return col
}

func openPianoRoll(a fyne.App) {
    w := a.NewWindow("Piano Roll")
    g := widget.NewLabel("Piano roll editor placeholder")
    w.SetContent(container.NewCenter(g))
    w.Resize(fyne.NewSize(800, 400))
    w.Show()
}

func main() {
    a := app.New()
    w := a.NewWindow("Daemon")
    engine := core.NewEngine()
    playBtn := widget.NewButton("Play", func() { engine.Play() })
    stopBtn := widget.NewButton("Stop", func() { engine.Stop() })
    pianoBtn := widget.NewButton("Piano Roll", func() { openPianoRoll(a) })
    tempoEntry := widget.NewEntry()
    tempoEntry.SetText("120")
    tempoEntry.OnChanged = func(t string) {
        if v, err := strconv.Atoi(t); err == nil {
            engine.SetTempo(v)
        }
    }

    patternSelect := widget.NewSelect([]string{"Pattern 1", "Pattern 2"}, func(s string) {})

    channelsRow := container.NewHBox()
    for i := 0; i < 6; i++ {
        col := makeChannelColumn(engine, i)
        channelsRow.Add(col)
    }

    mixerCols := container.NewVBox()
    for i := 0; i < 6; i++ {
        s := widget.NewSlider(0, 1)
        s.Value = 0.8
        idx := i
        s.OnChanged = func(v float64) { engine.SetChannelVol(idx, float32(v)) }
        mixerCols.Add(s)
    }

    left := container.NewBorder(container.NewHBox(playBtn, stopBtn, pianoBtn, tempoEntry, patternSelect), nil, nil, nil, container.NewVScroll(channelsRow))
    right := container.NewBorder(nil, nil, nil, nil, container.NewVScroll(mixerCols))
    content := container.New(layout.NewHSplitLayout(), left, right)
    w.SetContent(content)
    w.Resize(fyne.NewSize(1200, 800))
    go func() {
        for {
            engine.Tick()
            time.Sleep(10 * time.Millisecond)
        }
    }()
    saveBtn := widget.NewButton("Save .dmon", func() {
        p := Project{Name: "New Project", Tempo: 120}
        data, _ := json.Marshal(p)
        ioutil.WriteFile("project.dmon", data, 0644)
    })
    loadBtn := widget.NewButton("Load .dmon", func() {
        path := "project.dmon"
        if _, err := os.Stat(path); err == nil {
            b, _ := ioutil.ReadFile(path)
            var p Project
            json.Unmarshal(b, &p)
            _ = p
        }
    })
    menu := container.NewHBox(saveBtn, loadBtn)
    w.SetMainMenu(nil)
    w.SetContent(container.NewVBox(menu, content))
    w.ShowAndRun()
}

