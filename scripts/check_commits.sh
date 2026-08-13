#!/bin/sh
set -eu

# Configuration
MAX_COMMITS_BEHIND=100
JIRA_PATTERN='(DROP|TECH)-[0-9]{3,5}'
MIN_BODY_LINES=3

# Debug mode
DEBUG="${DEBUG:-false}"
debug() {
    if [ "$DEBUG" = "true" ]; then
        printf 'DEBUG: %s\n' "$*" >&2
    fi
}

# Colors
red() { printf '\033[0;31m%s\033[0m\n' "$*" >&2; }
purple() { printf '\033[0;35m%s\033[0m\n' "$*"; }
yellow() { printf '\033[0;33m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }

# Validate git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    red "Error: Not a git repository"
    exit 1
fi

# Get diff range
diff="${1:-origin/master...HEAD}"
debug "Using diff range: $diff"

purple "Executing commit validation for: $diff"
printf '\n'

# Ensure we have the target branch
target_branch=$(echo "${diff%%..*}" | sed 's|^origin/||')
if [ -n "$target_branch" ]; then
    debug "Fetching target branch: $target_branch"
    if git fetch origin "+refs/heads/$target_branch:refs/remotes/origin/$target_branch" >/dev/null 2>&1; then
        debug "Successfully fetched $target_branch"
    else
        yellow "Warning: Could not fetch $target_branch"
    fi
fi

# Verify branches exist
debug "Checking if branches exist..."
if git rev-parse --verify "origin/$target_branch" >/dev/null 2>&1; then
    debug "✓ origin/$target_branch exists"
else
    red "Error: origin/$target_branch does not exist"
    exit 1
fi

# Check if branch is too far behind
debug "Calculating commits behind..."
if num_commits_behind=$(git rev-list --left-only --count "$diff" 2>&1); then
    debug "Commits behind: $num_commits_behind"
else
    red "Error calculating commits behind"
    exit 1
fi

if [ "$num_commits_behind" -gt "$MAX_COMMITS_BEHIND" ]; then
    red "Branch is ${num_commits_behind} commits behind. Please sync (max: $MAX_COMMITS_BEHIND)"
    exit 1
fi

# Get commits to validate
debug "Fetching commits to validate..."
if commits=$(git log --no-merges --format="===%H%n%s%n%b;;;" "$diff" 2>&1); then
    debug "Successfully fetched commits"

    # Count commits
    commit_lines=$(echo "$commits" | grep -c "^===" || echo "0")
    debug "Found $commit_lines commits to validate"
else
    red "Error fetching commits"
    exit 1
fi

if [ -z "$commits" ]; then
    purple "No commits to validate"
    exit 0
fi

# Validation state - usando archivos temporales para evitar problema de subshell
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

echo "0" > "$tmp_dir/exit_code"
echo "0" > "$tmp_dir/commit_count"
echo "0" > "$tmp_dir/valid_count"
touch "$tmp_dir/failed_commits"
touch "$tmp_dir/failed_reasons"

# Validate a single commit
validate_commit() {
    commit_hash="$1"
    title="$2"
    body="$3"
    commit_num="$4"

    short_hash=$(echo "$commit_hash" | cut -c1-7)
    has_errors=0

    debug "Validating commit #$commit_num ($short_hash): $title"

    # Skip merge commits (Bitbucket PR merges, GitHub merges, manual merges)
    if echo "$title" | grep -qE "^(Merged in |Merge branch |Merge pull request )"; then
        debug "Skipping merge commit: $short_hash"
        return 0
    fi

    # Check JIRA in title (not allowed, when a ticket is present)
    if echo "$title" | grep -qE "$JIRA_PATTERN"; then
        printf '%s|%s\n' "$short_hash" "JIRA ticket should NOT be in title" >> "$tmp_dir/failed_reasons"
        has_errors=1
    fi

    # Check for Why: line
    if ! echo "$body" | grep -q "^Why:"; then
        printf '%s|%s\n' "$short_hash" "Missing 'Why:' line" >> "$tmp_dir/failed_reasons"
        has_errors=1
    fi

    # Check for What: line
    if ! echo "$body" | grep -q "^What:"; then
        printf '%s|%s\n' "$short_hash" "Missing 'What:' line" >> "$tmp_dir/failed_reasons"
        has_errors=1
    fi

    # Check minimum body length
    body_lines=$(echo "$body" | grep -cv "^$" || echo "0")
    if [ "$body_lines" -lt "$MIN_BODY_LINES" ]; then
        printf '%s|%s\n' "$short_hash" "Body too short ($body_lines lines, min: $MIN_BODY_LINES)" >> "$tmp_dir/failed_reasons"
        has_errors=1
    fi

    # If there are errors, store commit info
    if [ "$has_errors" -eq 1 ]; then
        printf '%s|%s\n' "$short_hash" "$title" >> "$tmp_dir/failed_commits"
        echo "1" > "$tmp_dir/exit_code"
        return 1
    fi

    debug "✓ Commit $short_hash passed validation"
    return 0
}

# Parse commits
debug "Starting commit parsing..."
current_hash=""
current_title=""
current_body=""

# FIX: Evitar subshell escribiendo el loop diferente
echo "$commits" > "$tmp_dir/commits.txt"

while IFS= read -r line; do
    case "$line" in
        ===*)
            # Process previous commit if exists
            if [ -n "$current_hash" ]; then
                count=$(cat "$tmp_dir/commit_count")
                count=$((count + 1))
                echo "$count" > "$tmp_dir/commit_count"

                debug "Processing commit $count: $current_hash"

                if validate_commit "$current_hash" "$current_title" "$current_body" "$count"; then
                    valid=$(cat "$tmp_dir/valid_count")
                    valid=$((valid + 1))
                    echo "$valid" > "$tmp_dir/valid_count"
                fi
            fi

            # Extract hash (remove ===)
            current_hash="${line#===}"
            current_title=""
            current_body=""
            debug "Starting new commit: $current_hash"
            ;;
        ";;;")
            debug "End of commit marker found"
            ;;
        *)
            # Capture title and body
            if [ -z "$current_title" ]; then
                current_title="$line"
                debug "Title: $current_title"
            else
                if [ -z "$current_body" ]; then
                    current_body="$line"
                else
                    current_body="$current_body
