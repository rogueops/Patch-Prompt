@echo off
REM Example cmd.exe AutoRun for PatchPrompt.
REM Register it once (per-user, no admin):
REM   reg add "HKCU\Software\Microsoft\Command Processor" /v AutoRun /d "%LOCALAPPDATA%\Programs\PatchPrompt\..\path\to\this.cmd" /f
REM Or simply paste the line below into your own AutoRun script.
REM
REM cmd cannot run a program on every prompt, so this sets a static colored PROMPT.
for /f "usebackq delims=" %%P in (`patchprompt init cmd 2^>nul`) do %%P
