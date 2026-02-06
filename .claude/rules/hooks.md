# Hooks System

## Hook Types

- **PreToolUse**: Before tool execution (validation, parameter modification)
- **PostToolUse**: After tool execution (auto-format, checks)
- **Stop**: When session ends (final verification)

## Current Hooks (in ~/.claude/settings.json)

**Note**: Some hooks are language-specific. For Go projects, consider adding Go-specific hooks.

### PreToolUse
- **tmux reminder**: Suggests tmux for long-running commands (npm, pnpm, yarn, cargo, etc.) - **Generic**
- **git push review**: Opens Zed for review before push - **Generic**
- **doc blocker**: Blocks creation of unnecessary .md/.txt files - **Generic**

### PostToolUse
- **PR creation**: Logs PR URL and GitHub Actions status - **Generic**
- **Prettier**: Auto-formats JS/TS files after edit - **JavaScript/TypeScript only** (not applicable for Go)
- **TypeScript check**: Runs tsc after editing .ts/.tsx files - **TypeScript only** (not applicable for Go)
- **console.log warning**: Warns about console.log in edited files - **JavaScript/TypeScript only** (for Go, consider checking for `fmt.Println`)

### Stop
- **console.log audit**: Checks all modified files for console.log before session ends - **JavaScript/TypeScript only**

## Go-Specific Hooks (Recommended)

Consider adding these hooks for Go projects:

### PostToolUse Hooks for Go
- **Go fmt auto-format**: Run `gofmt -w` after editing `.go` files
- **Go vet check**: Run `go vet` after editing `.go` files to catch issues
- **Go test run**: Run `go test ./...` after significant changes
- **Go lint check**: Run `golangci-lint` if configured

### Example hook configuration for Go:
```json
{
  "matcher": "tool == \"Edit\" && tool_input.file_path matches \"\\\\.go$\"",
  "hooks": [
    {
      "type": "command",
      "command": "gofmt -w \"${file_path}\""
    },
    {
      "type": "command",
      "command": "go vet \"${file_path}\" 2>&1 | head -20"
    }
  ]
}
```

## Auto-Accept Permissions

Use with caution:
- Enable for trusted, well-defined plans
- Disable for exploratory work
- Never use dangerously-skip-permissions flag
- Configure `allowedTools` in `~/.claude.json` instead

## TodoWrite Best Practices

Use TodoWrite tool to:
- Track progress on multi-step tasks
- Verify understanding of instructions
- Enable real-time steering
- Show granular implementation steps

Todo list reveals:
- Out of order steps
- Missing items
- Extra unnecessary items
- Wrong granularity
- Misinterpreted requirements
