# Daemon UI - Quick Start Guide

## Overview

The Daemon UI is a professional FL Studio-inspired Digital Audio Workstation (DAW) interface built entirely in Delphi/ObjectPascal. It features:

- **Daemon Logo** in the top-left corner (red chili pepper)
- **FL Studio-style dark interface** with cyan accents
- **Playlist panel** on the left for track management
- **16x8 Step Sequencer** in the center
- **8-channel mixer** on the right
- **Professional transport controls** (Play, Stop, Record)

## Layout

```
┌─────────────────────────────────────────────────────────────────────────┐
│ [🌶 Logo] DAEMON - FL Studio-style DAW  │ ▶ │ ⏹ │ ● │ [Tempo] [BPM] │
├──────────────┬───────────────────────────────────────┬──────────────────┤
│              │                                       │                  │
│  PLAYLIST    │          STEP SEQUENCER (16x8)       │  8 CHANNELS      │
│              │                                       │                  │
│ + Add Track  │    . . . . . . . . . . . . . . . .   │  CH 1 [===]      │
│ Track 1      │    . . . . . . . . . . . . . . . .   │                  │
│ Track 2      │    . . . . . . . . . . . . . . . .   │  CH 2 [===]      │
│ Track 3      │    . . . . . . . . . . . . . . . .   │                  │
│ Track 4      │    . . . . . . . . . . . . . . . .   │  CH 3 [===]      │
│              │    . . . . . . . . . . . . . . . .   │                  │
│              │    . . . . . . . . . . . . . . . .   │  ... (8 total)   │
│              │    . . . . . . . . . . . . . . . .   │                  │
│              │                                       │                  │
└──────────────┴───────────────────────────────────────┴──────────────────┘
```

## Installation & Compilation

### Windows (Primary Target)

1. **Install Delphi 10.4+** with FireMonkey support
2. **Navigate to the UI directory**:
   ```
   cd ui/
   ```
3. **Open in Delphi IDE** or compile from command line:
   ```
   dcc32 DaemonUI.dpr
   ```

### Running the Application

After compilation, execute:
```
DaemonUI.exe
```

## Menu Structure

### File Menu
- **New** - Create a new project
- **Open** - Load an existing project
- **Save** - Save current project
- **Exit** - Close application

### Edit Menu
- **Undo** - Undo last action
- **Redo** - Redo last action

### Help Menu
- **About** - Show application information

## Main Controls

### Transport Bar (Top Toolbar)
- **▶ Play** - Start playback (toggles to Pause when playing)
- **⏹ Stop** - Stop playback
- **● Record** - Toggle recording mode (turns red when active)

### Playlist Panel (Left Side)
- **Playlist label** - Shows "PLAYLIST"
- **+ Add Track** - Create new track
- **Track list** - Scrollable list of tracks
- Click to select tracks for editing

### Step Sequencer (Center)
- **16x8 Grid** - 16 steps × 8 channels
- **Click steps** to toggle on/off
- **Cyan color** indicates active steps
- **Dark gray** indicates inactive steps
- Real-time visualization during playback

### Channel Mixer (Right Side)
- **8 Channels** - One column per mixer channel
- **Volume slider** - Vertical fader for each channel
- **MUTE button** - Toggles mute on/off (turns red when muted)
- **CH# label** - Channel identifier

## Color Reference

| Element | Color | Hex Code |
|---------|-------|----------|
| Background | Dark Gray | #212121 |
| Panels | Lighter Gray | #303030 |
| Accent/Active | Cyan | #00BFFF |
| Alert/Record | Red | #FF4444 |
| Text | White | #FFFFFF |

## File Structure

```
daemon/
├── ui/
│   ├── DaemonUI.dpr              # Main project file
│   ├── DaemonUI.dproj            # Delphi project config
│   ├── Forms/
│   │   ├── MainForm.pas          # Main window
│   │   └── PlaylistForm.pas       # Playlist manager
│   ├── Components/
│   │   ├── SequencerPanel.pas     # Step sequencer
│   │   └── ChannelPanel.pas       # Channel control
│   ├── Utils/
│   │   └── StyleManager.pas       # Theme manager
│   └── README.md
├── images/
│   ├── daemonlogo01.PNG
│   └── daemonlogo02.png           # Used in UI
└── [other project files]
```

## Integration with Audio Engine

The UI components are designed to be easily connected to the Daemon audio engine:

### Event Handlers Available:

1. **Transport Controls**
   - `btnPlayClick()` - Connect to engine.Play()
   - `btnStopClick()` - Connect to engine.Stop()
   - `btnRecordClick()` - Connect to engine.Record()

2. **Sequencer Events**
   - `TSequencerPanel.OnStepToggle` - Receives (Channel, Step, Active)
   - `TSequencerPanel.AdvanceSequencer()` - Call from playback timer

3. **Channel Mixer Events**
   - `TChannelPanel.OnVolumeChange` - Receives (Channel, Volume 0.0-1.0)
   - `TChannelPanel.SetVolume()` - Update slider position

4. **Playlist Events**
   - `lstPlaylistClick()` - Receives selected track
   - `btnNewTrackClick()` - Create new track

## Customization

### Changing Colors

Edit `Utils/StyleManager.pas`:
```pascal
procedure TStyleManager.InitializeFLStudioTheme;
begin
  FTheme.DarkColor := $FF212121;      // Change this
  FTheme.DarkLightColor := $FF303030;
  FTheme.AccentColor := $FF00BFFF;
  FTheme.RedColor := $FFFF4444;
  FTheme.TextColor := $FFFFFFFF;
end;
```

### Adding Controls

1. Declare component in `MainForm.pas`
2. Create instance in `FormCreate()`
3. Set Parent, Alignment, and other properties
4. Add to initialization method

### Connecting to Audio Engine

In `MainForm.pas`, modify event handlers to call engine methods:
```pascal
procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  FIsPlaying := not FIsPlaying;
  UpdatePlayButtonState;
  FEngine.SetPlaying(FIsPlaying);  // Connect to engine
end;
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Logo not displaying | Verify image path relative to executable |
| Compilation errors | Ensure FireMonkey package is installed |
| Dark theme not applying | Check `ApplyFLStudioTheme()` is called in `FormShow()` |
| Steps not clickable | Verify `OnStepClick` handler is assigned |

## Performance Notes

- Step sequencer uses lightweight rectangle rendering
- 8×16 grid = 128 clickable elements (minimal memory footprint)
- Smooth animation at 60 FPS on modern hardware
- Optimized for real-time audio sync

## Next Steps

1. ✅ UI Framework - **COMPLETE**
2. ⏳ Connect to Daemon audio engine
3. ⏳ Implement MIDI input handling
4. ⏳ Add effects rack panel
5. ⏳ Piano roll editor window
6. ⏳ Sample browser

## Support

For issues or features, ensure you have:
- Delphi 10.4 or later installed
- FireMonkey framework enabled
- Windows 7 SP1 or later (for deployment)
