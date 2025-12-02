#!/usr/bin/env bash

set -euo pipefail
PATH="/opt/homebrew/bin:$PATH"

# fpm-based packager for KrankyBearClock
# - macOS: Packages bin/KrankyBearClock-macos-{arch} as .pkg installer to /Applications
# - Linux: Packages bin/KrankyBearClock-linux-amd64 as .deb and .rpm installers
# - macOS: Installs as KrankyBearClock.app bundle structure
# - Linux: Includes entire Resources/ tree -> /opt/local/bin/Resources
# - Outputs to ./installers (configurable)

usage() {
  cat <<EOF
Usage: ./package.sh [linux|mac|all] [ENV_VARS]

Arguments:
  linux       Build Linux packages (.deb and .rpm)
  mac         Build macOS package (.pkg)
  all         Build both Linux and macOS packages

Environment variables (optional):
  VERSION     Package version (default: 0.4.4)
  ITERATION   Package iteration/release (default: 1)
  ARCH        Target arch: amd64|arm64 (default: native for mac, amd64 for linux)
  OUTDIR      Output directory (default: ./installers)
  MAINTAINER  Maintainer (default: amarillier@gmail.com)
  VENDOR      Vendor (default: KrankyBear)
  URL         Project URL (default: https://github.com/amarillier/KrankyBearClock)
  LICENSE     License (default: GNU GPL v3)
  REPLACES    Comma-separated list of packages this replaces (default: KrankyBearClock)
              Allows overwriting shared files like /opt/local/bin/Resources/LICENSE

Examples:
  # Build macOS package
  ./package.sh mac
  ./package.sh mac VERSION=1.2.3 ARCH=amd64
  
  # Build Linux packages
  ./package.sh linux
  ./package.sh linux VERSION=1.2.3 ARCH=amd64
  
  # Build both
  ./package.sh all VERSION=1.2.3
EOF
}

# Parse command-line arguments
if [[ $# -eq 0 ]]
then
  usage
  exit 0
fi

TYPE_ARG="${1:-}"
shift  # Remove first argument, keep rest for env vars

# Validate TYPE argument
case "$TYPE_ARG" in
  linux|mac|all)
    # Valid argument
    ;;
  -h|--help|-?)
    usage
    exit 0
    ;;
  *)
    echo "Error: Invalid argument '$TYPE_ARG'. Must be 'linux', 'mac', or 'all'." >&2
    echo ""
    usage
    exit 1
    ;;
esac

# Process remaining arguments as environment variable assignments
# This allows: ./package.sh linux VERSION=1.2.3 ARCH=amd64
for arg in "$@"
do
  if [[ "$arg" == *"="* ]]
  then
    key="${arg%%=*}"
    value="${arg#*=}"
    export "$key=$value"
  fi
done

# Configurable via env vars
NAME=${NAME:-krankybearclock}
VERSION=${VERSION:-0.4.4}
ITERATION=${ITERATION:-1}
OUTDIR=${OUTDIR:-./installers}
MAINTAINER=${MAINTAINER:-"amarillier@gmail.com"}
VENDOR=${VENDOR:-"KrankyBear"}
URL=${URL:-"https://github.com/amarillier/KrankyBearClock"}
LICENSE=${LICENSE:-"GNU GPL v3"}
# Packages this can replace (allows overwriting shared files like LICENSE)
# Default covers every app that installs /opt/local/bin/Resources/*
# Add comma-separated values via REPLACES env to extend this list.
DEFAULT_REPLACES="KrankyBearClock,krankybearclock,KrankyBearNotify,krankybearnotify"
REPLACES=${REPLACES:-"$DEFAULT_REPLACES"}

# Efficiently copy directory trees while avoiding macOS extended attributes.
# Prefers rsync (to preserve permissions) but falls back to tar streams if unavailable.
copy_tree() {
  local src="$1"
  local dest="$2"
  if [[ ! -d "$src" ]]
  then
    echo "Error: Missing directory $src" >&2
    exit 1
  fi
  mkdir -p "$dest"
  if command -v rsync >/dev/null 2>&1
  then
    rsync -a --delete "$src"/ "$dest"/
  else
    (cd "$src" && tar -cf - .) | (cd "$dest" && tar -xf -)
  fi
}

