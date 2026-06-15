program DaemonUI;

uses
  System.StartUpCopy,
  FMX.Forms,
  MainForm in 'Forms\MainForm.pas' {frmMain},
  PlaylistForm in 'Forms\PlaylistForm.pas',
  SequencerPanel in 'Components\SequencerPanel.pas',
  ChannelPanel in 'Components\ChannelPanel.pas',
  StyleManager in 'Utils\StyleManager.pas';

{$R *.res}

begin
  Application.Initialize;
  Application.CreateForm(TfrmMain, frmMain);
  Application.Run;
end.
