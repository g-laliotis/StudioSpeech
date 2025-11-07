# AGENTS.md — Local Text‑to‑Speech (Greek & English) for YouTube (Mac & Windows)

**Goal:** Build a free, local (offline) TTS tool in Go that converts **.txt** and **.docx** files into high‑quality speech (WAV/MP3), suitable for **commercial/monetized YouTube** faceless channels. Runs smoothly on **macOS** and **Windows** without lag.

---

## 0) Ground Rules & Success Criteria

- **Local only:** No cloud APIs, no per‑character fees, works offline.
- **Commercial safety:** Only use voice models that allow **commercial use** (e.g., CC0/Apache/MIT or equivalent). Keep proof in repo.
- **Languages:** Greek (el‑GR) and English (en‑US/UK). Choose at least one production voice per language.
- **Quality:** Naturalness ≥ 4/5 in internal MOS; clear, stable prosody for ≥ 15‑minute scripts.
- **Performance:** ≤ 2.5× real‑time on CPU for medium voices (e.g., 10 min text renders ≤ 25 min) on a modest laptop; low memory footprint.
- **Output:** WAV (48 kHz mono) and MP3 (192 kbps). Loudness normalized to **−16 to −14 LUFS**.
- **Input:** .txt and .docx; Unicode safe; basic punctuation and paragraph handling.
- **UX:** Simple CLI (single command) + clear error messages. Optional small HTTP service later.
- **Portability:** macOS (Intel/Apple Silicon) and Windows (x86_64) with straightforward install steps.

**Non‑Goals (Phase 1):** voice cloning, cloud voices, phoneme‑level editing, GUI.

---

## 1) Platform & Dependencies (Decisions)

- **Language:** Go ≥ 1.21
- **Synth engine:** **Piper TTS** (CLI). Reason: fast, quality, offline, permissive licensing on many models.
- **Encoder/Post:** **ffmpeg** for WAV→MP3, resampling, and loudness normalization.
- **DOCX parsing:** Pure‑Go library (e.g., docx/OOXML reader). If any license risk appears, switch to a battle‑tested MIT/Apache lib.
- **Sentence splitting & normalization:** lightweight, rule‑based (numbers → words, abbreviations, punctuation pauses). Keep per‑language lists.

> **Commercial Safety Policy**
> - Only commit voice models whose metadata explicitly states **commercial use permitted**; store a copy of each model’s **license text** and checksums in `/voices/`.
> - Maintain `voices/catalog.json` with fields: `id`, `language`, `gender/style`, `sample_rate`, `license_name`, `license_url`, `source_url`, `sha256`.
> - Add a `LICENSES.md` appendix summarizing permissible use; link the upstream.

---

## 2) Architecture (High‑Level)

```
[text or .docx]
     │
     ▼
 TextIngestAgent ──► NormalizeAgent ──► SynthAgent(Piper) ──► PostProcessAgent(ffmpeg)
     │                     │                       │                 │
     └──► VoiceSelect ◄────┴── Prosody Params ◄────┘                 └──► Output (WAV/MP3)

                       ▲
                       └── CacheAgent (hash(text+voice+params) → reuse audio)
```

- **Interfaces:** Each agent has a simple input/output contract and a **Done checklist** with tests.
- **Orchestration:** CLI command calls agents in sequence; if any check fails, exit with actionable error.

---

## 3) Milestones & Iteration Flow

1. **M0 — Bootstrap & Voices**: scaffold repo, verify Piper/ffmpeg on macOS & Windows, curate 1 EN + 1 EL voice w/ licenses.
2. **M1 — TXT → WAV**: pipeline for .txt, normalization, WAV output; basic tests.
3. **M2 — DOCX → WAV/MP3**: add .docx ingest and MP3 export; loudness control.
4. **M3 — Performance & Caching**: real‑time factor targets; caching implemented.
5. **M4 — Release‑ready**: cross‑platform docs, packaging, license bundle, QA pass (15‑minute script).

> **Rule:** After each milestone, run the **Verification Checklist** at the bottom before moving on.

---

## 4) Agents (Tasks, Prompts, I/O, Tests)

