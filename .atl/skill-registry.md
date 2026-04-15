# Skill Registry — Arnela

Generated: 2026-04-15
Project: arnela
Stack: Go 1.24 (Gin) + Next.js 16 (React 19, TypeScript 5.9)

## Available Skills

| # | Skill | Source | Trigger |
|---|-------|--------|---------|
| 1 | go-testing | user | Writing Go tests, using teatest, adding test coverage |
| 2 | issue-creation | user | Creating GitHub issues, reporting bugs, requesting features |
| 3 | branch-pr | user | Creating PRs, opening PRs, preparing changes for review |
| 4 | skill-creator | user | Creating new skills, adding agent instructions, documenting AI patterns |
| 5 | judgment-day | user | "judgment day", adversarial review, "dual review", "juzgar" |

## Project Conventions

| File | Path | Description |
|------|------|-------------|
| copilot-instructions.md | .github/copilot-instructions.md | Coding conventions, architecture patterns, naming rules |
| Agent.md | Agent.md | Full technical definition and architecture details |

## Compact Rules

### go-testing
```
WHEN writing Go tests:
- Use table-driven tests with testify
- Follow AAA pattern (Arrange, Act, Assert)
- Use testify/assert and testify/require
- Mock interfaces with testify/mock
- Use miniredis for Redis tests
- Name: Test{Function}_{Scenario}
```

### issue-creation
```
WHEN creating GitHub issues:
- Use gh issue create with proper labels
- Include reproduction steps for bugs
- Follow issue template if available
- Link related issues/PRs
```

### branch-pr
```
WHEN creating pull requests:
- Follow issue-first enforcement (link to issue)
- Use gh pr create with structured body
- Include summary, test plan, and linked issue
- Branch naming: {type}/{issue-number}-{description}
```

### skill-creator
```
WHEN creating new AI skills:
- Use SKILL.md format with YAML frontmatter
- Include triggers, when-to-use, decision trees
- Add compact rules section for registry
- Place in ~/.cursor/skills/{name}/SKILL.md
```

### judgment-day
```
WHEN adversarial review requested:
- Launch two independent blind judge sub-agents in parallel
- Each reviews the same target independently
- Synthesize findings, apply fixes, re-judge
- Max 2 iterations, then escalate
```

### Project: copilot-instructions.md
```
WHEN working on Arnela:
- Backend: Go Clean Architecture (handler/service/domain/repository/integration)
- Frontend: Next.js App Router, Zustand for state, Shadcn UI
- API JSON keys: camelCase with json:"camelCase" tags
- PascalCase for exports, camelCase for private
- Use Docker for local dev (Postgres, Redis)
- TDD focus on backend
```