$line"
                fi
            fi
            ;;
    esac
done < "$tmp_dir/commits.txt"

# Validate last commit
if [ -n "$current_hash" ]; then
    count=$(cat "$tmp_dir/commit_count")
    count=$((count + 1))
    echo "$count" > "$tmp_dir/commit_count"

    debug "Processing final commit $count: $current_hash"

    if validate_commit "$current_hash" "$current_title" "$current_body" "$count"; then
        valid=$(cat "$tmp_dir/valid_count")
        valid=$((valid + 1))
        echo "$valid" > "$tmp_dir/valid_count"
    fi
fi

# Read final counts
exit_code=$(cat "$tmp_dir/exit_code")
commit_count=$(cat "$tmp_dir/commit_count")
valid_count=$(cat "$tmp_dir/valid_count")

debug "Parsing complete. Total commits: $commit_count, Valid: $valid_count"

# Summary
printf '\n'
purple "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$exit_code" -eq 0 ]; then
    green "✓ All $commit_count commits validated successfully"
else
    failed_count=$(grep -c "." "$tmp_dir/failed_commits" || echo "0")
    red "✗ Validation failed: $failed_count of $commit_count commits have errors"
    printf '\n'

    if [ -s "$tmp_dir/failed_commits" ]; then
        yellow "Failed commits:"
        while IFS='|' read -r hash title; do
            printf '  • %s - %s\n' "$hash" "$title"
        done < "$tmp_dir/failed_commits"
    fi

    printf '\n'
    if [ -s "$tmp_dir/failed_reasons" ]; then
        red "Errors found:"
        while IFS='|' read -r hash error; do
            printf '  [%s] %s\n' "$hash" "$error"
        done < "$tmp_dir/failed_reasons"
    fi

    printf '\n'
    yellow "To fix these commits:"
    echo "  git rebase -i $diff"
    echo "  # Then edit each failing commit with 'reword' or 'edit'"
fi

purple "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

exit "$exit_code"