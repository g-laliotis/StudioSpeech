#!/bin/bash

# StudioSpeech Voice Setup Script
# Downloads commercial-safe voice models

set -e

VOICES_DIR="voices"
mkdir -p "$VOICES_DIR"

echo "🎭 Setting up commercial-safe voice models..."

# English - LJSpeech (Public Domain)
echo "📥 Downloading English voice (LJSpeech)..."
if [ ! -f "$VOICES_DIR/en_US-ljspeech-medium.onnx" ]; then
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ljspeech/medium/en_US-ljspeech-medium.onnx" \
         -o "$VOICES_DIR/en_US-ljspeech-medium.onnx"
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ljspeech/medium/en_US-ljspeech-medium.onnx.json" \
         -o "$VOICES_DIR/en_US-ljspeech-medium.onnx.json"
fi

# Greek - Rapunzelina (CC0)
echo "📥 Downloading Greek voice (Rapunzelina)..."
if [ ! -f "$VOICES_DIR/el_GR-rapunzelina-low.onnx" ]; then
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/el/el_GR/rapunzelina/low/el_GR-rapunzelina-low.onnx" \
         -o "$VOICES_DIR/el_GR-rapunzelina-low.onnx"
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/el/el_GR/rapunzelina/low/el_GR-rapunzelina-low.onnx.json" \
         -o "$VOICES_DIR/el_GR-rapunzelina-low.onnx.json"
fi

# Update catalog with actual paths
echo "📝 Updating voice catalog..."
sed -i.bak 's|<local path>/|voices/|g' "$VOICES_DIR/catalog.json"

echo "✅ Voice setup complete!"
echo "Run './bin/ttscli check' to verify installation."