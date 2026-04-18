# Tasks: SDD bootstrap + Engram (Arnela)

## Phase 1: OpenSpec skeleton

- [x] 1.1 Add `openspec/config.yaml` with stack context, `strict_tdd`, `testing.capabilities`, phase `rules`.
- [x] 1.2 Add `openspec/specs/` and `openspec/changes/archive/` placeholders for future specs and archived changes.

## Phase 2: Skill registry + Cursor MCP

- [x] 2.1 Add `.atl/skill-registry.md` (delegator-oriented skill index; omit `sdd-*` from compact table per convention).
- [x] 2.2 Add `.cursor/mcp.json` registering `engram` with `engram mcp` (stdio).
- [x] 2.3 Add `.cursor/rules/engram.mdc` (`alwaysApply: true`) for Memory Protocol hints.

## Phase 3: Verification (this session)

- [x] 3.1 Backend: `go vet ./...`, `go build ./...`, `go test ./... -count=1` from `backend/` (Windows: `-race` requires `CGO_ENABLED=1`; CI runs `-race` on Linux).
- [x] 3.2 Frontend: `pnpm exec tsc --noEmit`, `pnpm lint` (warnings only), `pnpm run test -- --run` from `frontend/`.

## Phase 4: Git

- [x] 4.1 Commit tracked paths; omit `Agent-Orchestrator.md` unless explicitly requested. Push `main`.
