# Proposal: Skill Cursor para despliegue remoto vía SSH

## Intent

Dar a los agentes (y al usuario) un flujo **repetible y documentado** para desplegar Arnela (u otro monorepo similar) en un servidor Linux usando **SSH**, sin depender de una skill SDD genérica de “configurar servidor”. Hoy `docs/DEPLOYMENT.md` describe el objetivo final, pero no hay una **Agent Skill** que guíe paso a paso la recolección de datos, los comandos remotos y la verificación.

## Scope

### In Scope

- Nueva **Cursor Agent Skill** (`SKILL.md`) con triggers claros (p. ej. “deploy remoto”, “ssh deploy”, “desplegar servidor”).
- Flujo que **solicite** (en chat): host SSH, usuario, método de auth (ruta a clave privada vs contraseña con advertencia de riesgo), ruta del repo en el servidor, nombre/identificador del proyecto (carpeta o `COMPOSE_PROJECT_NAME`), y confirmación de rama/tag.
- Instrucciones para que el agente use la **terminal integrada** ejecutando `ssh` (sesión interactiva o `ssh … comando` no interactivo según política de seguridad documentada en la skill).
- Secuencia remota mínima: `git fetch` + `git pull` (o checkout explícito), **build y arranque** alineados a Arnela: `docker compose … up -d --build` (o los comandos que el usuario confirme si no usa Docker).
- Checklist de **verificación** post-deploy (`curl` health, `docker compose ps`, logs).

### Out of Scope

- Provisionar VPS desde cero (DNS, firewall, TLS Let’s Encrypt completo) — seguir `docs/DEPLOYMENT.md` o un cambio aparte.
- CI/CD en GitHub Actions sustituto del deploy manual.
- Almacenar secretos en el repo o en la skill (solo flujo y placeholders).
- Cambios al código de aplicación Arnela salvo que el deploy descubra un bug real (entonces otro change).

## Capabilities

### New Capabilities

- `remote-server-deploy`: Contrato de comportamiento de la skill: qué preguntar al usuario, qué comandos SSH ejecutar (plantillas), cómo manejar credenciales sin persistirlas en git, y criterios de éxito/rollback remotos.

### Modified Capabilities

- None (no hay `openspec/specs/*/spec.md` todavía fuera del bootstrap).

## Approach

1. Añadir un archivo **`SKILL.md`** siguiendo el formato de `skill-creator` / skills existentes (frontmatter `name`, `description` con **Trigger:**).
2. Contenido en secciones: prerequisitos (OpenSSH en cliente, acceso al servidor, Docker en remoto si aplica), **orden estricto** de preguntas, plantillas `ssh user@host 'bash -lc '\''…'\'''` para `git` + `docker compose`, y bloque **Seguridad** (no pegar claves privadas en el chat si se puede evitar; preferir `ssh-agent` / IdentityFile).
3. Referenciar **`docs/DEPLOYMENT.md`** y **`docker-compose.prod.yml`** como fuente de verdad de Arnela.
4. Opcional: enlace en `CONTRIBUTING.md` una línea “Deploy remoto asistido → ver skill X” (solo si cabe en el mismo change; si no, tarea en `sdd-tasks`).

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `~/.cursor/skills/…` o `skills/` en repo | New | Ubicación acordada del nuevo `SKILL.md` (preferencia del usuario: global vs proyecto). |
| `CONTRIBUTING.md` / `docs/DEPLOYMENT.md` | Optional | Cross-link a la skill si se aprueba. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Exponer contraseñas o claves en historial de chat | Med | Skill exige preferencia por clave en agente + `ssh -i`; advertir de pegar secretos. |
| `ssh` interactivo bloqueado en sandbox del agente | Med | Documentar modo no interactivo y “el usuario ejecuta este bloque en su terminal” como fallback. |
| Entornos sin Docker en remoto | Bajo | Paso de confirmación antes de asumir `docker compose`. |

## Rollback Plan

- La skill es un **archivo añadido**: revertir commit o borrar `SKILL.md`.
- Si se tocó documentación, revertir esas líneas.
- Ningún cambio en runtime del servidor por solo añadir la skill.

## Dependencies

- OpenSSH cliente en la máquina donde corre Cursor.
- Usuario con acceso SSH al servidor y permisos para `git` + Docker (o stack acordado).

## Success Criteria

- [ ] Existe `SKILL.md` instalable con triggers y pasos ordenados (preguntas → SSH → pull → build/up).
- [ ] La skill menciona explícitamente Arnela (`DEPLOYMENT.md`, compose prod) sin hardcodear secretos.
- [ ] Queda documentado el fallback si el agente no puede abrir TTY interactiva.
