# StudioSpeech

A free, local (offline) Text-to-Speech tool for creating high-quality speech from text files, designed for commercial YouTube content creation.

## Features

- **Local & Offline**: No cloud APIs, no per-character fees
- **Commercial Safe**: Only uses voice models with commercial licensing
- **Multi-language**: Supports Greek (el-GR) and English (en-US/UK)
- **High Quality**: Natural-sounding speech suitable for YouTube
- **Multiple Formats**: Outputs WAV (48kHz mono) and MP3 (192kbps)
- **Input Support**: Reads .txt and .docx files with Unicode support

## Quick Start

```bash
# Check system requirements
./ttscli check

# Convert text to speech with gender selection
./ttscli synth --in script.txt --lang en-US --gender female --out voice.mp3
./ttscli synth --in script.txt --lang el-GR --gender male --out voice.mp3
```

## Requirements

- Go 1.21+
- Piper TTS (for synthesis)
- FFmpeg (for audio processing)

## Installation

See [INSTALL.md](INSTALL.md) for detailed setup instructions.

## License

This project is licensed under the MIT License. See [LICENSES.md](LICENSES.md) for voice model licensing information.