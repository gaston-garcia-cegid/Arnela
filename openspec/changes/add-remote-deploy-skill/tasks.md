# Tasks: add-remote-deploy-skill

## Phase 1: Skill file

- [x] 1.1 Crear directorio `.cursor/skills/remote-deploy/`.
- [x] 1.2 Añadir `SKILL.md` con frontmatter (`name`, `description` + **Trigger:** en ES/EN: deploy remoto, ssh deploy, desplegar servidor, remote deploy).
- [x] 1.3 Documentar prerequisitos: OpenSSH cliente, servidor con Docker+Compose si Arnela, `docs/DEPLOYMENT.md` como referencia obligatoria.
- [x] 1.4 Sección “Datos a recoger” en orden: host, usuario, auth (clave `ssh -i` preferida; contraseña con advertencia), ruta remota del repo, rama/tag, confirmación Docker vs comandos custom.
- [x] 1.5 Plantilla remota Arnela: `cd` + `git fetch`/`git pull` (o checkout) + `docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d --build` (copiar de `docs/DEPLOYMENT.md` § Deploy).
- [x] 1.6 Bloque “Sin Docker”: pedir comandos explícitos de build/run antes de plantillar SSH.
- [x] 1.7 Modo agente vs usuario: `ssh … 'bash -lc …'` no interactivo; fallback “copiar y ejecutar en tu terminal” si no hay TTY.
- [x] 1.8 Post-deploy: `docker compose … ps`, `curl` health/readiness (URLs según doc).
- [x] 1.9 Seguridad: prohibir pegar `.env`/claves privadas en el repo; no persistir secretos en commits.

## Phase 2: Registry en repo

- [x] 2.1 Editar `.atl/skill-registry.md`: nueva fila en tabla User Skills — trigger + path **`.cursor/skills/remote-deploy/SKILL.md`** (relativo al workspace).
- [x] 2.2 Añadir subsección **### remote-deploy** bajo Compact Rules (≤10 líneas): SSH, compose Arnela, no secretos en chat, fallback terminal.

## Phase 3: Engram (post-implementación)

- [x] 3.1 Con MCP Engram activo: `mem_save` con `topic_key` estable (p. ej. `skill-registry`) y `content` = cuerpo actualizado de `.atl/skill-registry.md`, o observación dedicada `skills/arnela/remote-deploy` con triggers + path + reglas compactas.
- [x] 3.2 Verificar con `mem_search` que la observación es recuperable para el proyecto `Arnela`.

## Phase 4: Doc opcional y cierre

- [x] 4.1 (Opcional) Una línea en `CONTRIBUTING.md` apuntando a la skill y a `docs/DEPLOYMENT.md` para deploy remoto asistido.
- [x] 4.2 Revisión manual: leer `specs/remote-server-deploy/spec.md` y marcar escenarios S-1–S-3 cubiertos por el texto de `SKILL.md`.

## Phase 5: Verificación

- [x] 5.1 No hay secretos ni rutas absolutas de usuario en `SKILL.md` / registry (salvo ejemplos placeholder).
- [x] 5.2 `git status` limpio tras commit; PR listo para revisión.