# Function to build packages for a specific type
build_package() {
  local TYPE="$1"
  
  # Architecture handling
  local ARCH="${ARCH:-}"
  if [[ -z "$ARCH" ]]
  then
    if [[ "$TYPE" == "mac" ]]
    then
      # Default to native macOS arch
      ARCH=$(uname -m)
      [[ "$ARCH" == "x86_64" ]] && ARCH=amd64
    else
      ARCH=amd64
    fi
  fi

  case "$ARCH" in
    amd64|x86_64)
      DEB_ARCH=amd64
      RPM_ARCH=x86_64
      PKG_ARCH=amd64
      ;;
    arm64|aarch64)
      DEB_ARCH=arm64
      RPM_ARCH=aarch64
      PKG_ARCH=arm64
      ;;
    *)
      # Fall back to using the same string
      DEB_ARCH="$ARCH"
      RPM_ARCH="$ARCH"
      PKG_ARCH="$ARCH"
      ;;
  esac

  # Source assets - depends on TYPE
  # Determine binary name early for macOS staging
  BIN_NAME="KrankyBearClock"
  SRC_RESOURCES="Resources"
  SRC_INFO_PLIST="Info-plist.txt"
  SRC_README_PLIST="Readme-plist.txt"
  SRC_RELEASE_NOTES="ReleaseNotes.txt"
  if [[ -f "LICENSE" ]]
  then
    SRC_LICENSE_FILE="LICENSE"
  elif [[ -f "license.txt" ]]
  then
    SRC_LICENSE_FILE="license.txt"
  else
    SRC_LICENSE_FILE=""
  fi
  
  if [[ "$TYPE" == "mac" ]]
  then
    SRC_BIN="bin/KrankyBearClock-macos-${PKG_ARCH}"
  else
    SRC_BIN="bin/KrankyBearClock-linux-${DEB_ARCH}"
  fi

  # Validate sources
  if [[ ! -f "$SRC_BIN" ]]
  then
    if [[ "$TYPE" == "mac" ]]
    then
      echo "Error: Missing $SRC_BIN (build your macOS binary first with ./compile-mac.sh)." >&2
    else
      echo "Error: Missing $SRC_BIN (build your Linux binary first)." >&2
    fi
    exit 1
  fi
  if [[ "$TYPE" == "mac" ]]
  then
    if [[ ! -d "$SRC_RESOURCES" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES" >&2
      exit 1
    fi
    if [[ ! -d "$SRC_RESOURCES/Images" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES/Images" >&2
      exit 1
    fi
    if [[ ! -d "$SRC_RESOURCES/Sounds" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES/Sounds" >&2
      exit 1
    fi
    if [[ ! -f "$SRC_INFO_PLIST" ]]
    then
      echo "Error: Missing file $SRC_INFO_PLIST" >&2
      exit 1
    fi
  else
    if [[ ! -d "$SRC_RESOURCES" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES" >&2
      exit 1
    fi
    if [[ ! -d "$SRC_RESOURCES/Images" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES/Images" >&2
      exit 1
    fi
    if [[ ! -d "$SRC_RESOURCES/Sounds" ]]
    then
      echo "Error: Missing directory $SRC_RESOURCES/Sounds" >&2
      exit 1
    fi
    if [[ ! -f "$SRC_RELEASE_NOTES" ]]
    then
      echo "Error: Missing file $SRC_RELEASE_NOTES" >&2
      exit 1
    fi
    if [[ ! -f "$SRC_LICENSE_FILE" ]]
    then
      echo "Error: Missing file $SRC_LICENSE_FILE" >&2
      exit 1
    fi
  fi

  mkdir -p "$OUTDIR"

  # For macOS: Build complete .app bundle, then use it as source for fpm
  # For Linux packages built on macOS, create staging directory with files
  # stripped of extended attributes to avoid tar header compatibility issues
  STAGING_DIR=""
  APP_BUNDLE=""
  SRC_LICENSE_STAGED=""
  
  if [[ "$TYPE" == "mac" ]]
  then
    # Build the complete .app bundle structure
    APP_BUNDLE="KrankyBearClock.app"
    echo "Building .app bundle: $APP_BUNDLE..."
    
    # Remove existing bundle if present
    rm -rf "$APP_BUNDLE"
    
    # Create app bundle structure
    mkdir -p "$APP_BUNDLE/Contents/MacOS/Resources"
    
    # Copy binary to Contents/MacOS
    cp "$SRC_BIN" "$APP_BUNDLE/Contents/MacOS/$BIN_NAME"
    
    # Create symlink: KrankyBearClock -> KrankyBearClock in Contents/MacOS
    # Use -sf to force overwrite if it already exists
    ln -sf "$BIN_NAME" "$APP_BUNDLE/Contents/MacOS/clock"
    
    # Copy all Resources (Images, Sounds, data files, etc.)
    copy_tree "$SRC_RESOURCES" "$APP_BUNDLE/Contents/MacOS/Resources"
    # Ensure canonical files exist for downstream scripts
    if [[ -f "$APP_BUNDLE/Contents/MacOS/Resources/ReleaseNotes.txt" ]]
    then
      cp "$APP_BUNDLE/Contents/MacOS/Resources/ReleaseNotes.txt" \
         "$APP_BUNDLE/Contents/MacOS/Resources/ReleaseNotes-KrankyBearClock.txt"
    elif [[ -f "$SRC_RELEASE_NOTES" ]]
    then
      cp "$SRC_RELEASE_NOTES" \
         "$APP_BUNDLE/Contents/MacOS/Resources/ReleaseNotes-KrankyBearClock.txt"
    fi
    if [[ -f "$APP_BUNDLE/Contents/MacOS/Resources/license.txt" ]] && [[ ! -f "$APP_BUNDLE/Contents/MacOS/Resources/LICENSE" ]]
    then
      cp "$APP_BUNDLE/Contents/MacOS/Resources/license.txt" \
         "$APP_BUNDLE/Contents/MacOS/Resources/LICENSE"
    elif [[ -f "$SRC_LICENSE_FILE" ]]
    then
      cp "$SRC_LICENSE_FILE" \
         "$APP_BUNDLE/Contents/MacOS/Resources/LICENSE"
    fi
    
    # Set proper permissions on Resources directories (755 = rwxr-xr-x)
    chmod -R 755 "$APP_BUNDLE/Contents/MacOS/Resources"
    
    # Copy Info-plist.txt to Contents
    cp "$SRC_INFO_PLIST" "$APP_BUNDLE/Contents/Info-plist.txt"
    cp "$SRC_README_PLIST" "$APP_BUNDLE/Contents/Readme-plist.txt"
    
    # Verify files were copied correctly
    if [[ ! -d "$APP_BUNDLE/Contents/MacOS/Resources/Images" ]] || [[ -z "$(ls -A "$APP_BUNDLE/Contents/MacOS/Resources/Images" 2>/dev/null)" ]]
    then
      echo "Error: Images directory is empty or missing in app bundle" >&2
      exit 1
    fi
    if [[ ! -d "$APP_BUNDLE/Contents/MacOS/Resources/Sounds" ]] || [[ -z "$(ls -A "$APP_BUNDLE/Contents/MacOS/Resources/Sounds" 2>/dev/null)" ]]
    then
      echo "Error: Sounds directory is empty or missing in app bundle" >&2
      exit 1
    fi
    
    echo "App bundle created successfully."
    
  elif [[ "$TYPE" == "linux" ]]
  then
    echo "Creating staging directory for Linux payload..."
    STAGING_DIR=$(mktemp -d -t fpm-staging.XXXXXX)
    cleanup_staging() { rm -rf "$STAGING_DIR"; }
    trap cleanup_staging EXIT INT TERM
    
    # Copy binary and ensure executable permissions
    install -m 755 "$SRC_BIN" "$STAGING_DIR/KrankyBearClock"
    
    # Provide a convenience symlink for legacy launch scripts
    ln -sf "KrankyBearClock" "$STAGING_DIR/clock"
    
    # Copy the entire Resources tree (Images, Sounds, metadata, etc.)
    copy_tree "$SRC_RESOURCES" "$STAGING_DIR/Resources"
    
    # Ensure ReleaseNotes alias exists for backwards compatibility
    if [[ -f "$SRC_RELEASE_NOTES" ]]
    then
      cp "$SRC_RELEASE_NOTES" "$STAGING_DIR/Resources/ReleaseNotes-KrankyBearClock.txt"
      chmod 644 "$STAGING_DIR/Resources/ReleaseNotes-KrankyBearClock.txt"
    fi
    
    # Ensure LICENSE is present with expected casing
    SRC_LICENSE_STAGED=""
    if [[ -f "$STAGING_DIR/Resources/license.txt" ]] && [[ ! -f "$STAGING_DIR/Resources/LICENSE" ]]
    then
      cp "$STAGING_DIR/Resources/license.txt" "$STAGING_DIR/Resources/LICENSE"
    elif [[ -f "$SRC_LICENSE_FILE" ]]
    then
      cp "$SRC_LICENSE_FILE" "$STAGING_DIR/Resources/LICENSE"
    fi
    if [[ -f "$STAGING_DIR/Resources/LICENSE" ]]
    then
      chmod 644 "$STAGING_DIR/Resources/LICENSE"
      SRC_LICENSE_STAGED="$STAGING_DIR/LICENSE-config"
      mv "$STAGING_DIR/Resources/LICENSE" "$SRC_LICENSE_STAGED"
    fi
    
    # Remove macOS extended attributes when building from Darwin hosts
    if command -v xattr >/dev/null 2>&1
    then
      xattr -cr "$STAGING_DIR" 2>/dev/null || true
      find "$STAGING_DIR" -name '._*' -delete 2>/dev/null || true
    fi
    if [[ "$(uname -s)" == "Darwin" ]]
    then
      export COPYFILE_DISABLE=1
    fi
    
    # Update source paths to point to staging directory
    SRC_BIN="$STAGING_DIR/KrankyBearClock"
    SRC_SYMLINK="$STAGING_DIR/clock"
    SRC_RESOURCES_STAGED="$STAGING_DIR/Resources"
  fi

  COMMON_ARGS=(
    -s dir
    -n "$NAME"
    -v "$VERSION"
    --iteration "$ITERATION"
    --maintainer "$MAINTAINER"
    --vendor "$VENDOR"
    --url "$URL"
    --license "$LICENSE"
    --description "KrankyBear Clock - A cross-platform GUI desktop clock application"
    -f
  )
  if [[ -n "${FPM_DEBUG_WORKSPACE:-}" ]]
  then
    COMMON_ARGS+=(--debug-workspace)
    echo "FPM debug workspace enabled; temp build dirs will be preserved."
  fi

  if [[ "$TYPE" == "mac" ]]
  then
    # macOS .pkg installer - installs to /Applications as .app bundle
    PKG_OUTFILE="$OUTDIR/KrankyBearClock_${VERSION}-${ITERATION}_${PKG_ARCH}.pkg"
    echo "Building macOS .pkg ($PKG_ARCH) -> $PKG_OUTFILE..."
    
    # macOS .app bundle structure:
    # /Applications/KrankyBearClock.app/Contents/Info-plist.txt (sample Info.plist)
    # /Applications/KrankyBearClock.app/Contents/MacOS/KrankyBearClock (executable)
    # /Applications/KrankyBearClock.app/Contents/MacOS/Resources/Images (resources)
    # /Applications/KrankyBearClock.app/Contents/MacOS/Resources/Sounds (sounds)
    # Note: App looks for Resources/Images and Resources/Sounds relative to executable
    APP_NAME="KrankyBearClock.app"
    APP_DIR="/Applications/$APP_NAME"
    CONTENTS_DIR="$APP_DIR/Contents"
    MACOS_DIR="$CONTENTS_DIR/MacOS"
    RESOURCES_DIR="$MACOS_DIR/Resources"
    
    # Note: BIN_NAME is already defined earlier
    # Map bundle contents as directories/files
    FPM_FILES=(
      "$APP_BUNDLE/Contents/MacOS/$BIN_NAME=$MACOS_DIR/$BIN_NAME"
      "$APP_BUNDLE/Contents/MacOS/clock=$MACOS_DIR/clock"
      "$APP_BUNDLE/Contents/Info-plist.txt=$CONTENTS_DIR/Info-plist.txt"
      "$APP_BUNDLE/Contents/Readme-plist.txt=$CONTENTS_DIR/Readme-plist.txt"
    )

    while IFS= read -r -d '' resource_path
    do
      rel_path="${resource_path#$APP_BUNDLE/Contents/MacOS/Resources/}"
      FPM_FILES+=("$resource_path=$RESOURCES_DIR/$rel_path")
    done < <(find "$APP_BUNDLE/Contents/MacOS/Resources" -type f -print0)
    
    # Validate file mappings before running fpm
    if ! validate_fpm_files "${FPM_FILES[@]}"
    then
      echo "Aborting .pkg build due to validation errors." >&2
      exit 1
    fi
    
    fpm \
      "${COMMON_ARGS[@]}" \
      -t osxpkg \
      -a "$PKG_ARCH" \
      --directories "$APP_DIR" \
      --directories "$CONTENTS_DIR" \
      --directories "$MACOS_DIR" \
      --directories "$RESOURCES_DIR" \
      --package "$PKG_OUTFILE" \
      "${FPM_FILES[@]}"
    
    ./setIcon.sh Resources/Images/KrankyBearBeret.png "$PKG_OUTFILE"
    echo ""
    echo "Done. Package created:"
    echo "  $PKG_OUTFILE"
    
  else
    # Linux .deb and .rpm installers
    DEB_OUTFILE="$OUTDIR/KrankyBearClock_${VERSION}-${ITERATION}_${DEB_ARCH}.deb"
    RPM_OUTFILE="$OUTDIR/KrankyBearClock_${VERSION}-${ITERATION}_${RPM_ARCH}.rpm"
    
    echo "Building .deb ($DEB_ARCH) -> $DEB_OUTFILE..."
    # Build fpm args array for .deb
    DEB_ARGS=(
      "${COMMON_ARGS[@]}"
      -t deb
      -a "$DEB_ARCH"
      --deb-no-default-config-files
      --directories /opt/local/bin
      --directories /opt/local/bin/Resources
    )
    
    # Add --replaces if REPLACES is set (allows overwriting shared files)
    # fpm's --replaces works for both deb and rpm packages
    if [[ -n "$REPLACES" ]]
    then
      # Convert comma-separated list and add --replaces for each package
      # Save and restore IFS to avoid side effects
      OLD_IFS="$IFS"
      IFS=',' read -ra REPLACES_ARRAY <<< "$REPLACES"
      IFS="$OLD_IFS"
      for pkg in "${REPLACES_ARRAY[@]}"
      do
        # Trim whitespace from package name
        pkg=$(printf '%s' "$pkg" | sed -e 's/^ *//' -e 's/ *$//')
        [[ -n "$pkg" ]] && DEB_ARGS+=(--replaces "$pkg")
      done
      echo "  Note: Package will replace: $REPLACES (allows overwriting shared LICENSE file)"
    fi
    
    # Mark LICENSE file as a config file so it's preserved on uninstall
    # This allows the LICENSE to remain even when this package is removed,
    # since it's shared between multiple KrankyBear packages
    if [[ -n "$SRC_LICENSE_STAGED" ]]
    then
      DEB_ARGS+=(--config-files "/opt/local/bin/Resources/LICENSE")
      echo "  Note: LICENSE file will be preserved on uninstall (marked as config file)"
    fi
    
    # Build fpm file list (only include symlink if it exists)
    DEB_FPM_FILES=(
      "$SRC_BIN=/opt/local/bin/KrankyBearClock"
      "$SRC_RESOURCES_STAGED=/opt/local/bin"
    )
    # Only add symlink if it was created
    if [[ -n "$SRC_SYMLINK" ]]
    then
      DEB_FPM_FILES+=("$SRC_SYMLINK=/opt/local/bin/clock")
    fi
    if [[ -n "$SRC_LICENSE_STAGED" ]]
    then
      DEB_FPM_FILES+=("$SRC_LICENSE_STAGED=/opt/local/bin/Resources/LICENSE")
    fi
    
    # Validate file mappings before running fpm
    if ! validate_fpm_files "${DEB_FPM_FILES[@]}"
    then
      echo "Aborting .deb build due to validation errors." >&2
      exit 1
    fi
    
    fpm \
      "${DEB_ARGS[@]}" \
      --package "$DEB_OUTFILE" \
      "${DEB_FPM_FILES[@]}"
    
    echo ""
    echo "Building .rpm ($RPM_ARCH) -> $RPM_OUTFILE..."
    # Build RPM package with all files including symlink
    # Remove --directories flags to avoid "File listed twice" warnings
    # RPM will auto-create directories from file paths
    RPM_ARGS=(
      "${COMMON_ARGS[@]}"
      -t rpm
      -a "$RPM_ARCH"
      --rpm-os linux
      --rpm-auto-add-directories
    )
    
    # Add --replaces if REPLACES is set (allows overwriting shared files)
    # fpm's --replaces works for both deb and rpm packages
    if [[ -n "$REPLACES" ]]
    then
      # Convert comma-separated list and add --replaces for each package
      # Save and restore IFS to avoid side effects
      OLD_IFS="$IFS"
      IFS=',' read -ra REPLACES_ARRAY <<< "$REPLACES"
      IFS="$OLD_IFS"
      for pkg in "${REPLACES_ARRAY[@]}"
      do
        # Trim whitespace from package name
        pkg=$(printf '%s' "$pkg" | sed -e 's/^ *//' -e 's/ *$//')
        [[ -n "$pkg" ]] && RPM_ARGS+=(--replaces "$pkg")
      done
      echo "  Note: Package will replace: $REPLACES (allows overwriting shared LICENSE file)"
    fi
    
    # Mark LICENSE file as a config file so it's preserved on uninstall
    # This allows the LICENSE to remain even when this package is removed,
    # since it's shared between multiple KrankyBear packages
    if [[ -n "$SRC_LICENSE_STAGED" ]]
    then
      RPM_ARGS+=(--config-files "/opt/local/bin/Resources/LICENSE")
      echo "  Note: LICENSE file will be preserved on uninstall (marked as config file)"
    fi
    
    # Build fpm file list (only include symlink if it exists)
    RPM_FPM_FILES=(
      "$SRC_BIN=/opt/local/bin/KrankyBearClock"
      "$SRC_RESOURCES_STAGED=/opt/local/bin"
    )
    # Only add symlink if it was created
    if [[ -n "$SRC_SYMLINK" ]]
    then
      RPM_FPM_FILES+=("$SRC_SYMLINK=/opt/local/bin/clock")
    fi
    if [[ -n "$SRC_LICENSE_STAGED" ]]
    then
      RPM_FPM_FILES+=("$SRC_LICENSE_STAGED=/opt/local/bin/Resources/LICENSE")
    fi
    
    # Validate file mappings before running fpm
    if ! validate_fpm_files "${RPM_FPM_FILES[@]}"
    then
      echo "Aborting .rpm build due to validation errors." >&2
      exit 1
    fi
    
    fpm \
      "${RPM_ARGS[@]}" \
      --package "$RPM_OUTFILE" \
      "${RPM_FPM_FILES[@]}"
    
    echo ""
    echo "Done. Packages created:"
    echo "  $DEB_OUTFILE"
    echo "  $RPM_OUTFILE"
  fi
}

# Check for fpm before starting
if ! command -v fpm >/dev/null 2>&1
then
  echo "Error: fpm not found. Install with: gem install fpm" >&2
  exit 1
fi

# Pre-flight validation for fpm file mappings (Bash 3.x compatible)
# Detects: circular symlinks, duplicate destinations, symlinks overwriting files
validate_fpm_files() {
  local src dest target mapping
  local seen_dests=""
  local seen_srcs=""
  
  for mapping in "$@"
  do
    src="${mapping%%=*}"
    dest="${mapping#*=}"
    
    # Check for circular/broken symlinks
    if [[ -L "$src" ]]
    then
      target=$(readlink "$src" 2>/dev/null || echo "")
      if [[ -z "$target" ]] || [[ "$target" == "$(basename "$src")" ]]
      then
        echo "Error: Circular or broken symlink detected: $src" >&2
        echo "  This will cause fpm to fail with 'destination is a directory' error." >&2
        return 1
      fi
      # Check if symlink target exists
      if [[ ! -e "$src" ]]
      then
        echo "Error: Symlink target does not exist: $src -> $target" >&2
        return 1
      fi
    fi
    
    # Check for duplicate destinations (Bash 3.x compatible using string matching)
    # Use newline-delimited string for tracking seen destinations
    if echo "$seen_dests" | grep -qxF "$dest"
    then
      # Find the original source for this destination
      local orig_src=""
      local check_mapping check_dest
      for check_mapping in "$@"
      do
        check_dest="${check_mapping#*=}"
        if [[ "$check_dest" == "$dest" ]]
        then
          orig_src="${check_mapping%%=*}"
          break
        fi
      done
      echo "Error: Duplicate destination detected: $dest" >&2
      echo "  First source: $orig_src" >&2
      echo "  Second source: $src" >&2
      echo "  This will cause fpm to fail with 'destination is a directory' error." >&2
      return 1
    fi
    seen_dests="$seen_dests
$dest"
    seen_srcs="$seen_srcs
$src"
    
    # Check source exists (for non-symlinks)
    if [[ ! -e "$src" ]] && [[ ! -L "$src" ]]
    then
      echo "Error: Source file does not exist: $src" >&2
      return 1
    fi
  done
  
  return 0
}

# Build packages based on TYPE_ARG
case "$TYPE_ARG" in
  linux)
    build_package "linux"
    ;;
  mac)
    build_package "mac"
    ;;
  all)
    echo "Building all packages..."
    echo ""
    build_package "linux"
    echo ""
    build_package "mac"
    ;;
esac

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
