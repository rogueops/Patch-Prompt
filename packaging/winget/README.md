# Winget Packaging for PatchPrompt

> winget approval is **not** assumed. These steps prepare and submit the package.
> Local install never requires winget (use `scripts\install.ps1`).

## Steps

1. **Build the release artifacts**

   ```powershell
   .\scripts\build-release.ps1
   ```

   This produces `dist\PatchPromptSetup-1.0.0.exe` (requires Inno Setup) plus
   `dist\patchprompt-windows-amd64.zip` and `dist\checksums.txt`.

2. **Create a GitHub Release** and upload `PatchPromptSetup-1.0.0.exe`.
   Note the public download URL of the asset.

3. **Generate the SHA256** of the installer:

   ```powershell
   (Get-FileHash .\dist\PatchPromptSetup-1.0.0.exe -Algorithm SHA256).Hash
   ```

4. **Use wingetcreate** to produce the manifest:

   ```powershell
   winget install wingetcreate
   wingetcreate new <installer-download-url>
   ```

   Fill in PackageIdentifier `PatchPrompt.PatchPrompt`, version `1.0.0`,
   InstallerType `inno`, and the SHA256 from step 3. See
   `PatchPrompt.PatchPrompt.yaml.example` for the field layout.

5. **Submit** to the community repo:

   ```powershell
   wingetcreate submit --token <github-token> .\manifests\...
   ```

   This opens a PR against [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs).

6. **After approval**, anyone can install with:

   ```powershell
   winget install PatchPrompt.PatchPrompt
   ```
