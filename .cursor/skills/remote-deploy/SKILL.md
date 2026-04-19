---
name: remote-deploy
description: >
  Guía para desplegar Arnela (o monorepo similar) en un servidor Linux vía SSH: recoger credenciales/contexto,
  git pull remoto, docker compose build/up según docs/DEPLOYMENT.md, verificación post-deploy.
  Trigger: deploy remoto, desplegar servidor, ssh deploy, remote deploy, pull y compose en servidor,
  actualizar producción ssh, docker compose remoto arnela.
license: MIT
metadata:
  author: arnela
  version: "1.0"
---

## Purpose

Orquestar un **despliegue remoto seguro** sin escribir secretos en el repositorio. La fuente de verdad de Arnela en producción es **`docs/DEPLOYMENT.md`** y **`docker-compose.prod.yml`**.

## Prerequisites

- Cliente **OpenSSH** (`ssh`) en la máquina donde corre el agente (o el usuario ejecutará comandos manualmente).
- Servidor con **Git** y, si aplica Arnela estándar, **Docker** + **Docker Compose** y clone del repo.
- En el servidor ya debe existir **`.env.prod`** (no lo generes ni pegues contenido en el chat).

## Data to collect (strict order)

Ask the user for each item before running SSH:

1. **SSH host** — hostname or IP.
2. **SSH user** — login on server.
3. **Auth** — **Preferred:** path to **private key** on the client for `ssh -i <path>`. If the user insists on **password**: warn that pasting passwords in chat may leak to logs; recommend typing password only in a **local interactive** `ssh` session they open themselves.
4. **Remote repo path** — absolute path to the git clone on the server (e.g. `/srv/arnela`).
5. **Git ref** — branch or tag (default `main` only if the user confirms).
6. **Stack** — Confirm **Docker Compose** deploy per Arnela docs. If **not** Docker, stop and ask for exact remote **build** and **run** commands before proposing any script.

## Arnela + Docker: remote one-liner template

Use **non-interactive** SSH when possible. Replace placeholders; do not commit this string with real secrets.

```bash
ssh -i "<CLIENT_KEY_PATH>" "<USER>@<HOST>" 'bash -lc '\''set -e
cd "<REMOTE_REPO_PATH>"
git fetch --all
git checkout "<REF>" && git pull --ff-only
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
'\'''
```

Aligns with **Deploy** section of `docs/DEPLOYMENT.md`. If the compose files or env file path differ on the server, use the user’s confirmed paths.

## No Docker on server

If the user says Docker is not used:

1. Do **not** run the compose template above.
2. Collect explicit remote commands (build + process manager / binary run).
3. Wrap them in the same `ssh … 'bash -lc …'` pattern once confirmed.

## Agent vs user terminal

- **Primary:** agent runs `ssh` from the integrated terminal with the one-liner (key-based auth).
- **Fallback:** if the environment blocks TTY/SSH or password is required without key, output a **copy-paste block** for the user to run in their own terminal; do not ask them to paste `.env.prod` or private keys into chat.

## Post-deploy checks

From docs / server context, suggest:

```bash
ssh -i "<CLIENT_KEY_PATH>" "<USER>@<HOST>" 'bash -lc '\''cd "<REMOTE_REPO_PATH>" && docker compose -f docker-compose.prod.yml ps'\'''
```

Health (adjust URL: on server `localhost`, or public URL behind Nginx):

```bash
curl -sS http://127.0.0.1:8080/health
curl -sS http://127.0.0.1:8080/readiness
```

## Security rules

- **Never** commit `.env.prod`, private keys, or passwords into git.
- **Never** instruct storing secrets inside `SKILL.md` or markdown in the repo.
- Prefer `ssh -i` + key file on client; use `IdentitiesOnly=yes` in `~/.ssh/config` when helpful (document as optional tip).

## Optional: project name / compose name

If the user needs a specific **Compose project name**, they can export `COMPOSE_PROJECT_NAME` inside the remote script before `docker compose` — only after they confirm the value.

## References

- `docs/DEPLOYMENT.md` — full production guide, TLS, backups, troubleshooting.
- `docker-compose.prod.yml` — production overrides.
