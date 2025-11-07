# AGENTS SPECIFICATION

## PIPELINE_ORDER
1. EnvironmentAgent
2. VoiceCatalogAgent
3. TextIngestAgent
4. NormalizeAgent
5. SynthAgent
6. PostProcessAgent
7. CacheAgent

## AGENT_INTERFACE
```go
type Agent interface {
    Name() string
    Process(input interface{}) (interface{}, error)
    Validate() error
}
```

## ENVIRONMENT_AGENT
FILE: internal/agents/environment.go
ORDER: 1
FUNCTION: System requirements validation and dependency checking
DEPENDENCIES:
- Piper TTS (optional, fallback to macOS TTS)
- FFmpeg (required for audio processing)
- macOS 'say' command (macOS only)
VALIDATION_CHECKS:
- Operating system compatibility
- Required binaries availability
- System resource requirements
- Audio subsystem functionality
OUTPUTS:
- System capability report
- Installation recommendations
- Fallback options available
EXAMPLES:
- macOS: Validates 'say' command and FFmpeg
- Linux: Checks Piper TTS installation
- Windows: Validates SAPI availability
ERROR_HANDLING:
- Missing dependencies logged with install instructions
- Graceful degradation to available TTS engines

## VOICE_CATALOG_AGENT
FILE: internal/agents/voice_catalog.go
ORDER: 2
FUNCTION: Commercial-safe voice model management and selection
CATALOG_FILE: voices/catalog.json
VOICE_PROPERTIES:
- voice_id: Unique identifier
- name: Human-readable name
- language: ISO language code (en-US, el-GR)
- gender: male/female/neutral
- commercial_use_allowed: boolean
- quality_rating: 1-5 scale
- file_size_mb: Model size
- sample_rate: Audio sample rate
SELECTION_CRITERIA:
- Language matching (exact or fallback)
- Gender preference
- Commercial licensing requirement
- Quality rating priority
- System compatibility
EXAMPLES:
- "en-US + female" -> en_us_ljspeech_female
- "el-GR + male" -> el_gr_male_voice
- Auto-selection based on text language detection
FALLBACK_STRATEGY:
- Primary: Piper TTS models
- Secondary: macOS built-in voices
- Tertiary: System default TTS

## TEXT_INGEST_AGENT
FILE: internal/agents/text_ingest.go
ORDER: 3
FUNCTION: Multi-format text file processing and language detection
SUPPORTED_FORMATS:
- .txt: Plain text with UTF-8 encoding
- .docx: Microsoft Word documents
- .pdf: Portable Document Format
PROCESSING_PIPELINE:
- File format detection
- Content extraction
- Encoding validation (UTF-8)
- Language detection (en-US, el-GR)
- Paragraph segmentation
- Word count calculation
LANGUAGE_DETECTION:
- Greek: Unicode range detection (Ά-ώ)
- English: Default fallback
- Confidence scoring
- Manual override support
OUTPUT_STRUCTURE:
```go
type TextContent struct {
    Paragraphs []string
    Language   string
    WordCount  int
    Metadata   map[string]interface{}
}
```
ERROR_HANDLING:
- Unsupported file formats
- Corrupted or encrypted files
- Encoding issues
- Empty or invalid content

## NORMALIZE_AGENT
FILE: internal/agents/normalize.go
ORDER: 4
FUNCTION: Text cleanup and prosody preparation for natural speech
NORMALIZATION_RULES:
- Abbreviation expansion (language-specific)
- Number-to-words conversion (0-20)
- Punctuation normalization
- Sentence segmentation
- Length validation (max 1500 characters)
ABBREVIATION_MAPS:
- English: Dr. -> Doctor, etc. -> etcetera
- Greek: κ.λπ. -> και λοιπά, π.χ. -> παραδείγματος χάρη
NUMBER_CONVERSION:
- English: 1 -> one, 2 -> two
- Greek: 1 -> ένα, 2 -> δύο
PROSODY_IMPROVEMENTS:
- Comma spacing: ", " for natural pauses
- Colon emphasis: ": " for proper intonation
- Parenthetical handling: " (text) "
- Multiple space cleanup
SENTENCE_SPLITTING:
- Punctuation-based segmentation (.!?)
- Preserve original punctuation
- Ensure proper sentence endings
OUTPUT_STRUCTURE:
```go
type NormalizedText struct {
    Sentences []string
    Language  string
    Metadata  map[string]interface{}
}
```
VALIDATION:
- Sentence length limits (1500 chars)
- Non-empty sentence check
- Proper punctuation validation

