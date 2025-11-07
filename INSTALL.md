# Installation Guide

## Prerequisites

### macOS
```bash
# Install dependencies
brew install piper-tts ffmpeg

# Verify installation
piper --version
ffmpeg -version
```

### Windows
```bash
# Using Chocolatey
choco install piper-tts ffmpeg

# Or download manually:
# Piper: https://github.com/rhasspy/piper/releases
# FFmpeg: https://ffmpeg.org/download.html
```

### Linux
```bash
# Download Piper from releases
wget https://github.com/rhasspy/piper/releases/latest/download/piper_linux_x86_64.tar.gz
tar -xzf piper_linux_x86_64.tar.gz
sudo cp piper/piper /usr/local/bin/

# Install FFmpeg
sudo apt install ffmpeg  # Ubuntu/Debian
sudo yum install ffmpeg  # CentOS/RHEL
```

## Voice Models

1. Download commercial-safe models from [Piper Voices](https://huggingface.co/rhasspy/piper-voices)
2. Update `voices/catalog.json` with actual file paths
3. Verify with: `./ttscli check`

## Quick Test
```bash
# Build
make build

# Test system
./bin/ttscli check

# Synthesize speech
./bin/ttscli synth --in testdata/demo.txt --out test.mp3
```