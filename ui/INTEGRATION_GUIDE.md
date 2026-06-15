DAEMON UI - INTEGRATION ARCHITECTURE
====================================

This document describes how to connect the Delphi UI to the Daemon audio engine.

## ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────┐
│         DAEMON UI (Delphi/FireMonkey)                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  MainForm.pas - Main UI Window                      │   │
│  │  - Transport Controls (Play/Stop/Record)            │   │
│  │  - Playlist Manager                                 │   │
│  │  - Step Sequencer Display                           │   │
│  │  - Channel Mixer Display                            │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ Communication Layer
                       │ (DLL/TCP/IPC)
                       │
┌──────────────────────▼──────────────────────────────────────┐
│         DAEMON ENGINE (Go/C++)                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  core/engine.go - Audio Processing                 │   │
│  │  - Sequencer State                                  │   │
│  │  - Channel Mixer                                    │   │
│  │  - Playback Control                                 │   │
│  │  - Audio Output                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  cpp/engine.cpp - Audio Synthesis                   │   │
│  │  - Tone Generation                                  │   │
│  │  - DSP Effects                                      │   │
│  │  - MIDI Processing                                 │   │
│  └─────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

## COMMUNICATION INTERFACE

### Option 1: DLL Interface (Recommended for Windows)

**Pros:**
- No IPC overhead
- Direct memory access
- Simplest implementation

**Cons:**
- Platform-specific
- Must match architecture (32/64-bit)

**Implementation:**

1. **Create DLL in C++/Go:**

```cpp
// daemon_api.dll
extern "C" {
    __declspec(dllexport) void DaemonEngine_Init();
    __declspec(dllexport) void DaemonEngine_Play();
    __declspec(dllexport) void DaemonEngine_Stop();
    __declspec(dllexport) void DaemonEngine_SetChannelVolume(int channel, float volume);
    __declspec(dllexport) void DaemonEngine_SetStep(int channel, int step, bool active);
    __declspec(dllexport) void DaemonEngine_Shutdown();
}
```

2. **Load in Delphi:**

```pascal
// In MainForm.pas
type
  PDaemonEngine = pointer;

function DaemonEngine_Init: PDaemonEngine; external 'daemon_api.dll';
procedure DaemonEngine_Play(engine: PDaemonEngine); external 'daemon_api.dll';
procedure DaemonEngine_Stop(engine: PDaemonEngine); external 'daemon_api.dll';
procedure DaemonEngine_SetChannelVolume(engine: PDaemonEngine; 
                                       channel: Integer; 
                                       volume: Single); external 'daemon_api.dll';

var
  FEngine: PDaemonEngine;

procedure TfrmMain.FormCreate(Sender: TObject);
begin
  FEngine := DaemonEngine_Init();
end;

procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  DaemonEngine_Play(FEngine);
end;
```

### Option 2: TCP/IP Network (For Multi-Process)

**Pros:**
- Cross-platform
- Engine can run on different machine
- Supports remote DAW control

**Cons:**
- Network latency
- More complex debugging

**Implementation:**

```pascal
procedure TfrmMain.ConnectToEngine;
var
  LClient: TIdTCPClient;
begin
  LClient := TIdTCPClient.Create;
  LClient.Host := 'localhost';
  LClient.Port := 9090;
  LClient.Connect;
  FEngineConnection := LClient;
end;

procedure TfrmMain.SendCommand(ACommand: string);
begin
  FEngineConnection.IOHandler.WriteLn(ACommand);
end;

procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  SendCommand('PLAY');
end;
```

### Option 3: Named Pipes (Windows IPC)

**Pros:**
- Low latency
- Windows-specific optimization
- Good for same-machine communication

**Cons:**
- Windows only
- More complex implementation

## COMPONENT INTEGRATION

### Transport Controls Integration

**UI Component:**
```pascal
procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  FIsPlaying := not FIsPlaying;
  UpdatePlayButtonState;
  // TODO: Call engine
end;
```

**Integration Code:**
```pascal
procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  FIsPlaying := not FIsPlaying;
  UpdatePlayButtonState;
  
  if FIsPlaying then
    DaemonEngine_Play(FEngine)
  else
    DaemonEngine_Pause(FEngine);
end;
```

### Sequencer Integration

**UI Component:**
```pascal
// SequencerPanel.pas
procedure TSequencerPanel.OnStepClick(Sender: TObject);
var
  LChannel, LStep: Integer;
begin
  LChannel := TRectangle(Sender).Tag div 16;
  LStep := TRectangle(Sender).Tag mod 16;
  ToggleStep(LChannel, LStep);
  
  // Callback to engine
  if Assigned(FOnStepToggle) then
    FOnStepToggle(LChannel, LStep, FStepStates[LChannel, LStep]);
end;
```

**Connection in MainForm:**
```pascal
procedure TfrmMain.SetupSequencer;
begin
  pnlSequencer.OnStepToggle := SequencerStepChanged;
end;

procedure TfrmMain.SequencerStepChanged(AChannel, AStep: Integer; AActive: Boolean);
begin
  DaemonEngine_SetStep(FEngine, AChannel, AStep, AActive);
end;
```

### Channel Mixer Integration

**UI Component:**
```pascal
// ChannelPanel.pas
procedure TChannelPanel.OnVolumeChange(Sender: TObject);
begin
  if Assigned(FOnVolumeChange) then
    FOnVolumeChange(FChannelNumber, GetVolume);
end;
```