## SYNTH_AGENT
FILE: internal/agents/synth.go, internal/agents/macos_tts.go
ORDER: 5
FUNCTION: Speech synthesis with multiple TTS engine support
TTS_ENGINES:
- Primary: Piper TTS (high quality, commercial safe)
- Fallback: macOS built-in TTS (system integration)
- Future: SAPI (Windows), eSpeak (Linux)
SYNTHESIS_PARAMETERS:
- Voice selection (from VoiceCatalogAgent)
- Speech rate (175 WPM English, 160 WPM Greek)
- Audio quality (48kHz, mono)
- Output format (WAV native, MP3 converted)
MACOS_TTS_VOICES:
- English Male: Alex
- English Female: Samantha
- Greek Female: Melina (native pronunciation)
PIPER_TTS_INTEGRATION:
- Model loading and validation
- Batch sentence processing
- Memory management
- Error recovery
AUDIO_OUTPUT:
- Temporary WAV generation
- Format conversion via FFmpeg
- Quality validation
- Cleanup procedures
ERROR_HANDLING:
- TTS engine failures
- Audio generation errors
- Format conversion issues
- Graceful fallback between engines

## POSTPROCESS_AGENT
FILE: internal/agents/postprocess.go
ORDER: 6
FUNCTION: Audio format conversion and quality enhancement
AUDIO_PROCESSING:
- Format conversion (WAV -> MP3)
- Sample rate normalization (48kHz)
- Loudness normalization (-23 LUFS)
- Bitrate optimization (192kbps MP3)
FFMPEG_OPERATIONS:
- Input validation
- Format detection
- Conversion parameters
- Quality verification
SUPPORTED_FORMATS:
- Input: WAV, AIFF (macOS native)
- Output: MP3, WAV
- Future: OGG, FLAC, M4A
QUALITY_PARAMETERS:
- MP3: 192kbps CBR
- WAV: 48kHz 16-bit mono
- Loudness: -23 LUFS (broadcast standard)
- Dynamic range preservation
ERROR_HANDLING:
- FFmpeg execution errors
- Unsupported format combinations
- Quality validation failures
- Disk space limitations

## CACHE_AGENT
FILE: internal/agents/cache.go
ORDER: 7
FUNCTION: Result caching and performance optimization
CACHE_STRATEGY:
- SHA-256 content hashing
- Filesystem-based storage
- Automatic cleanup policies
- Cache hit/miss tracking
CACHE_STRUCTURE:
- Key: SHA-256(normalized_text + voice_id + parameters)
- Value: Generated audio file path
- Metadata: Creation time, access count, file size
CACHE_POLICIES:
- TTL: 30 days default
- Size limit: 1GB default
- LRU eviction strategy
- Manual cache clearing
PERFORMANCE_BENEFITS:
- Avoid re-synthesis of identical content
- Faster response for repeated requests
- Reduced computational overhead
- Bandwidth savings for large files
CACHE_VALIDATION:
- Content integrity checks
- File existence verification
- Metadata consistency
- Automatic repair mechanisms

## TESTING_FRAMEWORK
UNIT_TESTS: internal/agents/*_test.go
INTEGRATION_TESTS: cmd/ttscli/*_test.go
BENCHMARK_TESTS: benchmark_test.go
TEST_DATA: testdata/samples/*, testdata/comprehensive/*
COVERAGE_TARGET: 80%+ line coverage
TEST_COMMANDS:
- make test: Run all tests
- make test-verbose: Detailed test output
- make test-coverage: Coverage report
- make bench: Performance benchmarks

## IMPLEMENTATION_REQUIREMENTS
- Process data through sequential pipeline
- Maintain immutable agent state
- Handle errors gracefully with fallbacks
- Support concurrent processing where safe
- Provide detailed logging and metrics
- Follow Go best practices and conventions
- Ensure deterministic output
- Support configuration via environment variables

## EXTENSION_PROTOCOL
1. Create agent file in internal/agents/
2. Implement Agent interface methods
3. Add to pipeline in cmd/ttscli/pipeline.go
4. Create comprehensive unit tests
5. Update this specification
6. Add integration tests
7. Update documentation and examples

## CONFIGURATION_SCHEMA
```yaml
agents:
  environment:
    check_dependencies: true
    install_recommendations: true
  voice_catalog:
    catalog_file: "voices/catalog.json"
    commercial_only: true
  text_ingest:
    max_file_size_mb: 10
    supported_formats: [".txt", ".docx", ".pdf"]
  normalize:
    max_sentence_length: 1500
    expand_abbreviations: true
  synth:
    preferred_engine: "piper"
    fallback_engine: "macos"
    speech_rate_wpm: 175
  postprocess:
    output_format: "mp3"
    bitrate_kbps: 192
    sample_rate_hz: 48000
  cache:
    enabled: true
    ttl_days: 30
    max_size_gb: 1
```

## API_ENDPOINTS
```go
// Agent interface for all pipeline components
type Agent interface {
    Name() string
    Process(input interface{}) (interface{}, error)
    Validate() error
}

// Pipeline orchestration
func RunPipeline(inputFile string, config Config) (*AudioResult, error)

// Individual agent access
func NewEnvironmentAgent() *EnvironmentAgent
func NewVoiceCatalogAgent(catalogPath string) *VoiceCatalogAgent
func NewTextIngestAgent() *TextIngestAgent
func NewNormalizeAgent() *NormalizeAgent
func NewSynthAgent(tempDir string) *SynthAgent
func NewPostProcessAgent() *PostProcessAgent
func NewCacheAgent(cacheDir string) *CacheAgent
```