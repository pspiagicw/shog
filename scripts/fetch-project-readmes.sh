#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
URLS_FILE="${ROOT_DIR}/content/projects.urls"
OUTPUT_DIR="${ROOT_DIR}/content/projects"

if ! command -v curl >/dev/null 2>&1; then
  echo "error: curl is required" >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"
find "${OUTPUT_DIR}" -type f -name '*.md' -delete

AUTH_ARGS=()
if [[ -n "${GITHUB_TOKEN:-}" ]]; then
  AUTH_ARGS=(-H "Authorization: Bearer ${GITHUB_TOKEN}")
fi

normalize_url() {
  local url="${1#https://}"
  url="${url#http://}"
  url="${url%/}"
  printf '%s\n' "${url}"
}

slug_from_url() {
  local normalized
  normalized="$(normalize_url "$1")"
  printf '%s\n' "${normalized##*/}"
}

while IFS= read -r raw_line || [[ -n "${raw_line}" ]]; do
  line="$(printf '%s' "${raw_line}" | sed 's/[[:space:]]*$//')"
  if [[ -z "${line}" || "${line}" == \#* ]]; then
    continue
  fi

  normalized="$(normalize_url "${line}")"
  if [[ "${normalized}" != github.com/*/* ]]; then
    echo "warning: skipping unsupported URL: ${line}" >&2
    continue
  fi

  repo_path="${normalized#github.com/}"
  slug="$(slug_from_url "${normalized}")"
  output_path="${OUTPUT_DIR}/${slug}.md"

  echo "fetching README for ${repo_path}"
  curl -fsSL \
    -H "Accept: application/vnd.github.raw+json" \
    "${AUTH_ARGS[@]}" \
    "https://api.github.com/repos/${repo_path}/readme" \
    -o "${output_path}"
done < "${URLS_FILE}"
