# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Machine-readable AGENTS.md specification
- Professional project documentation structure
- GitHub Actions CI/CD workflows
- GitHub Pages documentation site
- Issue and PR templates
- Security policy and vulnerability reporting
- Code of Conduct (Contributor Covenant v2.1)
- Comprehensive contributing guidelines

### Changed
- Improved speech naturalness with better prosody handling
- Enhanced punctuation spacing and pause insertion
- Optimized speech rates (175 WPM English, 160 WPM Greek)
- Increased sentence length limit to 1500 characters for PDF processing

### Fixed
- Better sentence splitting preserving original punctuation
- Improved parenthetical expression handling
- Natural pause insertion after commas, colons, semicolons

## [1.0.0] - 2025-01-07

### Added
- Complete TTS pipeline with agent-based architecture
- Multi-format input support (.txt, .docx, .pdf)
- Commercial-safe voice catalog system
- Greek and English language support
- macOS native TTS integration with fallback support
- PDF processing capability with automatic language detection
- Gender-based voice selection (male/female)
- Audio format conversion (WAV/MP3) with FFmpeg
- Caching system for performance optimization
- CLI interface with comprehensive parameter validation

### Features
- **EnvironmentAgent**: System requirements validation
- **VoiceCatalogAgent**: Commercial voice model management
- **TextIngestAgent**: Multi-format file processing
- **NormalizeAgent**: Text cleanup and prosody preparation
- **SynthAgent**: Speech synthesis with Piper TTS and macOS fallback
- **PostProcessAgent**: Audio format conversion and quality enhancement
- **CacheAgent**: SHA-256 based result caching

### Technical Details
- Go 1.21+ compatibility
- Agent-based pipeline architecture
- Deterministic processing with no side effects
- UTF-8 text support with proper encoding handling
- Natural speech patterns with punctuation-aware pauses
- Commercial licensing validation for YouTube monetization
- Local/offline processing for privacy and cost savings

### Supported Platforms
- macOS (primary, with native TTS integration)
- Linux (with Piper TTS)
- Windows (planned)

### Voice Support
- English: Male (Alex), Female (Samantha, LJSpeech)
- Greek: Female (Melina) with native pronunciation
- Commercial-use-allowed voices only

### Audio Quality
- Output formats: MP3 (192kbps), WAV (48kHz mono)
- Loudness normalization (-23 LUFS broadcast standard)
- Natural speech rates optimized per language
- High-quality synthesis suitable for YouTube content

### CLI Commands
```bash
# System check
./ttscli check

# Basic synthesis
./ttscli synth --in input.txt --out output.mp3

# Advanced options
./ttscli synth --in script.pdf --lang el-GR --gender female --out greek_voice.mp3
```

### Dependencies
- Go 1.21+
- FFmpeg (for audio processing)
- Piper TTS (optional, with macOS fallback)
- github.com/ledongthuc/pdf (for PDF processing)

## [0.3.0] - 2025-01-07

### Added
- PDF file processing support
- Automatic language detection for Greek text
- Enhanced text normalization with abbreviation expansion
- Number-to-words conversion for natural speech

### Changed
- Improved sentence segmentation for better audio quality
- Enhanced error handling and user feedback
- Better file validation and format support

### Fixed
- Greek text processing and pronunciation
- PDF content extraction and paragraph handling
- Memory management for large documents

## [0.2.0] - 2025-01-07

### Added
- macOS native TTS integration using 'say' command
- Greek language support with Melina voice
- Gender-based voice selection
- Audio format conversion with FFmpeg
- Comprehensive error handling and fallback mechanisms

### Changed
- Switched from external TTS dependencies to macOS native support
- Improved voice selection algorithm
- Enhanced audio quality with proper sample rates

### Removed
- External TTS library dependencies (espeak-ng)
- Complex installation requirements

## [0.1.0] - 2025-01-07

### Added
- Initial project structure with Go modules
- Basic CLI framework using Cobra
- Agent-based architecture foundation
- Environment validation system
- Voice catalog management
- Text ingestion for .txt and .docx files
- Basic speech synthesis pipeline
- Caching mechanism for performance
- Makefile for build automation

### Technical Foundation
- Modular agent design for extensibility
- Pipeline orchestration system
- Configuration management
- Logging and error handling framework
- Test structure and basic coverage

---

## Release Notes

### Version Numbering
- **Major** (X.0.0): Breaking changes, major new features
- **Minor** (0.X.0): New features, backwards compatible
- **Patch** (0.0.X): Bug fixes, small improvements

### Supported Versions
- **1.x.x**: Active development and support
- **0.x.x**: Legacy versions, security fixes only

### Migration Guides
- **0.x to 1.0**: No breaking changes, direct upgrade
- **Future versions**: Migration guides will be provided for breaking changes

### Security Updates
Security updates are released as patch versions and documented in the [Security Policy](SECURITY.md).

### Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

### License
This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.