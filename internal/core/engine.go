package core

import (
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "math"
    "os"
    "path/filepath"
    "plugin"
    "strings"
    "sync"
    "time"

    "github.com/hajimehoshi/oto/v2"
)

const stepCount = 16
const defaultChannelCount = 6

var pianoKeyFrequencies = []float64{
    261.63, 277.18, 293.66, 311.13,
    329.63, 349.23, 369.99, 392.00,
    415.30, 440.00, 466.16, 493.88,
}

type Engine struct {
    ctx            *oto.Context
    p              oto.Player
    pw             *io.PipeWriter
    running        bool
    recording      bool
    recordFile     *os.File
    recordBytes    int
    mu             sync.Mutex
    plugins        []Plugin
    pluginFiles    []string
    pluginInfos    []PluginInfo
    recentProjects []string
    channels       []Channel
    channelPhases  []float64
    pianoPhases    []float64
    patterns       []Pattern
    currentPattern int
    bpm            int
    stepIndex      int
    sampleRate     float64
    samplesSince   int
    pianoRoll      PianoRoll
    projectName    string
    masterVolume   float32
    limiterEnabled bool
    undoStack      []Action
    redoStack      []Action
}

type Plugin interface {
    Process([]float32) []float32
}

type PluginFunc func([]float32) []float32

func (f PluginFunc) Process(b []float32) []float32 { return f(b) }

