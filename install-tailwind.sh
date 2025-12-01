#!/usr/bin/env bash

set -e

# ============================================
# TAILWIND CSS STANDALONE INSTALLER
# ============================================

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"
TAILWIND_VERSION="v4.0.0-alpha.32" # Ou latest stable
TAILWIND_BIN="$BIN_DIR/tailwindcss"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎨 Tailwind CSS Standalone Installer"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# ============================================
# 1. DÉTECTION SYSTÈME
# ============================================

OS="$(uname -s)"
ARCH="$(uname -m)"

echo "📋 Système détecté: $OS $ARCH"

case "$OS" in
    Linux*)
        PLATFORM="linux"
        ;;
    Darwin*)
        PLATFORM="macos"
        ;;
    *)
        echo "❌ OS non supporté: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCHITECTURE="x64"
        ;;
    aarch64 | arm64)
        ARCHITECTURE="arm64"
        ;;
    armv7l)
        ARCHITECTURE="armv7"
        ;;
    *)
        echo "❌ Architecture non supportée: $ARCH"
        exit 1
        ;;
esac

# ============================================
# 2. CONSTRUCTION URL DE TÉLÉCHARGEMENT
# ============================================

# Tailwind Standalone binary URLs
# https://github.com/tailwindlabs/tailwindcss/releases
DOWNLOAD_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-${PLATFORM}-${ARCHITECTURE}"

echo "📦 URL téléchargement:"
echo "   $DOWNLOAD_URL"

# ============================================
# 3. CRÉATION DOSSIERS
# ============================================

mkdir -p "$BIN_DIR"
mkdir -p "$PROJECT_ROOT/public/css"

# ============================================
# 4. TÉLÉCHARGEMENT
# ============================================

echo "⬇️  Téléchargement Tailwind CSS..."

if command -v curl &>/dev/null; then
    curl -sL "$DOWNLOAD_URL" -o "$TAILWIND_BIN"
elif command -v wget &>/dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TAILWIND_BIN"
else
    echo "❌ curl ou wget requis"
    exit 1
fi

# ============================================
# 5. PERMISSIONS EXÉCUTION
# ============================================

chmod +x "$TAILWIND_BIN"

echo "✅ Tailwind CSS installé: $TAILWIND_BIN"

# ============================================
# 6. VÉRIFICATION
# ============================================

if [ -f "$TAILWIND_BIN" ]; then
    VERSION=$("$TAILWIND_BIN" --help | head -n 1 || echo "tailwindcss")
    echo "✅ Version: $VERSION"
else
    echo "❌ Installation échouée"
    exit 1
fi

# ============================================
# 7. CRÉATION FICHIERS CONFIG
# ============================================

# Créer input.css si n'existe pas
INPUT_CSS="$PROJECT_ROOT/assets/css/input.css"
mkdir -p "$(dirname "$INPUT_CSS")"

if [ ! -f "$INPUT_CSS" ]; then
    cat >"$INPUT_CSS" <<'EOF'
@tailwind base;
@tailwind components;
@tailwind utilities;

/* === CUSTOM STYLES === */
@layer components {
  .btn {
    @apply px-4 py-2 rounded font-medium transition-colors;
  }
  
  .btn-primary {
    @apply bg-blue-600 text-white hover:bg-blue-700;
  }
  
  .btn-secondary {
    @apply bg-gray-600 text-white hover:bg-gray-700;
  }
}
EOF
    echo "✅ Créé: assets/css/input.css"
fi

# Créer tailwind.config.js si n'existe pas
TAILWIND_CONFIG="$PROJECT_ROOT/tailwind.config.js"

if [ ! -f "$TAILWIND_CONFIG" ]; then
    cat >"$TAILWIND_CONFIG" <<'EOF'
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/views/**/*.templ",
    "./internal/views/**/*.go",
    "./templates/**/*.html",
  ],
  theme: {
    extend: {
      colors: {
        terminal: {
          bg: '#0a0a0a',
          text: '#00ff00',
          border: '#00ff0044',
        }
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Courier New', 'monospace'],
      }
    },
  },
  plugins: [],
}
EOF
    echo "✅ Créé: tailwind.config.js"
fi

# ============================================
# 8. CRÉATION SCRIPTS HELPER
# ============================================

# Script build CSS
BUILD_SCRIPT="$PROJECT_ROOT/scripts/build-css.sh"
cat >"$BUILD_SCRIPT" <<'EOF'
#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TAILWIND_BIN="$PROJECT_ROOT/bin/tailwindcss"
INPUT_CSS="$PROJECT_ROOT/assets/css/input.css"
OUTPUT_CSS="$PROJECT_ROOT/public/css/style.css"

if [ ! -f "$TAILWIND_BIN" ]; then
    echo "❌ Tailwind CSS non installé. Run: ./scripts/install-tailwind.sh"
    exit 1
fi

echo "🎨 Building CSS..."
"$TAILWIND_BIN" -i "$INPUT_CSS" -o "$OUTPUT_CSS" --minify
echo "✅ CSS compilé: $OUTPUT_CSS"
EOF
chmod +x "$BUILD_SCRIPT"
echo "✅ Créé: scripts/build-css.sh"

# Script watch CSS
WATCH_SCRIPT="$PROJECT_ROOT/scripts/watch-css.sh"
cat >"$WATCH_SCRIPT" <<'EOF'
#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TAILWIND_BIN="$PROJECT_ROOT/bin/tailwindcss"
INPUT_CSS="$PROJECT_ROOT/assets/css/input.css"
OUTPUT_CSS="$PROJECT_ROOT/public/css/style.css"

if [ ! -f "$TAILWIND_BIN" ]; then
    echo "❌ Tailwind CSS non installé. Run: ./scripts/install-tailwind.sh"
    exit 1
fi

echo "👀 Watching CSS changes..."
echo "   Input:  $INPUT_CSS"
echo "   Output: $OUTPUT_CSS"
echo ""
"$TAILWIND_BIN" -i "$INPUT_CSS" -o "$OUTPUT_CSS" --watch
EOF
chmod +x "$WATCH_SCRIPT"
echo "✅ Créé: scripts/watch-css.sh"

# ============================================
# 9. PREMIÈRE BUILD
# ============================================

echo ""
echo "🎨 Compilation CSS initiale..."
"$TAILWIND_BIN" -i "$INPUT_CSS" -o "$PROJECT_ROOT/public/css/style.css" --minify

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Installation complète!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 Utilisation:"
echo "   Build CSS:  ./scripts/build-css.sh"
echo "   Watch CSS:  ./scripts/watch-css.sh"
echo ""
echo "📁 Fichiers:"
echo "   Binary:     $TAILWIND_BIN"
echo "   Input:      $INPUT_CSS"
echo "   Output:     $PROJECT_ROOT/public/css/style.css"
echo "   Config:     $TAILWIND_CONFIG"
echo ""