> Each section includes an **Amazon Q Prompt** you can paste in VS Code to drive implementation/refactors/tests without writing code here.

### 4.1 ProjectScaffoldAgent
**Goal:** Create the repo layout, modules, and CI basics.

**Deliverables**
- `/cmd/ttscli` (CLI entry)
- `/internal/agents/` (each agent’s code)
- `/voices/` (models not committed by default; only metadata)
- `/scripts/` (installers, test data)
- `/testdata/` (tiny sample texts and docx)
- `voices/catalog.json`, `LICENSES.md`, `README.md`

**Done Checklist**
- `go mod tidy` succeeds.
- Build empty CLI binary on macOS & Windows.

**Amazon Q Prompt**
> Create a Go module with the above folders and an empty Cobra‑style CLI in `/cmd/ttscli`. Add a basic `make build` for macOS and Windows. Do not include any proprietary code.

---

### 4.2 EnvironmentAgent
**Goal:** Ensure Piper and ffmpeg are present; detect OS/CPU; provide install hints.

**Inputs:** none. **Outputs:** struct with paths/versions.

**Done Checklist**
- `piper --version` and `ffmpeg -version` parsed.
- Helpful install message if missing (brew/choco + manual fallback).

**Tests**
- Simulate missing tools and assert clear guidance.

**Amazon Q Prompt**
> Implement a Go package that locates `piper` and `ffmpeg` in PATH, captures versions, and returns actionable install guidance for macOS (brew) and Windows (choco/manual). Include unit tests with exec stubs.

---

### 4.3 VoiceCatalogAgent
**Goal:** Maintain a safe set of voices (EN+EL) with license proofs.

**Inputs:** `voices/catalog.json`. **Outputs:** selected voice path + metadata.

**Done Checklist**
- Catalog schema validated.
- Refuses to run if license is unknown or non‑commercial.
- SHA‑256 check for local model files (if provided by user).

**Tests**
- Catalog with non‑commercial entry → program exits with clear error.
- Happy path picks language‑appropriate default.

**Amazon Q Prompt**
> Define a `catalog.json` schema and loader. Add validation for `license_name`, `license_url`, and `commercial_use_allowed: true`. Provide a selector that prefers language match, then fallback. Include tests.

---

### 4.4 TextIngestAgent
**Goal:** Read `.txt` and `.docx`, unify to UTF‑8 text, preserve paragraphs.

**Inputs:** path to file. **Outputs:** normalized string array (paragraphs/sentences).

**Done Checklist**
- Handles UTF‑8 `.txt` and basic Windows‑1253/1252 with auto‑detect.
- `.docx` paragraphs extracted; inline soft breaks handled.
- Max line length control (soft wrap) optional.

**Tests**
- Mixed Greek/English docx round‑trips correctly.

**Amazon Q Prompt**
> Add a reader that accepts .txt and .docx. For .docx use a permissive‑license Go lib. Return a slice of paragraphs. Include encoding sniffing for legacy .txt files and tests with Greek text.

---

### 4.5 NormalizeAgent
**Goal:** Light text cleanup and prosody hints, per language.

**Inputs:** paragraphs. **Outputs:** sentences with hints.

**Functions**
- Replace straight quotes, normalize spaces.
- Expand common numbers (10 → ten / δέκα); configurable.
- Abbreviation map (e.g., “Dr.”, “κ.λπ.”), do not split.
- Optional markup tags: `[PAUSE=300ms]` → punctuation insert for Piper (via sentence breaks).

**Done Checklist**
- Language detection (simple heuristic) or user flag `--lang`.
- Stable sentence boundaries for long scripts.

**Tests**
- Numbers/dates/symbols in EL/EN.

**Amazon Q Prompt**
> Implement a normalization pass with sentence splitting for Greek and English. Provide small pluggable rulesets and unit tests covering abbreviations and numbers. Keep it lightweight—no external heavy NLP.

---

### 4.6 SynthAgent (Piper)
**Goal:** Convert sentences to WAV via Piper with tunable prosody.

**Inputs:** sentences + voice metadata + params (speed, noise, length_scale, speaker).
**Outputs:** one WAV per block or single concatenated WAV.