type PluginInfo struct {
    File        string `json:"file"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Version     string `json:"version"`
    Author      string `json:"author"`
}

type Action struct {
    Kind         string
    PatternIndex int
    Channel      int
    Step         int
    Key          int
    OldState     bool
    NewState     bool
}

func NewEngine() *Engine {
    ctx, ready, err := oto.NewContext(44100, 2, 2)
    if err != nil {
        ctx = nil
    } else {
        <-ready
    }
    e := &Engine{
        ctx:            ctx,
        bpm:            120,
        sampleRate:     44100,
        currentPattern: 0,
        projectName:    "New Project",
        masterVolume:   1.0,
        limiterEnabled: true,
    }
    for i := 0; i < defaultChannelCount; i++ {
        c := Channel{Name: fmt.Sprintf("Chan %d", i+1), Volume: 0.8, Freq: 220.0 + float64(i)*110.0}
        e.channels = append(e.channels, c)
        e.channelPhases = append(e.channelPhases, 0)
    }
    e.pianoPhases = make([]float64, len(pianoKeyFrequencies))
    e.patterns = []Pattern{newPattern("Pattern 1", defaultChannelCount), newPattern("Pattern 2", defaultChannelCount)}
    e.pianoRoll = newPianoRoll()
    return e
}

func newPattern(name string, channelCount int) Pattern {
    p := Pattern{Name: name, Steps: make([][]bool, channelCount)}
    for i := 0; i < channelCount; i++ {
        p.Steps[i] = make([]bool, stepCount)
    }
    return p
}

func newPianoRoll() PianoRoll {
    p := PianoRoll{Notes: make([][]bool, len(pianoKeyFrequencies))}
    for i := range p.Notes {
        p.Notes[i] = make([]bool, stepCount)
    }
    return p
}

func (e *Engine) AddChannel(name string, freq float64) {
    e.mu.Lock()
    e.channels = append(e.channels, Channel{Name: name, Volume: 0.8, Freq: freq})
    e.channelPhases = append(e.channelPhases, 0)
    for i := range e.patterns {
        if len(e.patterns[i].Steps) < len(e.channels) {
            e.patterns[i].Steps = append(e.patterns[i].Steps, make([]bool, stepCount))
        }
    }
    e.mu.Unlock()
}

func (e *Engine) ToggleStep(ch, step int) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch < 0 || ch >= len(e.channels) || step < 0 || step >= stepCount {
        return
    }
    if e.currentPattern < 0 || e.currentPattern >= len(e.patterns) {
        return
    }
    e.patterns[e.currentPattern].Steps[ch][step] = !e.patterns[e.currentPattern].Steps[ch][step]
}

func (e *Engine) SetStep(ch, step int, active bool) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch < 0 || ch >= len(e.channels) || step < 0 || step >= stepCount {
        return
    }
    if e.currentPattern < 0 || e.currentPattern >= len(e.patterns) {
        return
    }
    e.patterns[e.currentPattern].Steps[ch][step] = active
}

func (e *Engine) GetStep(ch, step int) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch < 0 || ch >= len(e.channels) || step < 0 || step >= stepCount {
        return false
    }
    if e.currentPattern < 0 || e.currentPattern >= len(e.patterns) {
        return false
    }
    return e.patterns[e.currentPattern].Steps[ch][step]
}

func (e *Engine) SetPianoNoteState(key, step int, active bool) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if key < 0 || key >= len(e.pianoRoll.Notes) || step < 0 || step >= stepCount {
        return
    }
    e.pianoRoll.Notes[key][step] = active
}

func (e *Engine) GetPianoNoteState(key, step int) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    if key < 0 || key >= len(e.pianoRoll.Notes) || step < 0 || step >= stepCount {
        return false
    }
    return e.pianoRoll.Notes[key][step]
}

func (e *Engine) SetPattern(index int) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if index < 0 || index >= len(e.patterns) {
        return
    }
    e.currentPattern = index
}

func (e *Engine) SetPatternByName(name string) {
    e.mu.Lock()
    defer e.mu.Unlock()
    for i, p := range e.patterns {
        if p.Name == name {
            e.currentPattern = i
            return
        }
    }
}

func (e *Engine) PatternNames() []string {
    e.mu.Lock()
    defer e.mu.Unlock()
    names := make([]string, len(e.patterns))
    for i, p := range e.patterns {
        names[i] = p.Name
    }
    return names
}

func (e *Engine) CurrentPatternName() string {
    e.mu.Lock()
    defer e.mu.Unlock()
    if e.currentPattern < 0 || e.currentPattern >= len(e.patterns) {
        return ""
    }
    return e.patterns[e.currentPattern].Name
}

func (e *Engine) PatternCount() int {
    e.mu.Lock()
    defer e.mu.Unlock()
    return len(e.patterns)
}

func (e *Engine) ChannelCount() int {
    e.mu.Lock()
    defer e.mu.Unlock()
    return len(e.channels)
}

func (e *Engine) ProjectName() string {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.projectName
}

func (e *Engine) SetProjectName(name string) {
    e.mu.Lock()
    e.projectName = name
    e.mu.Unlock()
}

func (e *Engine) SetMasterVolume(v float32) {
    if v < 0 {
        v = 0
    }
    if v > 1 {
        v = 1
    }
    e.mu.Lock()
    e.masterVolume = v
    e.mu.Unlock()
}

func (e *Engine) MasterVolume() float32 {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.masterVolume
}

func (e *Engine) SetLimiterEnabled(enabled bool) {
    e.mu.Lock()
    e.limiterEnabled = enabled
    e.mu.Unlock()
}

func (e *Engine) LimiterEnabled() bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.limiterEnabled
}

func (e *Engine) AddPattern(name string) {
    e.mu.Lock()
    defer e.mu.Unlock()
    p := newPattern(name, len(e.channels))
    e.patterns = append(e.patterns, p)
}

func (e *Engine) clearHistoryLocked() {
    e.undoStack = nil
    e.redoStack = nil
}

func (e *Engine) PushAction(action Action) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if action.Kind == "" {
        return
    }
    e.undoStack = append(e.undoStack, action)
    if len(e.undoStack) > 128 {
        e.undoStack = e.undoStack[1:]
    }
    e.redoStack = nil
}

func (e *Engine) Undo() {
    e.mu.Lock()
    if len(e.undoStack) == 0 {
        e.mu.Unlock()
        return
    }
    action := e.undoStack[len(e.undoStack)-1]
    e.undoStack = e.undoStack[:len(e.undoStack)-1]
    e.redoStack = append(e.redoStack, action)
    if action.Kind == "step" {
        if action.PatternIndex >= 0 && action.PatternIndex < len(e.patterns) && action.Channel >= 0 && action.Channel < len(e.channels) {
            e.patterns[action.PatternIndex].Steps[action.Channel][action.Step] = action.OldState
        }
    }
    if action.Kind == "note" {
        if action.Key >= 0 && action.Key < len(e.pianoRoll.Notes) {
            e.pianoRoll.Notes[action.Key][action.Step] = action.OldState
        }
    }
    e.mu.Unlock()
}

func (e *Engine) Redo() {
    e.mu.Lock()
    if len(e.redoStack) == 0 {
        e.mu.Unlock()
        return
    }
    action := e.redoStack[len(e.redoStack)-1]
    e.redoStack = e.redoStack[:len(e.redoStack)-1]
    e.undoStack = append(e.undoStack, action)
    if action.Kind == "step" {
        if action.PatternIndex >= 0 && action.PatternIndex < len(e.patterns) && action.Channel >= 0 && action.Channel < len(e.channels) {
            e.patterns[action.PatternIndex].Steps[action.Channel][action.Step] = action.NewState
        }
    }
    if action.Kind == "note" {
        if action.Key >= 0 && action.Key < len(e.pianoRoll.Notes) {
            e.pianoRoll.Notes[action.Key][action.Step] = action.NewState
        }
    }
    e.mu.Unlock()
}

func (e *Engine) CanUndo() bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return len(e.undoStack) > 0
}

func (e *Engine) CanRedo() bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return len(e.redoStack) > 0
}

func (e *Engine) NewProject(name, template string) {
    e.mu.Lock()
    defer e.mu.Unlock()
    e.projectName = name
    e.bpm = 120
    e.currentPattern = 0
    e.stepIndex = 0
    e.samplesSince = 0
    e.masterVolume = 1.0
    e.limiterEnabled = true
    for i := range e.channels {
        e.channels[i].Muted = false
        e.channels[i].Volume = 0.8
        e.channels[i].Phase = 0
    }
    e.patterns = []Pattern{newPattern("Pattern 1", len(e.channels)), newPattern("Pattern 2", len(e.channels))}
    if template == "Starter Kit" {
        for i := 0; i < len(e.channels) && i < 2; i++ {
            for step := 0; step < stepCount; step += 2 {
                e.patterns[i].Steps[i][step] = true
            }
        }
    }
    if template == "Beat Machine" {
        for step := 0; step < stepCount; step += 4 {
            e.patterns[0].Steps[0][step] = true
            if step+2 < stepCount {
                e.patterns[1].Steps[1][step+2] = true
            }
        }
    }
    e.pianoRoll = newPianoRoll()
    e.clearHistoryLocked()
}

func (e *Engine) PluginNames() []string {
    e.mu.Lock()
    defer e.mu.Unlock()
    names := make([]string, len(e.pluginInfos))
    for i, info := range e.pluginInfos {
        if info.Name != "" {
            names[i] = info.Name
        } else {
            names[i] = info.File
        }
    }
    return names
}

func (e *Engine) PluginInfos() []PluginInfo {
    e.mu.Lock()
    defer e.mu.Unlock()
    infos := make([]PluginInfo, len(e.pluginInfos))
    copy(infos, e.pluginInfos)
    return infos
}

func (e *Engine) workspaceFilePath() string {
    cwd, err := os.Getwd()
    if err != nil {
        cwd = "."
    }
    return filepath.Join(cwd, ".daemon_workspace.json")
}

func (e *Engine) SaveWorkspace() error {
    e.mu.Lock()
    workspace := Workspace{RecentProjects: append([]string(nil), e.recentProjects...)}
    e.mu.Unlock()

    data, err := json.MarshalIndent(workspace, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(e.workspaceFilePath(), data, 0644)
}

func (e *Engine) LoadWorkspace() error {
    path := e.workspaceFilePath()
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    var workspace Workspace
    if err := json.Unmarshal(data, &workspace); err != nil {
        return err
    }
    e.mu.Lock()
    e.recentProjects = append([]string(nil), workspace.RecentProjects...)
    e.mu.Unlock()
    return nil
}

func (e *Engine) AddRecentProject(path string) {
    path = filepath.Clean(path)
    e.mu.Lock()
    defer e.mu.Unlock()
    idx := -1
    for i, p := range e.recentProjects {
        if p == path {
            idx = i
            break
        }
    }
    if idx >= 0 {
        e.recentProjects = append(append([]string(nil), e.recentProjects[:idx]...), e.recentProjects[idx+1:]...)
    }
    e.recentProjects = append([]string{path}, e.recentProjects...)
    if len(e.recentProjects) > 10 {
        e.recentProjects = e.recentProjects[:10]
    }
}

func (e *Engine) RecentProjects() []string {
    e.mu.Lock()
    defer e.mu.Unlock()
    recents := append([]string(nil), e.recentProjects...)
    return recents
}

func (e *Engine) SetTempo(bpm int) {
    if bpm < 20 {
        bpm = 20
    }
    if bpm > 300 {
        bpm = 300
    }
    e.mu.Lock()
    e.bpm = bpm
    e.mu.Unlock()
}

func (e *Engine) Tempo() int {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.bpm
}

func (e *Engine) SetChannelVol(ch int, v float32) {
    if v < 0 {
        v = 0
    }
    if v > 1 {
        v = 1
    }
    e.mu.Lock()
    if ch >= 0 && ch < len(e.channels) {
        e.channels[ch].Volume = v
    }
    e.mu.Unlock()
}

func (e *Engine) GetChannelVol(ch int) float32 {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch < 0 || ch >= len(e.channels) {
        return 0
    }
    return e.channels[ch].Volume
}

func (e *Engine) MuteChannel(ch int, mute bool) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch >= 0 && ch < len(e.channels) {
        e.channels[ch].Muted = mute
    }
}

func (e *Engine) IsChannelMuted(ch int) bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    if ch < 0 || ch >= len(e.channels) {
        return false
    }
    return e.channels[ch].Muted
}

func (e *Engine) IsPlaying() bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.running
}

func (e *Engine) IsRecording() bool {
    e.mu.Lock()
    defer e.mu.Unlock()
    return e.recording
}

func (e *Engine) Play() {
    e.mu.Lock()
    if e.running {
        e.mu.Unlock()
        return
    }
    if e.ctx != nil {
        pr, pw := io.Pipe()
        e.p = e.ctx.NewPlayer(pr)
        e.pw = pw
    }
    e.running = true
    e.mu.Unlock()
    go e.render()
}

func (e *Engine) Stop() {
    e.mu.Lock()
    if !e.running {
        e.mu.Unlock()
        return
    }
    e.running = false
    p := e.p
    pw := e.pw
    e.p = nil
    e.pw = nil
    e.mu.Unlock()
    if pw != nil {
        pw.Close()
    }
    if p != nil {
        p.Close()
    }
    e.StopRecording()
}

func (e *Engine) StartRecording(path string) error {
    e.mu.Lock()
    if e.recording {
        e.mu.Unlock()
        return nil
    }
    file, err := os.Create(path)
    if err != nil {
        e.mu.Unlock()
        return err
    }
    _, err = file.Write(make([]byte, 44))
    if err != nil {
        file.Close()
        e.mu.Unlock()
        return err
    }
    e.recordFile = file
    e.recordBytes = 0
    e.recording = true
    e.mu.Unlock()
    return nil
}

func (e *Engine) StopRecording() error {
    e.mu.Lock()
    if !e.recording || e.recordFile == nil {
        e.mu.Unlock()
        return nil
    }
    file := e.recordFile
    bytes := e.recordBytes
    e.recording = false
    e.recordFile = nil
    e.recordBytes = 0
    e.mu.Unlock()
    if err := writeWAVHeader(file, bytes); err != nil {
        file.Close()
        return err
    }
    return file.Close()
}

func (e *Engine) SaveProject(path string) error {
    e.mu.Lock()
    proj := Project{
        Name:           e.projectName,
        Tempo:          e.bpm,
        CurrentPattern: e.currentPattern,
        Channels:       e.channels,
        Patterns:       e.patterns,
        PianoRoll:      e.pianoRoll,
    }
    e.mu.Unlock()
    data, err := json.MarshalIndent(proj, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func (e *Engine) LoadProject(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
	return e.LoadProjectFromBytes(data)
}

func (e *Engine) LoadProjectFromBytes(data []byte) error {
	var proj Project
	if err := json.Unmarshal(data, &proj); err != nil {
		return err
	}
	e.mu.Lock()
	e.projectName = proj.Name
	e.bpm = proj.Tempo
	e.currentPattern = proj.CurrentPattern
	if len(proj.Channels) > 0 {
		e.channels = proj.Channels
	}
	if len(proj.Patterns) > 0 {
		e.patterns = proj.Patterns
	}
	if len(proj.PianoRoll.Notes) > 0 {
		e.pianoRoll = proj.PianoRoll
	}
	for len(e.channelPhases) < len(e.channels) {
		e.channelPhases = append(e.channelPhases, 0)
	}
	for i := range e.patterns {
		if len(e.patterns[i].Steps) < len(e.channels) {
			for len(e.patterns[i].Steps) < len(e.channels) {
				e.patterns[i].Steps = append(e.patterns[i].Steps, make([]bool, stepCount))
			}
		}
	}
	for len(e.pianoPhases) < len(e.pianoRoll.Notes) {
		e.pianoPhases = append(e.pianoPhases, 0)
	}
	e.mu.Unlock()
	return nil
}

func (e *Engine) LoadPlugins(dir string) {
    files, _ := os.ReadDir(dir)
    for _, f := range files {
        if f.IsDir() {
            continue
        }
        if filepath.Ext(f.Name()) == ".so" {
            p, err := plugin.Open(filepath.Join(dir, f.Name()))
            if err != nil {
                continue
            }
            sym, err := p.Lookup("Process")
            if err != nil {
                continue
            }
            info := PluginInfo{File: f.Name(), Name: strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))}
            if proc, ok := sym.(func([]float32) []float32); ok {
                e.mu.Lock()
                e.plugins = append(e.plugins, PluginFunc(proc))
                e.pluginFiles = append(e.pluginFiles, f.Name())
                e.pluginInfos = append(e.pluginInfos, info)
                e.mu.Unlock()
            }
            if metaSym, err := p.Lookup("Metadata"); err == nil {
                if metaFunc, ok := metaSym.(func() map[string]string); ok {
                    metadata := metaFunc()
                    if metadata != nil {
                        if v, ok := metadata["name"]; ok {
                            info.Name = v
                        }
                        if v, ok := metadata["description"]; ok {
                            info.Description = v
                        }
                        if v, ok := metadata["version"]; ok {
                            info.Version = v
                        }
                        if v, ok := metadata["author"]; ok {
                            info.Author = v
                        }
                        e.mu.Lock()
                        if len(e.pluginInfos) > 0 {
                            e.pluginInfos[len(e.pluginInfos)-1] = info
                        }
                        e.mu.Unlock()
                    }
                }
            }
        }
    }
}

func (e *Engine) AddPluginFunc(proc func([]float32) []float32) {
    e.mu.Lock()
    e.plugins = append(e.plugins, PluginFunc(proc))
    e.mu.Unlock()
}

func (e *Engine) applyPlugins(samples []float32) []float32 {
    e.mu.Lock()
    plugins := append([]Plugin(nil), e.plugins...)
    e.mu.Unlock()
    for _, plugin := range plugins {
        samples = plugin.Process(samples)
    }
    return samples
}

func (e *Engine) render() {
    if e.sampleRate <= 0 {
        return
    }
    buf := make([]byte, 4096)
    for {
        e.mu.Lock()
        if !e.running {
            e.mu.Unlock()
            return
        }
        bpm := e.bpm
        stepIndex := e.stepIndex
        samplesSince := e.samplesSince
        channels := append([]Channel(nil), e.channels...)
        channelPhases := append([]float64(nil), e.channelPhases...)
        patterns := make([]Pattern, len(e.patterns))
        copy(patterns, e.patterns)
        currentPatternIndex := e.currentPattern
        pianoRoll := e.pianoRoll
        pianoPhases := append([]float64(nil), e.pianoPhases...)
        recording := e.recording
        recordFile := e.recordFile
        pw := e.pw
        e.mu.Unlock()
        if bpm <= 0 {
            bpm = 120
        }
        stepSamples := int(e.sampleRate * 60.0 / float64(bpm) / 4.0)
        if stepSamples <= 0 {
            stepSamples = 1
        }
        if len(patterns) == 0 {
            patterns = []Pattern{newPattern("Pattern 1", len(channels))}
        }
        if currentPatternIndex < 0 || currentPatternIndex >= len(patterns) {
            currentPatternIndex = 0
        }
        currentPattern := patterns[currentPatternIndex]
        samples := make([]float32, len(buf)/4)
        for i := range samples {
            sample := float32(0)
            for ci := range channels {
                if ci < len(currentPattern.Steps) && currentPattern.Steps[ci][stepIndex] {
                    ch := channels[ci]
                    if !ch.Muted {
                        sample += float32(math.Sin(2*math.Pi*ch.Freq*channelPhases[ci]/e.sampleRate)) * ch.Volume * 0.2
                        channelPhases[ci]++
                    }
                }
            }
            for ki := range pianoRoll.Notes {
                if stepIndex < len(pianoRoll.Notes[ki]) && pianoRoll.Notes[ki][stepIndex] {
                    freq := pianoKeyFrequencies[ki%len(pianoKeyFrequencies)]
                    sample += float32(math.Sin(2*math.Pi*freq*pianoPhases[ki]/e.sampleRate)) * 0.15
                    pianoPhases[ki]++
                }
            }
            if sample > 1 {
                sample = 1
            }
            if sample < -1 {
                sample = -1
            }
            samples[i] = sample
            samplesSince++
            if samplesSince >= stepSamples {
                samplesSince = 0
                stepIndex = (stepIndex + 1) % stepCount
            }
        }
        samples = e.applyPlugins(samples)
        for i := 0; i < len(buf); i += 4 {
            s := int16(samples[i/4] * 32767)
            buf[i] = byte(s)
            buf[i+1] = byte(s >> 8)
            buf[i+2] = byte(s)
            buf[i+3] = byte(s >> 8)
        }
        if pw != nil {
            _, _ = pw.Write(buf)
        }
        if recording && recordFile != nil {
            _, _ = recordFile.Write(buf)
            e.mu.Lock()
            e.recordBytes += len(buf)
            e.mu.Unlock()
        }
        e.mu.Lock()
        e.stepIndex = stepIndex
        e.samplesSince = samplesSince
        e.channelPhases = channelPhases
        e.pianoPhases = pianoPhases
        e.mu.Unlock()
        time.Sleep(10 * time.Millisecond)
    }
}

func writeWAVHeader(file *os.File, dataBytes int) error {
    header := make([]byte, 44)
    copy(header[0:], []byte("RIFF"))
    binary.LittleEndian.PutUint32(header[4:], uint32(36+dataBytes))
    copy(header[8:], []byte("WAVE"))
    copy(header[12:], []byte("fmt "))
    binary.LittleEndian.PutUint32(header[16:], 16)
    binary.LittleEndian.PutUint16(header[20:], 1)
    binary.LittleEndian.PutUint16(header[22:], 2)
    binary.LittleEndian.PutUint32(header[24:], 44100)
    binary.LittleEndian.PutUint32(header[28:], 44100*2*2)
    binary.LittleEndian.PutUint16(header[32:], 4)
    binary.LittleEndian.PutUint16(header[34:], 16)
    copy(header[36:], []byte("data"))
    binary.LittleEndian.PutUint32(header[40:], uint32(dataBytes))
    _, err := file.Seek(0, 0)
    if err != nil {
        return err
    }
    _, err = file.Write(header)
    return err
}
