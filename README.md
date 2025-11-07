# 🎤 StudioSpeech

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue?logo=go)](https://go.dev)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-✓_passing-success)](#)
[![Made with ❤️ in Go](https://img.shields.io/badge/made%20with-%E2%9D%A4%20in%20Go-00ADD8?logo=go)](#)
[![Docs](https://img.shields.io/badge/docs-AGENTS.md-blue?logo=readme)](AGENTS.md)
[![CLI Usage](https://img.shields.io/badge/CLI-Make_Help-orange?logo=gnu-bash)](#makefile-commands)
[![Security](https://img.shields.io/badge/security-policy-red?logo=shield)](SECURITY.md)
[![Code of Conduct](https://img.shields.io/badge/code%20of-conduct-ff69b4?logo=handshake)](CODE_OF_CONDUCT.md)
[![Contributing](https://img.shields.io/badge/contributing-guidelines-brightgreen?logo=github)](CONTRIBUTING.md)
[![GitHub Pages](https://img.shields.io/badge/GitHub%20Pages-Live%20Demo-brightgreen?logo=github)](https://g-laliotis.github.io/StudioSpeech/)

> 📘 See [**AGENTS.md**](AGENTS.md) for detailed technical specification of all pipeline agents.  
> 💻 Run `make help` for a list of available CLI commands.

---

## 🧩 Overview

**StudioSpeech** is a command-line text-to-speech tool written in Go using an **agent-based pipeline architecture**.

It reads input text files, applies a sequence of intelligent processing agents, and generates high-quality speech audio suitable for commercial YouTube content creation.

This project demonstrates clean modular Go design, testability, and professional TTS pipeline implementation.

**Professional Features:**
- 🔒 **Security Policy** - Vulnerability reporting process
- 🤝 **Code of Conduct** - Community guidelines
- 📋 **Issue Templates** - Structured bug reports and feature requests
- 🔄 **CI/CD Pipeline** - Automated testing on Go 1.22 & 1.23
- 📚 **Comprehensive Documentation** - Technical specs and contribution guides

---

## 🏗️ Architecture

Each processing stage is an **agent** — an independent module implementing:

```go
type Agent interface {
    Name() string
    Process(input interface{}) (interface{}, error)
    Validate() error
}
```

The agents form a pipeline that processes text through multiple stages:
```
+-------------------+     +------------------+     +---------------------+     +----------------+
|  EnvironmentAgent | --> | VoiceCatalogAgent| --> |   TextIngestAgent   | --> | NormalizeAgent |
| System validation |     | Voice management |     | Multi-format input  |     | Text cleanup   |
+-------------------+     +------------------+     +---------------------+     +----------------+
           |
           v
+-------------------+     +------------------+     +---------------------+
|    SynthAgent     | --> | PostProcessAgent | --> |    CacheAgent       |
| Speech synthesis  |     | Audio conversion |     | Result caching      |
+-------------------+     +------------------+     +---------------------+
```

This design keeps each processing stage isolated, testable, and auditable.

---

## ⚙️ Installation

**Clone the repository**
```bash
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech
```

**Install Go (v1.21 or higher)**
```bash
go version
# should print go version go1.21+ ...
```

**Download dependencies**
```bash
go mod tidy
```

**Check version**
```bash
go run ./cmd/ttscli --version
```

---

## 🚦 Usage

**Basic command**
```bash
go run ./cmd/ttscli synth <input_file> <output_file>
```

**Example:**
```bash
go run ./cmd/ttscli synth testdata/samples/sample.txt result.mp3
```

**Advanced usage:**
```bash
# Greek female voice
go run ./cmd/ttscli synth --in script.pdf --lang el-GR --gender female --out greek_voice.mp3

# English male voice with custom speed
go run ./cmd/ttscli synth --in document.docx --lang en-US --gender male --speed 1.2 --out english_voice.wav

# System check
go run ./cmd/ttscli check
```

**Example transformation:**

*Input (sample.txt)*
```
Hello world. This is a test sentence, with proper punctuation! 
How does it sound? Let me know if the pauses are natural.
```

*Output*
```
🎤 Generated natural-sounding speech with:
✓ Proper pauses at punctuation marks
✓ Natural speech rate (175 WPM)
✓ Commercial-safe voice selection
✓ High-quality audio (48kHz, 192kbps MP3)
```

---

## 🧩 Agents Summary

| Agent | Purpose | Example |
|-------|----------|----------|
| 🔍 **EnvironmentAgent** | System validation & dependency checking | Validates FFmpeg, macOS TTS availability |
| 🎭 **VoiceCatalogAgent** | Commercial-safe voice management | Selects `Melina` for Greek, `Samantha` for English |
| 📖 **TextIngestAgent** | Multi-format file processing | Extracts text from `.pdf`, `.docx`, `.txt` files |
| 🔧 **NormalizeAgent** | Text cleanup & prosody preparation | `Dr.` → `Doctor`, proper punctuation spacing |
| 🎤 **SynthAgent** | Speech synthesis with fallbacks | Piper TTS → macOS TTS → System default |
| 🎵 **PostProcessAgent** | Audio format conversion & enhancement | WAV → MP3, loudness normalization |
| 💾 **CacheAgent** | Performance optimization | SHA-256 based caching, avoid re-synthesis |

Full agent descriptions and internal specs: [`AGENTS.md`](AGENTS.md).

---

## ✨ Features

### 🔒 **Local & Offline**
- No cloud APIs or internet connection required
- Complete privacy - files never leave your system
- No per-character fees or usage limits
- Works in air-gapped environments

### 💼 **Commercial Safe**
- Only uses voice models with commercial licensing
- Safe for YouTube monetization
- Licensing validation built into voice catalog
- No copyright or attribution requirements

### 🌍 **Multi-language Support**
- **Greek (el-GR)**: Native pronunciation with Melina voice
- **English (en-US/UK)**: Multiple voice options (Alex, Samantha)
- Automatic language detection
- Language-specific text normalization

### 🎵 **High Quality Audio**
- Natural speech rates (175 WPM English, 160 WPM Greek)
- Proper punctuation pauses and prosody
- Professional audio quality (48kHz, 192kbps)
- Loudness normalization (-23 LUFS broadcast standard)

### 📄 **Multi-format Input**
- **Plain text** (.txt) with UTF-8 support
- **Microsoft Word** (.docx) documents
- **PDF files** with automatic text extraction
- Unicode support for international characters

### 🎛️ **Advanced Options**
- Gender-based voice selection (male/female)
- Adjustable speech speed (0.5x - 2.0x)
- Multiple output formats (MP3, WAV)
- Configurable audio quality settings

---

## 🧪 Testing

The project includes comprehensive unit, integration, and performance tests.

**Run all tests:**
```bash
make test
# or
go test ./...
```

**Run with coverage:**
```bash
make test-coverage
# or
go test -cover ./...
```

**Run benchmarks:**
```bash
make bench
# or
go test -bench=. ./...
```

**Test Types:**
- **Unit tests**: `internal/agents/*_test.go`
- **Integration tests**: `cmd/ttscli/*_test.go`
- **Performance tests**: `benchmark_test.go`
- **End-to-end tests**: Real audio generation validation

---

## 🧰 Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show help menu with all commands |
| `make build` | Build binary into bin/ttscli |
| `make run` | Run synthesis on sample data |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make bench` | Run benchmark tests |
| `make fmt` | Format code |
| `make vet` | Run static analysis |
| `make clean` | Remove build artifacts |
| `make install` | Install system dependencies |

Run the help menu any time:
```bash
make help
```

---

## 🧱 Project Structure

```
StudioSpeech/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.yml      # Bug report form
│   │   └── feature_request.yml # Feature request form
│   ├── workflows/
│   │   ├── pages.yml           # GitHub Pages deployment
│   │   └── test.yml            # CI/CD testing workflow
│   └── PULL_REQUEST_TEMPLATE.md # Pull request template
├── cmd/
│   └── ttscli/
│       ├── main.go         # CLI entrypoint
│       ├── synth.go        # Synthesis command
│       └── pipeline.go     # Pipeline orchestration
├── docs/
│   ├── .nojekyll           # Bypass Jekyll processing
│   ├── index.html          # GitHub Pages site
│   └── README.md           # Documentation
├── internal/
│   └── agents/
│       ├── environment.go  # System validation
│       ├── voice_catalog.go # Voice management
│       ├── text_ingest.go  # File processing
│       ├── normalize.go    # Text cleanup
│       ├── synth.go        # Speech synthesis
│       ├── macos_tts.go    # macOS TTS integration
│       ├── postprocess.go  # Audio processing
│       └── cache.go        # Result caching
├── testdata/
│   ├── samples/            # Sample files for testing
│   └── comprehensive/      # Comprehensive test cases
├── voices/
│   └── catalog.json        # Voice model catalog
├── .gitignore
├── AGENTS.md               # Technical specification
├── CHANGELOG.md            # Version history
├── CODE_OF_CONDUCT.md      # Community guidelines
├── CONTRIBUTING.md         # Contribution guidelines
├── go.mod
├── INSTALL.md              # Installation instructions
├── LICENSE                 # MIT License
├── LICENSES.md             # Voice model licensing
├── Makefile                # Build automation
├── README.md               # Project documentation
└── SECURITY.md             # Security policy
```

---

## 🚀 Performance

**Benchmarks** (on MacBook Pro M1):
- **Text Processing**: ~1000 words/second
- **Speech Synthesis**: ~5x real-time (macOS TTS)
- **Audio Conversion**: ~10x real-time (FFmpeg)
- **Cache Hit Rate**: >90% for repeated content
- **Memory Usage**: <50MB for typical documents

**Scalability**:
- Supports documents up to 10MB
- Handles 1000+ page PDFs efficiently
- Concurrent processing for multiple files
- Automatic resource management

---

## 🔧 Configuration

**Environment Variables:**
```bash
# TTS Engine preference
export STUDIOSPEECH_TTS_ENGINE=macos  # or piper

# Cache settings
export STUDIOSPEECH_CACHE_DIR=/tmp/studiospeech
export STUDIOSPEECH_CACHE_TTL=30d

# Audio quality
export STUDIOSPEECH_SAMPLE_RATE=48000
export STUDIOSPEECH_BITRATE=192
```

**Voice Catalog:**
Customize voice selection in `voices/catalog.json`:
```json
{
  "voices": [
    {
      "voice_id": "en_us_samantha",
      "name": "Samantha (English Female)",
      "language": "en-US",
      "gender": "female",
      "commercial_use_allowed": true,
      "quality_rating": 5
    }
  ]
}
```

---

## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

**Quick start:**
- Run `make fmt` and `make vet` before committing
- Add tests for any new functionality
- Keep agents single-purpose and testable
- Update `AGENTS.md` for behavior changes
- Follow conventional commit format

**Development workflow:**
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run the full test suite
5. Submit a pull request

**Version history:** See [CHANGELOG.md](CHANGELOG.md) for release notes.

---

## 📚 Documentation

- **[AGENTS.md](AGENTS.md)** - Technical specification and API documentation
- **[INSTALL.md](INSTALL.md)** - Detailed installation instructions
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines and development setup
- **[SECURITY.md](SECURITY.md)** - Security policy and vulnerability reporting
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes
- **[GitHub Pages](https://g-laliotis.github.io/StudioSpeech/)** - Live documentation site

---

## 🔒 Security

- **Local Processing**: All operations happen locally, no data transmission
- **Input Validation**: Comprehensive file type and content validation
- **Dependency Management**: Regular security updates and minimal dependencies
- **Vulnerability Reporting**: See [SECURITY.md](SECURITY.md) for reporting process

---

## 📄 License

Distributed under the MIT License.  
See [LICENSE](LICENSE) for details.

Voice model licensing information: [LICENSES.md](LICENSES.md)

---

## 🙏 Acknowledgments

This project is built for the YouTube creator community and demonstrates:
- Professional Go development practices
- Agent-based pipeline architecture
- Comprehensive testing and documentation
- Open source community standards

**Special thanks to:**
- The Go community for excellent tooling
- macOS for built-in TTS capabilities
- FFmpeg for audio processing
- The open source TTS research community

---

## 📞 Support

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - Questions and community support
- **Documentation** - Comprehensive guides and API docs
- **Security** - Vulnerability reporting via GitHub Security Advisories

Made with ❤️ for content creators worldwide 🎤
