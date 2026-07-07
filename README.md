# Meso

**A project generator for Claude Code.** Bootstrap new Claude Code projects with pre-configured agents, workflows, patterns, quality rules, and examples—all set up and ready to go.

Think `create-react-app` or `cargo init`, but for multi-agent orchestration in Claude Code.

## What is Meso?

Meso scaffolds new Claude Code projects in seconds. When you run `bun meso`:

1. It fetches templates from `https://github.com/killallservers/templates`
2. You pick one
3. It copies the full `.claude/` scaffold into your current directory
4. You get a fully configured Claude Code project with:
   - ✅ Custom agents (architect, code-reviewer)
   - ✅ Reusable workflow patterns (audit, migrate, research, judge panel, etc.)
   - ✅ Quality rules (pre-commit hooks, linting, formatting, type checking)
   - ✅ Orchestration patterns (adversarial verify, dedup, cost-aware scaling, etc.)
   - ✅ Real-world examples for common tasks
   - ✅ Memory persistence for multi-phase workflows
   - ✅ Skills and tech-stack guidance

No setup. No configuration. Just `bun meso` → pick template → start building with Claude Code agents.

## Project Structure

**Meso (the generator):**
```
meso/
├── src/
│   └── cli.ts                 # Meso CLI entry point
├── package.json               # Bun project manifest
├── biome.jsonc                # Biome linter/formatter config
├── tsconfig.json              # TypeScript config
└── README.md                  # This file
```

**What you get after running `bun meso` (the template):**
```
my-project/
├── src/
│   └── cli.ts                 # Your app's entry point
├── .claude/
│   ├── agents/                # Custom agent definitions
│   │   ├── architect.md       # Systems architect for design decisions
│   │   └── code-reviewer.md   # Code review agent
│   ├── rules/                 # Operational guidelines
│   │   ├── default.md         # Security rules, access controls
│   │   ├── pre-commit-checks.md
│   │   ├── failure-scenarios.md
│   │   ├── workflow-composition.md
│   │   ├── anti-patterns.md
│   │   └── patterns/          # Orchestration patterns
│   │       ├── adversarial-verify.md
│   │       ├── dedup.md
│   │       ├── judge-panel.md
│   │       ├── perspective-diverse.md
│   │       ├── phase-orchestration.md
│   │       └── cost-aware.md
│   ├── examples/              # Real-world workflow examples
│   │   ├── 01-security-audit-workflow.md
│   │   ├── 02-drizzle-migration-workflow.md
│   │   ├── 03-bug-discovery-workflow.md
│   │   └── 04-architecture-decision.md
│   ├── memory/                # Persistent memory for multi-phase workflows
│   └── README.md              # Product template (customize for your project)
├── package.json               # Your project's Bun manifest
├── biome.jsonc                # Linting/formatting config
├── tsconfig.json              # TypeScript config
└── README.md                  # Your project's README
```

## Tech Stack

