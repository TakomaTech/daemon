unit StyleManager;

interface

uses
  System.SysUtils, FMX.Graphics;

type
  TFLStudioTheme = record
    DarkColor: TAlphaColor;
    DarkLightColor: TAlphaColor;
    AccentColor: TAlphaColor;
    RedColor: TAlphaColor;
    TextColor: TAlphaColor;
  end;

  TStyleManager = class
  private
    FTheme: TFLStudioTheme;
  public
    constructor Create;
    procedure InitializeFLStudioTheme;
    property Theme: TFLStudioTheme read FTheme;
  end;

implementation

constructor TStyleManager.Create;
begin
  InitializeFLStudioTheme;
end;

procedure TStyleManager.InitializeFLStudioTheme;
begin
  FTheme.DarkColor := $FF212121;
  FTheme.DarkLightColor := $FF303030;
  FTheme.AccentColor := $FF00BFFF;
  FTheme.RedColor := $FFFF4444;
  FTheme.TextColor := $FFFFFFFF;
end;

end.
