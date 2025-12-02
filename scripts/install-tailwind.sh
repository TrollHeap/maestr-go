#!/usr/bin/env bash

set -e

# ============================================
# TAILWIND CSS STANDALONE INSTALLER
# ============================================

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$PROJECT_ROOT/bin"
TAILWIND_VERSION="v4.1.17" # ✅ Dernière stable v4
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
    *)
        echo "❌ Architecture non supportée: $ARCH"
        exit 1
        ;;
esac

# ============================================
# 2. URL TÉLÉCHARGEMENT
# ============================================

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

echo "⬇️  Téléchargement Tailwind CSS v4.1.17..."

if command -v curl &>/dev/null; then
    curl -sL "$DOWNLOAD_URL" -o "$TAILWIND_BIN"
elif command -v wget &>/dev/null; then
    wget -q "$DOWNLOAD_URL" -O "$TAILWIND_BIN"
else
    echo "❌ curl ou wget requis"
    exit 1
fi

# ============================================
# 5. PERMISSIONS
# ============================================

chmod +x "$TAILWIND_BIN"

echo "✅ Tailwind CSS installé: $TAILWIND_BIN"

# ============================================
# 6. VÉRIFICATION
# ============================================

if [ -f "$TAILWIND_BIN" ]; then
    VERSION=$("$TAILWIND_BIN" --help 2>&1 | head -n 1 || echo "tailwindcss v4")
    echo "✅ Version: $VERSION"
else
    echo "❌ Installation échouée"
    exit 1
fi

# ============================================
# 7. INPUT CSS (V4 SYNTAX)
# ============================================

INPUT_CSS="$PROJECT_ROOT/public/css/input.css"

if [ ! -f "$INPUT_CSS" ]; then
    cat >"$INPUT_CSS" <<'EOF'
@import "tailwindcss";

/* === TERMINAL THEME === */
@theme {
  --color-terminal-bg: #0a0a0a;
  --color-terminal-text: #00ff00;
  --color-terminal-border: #00ff0044;
}

/* === CUSTOM UTILITIES === */
@layer components {
  .terminal-nav {
    @apply bg-black border-b-2 border-green-500/30 px-4 py-3;
  }
  
  .nav-link {
    @apply px-4 py-2 text-green-400 hover:bg-green-500/10 
           transition-colors border border-green-500/30 rounded;
  }
  
  .terminal-header {
    @apply bg-black border-b-2 border-green-500/30 px-6 py-4;
  }
  
  .btn {
    @apply px-4 py-2 rounded font-medium transition-all
           border-2 border-green-500/50 hover:border-green-500
           hover:shadow-[0_0_10px_rgba(0,255,0,0.5)];
  }
}

/* === ANIMATIONS === */
@keyframes scan {
  0% { transform: translateY(-100%); }
  100% { transform: translateY(100vh); }
}

.scan-line {
  animation: scan 8s linear infinite;
}
EOF
    echo "✅ Créé: public/css/input.css"
fi

# ============================================
# 8. TAILWIND CONFIG (V4)
# ============================================

TAILWIND_CONFIG="$PROJECT_ROOT/tailwind.config.js"

if [ ! -f "$TAILWIND_CONFIG" ]; then
    cat >"$TAILWIND_CONFIG" <<'EOF'
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./internal/views/**/*.{templ,go}",
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
    },
  },
}
EOF
    echo "✅ Créé: tailwind.config.js"
fi

# ============================================
# 9. SCRIPTS HELPER
# ============================================

BUILD_SCRIPT="$PROJECT_ROOT/scripts/build-css.sh"
cat >"$BUILD_SCRIPT" <<'EOF'
#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TAILWIND_BIN="$PROJECT_ROOT/bin/tailwindcss"
CSS_ROOT="$PROJECT_ROOT/public/css"
INPUT_CSS="$CSS_ROOT/input.css"
OUTPUT_CSS="$CSS_ROOT/style.css"


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

WATCH_SCRIPT="$PROJECT_ROOT/scripts/watch-css.sh"
cat >"$WATCH_SCRIPT" <<'EOF'
#!/usr/bin/env bash
set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TAILWIND_BIN="$PROJECT_ROOT/bin/tailwindcss"
CSS_ROOT="$PROJECT_ROOT/public/css"
INPUT_CSS="$CSS_ROOT/input.css"
OUTPUT_CSS="$CSS_ROOT/style.css"

if [ ! -f "$TAILWIND_BIN" ]; then
    echo "❌ Tailwind CSS non installé. Run: ./scripts/install-tailwind.sh"
    exit 1
fi

echo "👀 Watching CSS changes..."
"$TAILWIND_BIN" -i "$INPUT_CSS" -o "$OUTPUT_CSS" --watch
EOF
chmod +x "$WATCH_SCRIPT"
echo "✅ Créé: scripts/watch-css.sh"

# ============================================
# 10. BUILD INITIAL
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
echo ""
