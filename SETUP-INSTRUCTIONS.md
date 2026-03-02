# DockPulse — Agent Setup Instructions

> **Purpose**: Complete project setup before any development begins.
> **Prerequisites**: `gh` CLI installed and authenticated as HerbHall. Run all commands in CMD (not PowerShell) to avoid bracket escaping issues.
> **Shell**: CMD (`cmd /c` wrapper if calling from PowerShell)
> **Working directory**: `D:\devspace\DockPulse`

---

## Pre-flight Checks

Before running any commands, verify:

1. `gh auth status` returns authenticated as HerbHall
2. `git --version` returns a version
3. `D:\devspace\DockPulse\.git` directory exists (git init was already done)
4. No remote repo exists yet — `git remote -v` should return nothing

If `.git` does NOT exist, run `git init` before proceeding.

---

## Step 1 — Create GitHub Repository and Push

```cmd
cd /d D:\devspace\DockPulse
gh repo create HerbHall/DockPulse --public --source=. --remote=origin --description "Docker Desktop extension — check if your container images have updates available"
git add -A
git commit -m "chore: initial project scaffold"
git push -u origin main
```

**Verify**: `gh repo view HerbHall/DockPulse` shows the repo. `git log --oneline -1` shows the commit.

---

## Step 2 — Create Issue Labels

These labels must exist before the issue batch file runs.

```cmd
cd /d D:\devspace\DockPulse
gh label create mvp --color 0E8A16 --description "Minimum viable product" --repo HerbHall/DockPulse
gh label create feat --color 1D76DB --description "New feature" --repo HerbHall/DockPulse
gh label create enhancement --color A2EEEF --description "Enhancement" --repo HerbHall/DockPulse
gh label create chore --color FBCA04 --description "Maintenance task" --repo HerbHall/DockPulse
gh label create docs --color 0075CA --description "Documentation" --repo HerbHall/DockPulse
```

**Verify**: `gh label list --repo HerbHall/DockPulse` shows all five labels.

---

## Step 3 — Create GitHub Issues

```cmd
cd /d D:\devspace\DockPulse
create-issues.bat
```

This creates 10 issues covering the full backlog (5 MVP, 3 enhancement, 1 chore, 1 docs).

**Verify**: `gh issue list --repo HerbHall/DockPulse` shows 10 open issues.

---

## Step 4 — Clean Up Batch Artifacts

The batch file creates temporary `.body` and `.issuebody` files. Delete them if present and make sure they stay in `.gitignore` (they already are).

```cmd
cd /d D:\devspace\DockPulse
if exist .body del .body
if exist .issuebody del .issuebody
```

---

## Step 5 — Final Commit

If any files were created or modified during setup:

```cmd
cd /d D:\devspace\DockPulse
git add -A
git status
git commit -m "chore: complete project setup" --allow-empty
git push
```

---

## Completion Checklist

- [ ] GitHub repo `HerbHall/DockPulse` exists and is public
- [ ] All scaffold files pushed to `main` branch
- [ ] 5 labels created (mvp, feat, enhancement, chore, docs)
- [ ] 10 issues created with correct labels
- [ ] No temp files (.body, .issuebody) left in working directory
- [ ] `git status` is clean

---

## What NOT To Do

- Do NOT create `ui/` or `backend/` directories — those get scaffolded when development begins
- Do NOT modify CLAUDE.md, HANDOFF.md, or docs/ content — research phase is complete
- Do NOT run `docker extension init` — the scaffold is already in place and customized
- Do NOT use PowerShell directly for `gh` commands — use CMD or `cmd /c` wrapper
