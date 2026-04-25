#!/usr/bin/env bash
set -uo pipefail

# 無限ループ防止: 既に Stop hook がブロック済みなら再チェックせず素通り
stop_hook_active=$(jq -r '.stop_hook_active // false')
if [ "$stop_hook_active" = "true" ]; then
  terminal-notifier -title 'Claude Code' -message 'Claude has completed the task.' -sound 'Boop'
  exit 0
fi

if [ -f go.mod ] && git_status=$(git status --porcelain 2>/dev/null) && printf '%s\n' "$git_status" | grep -qE '\.go$'; then
  lint_output=$({ go vet ./... && go tool staticcheck ./...; } 2>&1)
  lint_status=$?
  test_output=$(go test -shuffle=on ./... 2>&1)
  test_status=$?
  if [ $lint_status -ne 0 ] || [ $test_status -ne 0 ]; then
    jq -n --arg l "$lint_output" --arg t "$test_output" \
      '{
        decision: "block",
        reason: (
          "静的解析またはテストが失敗しています。修正してください。\n\n"
          + "静的解析:\n" + $l + "\n\n"
          + "テスト:\n" + $t
        )
      }'
    exit 0
  fi
fi

terminal-notifier -title 'Claude Code' -message 'Claude has completed the task.' -sound 'Boop'
