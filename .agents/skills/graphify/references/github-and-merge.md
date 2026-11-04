# graphify reference: GitHub clone and cross-repo merge

Load this when the user passed one or more `https://github.com/...` URLs, or named several local subfolders to merge into one graph.

### Step 0 - Clone GitHub repo(s) (only if a GitHub URL was given)

**Single repo:**
```bash
LOCAL_PATH=$(scripts/graphify clone <github-url> [--branch <branch>])
# Use LOCAL_PATH as the target for all subsequent steps
```

**Multiple repos (cross-repo graph):**
```bash
# Clone each repo; keep graph commands in the active project root
scripts/graphify clone <url1>   # → ~/.graphify/repos/<owner1>/<repo1>
scripts/graphify clone <url2>   # → ~/.graphify/repos/<owner2>/<repo2>
# Extract into separate directories under this project's graph output
GRAPHIFY_OUT=docs/graphify/inputs/repo1 scripts/graphify extract ~/.graphify/repos/<owner1>/<repo1>
GRAPHIFY_OUT=docs/graphify/inputs/repo2 scripts/graphify extract ~/.graphify/repos/<owner2>/<repo2>
# Then merge into the project graph:
scripts/graphify merge-graphs \
  docs/graphify/inputs/repo1/graph.json \
  docs/graphify/inputs/repo2/graph.json \
  --out docs/graphify/graph.json
```

Graphify clones into `~/.graphify/repos/<owner>/<repo>` and reuses existing clones on repeat runs. Each node in the merged graph carries a `repo` attribute so you can filter by origin.

**Multiple local subfolders (monorepo or multi-service layout):**

The skill pipeline writes intermediate and final outputs to the active project's graph directory (`docs/graphify/` in these examples). Substitute the configured output directory consistently when different. Running extraction on each subfolder with that same output will overwrite the previous graph. For an explicitly requested multi-input merge, use command-local `GRAPHIFY_OUT` values under `docs/graphify/inputs/` when the project's wrapper supports overrides. Keep the working directory at the project root:

```bash
GRAPHIFY_OUT=docs/graphify/inputs/core scripts/graphify extract ./core/
GRAPHIFY_OUT=docs/graphify/inputs/service scripts/graphify extract ./service/
GRAPHIFY_OUT=docs/graphify/inputs/platform scripts/graphify extract ./platform/
# Add --backend gemini|kimi|openai|deepseek|claude-cli depending on which API key you have set

# Then merge at the project root:
scripts/graphify merge-graphs \
  docs/graphify/inputs/core/graph.json \
  docs/graphify/inputs/service/graph.json \
  docs/graphify/inputs/platform/graph.json \
  --out docs/graphify/graph.json
```

These overrides apply only to the extraction command; subsequent commands keep `GRAPHIFY_OUT=docs/graphify` from the repository setup. Do not change to a cloned repository or scanned subfolder to run the wrapper.

Once `docs/graphify/graph.json` exists, the fast path above takes over: any codebase question runs `scripts/graphify query` directly on the merged graph — no re-extraction, no size gate.
