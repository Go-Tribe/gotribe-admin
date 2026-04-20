#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEB_DIR="${ROOT_DIR}/web/admin"
RUN_FRONTEND="${RUN_FRONTEND:-1}"
VITE_API_BASE_URL="${VITE_API_BASE_URL:-/api}"

cd "${ROOT_DIR}"

section() {
  printf "\n>>> %s\n" "$1"
}

if [[ "${RUN_FRONTEND}" == "1" && -f "${WEB_DIR}/package.json" ]]; then
  section "Frontend dependency check"
  if ! command -v pnpm >/dev/null 2>&1; then
    printf "pnpm is required for frontend checks. Install pnpm or run RUN_FRONTEND=0 %s.\n" "$0"
    exit 1
  fi

  cd "${WEB_DIR}"
  if [[ ! -d node_modules ]]; then
    pnpm install --frozen-lockfile
  fi

  section "Frontend lint"
  VITE_API_BASE_URL="${VITE_API_BASE_URL}" pnpm run lint

  section "Frontend build"
  VITE_API_BASE_URL="${VITE_API_BASE_URL}" pnpm run build
  cd "${ROOT_DIR}"
elif [[ ! -f "${WEB_DIR}/dist/index.html" ]]; then
  printf "web/admin/dist is missing. Run with RUN_FRONTEND=1 or build the frontend before Go checks.\n"
  exit 1
fi

section "Go module verification"
go mod verify

section "Go format check"
unformatted="$(gofmt -s -l .)"
if [[ -n "${unformatted}" ]]; then
  printf "The following Go files are not gofmt -s formatted:\n%s\n" "${unformatted}"
  exit 1
fi

section "Go vet"
go vet ./...

section "Go tests with race detector and coverage"
go test -race -coverprofile="${ROOT_DIR}/coverage.out" ./...

section "All checks passed"
