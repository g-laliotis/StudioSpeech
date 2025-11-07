#!/bin/bash

# StudioSpeech Voice Setup Script
# Downloads commercial-safe voice models

set -e

VOICES_DIR="voices"
mkdir -p "$VOICES_DIR"

echo "🎭 Setting up commercial-safe voice models..."

# English Female - LJSpeech (Public Domain)
echo "📥 Downloading English female voice (LJSpeech)..."
if [ ! -f "$VOICES_DIR/en_US-ljspeech-medium.onnx" ]; then
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ljspeech/medium/en_US-ljspeech-medium.onnx" \
         -o "$VOICES_DIR/en_US-ljspeech-medium.onnx"
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ljspeech/medium/en_US-ljspeech-medium.onnx.json" \
         -o "$VOICES_DIR/en_US-ljspeech-medium.onnx.json"
fi

# English Male - Ryan (Public Domain)
echo "📥 Downloading English male voice (Ryan)..."
if [ ! -f "$VOICES_DIR/en_US-ryan-high.onnx" ]; then
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx" \
         -o "$VOICES_DIR/en_US-ryan-high.onnx"
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/ryan/high/en_US-ryan-high.onnx.json" \
         -o "$VOICES_DIR/en_US-ryan-high.onnx.json"
fi

# Greek Female - Rapunzelina (CC0)
echo "📥 Downloading Greek female voice (Rapunzelina)..."
if [ ! -f "$VOICES_DIR/el_GR-rapunzelina-low.onnx" ]; then
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/el/el_GR/rapunzelina/low/el_GR-rapunzelina-low.onnx" \
         -o "$VOICES_DIR/el_GR-rapunzelina-low.onnx"
    curl -L "https://huggingface.co/rhasspy/piper-voices/resolve/main/el/el_GR/rapunzelina/low/el_GR-rapunzelina-low.onnx.json" \
         -o "$VOICES_DIR/el_GR-rapunzelina-low.onnx.json"
fi

# Note: Greek male voice would need to be sourced from available commercial-safe models
echo "⚠️  Greek male voice: Check Piper voices repository for commercial-safe options"

# Update catalog with actual paths
echo "📝 Updating voice catalog..."
sed -i.bak 's|<local path>/|voices/|g' "$VOICES_DIR/catalog.json"

echo "✅ Voice setup complete!"
echo "Run './bin/ttscli check' to verify installation."