**Done Checklist**
- Map `speed` → Piper `--length_scale = 1/speed`.
- Safe defaults: speed 1.03 (slight pace), noise 0.667, noisew 0.8.
- Fail fast if model missing or wrong sample rate.

**Tests**
- Dry‑run mode constructs correct Piper command; no exec.
- Golden test: synth tiny sentence to WAV and assert headers.

**Amazon Q Prompt**
> Write a Piper wrapper that accepts text via stdin and writes a temp WAV. Expose params: model path, speed, noise, noisew, speaker. Include a dry‑run mode returning the command line for tests.

---

### 4.7 PostProcessAgent (ffmpeg)
**Goal:** Concatenate WAVs (if needed), resample to 48 kHz mono, normalize loudness, and export WAV/MP3.

**Inputs:** one or many WAVs. **Outputs:** final WAV and/or MP3.

**Done Checklist**
- Loudness normalized to −16…−14 LUFS (`ebur128` / `loudnorm`).
- MP3 export at 192 kbps CBR (or user‑set).
- Removes temps; deterministic output names.

**Tests**
- Verify sample rate/channel count and MP3 bitrate.

**Amazon Q Prompt**
> Implement a post‑processor using ffmpeg: concat list of WAVs, resample 48k mono, apply EBU R128 loudness normalization, export WAV and MP3. Include integration tests that check container/codec metadata.

---

### 4.8 CacheAgent
**Goal:** Avoid re‑synth for unchanged inputs.

**Inputs:** text + voice id + params. **Outputs:** cached file path or miss.

**Strategy:** SHA‑256 of canonicalized text + voice id + param JSON.

**Done Checklist**
- Cache hits logged; option `--no-cache` to bypass.

**Tests**
- Identical inputs reuse file; param change invalidates.

**Amazon Q Prompt**
> Implement a content‑addressed cache with SHA‑256 keys and an index JSON. Include unit tests for hit/miss and prune.

---

### 4.9 UXAgent (CLI)
**Goal:** One easy command for users.

**CLI Sketch**
```
ttscli synth --in myscript.docx --lang el-GR --voice greek_default \
  --out voice.mp3 --speed 1.05 --format mp3
```

**Done Checklist**
- Auto‑detect language if `--lang` omitted.
- Helpful errors with suggested fixes.
- `--check` prints environment (piper/ffmpeg versions, voices found).

**Tests**
- CLI flag parsing; invalid combos → errors.

**Amazon Q Prompt**
> Build a user‑friendly CLI with subcommands: `check`, `synth`. Provide examples in `--help`. Wire it to agents. Add unit tests for flag parsing and error messages.

---

### 4.10 CrossPlatformAgent
**Goal:** Ensure smooth macOS & Windows usage.

