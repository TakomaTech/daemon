unit PlaylistForm;

interface

uses
  System.SysUtils, System.Types, System.UITypes, System.Classes, System.Variants,
  FMX.Types, FMX.Controls, FMX.Forms, FMX.Graphics, FMX.Dialogs, FMX.ListBox,
  FMX.Controls.Presentation, FMX.StdCtrls, FMX.Layouts, FMX.Objects;

type
  TPlaylistItem = class
  public
    Name: string;
    Duration: Integer; // in beats
    IsActive: Boolean;
  end;

  TfrmPlaylist = class(TForm)
    pnlHeader: TPanel;
    lblTitle: TLabel;
    lstPlaylist: TListBox;
    pnlFooter: TPanel;
    btnAdd: TButton;
    btnRemove: TButton;
    btnEdit: TButton;
  private
    FPlaylistItems: TList<TPlaylistItem>;
    procedure RefreshPlaylist;
    procedure AddPlaylistItem(const AName: string; ADuration: Integer);
    procedure RemovePlaylistItem(AIndex: Integer);
  public
    constructor Create(AOwner: TComponent); override;
    destructor Destroy; override;
  end;

var
  frmPlaylist: TfrmPlaylist;

implementation

{$R *.fmx}

const
  FL_STUDIO_DARK = $FF212121;
  FL_STUDIO_DARK_LIGHT = $FF303030;
  FL_STUDIO_ACCENT = $FF00BFFF;

constructor TfrmPlaylist.Create(AOwner: TComponent);
begin
  inherited;
  FPlaylistItems := TList<TPlaylistItem>.Create;
end;

destructor TfrmPlaylist.Destroy;
var
  i: Integer;
begin
  for i := 0 to FPlaylistItems.Count - 1 do
    FPlaylistItems[i].Free;
  FPlaylistItems.Free;
  inherited;
end;

procedure TfrmPlaylist.AddPlaylistItem(const AName: string; ADuration: Integer);
var
  LItem: TPlaylistItem;
begin
  LItem := TPlaylistItem.Create;
  LItem.Name := AName;
  LItem.Duration := ADuration;
  LItem.IsActive := False;
  FPlaylistItems.Add(LItem);
  RefreshPlaylist;
end;

procedure TfrmPlaylist.RemovePlaylistItem(AIndex: Integer);
begin
  if (AIndex >= 0) and (AIndex < FPlaylistItems.Count) then
  begin
    FPlaylistItems[AIndex].Free;
    FPlaylistItems.Delete(AIndex);
    RefreshPlaylist;
  end;
end;

procedure TfrmPlaylist.RefreshPlaylist;
var
  i: Integer;
begin
  lstPlaylist.Clear;
  for i := 0 to FPlaylistItems.Count - 1 do
  begin
    lstPlaylist.Items.Add(FPlaylistItems[i].Name + 
                         ' (' + IntToStr(FPlaylistItems[i].Duration) + ' beats)');
  end;
end;

end.
