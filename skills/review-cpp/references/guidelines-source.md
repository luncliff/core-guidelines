# Guidelines source

## Choose a runtime tool

Use host file-reading tools to inspect rules. For downloads, prefer [curl](https://curl.se/docs/manpage.html) when installed: it saves raw bytes, reports HTTP errors, and sets timeout and size limits without Python. Check availability first. Each host controls network and write permissions.

| Runtime | Download alternative | Read and validate |
|---|---|---|
| Windows | `curl.exe`; else PowerShell `Invoke-WebRequest` | PowerShell or host file tools |
| macOS or Linux | `curl`; else permitted host download tool | Shell utilities or host file tools |
| Agent host with fetch but no shell | Permitted tool only if it returns the complete raw file; save only if writes are supported | Host file tools; state unavailable checks |
| Offline or read-only host | Existing or user-provided document | Read anchors directly |

PowerShell `curl` can be an alias. On Windows, use `curl.exe`. With PowerShell only, set `$url` to the fixed URL below and `$newFile` to an approved new path:

```powershell
$response = Invoke-WebRequest -Uri $url -OutFile $newFile -PassThru -UseBasicParsing -TimeoutSec 30 -ErrorAction Stop
$expected = $response.Headers['Content-Length']
if ($expected -and (Get-Item -LiteralPath $newFile).Length -ne [long]$expected) {
    throw 'Incomplete download'
}
```

[Invoke-WebRequest](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/invoke-webrequest) saves to `-OutFile`; size validation happens after the download. Validate the result before selecting it. A fetch tool that truncates or transforms the file is not a substitute. Tool installation is a separate user decision, not part of the fallback.

## Select a local file

- The skill bundles the index and review references. The complete Guidelines document is supplied or downloaded separately.
- Use the user-supplied path first. Correct an invalid explicit file.
- First, attempt an approved refresh. If unavailable or failed, inspect the preferred OS cache, legacy cache, then `CppCoreGuidelines.md` in the working directory. Validate each candidate before use. Record rejected candidates.
- The preferred cache is `review-cpp/CppCoreGuidelines.md`; also check the legacy `review-cpp-core-guidelines/CppCoreGuidelines.md` path under:
  - Windows: `%LOCALAPPDATA%`, or `%USERPROFILE%\AppData\Local` when available.
  - macOS: `~/Library/Caches` when a home directory is available.
  - Linux: absolute `$XDG_CACHE_HOME`, or `~/.cache` when a home directory is available.
- If those paths are inaccessible, use an existing readable file or request one. Choose a download destination only after approval for network access and writes.

## Download after approval

Use the host's download tool for this fixed URL: `https://raw.githubusercontent.com/isocpp/CppCoreGuidelines/master/CppCoreGuidelines.md`. Choose a new file in an approved writable directory. Keep existing copies intact. For curl (available on Windows as `curl.exe`):

```text
curl --disable --fail --proto "=https" --max-time 30 --max-filesize 16777216 --output "<new-file>" "https://raw.githubusercontent.com/isocpp/CppCoreGuidelines/master/CppCoreGuidelines.md"
```

Require exit status 0. Then validate below. In PowerShell, check immediately after curl:

```powershell
if ($LASTEXITCODE -ne 0) { throw 'Guidelines download failed' }
```

In a POSIX shell, append `|| exit 1` to the curl command. On failure, remove only the new partial file. Return to local selection. If no permitted downloader is available, request the document from the user.

Choose a user-approved reusable cache first. For automatic reuse, use the cache filename above if it does not exist. To refresh an existing cache, validate a new file first. Replace the cache only with approval and a host-supported replacement operation. Otherwise, keep the validated file and use its explicit path. If the cache is unavailable, get approval for another writable directory, such as temporary storage. Keep the selected absolute path for every rule read. Temporary storage allows reuse only for its lifetime.

## Validate and record

Set `$path` to the selected absolute file path. Then run in PowerShell:

```powershell
$ErrorActionPreference = 'Stop'
$file = Get-Item -LiteralPath $path
if ($file.PSIsContainer -or $file.Length -eq 0 -or $file.Length -gt 16777216) {
    throw 'Guidelines file must contain 1 to 16777216 bytes'
}
$utf8 = New-Object System.Text.UTF8Encoding($false, $true)
$bytes = [System.IO.File]::ReadAllBytes($file.FullName)
$text = $utf8.GetString($bytes)
$prefix = [System.Text.Encoding]::ASCII.GetString($bytes, 0, [Math]::Min(16384, $bytes.Length))
if (-not $prefix.Contains('C++ Core Guidelines') -or
    -not $text.Contains('<a name="s-introduction"></a>') -or
    -not $text.Contains('<a name="s-resource"></a>') -or
    -not $text.Contains('<a name="s-unclassified"></a>') -or
    $text -notmatch '(?m)^# Bibliography[ \t]*\r?$') {
    throw 'Document is missing expected Guidelines sections; inspect upstream changes'
}
Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256
```

On Unix, set `path` to the selected file. Use the installed shell utilities:

```sh
if ! {
    test -f "$path" && test -s "$path" &&
    test "$(wc -c < "$path")" -le 16777216 &&
    iconv -f UTF-8 -t UTF-8 "$path" >/dev/null &&
    head -c 16384 "$path" | grep -F 'C++ Core Guidelines' >/dev/null &&
    grep -qF '<a name="s-introduction"></a>' "$path" &&
    grep -qF '<a name="s-resource"></a>' "$path" &&
    grep -qF '<a name="s-unclassified"></a>' "$path" &&
    grep -qE '^# Bibliography[[:space:]]*$' "$path"
}; then
    printf '%s\n' 'Invalid or incomplete Guidelines document' >&2
    exit 1
fi
```

Record the path and `sha256sum "$path"` (or `shasum -a 256 "$path"` on macOS). Hosts with only file-reading tools must inspect the title and anchors directly. State which size, encoding, or digest checks remain unavailable. These checks detect common damage. Section markers alone cannot prove a complete download. Without a reported length or trusted digest, record completeness as unverified.

## Read rules

Use `toc.md` and host search to find each rule's anchor and ID.

- For a heading, read through the next heading at the same or higher level. Include Exceptions and Notes.
- For a profile rule in a list, read the full item, including nested and continuation lines, through the next sibling item or heading.

Cite the rule's anchor with the official Guidelines link from the report template.