**Connection in MainForm:**
```pascal
procedure TfrmMain.SetupChannels;
var
  i: Integer;
  LChannel: TChannelPanel;
begin
  for i := 1 to 8 do
  begin
    LChannel := TChannelPanel.Create(sbxChannels);
    LChannel.SetChannelNumber(i);
    LChannel.OnVolumeChange := ChannelVolumeChanged;
  end;
end;

procedure TfrmMain.ChannelVolumeChanged(AChannel: Integer; AVolume: Single);
begin
  DaemonEngine_SetChannelVolume(FEngine, AChannel - 1, AVolume);
end;
```

### Playlist Integration

**UI Component:**
```pascal
procedure TfrmMain.lstPlaylistClick(Sender: TObject);
begin
  if lstPlaylist.ItemIndex >= 0 then
  begin
    // Load selected track
  end;
end;
```

**Connection Code:**
```pascal
procedure TfrmMain.lstPlaylistClick(Sender: TObject);
var
  LTrackIndex: Integer;
begin
  LTrackIndex := lstPlaylist.ItemIndex;
  if LTrackIndex >= 0 then
    DaemonEngine_LoadTrack(FEngine, LTrackIndex);
end;
```

## REAL-TIME FEEDBACK

### Playback Position Sync

**Engine broadcasts current step to UI:**

```pascal
// Timer in MainForm to poll/receive updates
procedure TfrmMain.OnPlaybackTimer(Sender: TObject);
var
  LCurrentStep: Integer;
begin
  LCurrentStep := DaemonEngine_GetCurrentStep(FEngine);
  pnlStepGrid.UpdateCurrentStep(LCurrentStep);
end;
```

### Status Updates

**Listen for engine state changes:**

```pascal
procedure TfrmMain.OnEngineStatusChange(AStatus: string);
begin
  case AStatus of
    'PLAYING': lblStatus.Text := 'Playing...';
    'STOPPED': lblStatus.Text := 'Stopped';
    'RECORDING': lblStatus.Text := 'Recording...';
  end;
end;
```

## THREAD SAFETY

**Audio engines often use worker threads. Ensure thread safety:**

```pascal
procedure TfrmMain.ChannelVolumeChanged(AChannel: Integer; AVolume: Single);
begin
  // Use TThread.Synchronize to send from UI thread to engine
  TThread.Synchronize(nil, procedure
  begin
    DaemonEngine_SetChannelVolume(FEngine, AChannel, AVolume);
  end);
end;
```

## EVENT MODEL

### Callback-Based (Push)

Engine calls UI when state changes:

```pascal
type
  TEngineCallback = procedure(AEvent: string; AData: Pointer) of object;

procedure TfrmMain.EngineEventCallback(AEvent: string; AData: Pointer);
begin
  case AEvent of
    'STEP_ADVANCE': UpdateSequencerDisplay(Integer(AData));
    'PLAYBACK_STARTED': lblStatus.Text := 'Playing...';
    'PLAYBACK_STOPPED': lblStatus.Text := 'Stopped';
  end;
end;

procedure TfrmMain.FormCreate(Sender: TObject);
begin
  FEngine := DaemonEngine_Init();
  DaemonEngine_SetCallback(FEngine, EngineEventCallback);
end;
```

### Polling-Based (Pull)

UI polls engine for current state:

```pascal
procedure TfrmMain.OnUpdateTimer(Sender: TObject);
var
  LState: TEngineState;
begin
  LState := DaemonEngine_GetState(FEngine);
  UpdateUIFromState(LState);
end;
```

## GO/C++ ENGINE API SPECIFICATION

Required functions the engine must export:

```go
// core/engine.go

type EngineHandle struct {
    // private fields
}

func NewEngine() *EngineHandle { }
func (e *EngineHandle) Play() { }
func (e *EngineHandle) Stop() { }
func (e *EngineHandle) SetChannelVolume(ch int, vol float32) { }
func (e *EngineHandle) SetStep(ch, step int, active bool) { }
func (e *EngineHandle) GetCurrentStep() int { }
func (e *EngineHandle) GetState() EngineState { }
func (e *EngineHandle) Shutdown() { }

type EngineState struct {
    IsPlaying   bool
    IsRecording bool
    CurrentStep int
    Tempo       int
    Tracks      []TrackInfo
}
```

## TESTING CHECKLIST

- [ ] UI starts without engine (shows placeholder)
- [ ] Play button triggers engine.Play()
- [ ] Step clicks update engine state
- [ ] Volume changes applied to engine
- [ ] Playlist selection loads track
- [ ] Real-time position updates in sequencer
- [ ] Record button triggers recording
- [ ] Menu operations work correctly

## DEPLOYMENT CONSIDERATIONS

### Single Executable

UI and Engine in same process:
- Simpler distribution
- Shared memory access
- Potential UI blocking

### Separate Processes

UI and Engine as separate executables:
- Independent updates
- Better stability (UI crash doesn't stop audio)
- Network/IPC communication required

## FUTURE ENHANCEMENTS

- [ ] MIDI controller support
- [ ] VST plugin hosting
- [ ] Audio file import
- [ ] Pattern sequencer
- [ ] Arpeggiator
- [ ] Effect rack UI
- [ ] Piano roll editor
