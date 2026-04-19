# Delta Spec: remote-server-deploy (Cursor skill)

**Change:** `add-remote-deploy-skill`  
**Type:** New capability — behavior of the Agent Skill only (no runtime change to Arnela services until the skill is used).

---

## Requirements

### R-1 — Skill artifact and discovery

- The deliverable **MUST** be a single `SKILL.md` following the Agent Skills format (YAML frontmatter with `name` and `description` containing a **Trigger:** line).
- Triggers **SHOULD** include Spanish and English phrases equivalent to: remote deploy via SSH, `git pull` on server, Docker compose up on server.

### R-2 — Pre-flight information (agent MUST collect before SSH)

The skill **MUST** require the agent to obtain, in order and explicitly:

1. **SSH host** (hostname or IP).
2. **SSH user** (login on the server).
3. **Authentication method**: (a) path to private key on the **client** machine, or (b) password — if password, the skill **MUST** state that pasting passwords in chat is a security risk and **SHOULD** recommend key-based auth or user typing in a local terminal.
4. **Remote repository path** (absolute path to the clone on the server).
5. **Git ref** to deploy (branch name or tag); default **MAY** be `main` if the user accepts.
6. **Project / compose context**: for Arnela, confirm use of `docker compose -f docker-compose.prod.yml --env-file .env.prod` per **`docs/DEPLOYMENT.md`** (production server uses prod compose only; do not merge with dev `docker-compose.yml`); if the user does not use Docker, the skill **MUST** stop the Docker path and ask for the user’s build/run commands instead.

### R-3 — SSH execution model

- The skill **MUST** describe two modes: (A) agent runs `ssh` from the integrated terminal; (B) user copies a documented script block into their own terminal if the agent environment cannot run interactive SSH.
- Remote commands **SHOULD** use a single non-interactive form when possible, e.g. `ssh -i <key> <user>@<host> 'bash -lc '\''…'\'''` with a clear placeholder for the inner script.
- The skill **MUST NOT** instruct writing private keys, `.env.prod` contents, or passwords into the repository or into committed files.

### R-4 — Remote deployment sequence (Arnela + Docker)

When the user confirms Docker-based Arnela deploy, the inner remote script **MUST** include, in order:

1. `cd <remote-repo-path>`
2. `git fetch` and `git pull` (or `git checkout` + pull as agreed) for the chosen ref.
3. `docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build` (same flags as **`docs/DEPLOYMENT.md`** § Deploy).

The skill **SHOULD** remind that `.env.prod` **MUST** exist on the server and is out of scope for the skill to create.

### R-5 — Post-deploy verification

The skill **MUST** list verification steps the agent or user runs **after** compose succeeds, aligned with project docs, minimally:

- `docker compose … ps` (healthy containers).
- `curl` to backend health/readiness (URL as appropriate: localhost on server or public URL behind Nginx).

### R-6 — Documentation cross-links

- The skill **MUST** reference `docs/DEPLOYMENT.md` and `docker-compose.prod.yml` as the source of truth for Arnela production compose.
- A link from **`CONTRIBUTING.md`** to the skill **MAY** be added in the same change if maintainers agree (optional scenario S-4).

---

## Scenarios

### S-1 — Happy path: key auth, Docker Arnela, agent runs SSH

**Given** the user supplied host, user, key path, remote repo path, branch `main`, and confirmed Docker Arnela.  
**When** the agent follows the skill and runs the non-interactive `ssh` one-liner (or equivalent).  
**Then** the remote host performs `git pull` and `docker compose … up -d --build` without the skill storing secrets in the repo, and the agent proposes the verification commands from R-5.

### S-2 — Password auth (discouraged)

**Given** the user insists on password authentication.  
**When** the skill is applied.  
**Then** the skill **MUST** warn about chat/logging exposure and **SHOULD** instruct the user to type the password only in a local interactive `ssh` session, not pasted into the agent chat.

### S-3 — No Docker on server

**Given** the user states the server does not use Docker for this app.  
**When** the agent applies the skill.  
**Then** the agent **MUST NOT** run the Arnela compose command; the skill **MUST** require the user to supply explicit build and start commands before any remote execution template is produced.

### S-4 — CONTRIBUTING cross-link (optional)

**Given** maintainers want discoverability.  
**When** the change is merged.  
**Then** `CONTRIBUTING.md` contains a one-line pointer to the skill path and trigger phrase.

---

## Non-goals (this spec)

- Defining VPS provisioning, firewall rules, or TLS certificate issuance (remain in `docs/DEPLOYMENT.md` or future changes).
- Automating GitHub Actions or other CI deploy targets.
