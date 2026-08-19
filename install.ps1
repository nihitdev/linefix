$ErrorActionPreference = "Stop"

$repo = "nihitdev/linefix"
$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
switch ($architecture) {
    "x64" { $arch = "amd64" }
    "arm64" { $arch = "arm64" }
    default { throw "linefix: unsupported Windows architecture: $architecture" }
}

$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name.TrimStart("v")
$archive = "linefix_${version}_windows_${arch}.zip"
$baseUrl = "https://github.com/$repo/releases/download/$($release.tag_name)"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("linefix-" + [guid]::NewGuid())
$installDir = if ($env:LINEFIX_INSTALL_DIR) { $env:LINEFIX_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\linefix\bin" }

try {
    New-Item -ItemType Directory -Path $tempDir | Out-Null
    Invoke-WebRequest "$baseUrl/$archive" -OutFile (Join-Path $tempDir $archive)
    Invoke-WebRequest "$baseUrl/SHA256SUMS" -OutFile (Join-Path $tempDir "SHA256SUMS")
    $checksumLine = Get-Content (Join-Path $tempDir "SHA256SUMS") | Where-Object { $_ -match "\s$([regex]::Escape($archive))$" }
    if (-not $checksumLine) { throw "linefix: checksum not found for $archive" }
    $expected = ($checksumLine -split "\s+")[0].ToLowerInvariant()
    $actual = (Get-FileHash (Join-Path $tempDir $archive) -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $expected) { throw "linefix: checksum verification failed" }

    Expand-Archive (Join-Path $tempDir $archive) -DestinationPath $tempDir -Force
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Copy-Item (Join-Path $tempDir "linefix.exe") (Join-Path $installDir "linefix.exe") -Force

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $parts = @($userPath -split ";" | Where-Object { $_ })
    if ($parts -notcontains $installDir) {
        $newPath = (@($parts) + $installDir) -join ";"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = "$env:Path;$installDir"
        Write-Host "Added $installDir to your user PATH. Open a new terminal to use it."
    }
    Write-Host "linefix $version installed to $installDir\linefix.exe"
}
finally {
    Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}
