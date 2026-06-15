unit MainForm;

interface

uses
  System.SysUtils, System.Types, System.UITypes, System.Classes, System.Variants,
  FMX.Types, FMX.Controls, FMX.Forms, FMX.Graphics, FMX.Dialogs, FMX.Layouts,
  FMX.Objects, FMX.Controls.Presentation, FMX.StdCtrls, FMX.ListBox,
  FMX.ExtCtrls, FMX.Menus;

type
  TfrmMain = class(TForm)
    pnlMain: TPanel;
    pnlToolbar: TPanel;
    pnlPlaylist: TPanel;
    pnlSequencer: TPanel;
    pnlChannels: TPanel;
    imgLogo: TImage;
    lblAppName: TLabel;
    btnPlay: TButton;
    btnStop: TButton;
    btnRecord: TButton;
    lstPlaylist: TListBox;
    lblPlaylist: TLabel;
    btnNewTrack: TButton;
    pnlStepGrid: TPanel;
    sbxChannels: TScrollBox;
    pnlChannelContainer: TPanel;
    MainMenu: TMainMenu;
    mitFile: TMenuItem;
    mitNew: TMenuItem;
    mitOpen: TMenuItem;
    mitSave: TMenuItem;
    mitSep1: TMenuItem;
    mitExit: TMenuItem;
    mitEdit: TMenuItem;
    mitUndo: TMenuItem;
    mitRedo: TMenuItem;
    mitHelp: TMenuItem;
    mitAbout: TMenuItem;
    procedure FormCreate(Sender: TObject);
    procedure FormShow(Sender: TObject);
    procedure btnPlayClick(Sender: TObject);
    procedure btnStopClick(Sender: TObject);
    procedure btnRecordClick(Sender: TObject);
    procedure btnNewTrackClick(Sender: TObject);
    procedure lstPlaylistClick(Sender: TObject);
    procedure mitNewClick(Sender: TObject);
    procedure mitOpenClick(Sender: TObject);
    procedure mitSaveClick(Sender: TObject);
    procedure mitExitClick(Sender: TObject);
    procedure mitAboutClick(Sender: TObject);
  private
    FIsPlaying: Boolean;
    FIsRecording: Boolean;
    procedure InitializeUI;
    procedure SetupToolbar;
    procedure SetupPlaylist;
    procedure SetupSequencer;
    procedure SetupChannels;
    procedure CreateStepGrid;
    procedure CreateChannelControls;
    procedure UpdatePlayButtonState;
    procedure ApplyFLStudioTheme;
  public
    { Public declarations }
  end;

var
  frmMain: TfrmMain;

implementation

{$R *.fmx}
{$R *.LgXhdpi.fmx}

const
  FL_STUDIO_DARK = $FF212121;
  FL_STUDIO_DARK_LIGHT = $FF303030;
  FL_STUDIO_ACCENT = $FF00BFFF;
  FL_STUDIO_RED = $FFFF4444;
  FL_STUDIO_TEXT = $FFFFFFFF;

procedure TfrmMain.FormCreate(Sender: TObject);
begin
  FIsPlaying := False;
  FIsRecording := False;
  InitializeUI;
end;

procedure TfrmMain.FormShow(Sender: TObject);
begin
  ApplyFLStudioTheme;
end;

procedure TfrmMain.InitializeUI;
begin
  SetupToolbar;
  SetupPlaylist;
  SetupSequencer;
  SetupChannels;
end;

procedure TfrmMain.ApplyFLStudioTheme;
begin
  // Apply FL Studio dark theme
  pnlMain.Fill.Color := FL_STUDIO_DARK;
  pnlToolbar.Fill.Color := FL_STUDIO_DARK_LIGHT;
  pnlPlaylist.Fill.Color := FL_STUDIO_DARK;
  pnlSequencer.Fill.Color := FL_STUDIO_DARK_LIGHT;
  pnlChannels.Fill.Color := FL_STUDIO_DARK;
  
  lblAppName.TextSettings.FontColor := FL_STUDIO_ACCENT;
  lblPlaylist.TextSettings.FontColor := FL_STUDIO_TEXT;
  
  // Style buttons
  btnPlay.TextSettings.FontColor := FL_STUDIO_TEXT;
  btnStop.TextSettings.FontColor := FL_STUDIO_TEXT;
  btnRecord.TextSettings.FontColor := FL_STUDIO_TEXT;
  btnNewTrack.TextSettings.FontColor := FL_STUDIO_TEXT;
end;

