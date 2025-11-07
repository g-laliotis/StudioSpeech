# Contributing to StudioSpeech

Thank you for your interest in contributing to StudioSpeech! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- Git
- Basic understanding of text-to-speech systems

### Development Setup
```bash
# Clone the repository
git clone https://github.com/g-laliotis/StudioSpeech.git
cd StudioSpeech

# Install dependencies
go mod tidy

# Run tests
make test

# Build the project
make build
```

## 🏗️ Project Structure

```
StudioSpeech/
├── cmd/ttscli/           # CLI application entry point
├── internal/agents/      # Core TTS pipeline agents
├── voices/              # Voice model catalog
├── testdata/            # Test files and samples
├── docs/                # Documentation and GitHub Pages
└── scripts/             # Build and utility scripts
```

## 🧩 Agent Architecture

StudioSpeech uses a pipeline of agents for text-to-speech processing:

1. **EnvironmentAgent** - System validation
2. **VoiceCatalogAgent** - Voice management
3. **TextIngestAgent** - File processing
4. **NormalizeAgent** - Text cleanup
5. **SynthAgent** - Speech synthesis
6. **PostProcessAgent** - Audio processing
7. **CacheAgent** - Result caching

Each agent implements the `Agent` interface and should be:
- **Single-purpose** - One responsibility per agent
- **Testable** - Unit tests for all functionality
- **Deterministic** - Same input produces same output
- **Side-effect free** - No global state modifications

## 📝 Code Style

### Go Guidelines
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

### Commit Messages
Use conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
```
feat(agents): add PDF processing support
fix(synth): resolve macOS TTS rate control
docs: update installation instructions
test(normalize): add edge case coverage
```

## 🧪 Testing

### Running Tests
```bash
# All tests
make test

# Specific package
go test ./internal/agents/

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Test Categories
- **Unit tests** - Individual agent functionality
- **Integration tests** - End-to-end pipeline testing
- **Performance tests** - Benchmarking and memory usage

### Writing Tests
- Test both success and error cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Ensure tests are deterministic and fast

## 🐛 Bug Reports

Use the GitHub issue template and include:
- StudioSpeech version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant log output

## ✨ Feature Requests

Before submitting:
1. Check existing issues for duplicates
2. Consider if it fits the project scope
3. Provide clear use cases and benefits
4. Suggest implementation approach if possible

## 🔍 Code Review Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'feat: add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Review Criteria
- Code follows project conventions
- Tests pass and coverage is maintained
- Documentation is updated
- No breaking changes without discussion
- Performance impact is considered

## 📚 Documentation

### Types of Documentation
- **README.md** - User-facing documentation
- **AGENTS.md** - Machine-readable technical specs
- **Code comments** - Inline documentation
- **GitHub Pages** - Web documentation

### Documentation Standards
- Use clear, concise language
- Include code examples
- Keep documentation up-to-date with code changes
- Use proper markdown formatting

## 🔒 Security

- Report security vulnerabilities via GitHub Security Advisories
- Do not include credentials or sensitive data in code
- Follow secure coding practices
- Validate all user inputs

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🤝 Community

- Be respectful and inclusive
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)
- Help others learn and grow
- Share knowledge and best practices

## 📞 Getting Help

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - Questions and community support
- **Documentation** - Check existing docs first

Thank you for contributing to StudioSpeech! 🎤