# Installation Guide

## 🚀 Quick Start

StudioSpeech works **locally and offline** with no external dependencies on most systems. Choose your platform:

### 📱 macOS (Recommended)
```bash
# 1. Clone the repository
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech

# 2. Install FFmpeg (for audio conversion)
brew install ffmpeg

# 3. Build and run
make build
make run "your-file.txt"
```

**That's it!** macOS has built-in TTS with excellent Greek support.

### 🐧 Linux (Ubuntu/Debian)
```bash
# 1. Clone the repository
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech

# 2. Install dependencies
sudo apt update
sudo apt install ffmpeg golang-go espeak espeak-data

# 3. Build and test
make build
make run "testdata/samples/sample.txt"
```

**Note**: Linux uses eSpeak for TTS. For better quality, install Festival:
```bash
sudo apt install festival festvox-kallpc16k
```

### 🪟 Windows
```powershell
# 1. Install Go from https://golang.org/dl/
# 2. Install FFmpeg from https://ffmpeg.org/download.html
# 3. Install Git from https://git-scm.com/download/win

# 4. Clone and build
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech
go build -o ttscli.exe ./cmd/ttscli

# 5. Test installation
.\ttscli.exe version
.\ttscli.exe synth --in "testdata\samples\sample.txt" --out "test.mp3"
```

**Note**: Windows uses SAPI for TTS. Greek support is limited compared to macOS.

---

## 📋 System Requirements

### Minimum Requirements
- **Go**: 1.21 or higher
- **FFmpeg**: For audio format conversion
- **Disk Space**: 50MB for installation
- **Memory**: 100MB RAM during processing

### Platform-Specific TTS Engines

| Platform | Primary TTS | Voice Quality | Greek Support |
|----------|-------------|---------------|---------------|
| **macOS** | Built-in `say` | ⭐⭐⭐⭐⭐ | ✅ Native (Melina female only) |
| **Linux** | Built-in TTS | ⭐⭐⭐ | ⚠️ Limited |
| **Windows** | SAPI | ⭐⭐⭐ | ⚠️ Limited |

---

## 🛠️ Detailed Installation

### Step 1: Install Go
**macOS:**
```bash
brew install go
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang
```

**Windows:**
- Download from [golang.org](https://golang.org/dl/)
- Run installer and follow prompts

### Step 2: Install FFmpeg
**macOS:**
```bash
brew install ffmpeg
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# CentOS/RHEL
sudo yum install ffmpeg
```

**Windows:**
- Download from [ffmpeg.org](https://ffmpeg.org/download.html)
- Extract to `C:\ffmpeg`
- Add `C:\ffmpeg\bin` to PATH

### Step 3: Clone and Build
**macOS/Linux:**
```bash
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech
make deps
make build
```

**Windows:**
```powershell
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech
go mod tidy
go build -o ttscli.exe ./cmd/ttscli
```

### Step 4: Verify Installation
**macOS/Linux:**
```bash
make check
./bin/ttscli version
```

**Windows:**
```powershell
.\ttscli.exe version
.\ttscli.exe check
```

---

## 🎯 Simple Usage

### Basic Commands
```bash
# Convert text file to speech
make run "script.txt"
# Output: script.mp3

# Convert with specific language
make run-greek "greek-script.txt"
# Output: greek-script.mp3

# Convert with specific gender
make run-male "script.txt"
# Output: script.mp3 (male voice)
```

### Advanced Usage
```bash
# Full control
./bin/ttscli synth --in "input.pdf" --lang el-GR --gender female --out "output.mp3"

# System check
./bin/ttscli check

# Help
./bin/ttscli --help
```

---

## 🔧 Configuration

### Environment Variables
**macOS/Linux:**
```bash
# TTS Engine
export STUDIOSPEECH_TTS_ENGINE=macos  # or system

# Cache settings
export STUDIOSPEECH_CACHE_DIR=/tmp/studiospeech
export STUDIOSPEECH_CACHE_TTL=30d

# Audio quality
export STUDIOSPEECH_SAMPLE_RATE=48000
export STUDIOSPEECH_BITRATE=192
```

**Windows:**
```powershell
# TTS Engine
$env:STUDIOSPEECH_TTS_ENGINE="sapi"

# Cache settings
$env:STUDIOSPEECH_CACHE_DIR="C:\temp\studiospeech"
$env:STUDIOSPEECH_CACHE_TTL="30d"

# Audio quality
$env:STUDIOSPEECH_SAMPLE_RATE="48000"
$env:STUDIOSPEECH_BITRATE="192"
```

### Voice Customization
Edit `voices/catalog.json` to customize voice selection:
```json
{
  "voices": [
    {
      "voice_id": "en_us_samantha",
      "name": "Samantha (English Female)",
      "language": "en-US",
      "gender": "female",
      "commercial_use_allowed": true
    }
  ]
}
```

---

## 🚨 Troubleshooting

### Common Issues

**"Go not found"**
```bash
# Check Go installation
go version
# Should show: go version go1.21+ ...

# If not installed, follow Step 1 above
```

**"FFmpeg not found"**
```bash
# Check FFmpeg installation
ffmpeg -version
# Should show FFmpeg version info

# If not installed, follow Step 2 above
```

**"Permission denied"**
```bash
# Make binary executable
chmod +x bin/ttscli

# Or run with make
make run "your-file.txt"
```

**"No audio output"**
- Check file permissions in output directory
- Verify input file contains text
- Try with a simple .txt file first

### Platform-Specific Issues

**macOS:**
- If `say` command fails, check System Preferences > Accessibility > Speech
- For Greek voices, ensure language pack is installed

**Linux:**
- Install system TTS: `sudo apt install espeak espeak-data`
- For better quality: `sudo apt install festival festvox-kallpc16k`

**Windows:**
- Ensure SAPI voices are installed
- Check Windows Speech Platform Runtime

---

## 🔄 Updates

### Updating StudioSpeech
```bash
cd StudioSpeech
git pull origin main
make clean
make build
```

### Checking for Updates
```bash
# Check current version
./bin/ttscli version

# Check latest release
curl -s https://api.github.com/repos/g-laliotis/StudioSpeech/releases/latest | grep tag_name
```

---

## 🆘 Getting Help

- **Documentation**: [README.md](README.md)
- **Technical Specs**: [AGENTS.md](AGENTS.md)
- **Issues**: [GitHub Issues](https://github.com/g-laliotis/StudioSpeech/issues)
- **Discussions**: [GitHub Discussions](https://github.com/g-laliotis/StudioSpeech/discussions)

---

## 🎉 Success!

If you can run this command successfully, you're all set:
```bash
make run "testdata/samples/sample.txt"
```

You should see:
```
🎤 Generated natural-sounding speech with:
✓ Proper pauses at punctuation marks
✓ Natural speech rate (175 WPM)
✓ Commercial-safe voice selection
✓ High-quality audio (48kHz, 192kbps MP3)
```

**Happy creating! 🎤**