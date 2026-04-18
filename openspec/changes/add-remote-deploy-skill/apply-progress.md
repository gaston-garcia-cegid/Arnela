# Apply progress: add-remote-deploy-skill

**Mode:** Standard (no código con test runner para esta skill en sí).  
**Batch:** 1

## Completed

- Phase 1: `.cursor/skills/remote-deploy/SKILL.md` creado.
- Phase 2: `.atl/skill-registry.md` actualizado (fila + compact rules).
- Phase 4: `CONTRIBUTING.md` — sección deploy remoto asistido.
- Spec cross-check: S-1 (Docker+clave/ssh), S-2 (contraseña advertida), S-3 (sin Docker) cubiertos en `SKILL.md`.

## Pending (manual / otro actor)

- Ninguno para este change (Engram: `mem_save` + `mem_search` vía MCP `user-engram`; commit aplicado en repo).

## Verification (sdd-verify)

- R-1–R-6 y S-1–S-4: cubiertos por `SKILL.md`, registry, `CONTRIBUTING.md` (ver `apply-progress` y `spec.md`).

## Deviations from design

None — paths usan `.cursor/skills/remote-deploy/SKILL.md` relativo al workspace en registry; design allowed Opción A en repo.
