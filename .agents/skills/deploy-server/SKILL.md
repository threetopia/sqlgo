---
name: deploy-server
description: Use when deploying, releasing, rolling back, or restarting Rakkume apps on the systemd production server (build binary, rsync to server, restart unit, run migration, verify health).
---

# Deploy Server

## Overview

Production runs as a **native Go binary under systemd**, not Docker Compose. Deploy = build → transfer → restart → verify.

Full runbook: **[docs/ops/deployment.md](../../../docs/ops/deployment.md)**. Run **[deploy/deploy.sh](../../../deploy/deploy.sh)**; target values live in [deploy/targets.env](../../../deploy/targets.env). This skill is only a trigger and guardrail — do not duplicate the runbook steps here.

## When to Use

- A request to deploy/release to the server, roll back, or restart a production service.
- Changes that touch `migration/`, `web/template/<theme>/{assets,html}`, or `.env`.

Not for: running the local Docker stack (`docker compose up`) or `go run` for development.

## Configuration Layers

Keep these apart — they need different visibility:

| Layer | Content | Location |
|---|---|---|
| SSH connection | Host alias, User, Port, IdentityFile | `~/.ssh/config` |
| Server layout | `SSH_HOST`, `APP_BASE_DIR`, `BINARY_DIR`, `ENV_FILE`, `SERVICE_NAME_PREFIX`, `WEB_ASSET_ROOT` | `deploy/targets.env` |
| Runtime + secrets | DSN, ports, `*_MIGRATION_PATH`, `WEB_*_TEMPLATE_PATH` | `.env` on the server |

Never put secrets, IPs, or SSH users in `deploy/targets.env`; never put server paths in the app `.env`.

## Guardrails

- Run `--dry-run` first and check the output before a real deploy.
- Never fill `CHANGE_ME` in `deploy/targets.env` on your own — those values must come from the operator.
- Back up `.env` and the old binary (`.prev`) before overwriting anything; the script does both.
- Restart **only** the app unit being deployed, one at a time — do not restart all units.
- Do not deploy the `stress` entrypoint to production; the script rejects it.
- `migration/` must be available on the server: the API runs migrations automatically on startup (`k.DAO().DAOMigration().Run()`), and a failed migration means a failed start.
- Templates, static assets, and migration SQL are **not** embedded — they are read from an absolute env path, so they must be uploaded separately. Sync the whole `web/template/` tree: each theme carries `assets/` next to `html/`.
- All build artifacts go into `./.tmp/` and are cleaned up afterward.

## Quick Reference

| Need | Command |
|---|---|
| Show targets + app list | `./deploy/deploy.sh --list-targets` |
| Preview plan | `./deploy/deploy.sh --dry-run <app>` |
| Deploy an app | `./deploy/deploy.sh <app>` |
| Build only (no server) | `./deploy/deploy.sh --build-only <app>` |
| Binary only | `./deploy/deploy.sh --binary-only <app>` |
| Web assets only | `./deploy/deploy.sh --assets-only` |
| App + web assets | `./deploy/deploy.sh --with-assets <app>` |
| Build on the server | `./deploy/deploy.sh --local <app>` |
| Rollback | `<app>.prev` → `<app>` then restart, see [§11 runbook](../../../docs/ops/deployment.md#11-rollback) |

## Common Mistakes

- Putting an IP, SSH user, or secret into `deploy/targets.env` — that is not its place.
- Restarting all units at once, widening downtime for no reason.
- Forgetting to transfer `migration/` even though the migration path is read from an absolute env value.
- Forgetting to upload the HTML/static assets because templates are not embedded into the binary. Each theme keeps its static assets in `template/<theme>/assets`, next to `template/<theme>/html` — sync the whole `web/template/` tree, not a separate `web/static`.
- Deploying an app without assets while expecting CSS/JS changes to appear: an app-only deploy never touches assets; use `--with-assets`.
- Assuming the synced theme is the one being served: check `.env.prod` `WEB_*_TEMPLATE_PATH`/`*_STATIC_PATH` first (runbook §7). They currently point at this repo's `default` theme, but `web_dummy` still points at a non-existent `web/static/core`.
- Assuming `systemctl restart` succeeded without checking `is-active` and `journalctl` (the script does this for you).
- Thinking a binary rollback automatically reverts the database migration.
