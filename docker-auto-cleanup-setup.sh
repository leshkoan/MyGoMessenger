#!/bin/bash

# Docker Auto-Cleanup Setup Script
# This script sets up automatic Docker cleanup using cron jobs
# Author: MyGoMessenger Project

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if script is run as root (needed for cron setup)
if [ "$EUID" -ne 0 ]; then
    print_warning "Please run as root to setup cron jobs"
    print_info "Usage: sudo $0"
    exit 1
fi

# Check if docker-cleanup.sh exists
if [ ! -f "docker-cleanup.sh" ]; then
    print_error "docker-cleanup.sh not found in current directory"
    exit 1
fi

# Make script executable
chmod +x docker-cleanup.sh

# Setup cron job for daily cleanup at 2 AM
print_info "Setting up daily Docker cleanup at 2:00 AM..."

# Create cron job script
CRON_JOB="0 2 * * * cd $(pwd) && ./docker-cleanup.sh --all --force > /var/log/docker-cleanup.log 2>&1"

# Check if cron job already exists
if crontab -l 2>/dev/null | grep -F "docker-cleanup.sh" >/dev/null 2>&1; then
    print_warning "Docker cleanup cron job already exists"
else
    # Add to crontab
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    print_success "Added daily Docker cleanup cron job"
fi

# Setup log rotation for cleanup logs
print_info "Setting up log rotation..."
LOGROTATE_CONFIG="/etc/logrotate.d/docker-cleanup"

cat > "$LOGROTATE_CONFIG" << EOF
/var/log/docker-cleanup.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        systemctl reload cron || true
    endscript
}
EOF

print_success "Created log rotation configuration"

# Create log directory if it doesn't exist
mkdir -p /var/log
touch /var/log/docker-cleanup.log
chmod 644 /var/log/docker-cleanup.log

print_success "Docker auto-cleanup setup completed!"
print_info ""
print_info "Cleanup schedule:"
print_info "  - Daily at 2:00 AM"
print_info "  - Removes stopped containers, dangling images, and unused networks/volumes"
print_info "  - Logs are saved to /var/log/docker-cleanup.log"
print_info "  - Log rotation: 30 days"
print_info ""
print_info "To view logs: tail -f /var/log/docker-cleanup.log"
print_info "To test cleanup: ./docker-cleanup.sh --dry-run"
print_info "To run cleanup manually: ./docker-cleanup.sh --force"

# Show current cron jobs
print_info ""
print_info "Current cron jobs:"
crontab -l