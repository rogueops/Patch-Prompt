@echo off
REM PatchPrompt cmd.exe helper.
REM cmd cannot run an external program on every prompt, so this applies a
REM static, ANSI-colored PROMPT. Call it from your cmd AutoRun (see
REM examples\cmd-autorun-example.cmd) or run it manually in a session.
for /f "usebackq delims=" %%P in (`patchprompt init cmd 2^>nul`) do %%P
