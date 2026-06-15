package core

type Channel struct {
    Name   string  `json:"name"`
    Volume float32 `json:"volume"`
    Freq   float64 `json:"freq"`
    Muted  bool    `json:"muted"`
    Phase  float64 `json:"-"`
}

type Pattern struct {
    Name  string     `json:"name"`
    Steps [][]bool   `json:"steps"`
}

type PianoRoll struct {
    Notes [][]bool `json:"notes"`
}

type Project struct {
    Name           string    `json:"name"`
    Tempo          int       `json:"tempo"`
    CurrentPattern int       `json:"current_pattern"`
    Channels       []Channel `json:"channels"`
    Patterns       []Pattern `json:"patterns"`
    PianoRoll      PianoRoll `json:"piano_roll"`
}