**Meso itself:**
- **Runtime**: [Bun](https://bun.sh) (fast JavaScript runtime)
- **Language**: TypeScript
- **Code Quality**: Biome (linting + formatting)

**Generated projects include:**
- Bun + TypeScript starter
- Pre-configured Biome for linting/formatting
- Custom Claude Code agents and workflows
- Git hooks for quality enforcement
- Persistent memory for multi-phase orchestration

## What Gets Generated

When you run `bun meso`, every generated project includes:

### Custom Agents

Specialized Claude Code workers for your project:

- **architect**: Systems architect for design reviews and architecture decisions
- **code-reviewer**: Code reviewer for PR audits and implementation checks

Defined in `.claude/agents/` and ready to use with `/architect` and `/code-reviewer` commands.

### Workflow Patterns

Reusable orchestration techniques for common multi-agent patterns:

| Pattern | Use Case | Cost |
|---------|----------|------|
| **Adversarial Verify** | High-stakes findings (security, breaking changes) | 3× per finding |
| **Deduplication** | Iterative discovery (loop-until-dry) | O(1) per check |
| **Judge Panel** | Multi-perspective decisions | N drafts + N judges |
| **Perspective-Diverse** | Multi-faceted quality reviews | N lenses per item |
| **Phase Orchestration** | Structuring multi-stage workflows | Optimizes wall-clock |
| **Cost-Aware Scaling** | Open-ended discovery with budget limits | Adapts to token budget |

Documented in `.claude/rules/patterns/` with implementation guides.

### Enforced Quality Standards

Pre-commit hooks ensure code quality:

- ✅ Type checking: `bun run typecheck`
- ✅ Linting: `bun run lint`
- ✅ Formatting: `bun run format`
- ✅ Comprehensive check: `bun run check`

All commits must pass all four checks. No exceptions.

## Getting Started

### Prerequisites

- [Bun](https://bun.sh) (v1.0+)
- [Claude Code](https://claude.ai/code)

### Installation

```bash
git clone https://github.com/killallservers/meso.git
cd meso
bun install
```

### Running Meso

```bash
# Generate a new Claude Code project
bun run src/cli.ts

# Or once installed globally
bun meso
```

This will:
1. Fetch available templates from GitHub
2. Show you a list to choose from
3. Copy the template into your current directory
4. You're ready to start using Claude Code agents!

## Development

```bash
# Run type checking
bun run typecheck

# Run linter
bun run lint

# Format code
bun run format

# Comprehensive check (all of the above)
bun run check
```

### After Generating a Project

Once you've run `bun meso` and selected a template, your new project is ready. Launch Claude Code:

```bash
cd my-project
claude code
```

You'll have access to custom agents configured in `.claude/agents/`:

```
/architect
Design the architecture for a new feature.

/code-reviewer
Review the current diff for bugs and improvements.
```

Plus pre-built workflows and patterns for orchestrating multi-agent work.

## What You Get: Workflow Examples

Every generated project includes real-world examples in `.claude/examples/`:

### Example 1: Security Audit
Find security issues, verify each finding, report confirmed ones.
→ See `.claude/examples/01-security-audit-workflow.md`

### Example 2: Large Refactor
Migrate code across 100+ files in parallel without conflicts.
→ See `.claude/examples/02-drizzle-migration-workflow.md`

### Example 3: Bug Discovery
Iteratively find flaky tests, deduplicate, stop when converged.
→ See `.claude/examples/03-bug-discovery-workflow.md`

### Example 4: Architecture Decision
Draft solutions from different angles, judge each, synthesize best approach.
→ See `.claude/examples/04-architecture-decision.md`

These examples show you how to use the agents, patterns, and rules in your project.

## Roadmap

**Current phase**: MVP — Fetches and scaffolds templates from remote GitHub repo.

### Planned Improvements

- [ ] Support local template discovery
- [ ] Interactive template selection UI
- [ ] Project-specific scaffolding (API, web app, CLI templates with additional files)
- [ ] Template customization options during generation
- [ ] Validation after template copy (check dependencies, config)
- [ ] Built-in template creation/contribution workflow
- [ ] Documentation site with template gallery

## Understanding Workflows vs Subagents vs Skills

| | Subagents | Skills | Workflows |
|---|-----------|--------|-----------|
| **What** | Single specialized worker | Documentation/instructions | Multi-agent orchestration script |
| **Control Flow** | Claude decides each step | Claude follows instructions | Script executes deterministically |
| **Scale** | Few per turn | N/A | Dozens to hundreds of agents |
| **Use when** | Offloading a side task | Teaching Claude a skill | Building repeatable multi-phase work |

## Contributing to Meso

Want to contribute to the generator itself or improve the template? Here's how:

1. **Read the rules**: Review `.claude/rules/default.md` for code quality requirements
2. **Run checks**: Ensure `bun run check` passes before committing
3. **Update the template**: Changes to `.claude/` in this repo become part of every generated project
4. **Add new patterns**: Extend `.claude/rules/patterns/` with reusable orchestration techniques
5. **Improve agents**: Enhance `.claude/agents/` definitions for better specialized workers

All commits are checked by pre-commit hooks for quality. See `.claude/rules/pre-commit-checks.md`.

## Key Files in Generated Projects

Every generated project includes:

- **`.claude/README.md`** — Product template (customize for your project)
- **`.claude/AGENTS.md`** — Available agents and when to use them
- **`.claude/WORKFLOW_SELECTION.md`** — Choosing the right workflow for your task
- **`.claude/rules/default.md`** — Security and access control rules
- **`.claude/rules/pre-commit-checks.md`** — Quality check requirements
- **`.claude/rules/failure-scenarios.md`** — Common failures and recovery patterns
- **`.claude/rules/patterns/*.md`** — Reusable orchestration patterns

## Troubleshooting

**Type errors after commit?**
```bash
bun run typecheck
# Fix errors and commit again
```

**Linting failures?**
```bash
bun run lint   # Auto-fix many issues
bun run check  # See remaining issues
```

**Pre-commit hook failing?**
The hook runs all quality checks automatically. If it fails:
1. Review the error message
2. Run the failing check manually
3. Fix issues and retry commit

To bypass (not recommended):
```bash
git commit --no-verify
```

## Philosophy

Meso (the generator) and the projects it creates are built on these principles:

1. **Zero Setup**: One command gets you a fully configured Claude Code project
2. **Best Practices Baked In**: Quality rules, agents, and patterns included from day one
3. **Determinism**: Workflows execute the same way every time
4. **Parallelism**: Fan out work to agents without manual coordination
5. **Verification**: High-stakes findings are independently verified
6. **Cost Awareness**: Scale to the user's token budget, degrade gracefully
7. **Reusability**: Patterns and agents are composable across projects and domains

## License

Proprietary.

## Contact

Questions or feedback? See `.claude/ROADMAP.md` for planned improvements.
