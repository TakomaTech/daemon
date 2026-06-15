# Daemon UI - Delphi/ObjectPascal

A professional FL Studio-inspired DAW interface built in pure Delphi/ObjectPascal using FireMonkey (FMX).

## Project Structure

```
ui/
├── DaemonUI.dpr              # Main project file
├── DaemonUI.dproj            # Project configuration
├── Forms/
│   ├── MainForm.pas          # Main application window
│   └── PlaylistForm.pas       # Playlist management
├── Components/
│   ├── SequencerPanel.pas     # Step sequencer component
│   └── ChannelPanel.pas       # Channel mixer component
├── Utils/
│   └── StyleManager.pas       # FL Studio theme management
└── README.md
```

## Features

### Main Window
- **Daemon Logo**: Top-left branding with logo from images/
- **Toolbar**: Transport controls (Play, Stop, Record)
- **Playlist Panel**: Left sidebar with playlist management
- **Sequencer Grid**: 16x8 step sequencer (FL Studio-style)
- **Channel Mixer**: Right sidebar with 8 channels
- **Menu Bar**: File, Edit, and Help menus

### FL Studio-Inspired Design
- Dark theme with cyan accents (#00BFFF)
- Professional color scheme matching FL Studio
- Responsive layout with proper panel organization
- Keyboard-friendly interface

### UI Components
- **SequencerPanel**: Interactive 16x8 step grid
- **ChannelPanel**: Per-channel volume control with mute button
- **PlaylistForm**: Track management and organization

## Building

### Requirements
- Delphi 10.4 or later
- FireMonkey framework
- Windows 64-bit target (configurable)

### Compilation
```bash
dcc32 DaemonUI.dpr -DDEBUG  # Debug build
dcc32 DaemonUI.dpr          # Release build
```

Or use the Delphi IDE:
1. Open `DaemonUI.dproj`
2. Select Debug/Release configuration
3. Press F9 to compile and run

## Color Scheme

```
Background (Dark):          #212121
Background (Light):         #303030
Accent (Cyan):             #00BFFF
Alert/Recording (Red):     #FF4444
Text (White):              #FFFFFF
```

## Integration Points

The UI is ready to be connected to the Daemon audio engine:

1. **Playback Control**: `btnPlayClick`, `btnStopClick` - Connect to engine playback
2. **Step Sequencer**: `TSequencerPanel.OnStepToggle` - Connect to step events
3. **Volume Control**: `TChannelPanel.OnVolumeChange` - Connect to channel mixer
4. **Playlist**: `TPlaylistForm` - Connect to project management

## Future Enhancements

- [ ] Piano roll editor window
- [ ] Plugin effects rack
- [ ] Sample browser
- [ ] Real-time waveform display
- [ ] Master volume control
- [ ] BPM/Tempo display and adjustment
- [ ] Project save/load dialogs
- [ ] MIDI controller support
- [ ] Customizable keyboard shortcuts

## Notes

- The UI uses FMX for cross-platform compatibility (Windows, macOS, iOS, Android)
- All theme colors are centralized in `StyleManager.pas`
- Custom components can be easily styled through the theme manager
- The sequencer supports dynamic step animation during playback
