unit ChannelPanel;

interface

uses
  System.SysUtils, System.Types, System.UITypes, System.Classes, System.Variants,
  FMX.Types, FMX.Controls, FMX.Objects, FMX.Graphics, FMX.Layouts,
  FMX.Controls.Presentation, FMX.StdCtrls;

type
  TChannelEventHandler = procedure(AChannel: Integer; AValue: Single) of object;

  TChannelPanel = class(TPanel)
  private
    FChannelNumber: Integer;
    FVolumeSlider: TTrackBar;
    FChannelLabel: TLabel;
    FMuteButton: TButton;
    FOnVolumeChange: TChannelEventHandler;
    procedure OnVolumeChange(Sender: TObject);
    procedure OnMuteClick(Sender: TObject);
  public
    constructor Create(AOwner: TComponent); override;
    procedure SetChannelNumber(ANumber: Integer);
    procedure SetVolume(AVolume: Single);
    function GetVolume: Single;
    property OnVolumeChange: TChannelEventHandler read FOnVolumeChange write FOnVolumeChange;
  end;

implementation

const
  FL_STUDIO_DARK = $FF212121;
  FL_STUDIO_DARK_LIGHT = $FF303030;
  FL_STUDIO_ACCENT = $FF00BFFF;
  FL_STUDIO_RED = $FFFF4444;

constructor TChannelPanel.Create(AOwner: TComponent);
begin
  inherited;
  Fill.Color := FL_STUDIO_DARK_LIGHT;
  Width := 140;
  Height := 200;
  
  // Channel label
  FChannelLabel := TLabel.Create(Self);
  FChannelLabel.Parent := Self;
  FChannelLabel.Text := 'CH 1';
  FChannelLabel.Height := 30;
  FChannelLabel.Align := TAlignLayout.Top;
  FChannelLabel.TextSettings.FontColor := FL_STUDIO_ACCENT;
  FChannelLabel.Font.Size := 12;
  FChannelLabel.Font.Style := [TFontStyle.fsBold];
  FChannelLabel.Padding.Top := 5;
  
  // Volume slider
  FVolumeSlider := TTrackBar.Create(Self);
  FVolumeSlider.Parent := Self;
  FVolumeSlider.Min := 0;
  FVolumeSlider.Max := 100;
  FVolumeSlider.Value := 80;
  FVolumeSlider.Orientation := TOrientation.Vertical;
  FVolumeSlider.Height := 120;
  FVolumeSlider.Align := TAlignLayout.Top;
  FVolumeSlider.OnChange := OnVolumeChange;
  FVolumeSlider.Margins.Top := 5;
  
  // Mute button
  FMuteButton := TButton.Create(Self);
  FMuteButton.Parent := Self;
  FMuteButton.Text := 'MUTE';
  FMuteButton.Height := 35;
  FMuteButton.Align := TAlignLayout.Top;
  FMuteButton.OnClick := OnMuteClick;
  FMuteButton.TextSettings.FontColor := FL_STUDIO_ACCENT;
  FMuteButton.Margins.Top := 5;
end;

procedure TChannelPanel.SetChannelNumber(ANumber: Integer);
begin
  FChannelNumber := ANumber;
  FChannelLabel.Text := 'CHANNEL ' + IntToStr(ANumber);
end;

procedure TChannelPanel.SetVolume(AVolume: Single);
begin
  FVolumeSlider.Value := AVolume * 100;
end;

function TChannelPanel.GetVolume: Single;
begin
  Result := FVolumeSlider.Value / 100;
end;

procedure TChannelPanel.OnVolumeChange(Sender: TObject);
begin
  if Assigned(FOnVolumeChange) then
    FOnVolumeChange(FChannelNumber, GetVolume);
end;

procedure TChannelPanel.OnMuteClick(Sender: TObject);
begin
  if FMuteButton.TextSettings.FontColor = FL_STUDIO_ACCENT then
  begin
    FMuteButton.TextSettings.FontColor := FL_STUDIO_RED;
    FMuteButton.Text := 'MUTED';
  end
  else
  begin
    FMuteButton.TextSettings.FontColor := FL_STUDIO_ACCENT;
    FMuteButton.Text := 'MUTE';
  end;
end;

end.