**Tasks**
- Path handling (`\` vs `/`), temp dirs, quoting.
- Install docs: Homebrew (macOS), Chocolatey/MSI (Windows), manual fallback.
- Optional static builds or zip releases.

**Done Checklist**
- Smoke test on both OSes returns valid MP3.

**Amazon Q Prompt**
> Add OS‑specific helpers for quoting paths and temp files. Generate install docs for macOS (brew) and Windows (choco/manual). Provide a simple PowerShell example.

---

### 4.11 QAAgent (Verification & Benchmarks)
**Goal:** Prevent regressions and ensure YouTube readiness.

**Checks**
- **Quality:** Internal MOS on 10× sentences (EL/EN). No clipping, stable pacing.
- **Loudness:** −16…−14 LUFS. Peak < −1 dBFS.
- **Performance:** Measure real‑time factor on a 5‑minute script.
- **Stability:** 20‑minute long‑form render without crash or drift.
- **Licensing:** Catalog entries contain license proofs.

**Amazon Q Prompt**
> Create a test suite that renders small fixtures and extracts audio metadata with ffprobe to assert sample rate, channels, and loudness filters applied. Add a benchmark that times a synthetic 5‑minute text.

---

### 4.12 ReleaseAgent
**Goal:** Package and document.

**Additions:** Ensure a generated `LICENSES.md` is included with attribution blocks.

### LICENSES.md (Template)
```
# LICENSES for Voices Used in This Project

## Public Domain Voices
### en_us_ljspeech_medium
Dataset: **LJ Speech** (Public Domain)
Source: https://keithito.com/LJ-Speech-Dataset/
Commercial Use: **Allowed**
Attribution Required: **No**

### el_gr_rapunzelina_low
Dataset: Greek Single Speaker Dataset (**CC0**)
Source: https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset
Commercial Use: **Allowed**
Attribution Required: **No**

---

## Attribution-Required Voices
### en_us_libritts_high_multi
Dataset: **LibriTTS** (CC BY 4.0)
License: http://www.openslr.org/60/
Commercial Use: **Allowed**
Attribution Required: **Yes**

#### Required Attribution Snippet
Put this in your README, website, OR YouTube video description:

"This project uses the LibriTTS dataset (CC BY 4.0). © Original contributors. Licensed under CC BY 4.0 (http://www.openslr.org/60/). No endorsement implied."

---

## Non-Commercial Voices (Blocked)
These voices **must not** be used for monetized YouTube or commercial outputs.
The program must refuse to load them.

### en_us_ryan
Dataset License: **CC BY-NC-SA 4.0** (Non-Commercial)
Commercial Use: **NOT allowed**
Action: Blocked by VoiceCatalogAgent
```

### VoiceCatalogAgent updates
- If `commercial_use_allowed` is false → **fatal error**.
- If license contains `NC` or `Non-Commercial` → **auto-block**.
- If `attribution_required` is true → display required attribution in `ttscli check`.

---
**Goal:** Package and document.

**Deliverables**
- `README.md`: quick start, examples, troubleshooting.
- `INSTALL.md`: macOS/Windows steps.
- `LICENSES.md`: voice model licenses + checksums.
- Example `voices/catalog.json` (placeholders; users add files themselves).
- Changelog and versioning.

**Done Checklist**
- Fresh machine can follow docs and produce an MP3 in <10 minutes.

**Amazon Q Prompt**
> Generate clear docs for installation and first use. Add a sample `catalog.json` with placeholders and comments about commercial use. Do not bundle any third‑party models.

---

## 5) Configuration Files (Proposed)

- `voices/catalog.json` (**pre-filled with safe, high‑quality options**)
```json
{
  "voices": [
    {
      "id": "en_us_ljspeech_medium",
      "language": "en-US",
      "style": "female-neutral",
      "sample_rate": 22050,
      "commercial_use_allowed": true,
      "attribution_required": false,
      "license_name": "Public Domain (LJ Speech)",
      "license_url": "https://keithito.com/LJ-Speech-Dataset/",
      "source_url": "https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/ljspeech/medium",
      "sha256": "<fill after download>",
      "path": "<local path>/en_US-ljspeech-medium.onnx"
    },
    {
      "id": "en_us_libritts_high_multi",
      "language": "en-US",
      "style": "multi-speaker",
      "sample_rate": 22050,
      "commercial_use_allowed": true,
      "attribution_required": true,
      "license_name": "CC BY 4.0 (LibriTTS)",
      "license_url": "http://www.openslr.org/60/",
      "source_url": "https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/libritts/high",
      "sha256": "<fill after download>",
      "path": "<local path>/en_US-libritts-high.onnx"
    },
    {
      "id": "el_gr_rapunzelina_low",
      "language": "el-GR",
      "style": "female-neutral",
      "sample_rate": 16000,
      "commercial_use_allowed": true,
      "attribution_required": false,
      "license_name": "CC0 (Greek Single Speaker Dataset)",
      "license_url": "https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset",
      "source_url": "https://huggingface.co/rhasspy/piper-voices/tree/main/el/el_GR/rapunzelina/low",
      "sha256": "<fill after download>",
      "path": "<local path>/el_GR-rapunzelina-low.onnx"
    }
  ]
}
```

> **Why these?**
> - **English (LJSpeech, medium):** Model card says dataset is **public domain** → safe for commercial use without attribution. citeturn16view0
> - **English (LibriTTS, high, multi‑speaker):** Dataset license **CC BY 4.0** → commercial allowed **with attribution** (include a credit in your README/ABOUT). citeturn25view0
> - **Greek (rapunzelina, low):** Model card lists dataset **CC0** → commercial use allowed, no attribution required. (Lower samplerate, but clean.) citeturn11view0
>
> **Note:** Avoid **en_US/ryan** (dataset **CC BY‑NC‑SA 4.0**, non‑commercial). citeturn19view0

- `config.example.json` (defaults for speed/noise and output)
```json
{
  "defaults": {
    "lang": "auto",
    "voice": "auto",
    "speed": 1.03,
    "noise": 0.667,
    "noisew": 0.8,
    "format": "mp3",
    "sample_rate": 48000,
    "bitrate_kbps": 192
  }
}
```

- **Attribution template** (use when a voice requires it, e.g., LibriTTS):
```
This project uses the "LibriTTS" dataset (CC BY 4.0). © Original contributors. Licensed under CC BY 4.0 (http://www.openslr.org/60/). No endorsement implied.
```

---

## 6) Developer Runbook

**Quick local check**
1) `ttscli check` → prints Piper/ffmpeg availability and voices.
2) `ttscli synth --in demo.txt --lang en-US --voice en_us_default --out voice.mp3`

**Typical YouTube export**
- Use WAV during editing; export final MP3 only for delivery: `--format wav` during creation, convert to MP3 last.

**Performance tips**
- Prefer medium voices; keep speed ≤ 1.08.
- Use caching for re‑renders.

**Troubleshooting**
- Missing `piper`: follow OS‑specific installer message.
- Clipping: lower speed or apply limiter in PostProcessAgent.
- Wrong language detection: pass `--lang` explicitly.

---

## 7) Verification Checklist (Run After Each Milestone)

- [ ] **Licensing**: All selected voices have `commercial_use_allowed: true` with proof (license text + URL + checksum).
- [ ] **Input**: .txt and .docx parsed correctly; Greek diacritics preserved.
- [ ] **Output**: WAV 48 kHz mono and MP3 192 kbps play in OS default player.
- [ ] **Quality**: Internal MOS ≥ 4/5 on EL+EN samples; no robotic artifacts.
- [ ] **Loudness**: −16…−14 LUFS, peak < −1 dBFS.
- [ ] **Performance**: Real‑time factor ≤ 2.5× on reference laptop.
- [ ] **Stability**: 15–20 min script renders without crash or drift.
- [ ] **Cross‑platform**: Smoke test passed on macOS and Windows.

---

## 8) Security & Privacy Notes

- Runs offline; no text leaves the machine.
- Avoid bundling third‑party models without licenses.
- Validate file inputs to prevent path traversal when adding HTTP mode later.

---

## 9) Future Extensions (Post‑MVP)

- Small HTTP server with `/synthesize` endpoint (same agents under the hood).
- Pronunciation dictionary (per‑project JSON).
- SSML subset parser for fine prosody control.
- Basic GUI wrapper (Electron/Tauri) reusing the CLI.

---

## 10) VS Code Tasks & Amazon Q Workflows

This section gives you a **copy‑paste** set of VS Code tasks and **Amazon Q prompts** to scaffold the repo step‑by‑step on macOS and Windows.

### A) `.vscode/tasks.json`
Create `.vscode/tasks.json` with the following content:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "00: Check prerequisites (piper & ffmpeg)",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "command -v piper && piper --version && command -v ffmpeg && ffmpeg -version | head -n 1"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "if (-not (Get-Command piper -ErrorAction SilentlyContinue)) { throw 'piper not found' } ; piper --version ; if (-not (Get-Command ffmpeg -ErrorAction SilentlyContinue)) { throw 'ffmpeg not found' } ; ffmpeg -version | Select-Object -First 1"]
      },
      "problemMatcher": []
    },
    {
      "label": "01: Go module init",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "test -f go.mod || go mod init ttslocal && go mod tidy"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "if (-not (Test-Path go.mod)) { go mod init ttslocal } ; go mod tidy"]
      },
      "problemMatcher": []
    },
    {
      "label": "02: Scaffold folders",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "mkdir -p cmd/ttscli internal/agents voices scripts testdata .vscode && touch README.md"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "mkdir cmd/ttscli, internal/agents, voices, scripts, testdata, .vscode -Force ; if (-not (Test-Path README.md)) { '' | Out-File README.md -Encoding utf8 } "]
      },
      "problemMatcher": []
    },
    {
      "label": "03: Create sample catalog & config",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "cat > voices/catalog.json <<'JSON'
{
  \"voices\": [
    {
      \"id\": \"en_us_ljspeech_medium\",
      \"language\": \"en-US\",
      \"style\": \"female-neutral\",
      \"sample_rate\": 22050,
      \"commercial_use_allowed\": true,
      \"attribution_required\": false,
      \"license_name\": \"Public Domain (LJ Speech)\",
      \"license_url\": \"https://keithito.com/LJ-Speech-Dataset/\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/ljspeech/medium\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/en_US-ljspeech-medium.onnx\"
    },
    {
      \"id\": \"en_us_libritts_high_multi\",
      \"language\": \"en-US\",
      \"style\": \"multi-speaker\",
      \"sample_rate\": 22050,
      \"commercial_use_allowed\": true,
      \"attribution_required\": true,
      \"license_name\": \"CC BY 4.0 (LibriTTS)\",
      \"license_url\": \"http://www.openslr.org/60/\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/libritts/high\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/en_US-libritts-high.onnx\"
    },
    {
      \"id\": \"el_gr_rapunzelina_low\",
      \"language\": \"el-GR\",
      \"style\": \"female-neutral\",
      \"sample_rate\": 16000,
      \"commercial_use_allowed\": true,
      \"attribution_required\": false,
      \"license_name\": \"CC0 (Greek Single Speaker Dataset)\",
      \"license_url\": \"https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/el/el_GR/rapunzelina/low\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/el_GR-rapunzelina-low.onnx\"
    }
  ]
}
JSON
"] ,
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "$json = @'
{
  \"voices\": [
    {
      \"id\": \"en_us_ljspeech_medium\",
      \"language\": \"en-US\",
      \"style\": \"female-neutral\",
      \"sample_rate\": 22050,
      \"commercial_use_allowed\": true,
      \"attribution_required\": false,
      \"license_name\": \"Public Domain (LJ Speech)\",
      \"license_url\": \"https://keithito.com/LJ-Speech-Dataset/\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/ljspeech/medium\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/en_US-ljspeech-medium.onnx\"
    },
    {
      \"id\": \"en_us_libritts_high_multi\",
      \"language\": \"en-US\",
      \"style\": \"multi-speaker\",
      \"sample_rate\": 22050,
      \"commercial_use_allowed\": true,
      \"attribution_required\": true,
      \"license_name\": \"CC BY 4.0 (LibriTTS)\",
      \"license_url\": \"http://www.openslr.org/60/\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/en/en_US/libritts/high\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/en_US-libritts-high.onnx\"
    },
    {
      \"id\": \"el_gr_rapunzelina_low\",
      \"language\": \"el-GR\",
      \"style\": \"female-neutral\",
      \"sample_rate\": 16000,
      \"commercial_use_allowed\": true,
      \"attribution_required\": false,
      \"license_name\": \"CC0 (Greek Single Speaker Dataset)\",
      \"license_url\": \"https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset\",
      \"source_url\": \"https://huggingface.co/rhasspy/piper-voices/tree/main/el/el_GR/rapunzelina/low\",
      \"sha256\": \"<fill after download>\",
      \"path\": \"<local path>/el_GR-rapunzelina-low.onnx\"
    }
  ]
}
'@ ; $json | Out-File voices/catalog.json -Encoding utf8 "]
      },
      "problemMatcher": []
    },
    {
      "label": "04: Create main CLI (stub)",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "cat > cmd/ttscli/main.go <<'GO'
package main
import (
  \"fmt\"
)
func main(){
  fmt.Println(\"ttscli stub — run 'ttscli check' after implementing agents.\")
}
GO
"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "$code = @'
package main
import (
  \"fmt\"
)
func main(){
  fmt.Println(\"ttscli stub — run 'ttscli check' after implementing agents.\")
}
'@ ; $code | Out-File cmd/ttscli/main.go -Encoding utf8 "]
      },
      "problemMatcher": []
    },
    {
      "label": "05: Build",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "go build -o bin/ttscli ./cmd/ttscli"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "if (-not (Test-Path bin)) { mkdir bin | Out-Null } ; go build -o bin/ttscli.exe ./cmd/ttscli"]
      },
      "problemMatcher": []
    },
    {
      "label": "06: Create LICENSES.md",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "cat > LICENSES.md <<'MD'
# LICENSES for Voices Used in This Project

## Public Domain Voices
### en_us_ljspeech_medium
Dataset: LJ Speech (Public Domain)
Source: https://keithito.com/LJ-Speech-Dataset/
Commercial Use: Allowed
Attribution Required: No

### el_gr_rapunzelina_low
Dataset: Greek Single Speaker Dataset (CC0)
Source: https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset
Commercial Use: Allowed
Attribution Required: No

---

## Attribution-Required Voices
### en_us_libritts_high_multi
Dataset: LibriTTS (CC BY 4.0)
License: http://www.openslr.org/60/
Commercial Use: Allowed
Attribution Required: Yes

**Required attribution to place in README or video description:**
\"This project uses the LibriTTS dataset (CC BY 4.0). © Original contributors. Licensed under CC BY 4.0 (http://www.openslr.org/60/). No endorsement implied.\"

---

## Non-Commercial Voices (Blocked)
These must not be used for monetized outputs.

### en_us_ryan
Dataset License: CC BY-NC-SA 4.0 (Non-Commercial)
Commercial Use: NOT allowed
MD
"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "$md = @'
# LICENSES for Voices Used in This Project

## Public Domain Voices
### en_us_ljspeech_medium
Dataset: LJ Speech (Public Domain)
Source: https://keithito.com/LJ-Speech-Dataset/
Commercial Use: Allowed
Attribution Required: No

### el_gr_rapunzelina_low
Dataset: Greek Single Speaker Dataset (CC0)
Source: https://www.kaggle.com/datasets/bryanpark/greek-single-speaker-speech-dataset
Commercial Use: Allowed
Attribution Required: No

---

## Attribution-Required Voices
### en_us_libritts_high_multi
Dataset: LibriTTS (CC BY 4.0)
License: http://www.openslr.org/60/
Commercial Use: Allowed
Attribution Required: Yes

**Required attribution to place in README or video description:**
\"This project uses the LibriTTS dataset (CC BY 4.0). © Original contributors. Licensed under CC BY 4.0 (http://www.openslr.org/60/). No endorsement implied.\"

---

## Non-Commercial Voices (Blocked)
These must not be used for monetized outputs.

### en_us_ryan
Dataset License: CC BY-NC-SA 4.0 (Non-Commercial)
Commercial Use: NOT allowed
'@ ; $md | Out-File LICENSES.md -Encoding utf8 "]
      },
      "problemMatcher": []
    },
    {
      "label": "07: Smoke test placeholders",
      "type": "shell",
      "command": "bash",
      "args": ["-lc", "echo 'Hello world' > testdata/demo.txt && echo 'Created testdata/demo.txt'"],
      "windows": {
        "command": "powershell",
        "args": ["-NoProfile", "-ExecutionPolicy", "Bypass", "'Hello world' | Out-File testdata/demo.txt -Encoding utf8 ; Write-Host 'Created testdata/demo.txt'"]
      },
      "problemMatcher": []
    }
  ]
}
```

**How to use**: Press `⇧⌘B` / `Ctrl+Shift+B` to run tasks, or open the **Terminal → Run Task…** menu and execute tasks **in order** (00 → 07).

---

### B) `q-workflows/` — Amazon Q step prompts
Create a folder `q-workflows/` and add the following files. Each is a copy‑paste **goal prompt** for Amazon Q to implement the corresponding agent.

1. **`00-bootstrap.md`**
```
Goal: Create project scaffold per AGENTS.md §4.1 ProjectScaffoldAgent.
- Folders: /cmd/ttscli, /internal/agents, /voices, /scripts, /testdata
- Init Cobra-like CLI in /cmd/ttscli with subcommands: check, synth
- Add Makefile targets: build, test
- Add README skeleton with quick start
- Do not add any third-party voice files
Acceptance: `go build ./cmd/ttscli` succeeds on macOS and Windows.
```

2. **`01-environment.md`**
```
Goal: Implement EnvironmentAgent per §4.2.
- Detect piper and ffmpeg in PATH; return versions
- Provide OS-specific install guidance if missing
- Unit tests use exec stubs/mocks
Acceptance: `ttscli check` prints availability and friendly guidance.
```

3. **`02-voice-catalog.md`**
```
Goal: Implement VoiceCatalogAgent per §4.3.
- Load voices/catalog.json; validate schema
- Enforce commercial_use_allowed=true; block Non-Commercial
- Support attribution_required flag and print required snippet in `check`
- Provide selection by --lang/--voice with sensible defaults
Acceptance: Catalog with non-commercial voice fails fast with clear error.
```

4. **`03-text-ingest.md`**
```
Goal: Implement TextIngestAgent per §4.4.
- Read .txt (encoding sniff) and .docx using a permissive Go lib
- Return []string paragraphs; preserve basic formatting
- Tests include mixed Greek/English
Acceptance: Given testdata/demo.docx and demo.txt, returns >0 paragraphs.
```

5. **`04-normalize.md`**
```
Goal: Implement NormalizeAgent per §4.5.
- Sentence splitting for EL and EN
- Number expansion (configurable), abbreviation list
- Optional [PAUSE=ms] markup → sentence breaks
Acceptance: Unit tests for abbreviations and numbers pass.
```

6. **`05-synth-piper.md`**
```
Goal: Implement SynthAgent per §4.6.
- Wrap piper CLI; map speed→length_scale, expose noise/noisew/speaker
- Dry-run builds command without executing (for tests)
- Golden test: small sentence → valid WAV header
Acceptance: `ttscli synth --in testdata/demo.txt --voice en_us_ljspeech_medium --format wav` creates a WAV.
```

7. **`06-postprocess-ffmpeg.md`**
```
Goal: Implement PostProcessAgent per §4.7.
- Concat WAV blocks; resample 48k mono
- Apply loudnorm to target -16 to -14 LUFS
- Export WAV and MP3 (192 kbps default)
Acceptance: ffprobe shows 48k mono; MP3 bitrate 192 kbps.
```

8. **`07-cache.md`**
```
Goal: Implement CacheAgent per §4.8.
- SHA-256 of canonical text+voice+params → file key
- `--no-cache` bypass flag; prune command optional
Acceptance: Re-running same input is a cache hit.
```

9. **`08-cli-ux.md`**
```
Goal: Implement UXAgent/CLI per §4.9.
- Subcommands: check, synth
- Flags: --in, --lang, --voice, --speed, --format, --out
- Helpful errors with suggestions
Acceptance: `ttscli check` and a simple synth both succeed.
```

10. **`09-cross-platform.md`**
```
Goal: CrossPlatformAgent per §4.10.
- Handle temp dirs, quoting, path separators
- Update INSTALL docs for macOS (brew) and Windows (choco/manual)
Acceptance: Smoke test on both OSes produces a valid MP3.
```

11. **`10-qa-bench.md`**
```
Goal: QAAgent per §4.11.
- Minimal integration tests rendering fixtures
- ffprobe-based assertions; simple benchmark of 5-minute text
Acceptance: CI passes; real-time factor recorded in logs.
```

12. **`11-release.md`**
```
Goal: ReleaseAgent per §4.12.
- Finalize README, INSTALL, LICENSES; ensure voices not bundled
- Provide sample commands for YouTube-friendly output
Acceptance: Fresh machine can follow docs to produce an MP3 in <10 minutes.
```

---

**Usage tip:** Open each `q-workflows/*.md` file, select all, and ask **Amazon Q** to “apply this goal” in the current repo. After each goal completes, run the **Verification Checklist** in AGENTS.md §7 before proceeding.

---

**End of AGENTS.md**

