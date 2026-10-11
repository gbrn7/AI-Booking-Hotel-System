#!/usr/bin/env bash
# ==============================================================================
# Unified Container Management Helper for Hotel AI Booking System
# Automatically detects and supports Docker or Podman
# ==============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# ------------------------------------------------------------------------------
# 1. Detect Container Engine & Compose Tool
# ------------------------------------------------------------------------------
detect_tools() {
    if command -v podman-compose >/dev/null 2>&1; then
        COMPOSE_CMD="podman-compose"
        ENGINE_CMD="podman"
    elif command -v podman >/dev/null 2>&1 && podman compose version >/dev/null 2>&1; then
        COMPOSE_CMD="podman compose"
        ENGINE_CMD="podman"
    elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
        COMPOSE_CMD="docker compose"
        ENGINE_CMD="docker"
    elif command -v docker-compose >/dev/null 2>&1; then
        COMPOSE_CMD="docker-compose"
        ENGINE_CMD="docker"
    elif command -v podman >/dev/null 2>&1; then
        COMPOSE_CMD="podman"
        ENGINE_CMD="podman"
    elif command -v docker >/dev/null 2>&1; then
        COMPOSE_CMD="docker"
        ENGINE_CMD="docker"
    else
        echo "❌ Error: Tidak ditemukan 'docker' maupun 'podman' di sistem Anda."
        echo "   Silakan instal Docker atau Podman terlebih dahulu."
        exit 1
    fi
}

detect_tools

# ------------------------------------------------------------------------------
# 2. Print Header & Prerequisite Checks
# ------------------------------------------------------------------------------
print_banner() {
    echo "=================================================================="
    echo "  🏨 Hotel AI Booking Platform — Container Management"
    echo "  Engine  : $ENGINE_CMD"
    echo "  Compose : $COMPOSE_CMD"
    echo "=================================================================="
}

ensure_env_file() {
    if [ ! -f "$SCRIPT_DIR/.env" ]; then
        echo "⚠️  File .env tidak ditemukan. Menyalin dari .env.example..."
        cp "$SCRIPT_DIR/.env.example" "$SCRIPT_DIR/.env"
    fi
}

ensure_data_dirs() {
    mkdir -p "$SCRIPT_DIR/data/postgres" "$SCRIPT_DIR/data/evolution"
    ensure_env_file
}

# ------------------------------------------------------------------------------
# 3. Command Handlers
# ------------------------------------------------------------------------------
cmd_help() {
    print_banner
    echo "Penggunaan: ./container.sh <perintah>"
    echo ""
    echo "Perintah yang Tersedia:"
    echo "  build        Membangun container image (Backend & Frontend)"
    echo "  up           Jalankan DB + Backend + Frontend"
    echo "  up:all       Jalankan SEMUA service (termasuk WhatsApp Evolution API)"
    echo "  down         Hentikan semua container dan jaringan"
    echo "  logs         Pantau log real-time container"
    echo "  cli          Buka interactive CLI simulator di backend"
    echo "  status       Periksa status dan healthcheck container"
    echo "  help         Tampilkan pesan bantuan ini"
    echo ""
}

cmd_build() {
    print_banner
    echo "🔨 Membangun container image..."
    $COMPOSE_CMD build
    echo "✅ Build selesai!"
}

cmd_up() {
    print_banner
    ensure_data_dirs
    echo "🚀 Menjalankan Database, Backend Go, dan Frontend Next.js..."
    $COMPOSE_CMD up -d db backend frontend
    echo ""
    echo "✅ Layanan berjalan!"
    echo "   - Landing & Web Chat : http://localhost:3000"
    echo "   - Admin Dashboard    : http://localhost:3000/admin"
    echo "   - Go API Server      : http://localhost:8080"
}

cmd_up_all() {
    print_banner
    ensure_data_dirs
    echo "🚀 Menjalankan SEMUA layanan (termasuk WhatsApp Evolution API)..."
    $COMPOSE_CMD --profile full up -d
    echo ""
    echo "✅ Semua layanan berjalan!"
    echo "   - Frontend : http://localhost:3000"
    echo "   - Backend  : http://localhost:8080"
    echo "   - WhatsApp : http://localhost:8083"
}

cmd_down() {
    print_banner
    echo "🛑 Menghentikan seluruh container..."
    $COMPOSE_CMD down
    echo "✅ Semua container dihentikan."
}

cmd_logs() {
    print_banner
    SERVICE="${1:-}"
    if [ -n "$SERVICE" ]; then
        $COMPOSE_CMD logs -f "$SERVICE"
    else
        $COMPOSE_CMD logs -f backend frontend
    fi
}

cmd_cli() {
    print_banner
    echo "💬 Menjalankan simulator AI Sales Agent CLI..."
    if $ENGINE_CMD ps | grep -q "hotel-backend"; then
        $ENGINE_CMD exec -it hotel-backend /app/cli
    else
        echo "ℹ️  Container backend tidak aktif, mencoba menjalankan CLI lokal..."
        (cd backend && go run ./cmd/cli)
    fi
}

cmd_status() {
    print_banner
    $COMPOSE_CMD ps
}

# ------------------------------------------------------------------------------
# 4. Entrypoint
# ------------------------------------------------------------------------------
case "${1:-help}" in
    build)
        cmd_build
        ;;
    up)
        cmd_up
        ;;
    up:all)
        cmd_up_all
        ;;
    down)
        cmd_down
        ;;
    logs)
        cmd_logs "${2:-}"
        ;;
    cli)
        cmd_cli
        ;;
    status)
        cmd_status
        ;;
    help|--help|-h)
        cmd_help
        ;;
    *)
        echo "❌ Perintah tidak dikenal: $1"
        cmd_help
        exit 1
        ;;
esac