procedure TfrmMain.SetupToolbar;
var
  LHeight: Single;
begin
  LHeight := 60;
  pnlToolbar.Height := LHeight;
  pnlToolbar.Align := TAlignLayout.Top;
  pnlToolbar.Fill.Color := FL_STUDIO_DARK_LIGHT;
  
  // Logo
  imgLogo.Height := LHeight - 10;
  imgLogo.Width := LHeight - 10;
  imgLogo.Position.X := 5;
  imgLogo.Position.Y := 5;
  imgLogo.Align := TAlignLayout.Left;
  imgLogo.WrapMode := TImageWrapMode.Fit;
  
  try
    imgLogo.Bitmap.LoadFromFile('..\..\..\images\daemonlogo02.png');
  except
    // Logo not found, will show placeholder
  end;
  
  // App name
  lblAppName.Text := 'DAEMON - FL Studio-style DAW';
  lblAppName.Font.Size := 18;
  lblAppName.Position.X := 70;
  lblAppName.Position.Y := 10;
  lblAppName.Align := TAlignLayout.Left;
  lblAppName.Width := 300;
  
  // Transport controls
  btnPlay.Text := '▶ Play';
  btnPlay.Width := 80;
  btnPlay.Position.X := 350;
  btnPlay.OnClick := btnPlayClick;
  
  btnStop.Text := '⏹ Stop';
  btnStop.Width := 80;
  btnStop.Position.X := 435;
  btnStop.OnClick := btnStopClick;
  
  btnRecord.Text := '● Record';
  btnRecord.Width := 80;
  btnRecord.Position.X := 520;
  btnRecord.OnClick := btnRecordClick;
end;

procedure TfrmMain.SetupPlaylist;
var
  i: Integer;
begin
  pnlPlaylist.Width := 200;
  pnlPlaylist.Align := TAlignLayout.Left;
  pnlPlaylist.Fill.Color := FL_STUDIO_DARK;
  
  lblPlaylist.Text := 'PLAYLIST';
  lblPlaylist.Font.Size := 14;
  lblPlaylist.Font.Style := [TFontStyle.fsBold];
  lblPlaylist.Height := 30;
  lblPlaylist.Align := TAlignLayout.Top;
  lblPlaylist.TextSettings.FontColor := FL_STUDIO_ACCENT;
  
  btnNewTrack.Text := '+ Add Track';
  btnNewTrack.Height := 40;
  btnNewTrack.Align := TAlignLayout.Top;
  btnNewTrack.OnClick := btnNewTrackClick;
  
  lstPlaylist.Align := TAlignLayout.Client;
  lstPlaylist.Fill.Color := FL_STUDIO_DARK_LIGHT;
  lstPlaylist.OnChange := lstPlaylistClick;
  
  // Add sample tracks
  for i := 1 to 8 do
    lstPlaylist.Items.Add('Track ' + IntToStr(i));
  
  if lstPlaylist.Count > 0 then
    lstPlaylist.ItemIndex := 0;
end;

procedure TfrmMain.SetupSequencer;
begin
  pnlSequencer.Align := TAlignLayout.Client;
  pnlSequencer.Fill.Color := FL_STUDIO_DARK_LIGHT;
  
  pnlStepGrid.Align := TAlignLayout.Top;
  pnlStepGrid.Height := 300;
  pnlStepGrid.Fill.Color := FL_STUDIO_DARK;
  
  CreateStepGrid;
end;

procedure TfrmMain.SetupChannels;
begin
  pnlChannels.Width := 150;
  pnlChannels.Align := TAlignLayout.Right;
  pnlChannels.Fill.Color := FL_STUDIO_DARK;
  
  sbxChannels.Align := TAlignLayout.Client;
  sbxChannels.ShowScrollBars := TScrollBars.sbVertical;
  
  pnlChannelContainer.Align := TAlignLayout.Top;
  pnlChannelContainer.Height := 0;
  
  CreateChannelControls;
end;

procedure TfrmMain.CreateStepGrid;
var
  i, j: Integer;
  LBtn: TButton;
  LPanel: TPanel;
  LStepsPerRow: Integer;
