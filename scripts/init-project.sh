#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/init-project.sh --app-name <new-app-name> --module-path <new-module-path> [--dry-run]

Example:
  ./scripts/init-project.sh \
    --app-name xxxxx-admin \
    --module-path github.com/your-org/xxxxx-admin

This script is intended to be run once when bootstrapping a new project from the template.
It updates:
  - go.mod module path
  - Go import paths
  - project/binary names in docs and build files
  - entrypoint filename (<app-name>.go)
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

replace_in_file() {
  local file="$1"
  local old="$2"
  local new="$3"

  if [[ "$DRY_RUN" == "true" ]]; then
    if grep -Fq "$old" "$file"; then
      echo "[dry-run] update $file"
    fi
    return 0
  fi

  OLD_VALUE="$old" NEW_VALUE="$new" perl -0pi -e 's/\Q$ENV{OLD_VALUE}\E/$ENV{NEW_VALUE}/g' "$file"
}

is_text_file() {
  local file="$1"
  [[ -f "$file" ]] || return 1
  LC_ALL=C grep -Iq . "$file"
}

collect_files() {
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git ls-files -z --cached --others --exclude-standard
  else
    find . \
      -path './.git' -prune -o \
      -type f \
      -print0
  fi
}

rename_file() {
  local from="$1"
  local to="$2"

  if [[ ! -e "$from" || "$from" == "$to" ]]; then
    return 0
  fi

  if [[ "$DRY_RUN" == "true" ]]; then
    echo "[dry-run] rename $from -> $to"
    return 0
  fi

  mv "$from" "$to"
}

NEW_APP_NAME=""
NEW_MODULE_PATH=""
DRY_RUN="false"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --app-name)
      NEW_APP_NAME="${2:-}"
      shift 2
      ;;
    --module-path)
      NEW_MODULE_PATH="${2:-}"
      shift 2
      ;;
    --dry-run)
      DRY_RUN="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$NEW_APP_NAME" || -z "$NEW_MODULE_PATH" ]]; then
  usage
  exit 1
fi

require_cmd git
require_cmd perl
require_cmd sed

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
cd "$ROOT_DIR"

OLD_MODULE_PATH="$(sed -n '1s/^module //p' go.mod)"
if [[ -z "$OLD_MODULE_PATH" ]]; then
  echo "failed to detect current module path from go.mod" >&2
  exit 1
fi

OLD_APP_NAME="$(sed -n 's/^PROJECT_NAME := //p' Makefile | head -n 1 | tr -d '[:space:]')"
if [[ -z "$OLD_APP_NAME" ]]; then
  OLD_APP_NAME="$(basename "$OLD_MODULE_PATH")"
fi

if [[ "$OLD_APP_NAME" == "$NEW_APP_NAME" && "$OLD_MODULE_PATH" == "$NEW_MODULE_PATH" ]]; then
  echo "project already matches the requested app name and module path"
  exit 0
fi

echo "Bootstrapping project"
echo "  old app name:    $OLD_APP_NAME"
echo "  new app name:    $NEW_APP_NAME"
echo "  old module path: $OLD_MODULE_PATH"
echo "  new module path: $NEW_MODULE_PATH"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "  mode:            dry-run"
fi

while IFS= read -r -d '' file; do
  if [[ "$file" == ./.git/* ]]; then
    continue
  fi
  if ! is_text_file "$file"; then
    continue
  fi
  replace_in_file "$file" "$OLD_MODULE_PATH" "$NEW_MODULE_PATH"
  replace_in_file "$file" "$OLD_APP_NAME" "$NEW_APP_NAME"
done < <(collect_files)

rename_file "${OLD_APP_NAME}.go" "${NEW_APP_NAME}.go"

if [[ "$DRY_RUN" == "false" ]]; then
  if command -v go >/dev/null 2>&1; then
    gofmt -w "${NEW_APP_NAME}.go" ./internal ./pkg ./config ./docs >/dev/null 2>&1 || true
  fi
fi

echo "Initialization complete."
echo "Next steps:"
echo "  1. Review the diff"
echo "  2. Run: go test ./..."
echo "  3. Run: make run"
