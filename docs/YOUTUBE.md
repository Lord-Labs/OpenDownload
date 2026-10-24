# Public YouTube downloads

Current-source builds enable public YouTube video/audio through the existing
yt-dlp engine. Tagged v0.1.0 images remain unchanged and exclude YouTube. Build
from `main` for this addition. A public deployment requires its own verification;
local success does not establish cloud availability.

## Upstream research — 2026-10-24

- [yt-dlp 2026.10.21](https://github.com/yt-dlp/yt-dlp/releases/tag/2026.10.21)
  includes the new `tv_samsung` client and makes it a default YouTube client.
- [Default-client PR #17805](https://github.com/yt-dlp/yt-dlp/pull/17805)
  addresses failures of older TV clients. Keep upstream client selection rather
  than pinning a list of clients that may stop working independently.
- [EJS setup](https://github.com/yt-dlp/yt-dlp/wiki/EJS) requires a supported
  JavaScript runtime and matching solver scripts. The existing Node 24 image,
  `yt-dlp[default,curl-cffi]` and `--js-runtimes node` already satisfy this.
  Solver scripts are installed at build time; `--no-remote-components` remains.
- [SABR PR #13515](https://github.com/yt-dlp/yt-dlp/pull/13515) is still open.
  Its experimental branch and extra protobuf dependency are unnecessary for
  the verified stable-client path and are not included.
- [PO Token guide](https://github.com/yt-dlp/yt-dlp/wiki/PO-Token-Guide)
  explains additional per-client requirements. No token service or browser
  provider is installed. If a source requires verification or denies access,
  OpenDownload reports the existing safe error rather than importing credentials.

## Behavior and limits

The API applies its usual public URL/DNS validation, bounded analysis, guarded
proxy and owner-scoped job handling. The engine exposes only actual supported
formats, merges separate video/audio with FFmpeg, and offers source audio/MP3.
Private, account-only, paywalled, live and DRM media remain excluded. Playlists
are outside scope. There are no user-supplied cookies or extractor flags.

No Chromium, anonymous attestation provider, remote solver downloads, retry
fallback tree or parallel extraction service is added. The older experimental
YouTube branch remains unmerged.

## Verification

The stable build analyzed the previously reported public URL
`https://youtu.be/2cUkUbB3Gu4` and returned 49 formats locally. A real format-18
MP4 download completed with the new default TV client. This is a sample result,
not a platform-wide guarantee. API download evidence is recorded in
[VERIFICATION.md](VERIFICATION.md).

Repeat a small real video/MP3 test against an instance with media tools enabled:

```sh
python scripts/media-smoke.py \
  --source 'https://youtu.be/2cUkUbB3Gu4' \
  --video-option video-230
```

The smoke verifies actual stream types/duration using FFprobe, HTTP ranges and
file revocation after deletion. Use only media you own or may save. The selected
option must still be offered by current analysis. Live third-party availability
is opt-in and is not a required deterministic CI gate.