begin
  LStepsPerRow := 16;
  
  LPanel := TPanel.Create(pnlStepGrid);
  LPanel.Parent := pnlStepGrid;
  LPanel.Align := TAlignLayout.Client;
  LPanel.Fill.Color := FL_STUDIO_DARK_LIGHT;
  
  // Create 16x8 step grid
  for i := 0 to 7 do
  begin
    for j := 0 to 15 do
    begin
      LBtn := TButton.Create(LPanel);
      LBtn.Parent := LPanel;
      LBtn.Text := '.';
      LBtn.Width := 30;
      LBtn.Height := 30;
      LBtn.Position.X := j * 32;
      LBtn.Position.Y := i * 32;
      LBtn.Tag := i * LStepsPerRow + j;
      LBtn.OnClick := nil; // Will be connected to engine
    end;
  end;
end;

procedure TfrmMain.CreateChannelControls;
var
  i: Integer;
  LChannelPanel: TPanel;
  LLabel: TLabel;
  LSlider: TTrackBar;
begin
  // Create controls for 8 channels
  for i := 1 to 8 do
  begin
    LChannelPanel := TPanel.Create(pnlChannelContainer);
    LChannelPanel.Parent := pnlChannelContainer;
    LChannelPanel.Height := 80;
    LChannelPanel.Width := 140;
    LChannelPanel.Align := TAlignLayout.Top;
    LChannelPanel.Fill.Color := FL_STUDIO_DARK_LIGHT;
    LChannelPanel.Margins.Top := 5;
    LChannelPanel.Margins.Left := 5;
    LChannelPanel.Margins.Right := 5;
    
    // Channel label
    LLabel := TLabel.Create(LChannelPanel);
    LLabel.Parent := LChannelPanel;
    LLabel.Text := 'CH ' + IntToStr(i);
    LLabel.Height := 20;
    LLabel.Align := TAlignLayout.Top;
    LLabel.TextSettings.FontColor := FL_STUDIO_ACCENT;
    
    // Volume slider
    LSlider := TTrackBar.Create(LChannelPanel);
    LSlider.Parent := LChannelPanel;
    LSlider.Min := 0;
    LSlider.Max := 100;
    LSlider.Value := 80;
    LSlider.Align := TAlignLayout.Client;
    LSlider.Orientation := TOrientation.Vertical;
  end;
  
  pnlChannelContainer.Height := 8 * 85;
end;

procedure TfrmMain.UpdatePlayButtonState;
begin
  if FIsPlaying then
  begin
    btnPlay.Text := '⏸ Pause';
    btnPlay.TextSettings.FontColor := FL_STUDIO_RED;
  end
  else
  begin
    btnPlay.Text := '▶ Play';
    btnPlay.TextSettings.FontColor := FL_STUDIO_TEXT;
  end;
end;

procedure TfrmMain.btnPlayClick(Sender: TObject);
begin
  FIsPlaying := not FIsPlaying;
  UpdatePlayButtonState;
  // TODO: Connect to engine to start/pause playback
end;

procedure TfrmMain.btnStopClick(Sender: TObject);
begin
  FIsPlaying := False;
  UpdatePlayButtonState;
  // TODO: Connect to engine to stop playback
end;

procedure TfrmMain.btnRecordClick(Sender: TObject);
begin
  FIsRecording := not FIsRecording;
  if FIsRecording then
    btnRecord.TextSettings.FontColor := FL_STUDIO_RED
  else
    btnRecord.TextSettings.FontColor := FL_STUDIO_TEXT;
  // TODO: Connect to engine to start/stop recording
end;

procedure TfrmMain.btnNewTrackClick(Sender: TObject);
begin
  lstPlaylist.Items.Add('Track ' + IntToStr(lstPlaylist.Items.Count + 1));
  // TODO: Create new track in engine
end;

procedure TfrmMain.lstPlaylistClick(Sender: TObject);
begin
  // TODO: Load selected track
  if lstPlaylist.ItemIndex >= 0 then
    ShowMessage('Selected: ' + lstPlaylist.Items[lstPlaylist.ItemIndex]);
end;

procedure TfrmMain.mitNewClick(Sender: TObject);
begin
  // TODO: New project
  ShowMessage('New Project');
end;

procedure TfrmMain.mitOpenClick(Sender: TObject);
begin
  // TODO: Open project
  ShowMessage('Open Project');
end;

procedure TfrmMain.mitSaveClick(Sender: TObject);
begin
  // TODO: Save project
  ShowMessage('Project Saved');
end;

procedure TfrmMain.mitExitClick(Sender: TObject);
begin
  Close;
end;

procedure TfrmMain.mitAboutClick(Sender: TObject);
begin
  ShowMessage('DAEMON - FL Studio-style DAW' + sLineBreak +
              'Version 1.0' + sLineBreak +
              'A modern music production station');
end;

end.
