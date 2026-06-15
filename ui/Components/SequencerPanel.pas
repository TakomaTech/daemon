unit SequencerPanel;

interface

uses
  System.SysUtils, System.Types, System.UITypes, System.Classes, System.Variants,
  FMX.Types, FMX.Controls, FMX.Objects, FMX.Graphics, FMX.Layouts;

type
  TStepEventHandler = procedure(AChannel, AStep: Integer; AActive: Boolean) of object;

  TSequencerPanel = class(TPanel)
  private
    FSteps: array[0..7, 0..15] of TRectangle;
    FStepStates: array[0..7, 0..15] of Boolean;
    FOnStepToggle: TStepEventHandler;
    FCurrentStep: Integer;
    procedure InitializeSteps;
    procedure OnStepClick(Sender: TObject);
    procedure UpdateStepColor(AChannel, AStep: Integer);
  public
    constructor Create(AOwner: TComponent); override;
    destructor Destroy; override;
    procedure SetStepActive(AChannel, AStep: Integer; AActive: Boolean);
    procedure ToggleStep(AChannel, AStep: Integer);
    procedure AdvanceSequencer;
    property OnStepToggle: TStepEventHandler read FOnStepToggle write FOnStepToggle;
    property CurrentStep: Integer read FCurrentStep write FCurrentStep;
  end;

implementation

const
  FL_STUDIO_DARK = $FF212121;
  FL_STUDIO_DARK_LIGHT = $FF303030;
  FL_STUDIO_ACCENT = $FF00BFFF;
  FL_STUDIO_RED = $FFFF4444;
  STEP_OFF_COLOR = $FF404040;
  STEP_ON_COLOR = FL_STUDIO_ACCENT;

constructor TSequencerPanel.Create(AOwner: TComponent);
begin
  inherited;
  Fill.Color := FL_STUDIO_DARK_LIGHT;
  FCurrentStep := 0;
  InitializeSteps;
end;

destructor TSequencerPanel.Destroy;
begin
  inherited;
end;

procedure TSequencerPanel.InitializeSteps;
var
  i, j: Integer;
  LRect: TRectangle;
  LStepSize: Single;
begin
  LStepSize := 30;
  
  for i := 0 to 7 do
  begin
    for j := 0 to 15 do
    begin
      LRect := TRectangle.Create(Self);
      LRect.Parent := Self;
      LRect.Width := LStepSize;
      LRect.Height := LStepSize;
      LRect.Position.X := j * (LStepSize + 2);
      LRect.Position.Y := i * (LStepSize + 2);
      LRect.Fill.Color := STEP_OFF_COLOR;
      LRect.Stroke.Color := $FF505050;
      LRect.Stroke.Thickness := 1;
      LRect.OnClick := OnStepClick;
      LRect.Tag := i * 16 + j; // Store channel and step in tag
      
      FSteps[i, j] := LRect;
      FStepStates[i, j] := False;
    end;
  end;
end;

procedure TSequencerPanel.OnStepClick(Sender: TObject);
var
  LChannel, LStep: Integer;
begin
  LChannel := TRectangle(Sender).Tag div 16;
  LStep := TRectangle(Sender).Tag mod 16;
  ToggleStep(LChannel, LStep);
end;

procedure TSequencerPanel.SetStepActive(AChannel, AStep: Integer; AActive: Boolean);
begin
  if (AChannel >= 0) and (AChannel < 8) and (AStep >= 0) and (AStep < 16) then
  begin
    FStepStates[AChannel, AStep] := AActive;
    UpdateStepColor(AChannel, AStep);
  end;
end;

procedure TSequencerPanel.ToggleStep(AChannel, AStep: Integer);
begin
  SetStepActive(AChannel, AStep, not FStepStates[AChannel, AStep]);
  if Assigned(FOnStepToggle) then
    FOnStepToggle(AChannel, AStep, FStepStates[AChannel, AStep]);
end;

procedure TSequencerPanel.UpdateStepColor(AChannel, AStep: Integer);
begin
  if FStepStates[AChannel, AStep] then
    FSteps[AChannel, AStep].Fill.Color := STEP_ON_COLOR
  else
    FSteps[AChannel, AStep].Fill.Color := STEP_OFF_COLOR;
end;

procedure TSequencerPanel.AdvanceSequencer;
var
  i: Integer;
begin
  // Visual feedback for current playing step
  for i := 0 to 7 do
  begin
    if FStepStates[i, FCurrentStep] then
      FSteps[i, FCurrentStep].Stroke.Color := FL_STUDIO_RED;
  end;
  
  FCurrentStep := (FCurrentStep + 1) mod 16;
end;

end.
