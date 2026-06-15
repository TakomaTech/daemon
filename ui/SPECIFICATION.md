DAEMON UI - COMPLETE SPECIFICATION
==================================

## PRODUCT OVERVIEW

**Name:** DAEMON - FL Studio-Inspired DAW UI
**Platform:** Windows 64-bit (Delphi/FireMonkey)
**Version:** 1.0
**Status:** Complete UI Framework (Awaiting Audio Engine Integration)

## VISUAL DESIGN SPECIFICATION

### Color Palette (FL Studio Dark Theme)

```
┌─────────────────────────────────────┐
│ Color              Hex Code   RGB    │
├─────────────────────────────────────┤
│ Background Dark    #212121   (33,33,33)   │
│ Background Light   #303030   (48,48,48)   │
│ Accent (Cyan)      #00BFFF   (0,191,255)  │
│ Alert (Red)        #FF4444   (255,68,68)  │
│ Text (White)       #FFFFFF   (255,255,255)│
│ Border             #404040   (64,64,64)   │
│ Grid Inactive      #404040   (64,64,64)   │
└─────────────────────────────────────┘
```

### Window Dimensions

**Minimum:** 1024 × 600 pixels
**Recommended:** 1280 × 800 pixels
**Maximum:** Unlimited (scales to screen)

### DPI Support

- 96 DPI (100%) - Default
- 120 DPI (125%) - Supported
- 144 DPI (150%) - Supported
- 192 DPI (200%) - Supported

## LAYOUT SPECIFICATION

### Main Window Breakdown

```
╔════════════════════════════════════════════════════════════════════════════╗
║                                TOP TOOLBAR (60px)                         ║
║  [Logo] Daemon - FL Studio DAW  |  ▶ Play  ⏹ Stop  ● Record              ║
╠═════════════════╦════════════════════════════════════════════════════════╦═╣
║                 ║                                                        ║ ║
║   PLAYLIST      ║         STEP SEQUENCER (16 Steps × 8 Channels)        ║ ║
║   Panel         ║                                                        ║ ║
║   (200px)       ║     . . . . . . . . . . . . . . . .                   ║ ║  CHANNEL
║                 ║     . . . . . . . . . . . . . . . .                   ║ ║  MIXER
║  [+ New Track]  ║     . . . . . . . . . . . . . . . .                   ║ ║  (160px)
║                 ║     . . . . . . . . . . . . . . . .                   ║ ║
║  Track 1        ║     . . . . . . . . . . . . . . . .                   ║ ║
║  Track 2        ║     . . . . . . . . . . . . . . . .                   ║ ║
║  Track 3        ║     . . . . . . . . . . . . . . . .                   ║ ║
║  Track 4        ║     . . . . . . . . . . . . . . . .                   ║ ║
║                 ║                                                        ║ ║
║  [Scrollbar]    ║                                                        ║ ║
║                 ║                                                        ║ ║
╚═════════════════╩════════════════════════════════════════════════════════╩═╝
```

### Component Specifications

#### 1. Toolbar (Top Panel)

