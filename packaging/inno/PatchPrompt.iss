; PatchPrompt Inno Setup script - per-user install, no admin required.
; Build: ISCC.exe /DAppVersion=1.0.0 PatchPrompt.iss
; (build-release.ps1 invokes this automatically when Inno Setup is installed.)

#ifndef AppVersion
  #define AppVersion "1.0.0"
#endif

[Setup]
AppId={{B2F7C9E1-4A3D-4C2B-9E6A-PATCHPROMPT01}
AppName=PatchPrompt
AppVersion={#AppVersion}
AppPublisher=PatchPrompt
DefaultDirName={localappdata}\Programs\PatchPrompt
DefaultGroupName=PatchPrompt
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
OutputDir=..\..\dist
OutputBaseFilename=PatchPromptSetup-{#AppVersion}
Compression=lzma2
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
UninstallDisplayName=PatchPrompt

[Files]
; Binary is expected at dist\patchprompt.exe (produced by build-release.ps1).
Source: "..\..\dist\patchprompt.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\themes\*.json"; DestDir: "{app}\themes"; Flags: ignoreversion
Source: "..\..\README.md"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: "addtopath"; Description: "Add PatchPrompt to my user PATH"; GroupDescription: "Setup:"

[Registry]
; Add install dir to the per-user PATH when the task is selected.
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; \
  ValueData: "{olddata};{app}"; Check: NeedsAddPath('{app}'); Tasks: addtopath; \
  Flags: preservestringtype

[Run]
; Generate default config + themes into the user profile after install.
Filename: "{app}\patchprompt.exe"; Parameters: "config generate"; \
  Flags: runhidden nowait skipifsilent

[UninstallRun]
; Best-effort cleanup of the prompt is left to scripts\uninstall.ps1; the
; installer only removes files it placed.

[Code]
function NeedsAddPath(Param: string): Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', OrigPath) then
  begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + ExpandConstant(Param) + ';', ';' + OrigPath + ';') = 0;
end;
