package core

import (
    "math"
    "os"
    "path/filepath"
    "plugin"
    "sync"
    "time"

    "github.com/hajimehoshi/oto/v2"
)

type Engine struct {
    ctx          *oto.Context
    p            *oto.Player
    running      bool
    mu           sync.Mutex
    plugins      []Plugin
    channels     []Channel
    bpm          int
    stepIndex    int
    sampleRate   float64
    samplesSince int
}

type Channel struct {
    Name   string
    Steps  [16]bool
    Volume float32
    Freq   float64
    Phase  float64
}

type Plugin interface {
    Process([]float32) []float32
}

func NewEngine() *Engine {
    ctx, _ := oto.NewContext(44100, 2, 2, 32768)
    e := &Engine{ctx: ctx, bpm: 120, sampleRate: 44100}
    for i := 0; i < 6; i++ {
        c := Channel{Name: "Chan" + string('1'+i), Volume: 0.8, Freq: 220.0 + float64(i)*110}
        e.channels = append(e.channels, c)
    }
    return e
}

func (e *Engine) AddChannel(name string, freq float64) {
    e.mu.Lock()
    e.channels = append(e.channels, Channel{Name: name, Volume: 0.8, Freq: freq})
    e.mu.Unlock()
}

func (e *Engine) ToggleStep(ch, step int) {
    e.mu.Lock()
    if ch >= 0 && ch < len(e.channels) && step >= 0 && step < 16 {
        e.channels[ch].Steps[step] = !e.channels[ch].Steps[step]
    }
    e.mu.Unlock()
}

func (e *Engine) SetTempo(bpm int) {
    e.mu.Lock()
    e.bpm = bpm
    e.mu.Unlock()
}

func (e *Engine) SetChannelVol(ch int, v float32) {
    e.mu.Lock()
    if ch >= 0 && ch < len(e.channels) {
        e.channels[ch].Volume = v
    }
    e.mu.Unlock()
}

func (e *Engine) Play() {
    e.mu.Lock()
    if e.running {
        e.mu.Unlock()
        return
    }
    e.p, _ = e.ctx.NewPlayer()
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
    e.p.Close()
    e.mu.Unlock()
}

func (e *Engine) render() {
    sr := e.sampleRate
    buf := make([]byte, 4096)
    for {
        e.mu.Lock()
        if !e.running {
            e.mu.Unlock()
            return
        }
        bpm := e.bpm
        e.mu.Unlock()
        stepSamples := int(sr * 60.0 / float64(bpm) / 4.0)
        for i := 0; i < len(buf); i += 4 {
            sample := float32(0)
            e.mu.Lock()
            for ci := range e.channels {
                ch := &e.channels[ci]
                if ch.Steps[e.stepIndex] {
                    v := float32(math.Sin(2*math.Pi*ch.Freq*ch.Phase/sr)) * ch.Volume * 0.2
                    sample += v
                    ch.Phase += 1
                }
            }
            e.mu.Unlock()
            if sample > 1 { sample = 1 }
            if sample < -1 { sample = -1 }
            s := int16(sample * 32767)
            buf[i] = byte(s)
            buf[i+1] = byte(s >> 8)
            buf[i+2] = byte(s)
            buf[i+3] = byte(s >> 8)
            e.samplesSince++
            if e.samplesSince >= stepSamples {
                e.samplesSince = 0
                e.stepIndex = (e.stepIndex + 1) % 16
            }
        }
        e.p.Write(buf)
        time.Sleep(10 * time.Millisecond)
    }
}

func (e *Engine) Tick() {}

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
            if proc, ok := sym.(func([]float32) []float32); ok {
                e.plugins = append(e.plugins, PluginFunc(proc))
            }
        }
    }
}

type PluginFunc func([]float32) []float32

func (f PluginFunc) Process(b []float32) []float32 { return f(b) }