**Height:** 60 pixels
**Background:** FL_STUDIO_DARK_LIGHT (#303030)
**Contents:**

- **Logo** (55×55 px)
  - Image: daemonlogo02.png
  - Position: Top-left
  - Margin: 5px all sides
  - Fit mode: Proportional

- **App Title** (350×45 px)
  - Text: "DAEMON - FL Studio-style DAW"
  - Font: Bold, 18pt
  - Color: FL_STUDIO_ACCENT (#00BFFF)
  - Vertical align: Center

- **Transport Controls** (240×50 px combined)
  - Play Button: "▶ Play" (80×50 px)
  - Stop Button: "⏹ Stop" (80×50 px)
  - Record Button: "● Record" (80×50 px)
  - Margin: 5px between buttons

#### 2. Playlist Panel (Left)

**Width:** 200 pixels
**Background:** FL_STUDIO_DARK (#212121)
**Border:** FL_STUDIO_DARK_LIGHT (#303030), 1px dashed

**Contents:**

- **Header Label** (200×30 px)
  - Text: "PLAYLIST"
  - Font: Bold, 14pt
  - Color: FL_STUDIO_ACCENT (#00BFFF)

- **Add Track Button** (190×40 px)
  - Text: "+ Add Track"
  - Color: FL_STUDIO_ACCENT (#00BFFF)
  - Margin: 5px all sides

- **Track List** (Scrollable)
  - Font: Regular, 11pt
  - Color: FL_STUDIO_TEXT (#FFFFFF)
  - Item height: 25px
  - Selection highlight: FL_STUDIO_ACCENT
  - Capacity: 8-50 tracks (scrollable)

#### 3. Step Sequencer (Center)

**Dimensions:** Flexible (fills remaining space)
**Background:** FL_STUDIO_DARK_LIGHT (#303030)

**Step Grid:**
- **Columns:** 16 (steps per pattern)
- **Rows:** 8 (channels)
- **Cell size:** 32×32 px
- **Padding:** 2px between cells

**Step Cell Behavior:**
- **Inactive:** FL_STUDIO_DARK (#212121)
- **Active:** FL_STUDIO_ACCENT (#00BFFF)
- **Playing:** FL_STUDIO_RED (#FF4444) - pulsing border
- **Border:** FL_STUDIO_DARK_LIGHT, 1px

**Interaction:**
- Click to toggle active/inactive
- Updates in real-time
- Sends event to engine

#### 4. Channel Mixer (Right)

**Width:** 160 pixels
**Background:** FL_STUDIO_DARK (#212121)
**Border:** FL_STUDIO_DARK_LIGHT (#303030), 1px dashed
**Scrollable:** Vertical

**Per-Channel Control (×8):**

- **Channel Label** (140×25 px)
  - Text: "CHANNEL #"
  - Font: Bold, 11pt
  - Color: FL_STUDIO_ACCENT (#00BFFF)

- **Volume Slider** (140×120 px)
  - Orientation: Vertical
  - Range: 0-100 (0.0-1.0 normalized)
  - Default: 80
  - Type: Smooth fader

- **Mute Button** (140×35 px)
  - Text: "MUTE" or "MUTED"
  - Normal Color: FL_STUDIO_ACCENT
  - Active Color: FL_STUDIO_RED
  - Font: Bold, 10pt

**Total Height:** 8 × 85px = 680px (scrollable if window < 680px)

## INTERACTION SPECIFICATION

### Transport Controls

#### Play Button
- **Default State:** "▶ Play" (text color: WHITE)
- **Active State:** "⏸ Pause" (text color: RED)
- **Click Action:**
  - Toggle between playing and paused
  - Update button text and color
  - Send PLAY command to engine
- **Keyboard:** Space bar (if focused)

#### Stop Button
- **State:** Always "⏹ Stop"
- **Color:** WHITE (always)
- **Click Action:**
  - Stop playback
  - Reset playback position to beginning
  - Reset play button to "▶ Play"
  - Send STOP command to engine
- **Keyboard:** Ctrl+Space

#### Record Button
- **Default State:** "● Record" (text color: WHITE)
- **Active State:** "● RECORDING" (text color: RED, pulsing)
- **Click Action:**
  - Toggle recording mode
  - Update button color
  - Send RECORD command to engine
- **Keyboard:** R key

### Sequencer Interaction

**Step Cell Click:**
- Toggle between active/inactive
- Change color immediately
- Send UPDATE command to engine with (Channel, Step, Active)
- Visual feedback: brief highlight animation

**During Playback:**
- Display current playing step with red border
- Advance to next step on beat
- Allow editing while playing (live mode)

### Playlist Interaction

**Track Selection:**
- Click track in list
- Load track data into sequencer view
- Highlight selected track with cyan background

**Add Track:**
- Click "+ Add Track" button
- Append new track to list
- Initialize with default pattern
- Auto-select new track

**Remove Track:**
- Right-click track → Delete
- Confirm removal dialog
- Update track numbering

### Mixer Interaction

**Volume Slider:**
- Drag vertically (up = louder, down = quieter)
- Real-time value display in tooltip
- Range: 0-100%
- Send VOLUME command to engine

**Mute Button:**
- Click to toggle mute state
- Change color (cyan→red)
- Mute audio output for channel
- Send MUTE command to engine

## MENU STRUCTURE

### File Menu
```
File
├─ New              (Ctrl+N)  → Create new project
├─ Open             (Ctrl+O)  → Load project file
├─ Save             (Ctrl+S)  → Save current project
├─ ─────────────────────────
└─ Exit             (Alt+F4)  → Close application
```

### Edit Menu
```
Edit
├─ Undo             (Ctrl+Z)  → Undo last action
└─ Redo             (Ctrl+Y)  → Redo last undone action
```

### Help Menu
```
Help
└─ About            (F1)      → Show version and credits
```

## STATE DIAGRAM

```
┌─────────────┐
│   STOPPED   │
└──────┬──────┘
       │ Play
       ▼
┌─────────────┐       Record
│   PLAYING   ├──────────┐
└──────┬──────┘          │
       │ Stop/Pause      ▼
       │          ┌──────────────┐
       │          │   RECORDING  │
       │          └──────────────┘
       │                 ▲
       └─────────────────┘
```

## KEYBOARD SHORTCUTS

| Shortcut | Action |
|----------|--------|
| Space | Play/Pause |
| Ctrl+Space | Stop |
| R | Toggle Record |
| Ctrl+N | New Project |
| Ctrl+O | Open Project |
| Ctrl+S | Save Project |
| Ctrl+Z | Undo |
| Ctrl+Y | Redo |
| Alt+F4 | Exit |
| F1 | Help/About |

## WINDOW BEHAVIOR

### Resizing
- Minimum size: 1024×600 pixels
- Maximum size: Unlimited
- Components scale proportionally
- Playlist width: fixed at 200px
- Mixer width: fixed at 160px
- Sequencer: expands/contracts with window

### Fullscreen
- Supported (F11 or Alt+Enter)
- All controls remain accessible
- Recommended for performances

### Restore from Minimized
- All UI state preserved
- Playback continues if playing
- No glitches or redraw issues

## ANIMATION SPECIFICATIONS

### Step Grid Animation
- **Inactive→Active:** 100ms ease-in
- **Active→Inactive:** 100ms ease-out
- **Current Step:** 300ms pulse (repeating during playback)

### Button Hover Effects
- Scale: 1.0 → 1.05
- Duration: 100ms ease-out
- Color brightening: +10% luminance

### Slider Interaction
- Smooth tracking while dragging
- Tooltip shows percentage during drag

## ACCESSIBILITY REQUIREMENTS

- [ ] High contrast mode support
- [ ] Keyboard navigation
- [ ] Screen reader compatibility
- [ ] Focus indicators visible
- [ ] Minimum 12pt font size
- [ ] Tab navigation order logical

## PERFORMANCE TARGETS

- **Startup time:** < 2 seconds
- **Click response:** < 50ms
- **FPS during idle:** 60 FPS (smooth animations)
- **CPU usage:** < 5% when idle
- **Memory footprint:** < 100MB baseline

## FILE FORMATS

### Project File (.daemon)
```json
{
  "version": "1.0",
  "name": "Project Name",
  "tempo": 120,
  "tracks": [
    {
      "name": "Track 1",
      "channel": 1,
      "volume": 0.8,
      "steps": [true, false, true, ...]
    }
  ]
}
```

## DEPLOYMENT SPECIFICATIONS

**Executable:** DaemonUI.exe (64-bit Windows)
**Runtime Requirements:**
- Windows 7 SP1 or later
- .NET Framework 4.5+ (for some components)
- 4GB RAM minimum

**Dependencies:**
- FireMonkey libraries (included in Delphi)
- Standard Delphi RTL

## TESTING REQUIREMENTS

### Unit Tests
- [ ] Sequencer state management
- [ ] Channel volume calculations
- [ ] Playlist operations
- [ ] File I/O operations

### Integration Tests
- [ ] UI responds to engine commands
- [ ] Engine receives UI commands
- [ ] Real-time sync between UI and engine

### UAT Tests
- [ ] All buttons clickable
- [ ] All menus functional
- [ ] Smooth animations
- [ ] No memory leaks (30min continuous use)
- [ ] Crash recovery (project auto-save)

## FUTURE FEATURES (v1.1+)

- [ ] Piano roll editor
- [ ] Drum machine step editor
- [ ] Effects rack with VST support
- [ ] Sample browser
- [ ] MIDI keyboard support
- [ ] BPM/Tempo visual feedback
- [ ] Master volume/limiter
- [ ] Project templates
- [ ] Undo/redo panel
- [ ] Theme customizer
- [ ] Keyboard mapping editor
- [ ] ASIO driver support
- [ ] Audio waveform preview
- [ ] Pattern sequencer
- [ ] Arpeggiator
- [ ] Swing/shuffle controls
