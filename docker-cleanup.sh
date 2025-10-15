#!/bin/bash

# Docker Cleanup Script
# This script removes unused Docker containers, images, networks, and volumes
# Author: MyGoMessenger Project
# Usage: ./docker-cleanup.sh [options]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to count resources
count_resources() {
    local containers=$(docker ps -aq | wc -l)
    local images=$(docker images -q | wc -l)
    local volumes=$(docker volume ls -q | wc -l)
    local networks=$(docker network ls -q | wc -l)

    echo "Current Docker resources:"
    echo "  Containers: $containers"
    echo "  Images: $images"
    echo "  Volumes: $volumes"
    echo "  Networks: $networks"
}

# Function to remove stopped containers
cleanup_stopped_containers() {
    print_info "Removing stopped containers..."
    local stopped_containers=$(docker ps -aq --filter "status=exited" --filter "status=created")
    if [ -n "$stopped_containers" ]; then
        echo "Removing containers: $stopped_containers"
        docker rm $stopped_containers
        print_success "Removed $(echo "$stopped_containers" | wc -w) stopped containers"
    else
        print_info "No stopped containers to remove"
    fi
}

# Function to remove dangling images
cleanup_dangling_images() {
    print_info "Removing dangling images..."
    local dangling_images=$(docker images -f "dangling=true" -q)
    if [ -n "$dangling_images" ]; then
        echo "Removing dangling images: $dangling_images"
        docker rmi $dangling_images
        print_success "Removed $(echo "$dangling_images" | wc -w) dangling images"
    else
        print_info "No dangling images to remove"
    fi
}

# Function to remove unused networks
cleanup_unused_networks() {
    print_info "Removing unused networks..."
    local project_name=$(basename $(pwd) | tr '[:upper:]' '[:lower:]')
    local unused_networks=$(docker network ls -q --filter "label=com.docker.compose.project=${project_name}" | xargs docker network inspect 2>/dev/null | jq -r '.[]? | select(.Containers | length == 0) | .Id' 2>/dev/null || true)

    if [ -n "$unused_networks" ]; then
        echo "Removing unused networks: $unused_networks"
        echo "$unused_networks" | xargs docker network rm
        print_success "Removed $(echo "$unused_networks" | wc -w) unused networks"
    else
        print_info "No unused networks to remove"
    fi
}

# Function to remove unused volumes
cleanup_unused_volumes() {
    print_info "Removing unused volumes..."
    local project_name=$(basename $(pwd) | tr '[:upper:]' '[:lower:]')
    local used_volumes=$(docker ps -aq --filter "label=com.docker.compose.project=${project_name}" | xargs docker inspect 2>/dev/null | jq -r '.[].Mounts[].Name' 2>/dev/null | sort | uniq || true)
    local all_volumes=$(docker volume ls -q)

    local unused_volumes=""
    for volume in $all_volumes; do
        if ! echo "$used_volumes" | grep -q "^${volume}$"; then
            unused_volumes="$unused_volumes $volume"
        fi
    done

    if [ -n "$unused_volumes" ]; then
        echo "Removing unused volumes: $unused_volumes"
        echo "$unused_volumes" | xargs docker volume rm
        print_success "Removed $(echo "$unused_volumes" | wc -w) unused volumes"
    else
        print_info "No unused volumes to remove"
    fi
}

# Function to show system disk usage
show_disk_usage() {
    print_info "Docker disk usage:"
    docker system df
}

# Function to show help
show_help() {
    echo "Docker Cleanup Script"
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -d, --dry-run       Show what would be cleaned up without doing it"
    echo "  -f, --force         Force cleanup without confirmation"
    echo "  -c, --containers    Clean only stopped containers"
    echo "  -i, --images        Clean only dangling images"
    echo "  -n, --networks      Clean only unused networks"
    echo "  -v, --volumes       Clean only unused volumes"
    echo "  --all              Clean everything (default)"
    echo ""
    echo "Examples:"
    echo "  $0                  # Clean everything"
    echo "  $0 --dry-run        # Show what would be cleaned"
    echo "  $0 -c -i           # Clean only containers and images"
}

# Parse command line arguments
DRY_RUN=false
FORCE=false
CLEAN_CONTAINERS=false
CLEAN_IMAGES=false
CLEAN_NETWORKS=false
CLEAN_VOLUMES=false
CLEAN_ALL=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -c|--containers)
            CLEAN_CONTAINERS=true
            CLEAN_ALL=false
            shift
            ;;
        -i|--images)
            CLEAN_IMAGES=true
            CLEAN_ALL=false
            shift
            ;;
        -n|--networks)
            CLEAN_NETWORKS=true
            CLEAN_ALL=false
            shift
            ;;
        -v|--volumes)
            CLEAN_VOLUMES=true
            CLEAN_ALL=false
            shift
            ;;
        --all)
            CLEAN_ALL=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
main() {
    print_info "Starting Docker cleanup process..."

    # Show initial state
    count_resources
    echo ""

    # Confirmation unless forced
    if [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
        print_warning "This will remove unused Docker resources. Continue? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_info "Cleanup cancelled"
            exit 0
        fi
    fi

    # Perform cleanup
    if [ "$DRY_RUN" = true ]; then
        print_info "DRY RUN - Showing what would be cleaned:"
        cleanup_stopped_containers
        cleanup_dangling_images
        cleanup_unused_networks
        cleanup_unused_volumes
    else
        if [ "$CLEAN_ALL" = true ] || [ "$CLEAN_CONTAINERS" = true ]; then
            cleanup_stopped_containers
        fi

        if [ "$CLEAN_ALL" = true ] || [ "$CLEAN_IMAGES" = true ]; then
            cleanup_dangling_images
        fi

        if [ "$CLEAN_ALL" = true ] || [ "$CLEAN_NETWORKS" = true ]; then
            cleanup_unused_networks
        fi

        if [ "$CLEAN_ALL" = true ] || [ "$CLEAN_VOLUMES" = true ]; then
            cleanup_unused_volumes
        fi

        # Show final state
        echo ""
        print_info "Final Docker resources:"
        count_resources

        # Show disk usage after cleanup
        echo ""
        show_disk_usage
    fi

    print_success "Docker cleanup completed!"
}

# Run main function
main "$@"