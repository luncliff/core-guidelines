---
name: review-cpp
description: Reviews C++ code using the C++ Core Guidelines and project-specific build evidence. Use for C++ changes, module audits, ownership and lifetime questions, analyzer warnings, or Guidelines adoption plans.
tools: [read, search]
---

# C++ Core Guidelines review

Produce actionable findings grounded in code, applicable Guidelines, and observed tool results. Keep the reviewed project unchanged. Scale the review to its scope.

## Workflow

### 1. Establish scope

- Identify the requested diff, files, question, analyzer warnings, or adoption goal.
- For a change review, report defects introduced or worsened by the change. For an audit, inspect the named files or module.
- With no stated scope, start with uncommitted changes. If none exist, establish the intended comparison before choosing a branch or commit range.
- Read [toc.md](references/toc.md) and [tools.md](references/tools.md) to plan the review. Batch these reads when the host supports it.

### 2. Select the Guidelines source

- Resolve `<skill-dir>` from this loaded skill's location. Use absolute paths, quote paths with spaces, and adapt command syntax to the host shell.
- The full Guidelines document is supplied separately. Use an explicit user-provided file when given.
- Otherwise, retrieve the latest upstream document when the user has approved both network access and destination writes. Use curl when available; the source reference provides alternatives and fallback paths. Host tool approval still applies.
- Follow [guidelines-source.md](references/guidelines-source.md) to validate every candidate, including a user-provided file, and record its digest.
- Resolve cache roots from the host environment at runtime, not from a guessed user or home path.
- If refresh is unavailable, use a validated cache or current-directory copy. In non-interactive reviews, record a skipped refresh rather than assuming permission.
- If no valid document is available, disclose that rule-grounded review is blocked. In an interactive review, ask for a file or download approval.
- Report the selected absolute path, digest, refresh status, and completeness. Use "completeness unverified" until an expected length or a trusted digest confirms it. On a failed refresh, use only a validated local copy and disclose its freshness limit.

### 3. Understand code and build context

- Read the scoped code, complete affected functions, and relevant callers. Trace ownership, borrowed lifetimes, shared state, and error paths.
- Use [build-discovery.md](references/build-discovery.md) to inspect relevant build instructions, presets, toolchain files, and CI jobs.
- Establish the language standard, compiler flags, host, and target from evidence. A host OS, generator, or missing toolchain alone does not establish a native target. Keep unknown facts explicit; select target-compatible checks.

### 4. Select and read rules

- Use the index for sections and [review-lenses.md](references/review-lenses.md) for code signals. Apply lens 0's calibration guards and the lenses relevant to the code.
- List rules in selected sections. Before applying a rule:
  - Locate its starting line and the next heading at the same or higher level, or the end of the file. For a profile list item, locate the next sibling item or heading instead.
  - Read the entire range from the start through that boundary, including Exceptions and Notes. Request the range in one read when possible. If output is truncated, continue with contiguous reads until all lines are visible.
  - Record the source line range actually read alongside each cited rule.
- Cite only rules whose full ranges were read. State any incomplete range as a limitation and leave that rule out of findings.

### 5. Gather evidence

- Use the tool reference to select checks. Prefer the project's configured analysis and tests over ad hoc commands.
- Follow the action boundaries in the build discovery reference. Obtain approval for new build trees, long builds, tool installation, devices, and network access.
- Host permissions govern every action, including the refresh and cache writes. When approval or tools are unavailable, record the skipped check and reason.
- Treat analyzer diagnostics as leads. Distinguish tool/environment failures from reachable code defects; report only commands actually run.

### 6. Judge and report

- Read [review-report-template.md](references/review-report-template.md) before drafting. Preserve its required fields in focused answers.
- Admit a finding when it has an in-scope location, an applicable rule, a concrete consequence, and an actionable fix or check. Label tool-only findings.
- Distinguish a demonstrated defect from an interface risk. Put unresolved contracts in Questions. Check exceptions and counterexamples before reporting.
- For a fix returning a borrowed value, establish that its owner outlives every use. A view parameter does not extend a temporary owner's lifetime. Otherwise propose an owning return or ask about the lifetime contract.
- For each rule-based finding, cite a `file:line`, exactly one applicable rule, its canonical `isocpp.github.io` link, complete source line range, and verbatim quote. State evidence, consequence, and one fix or check that accounts for affected callers. Keep hypothetical effects out of findings.
- Use headings through level 3 and nested lists. Use colons, commas, or ASCII hyphens in place of U+00B7 and U+2014; choose quotations that satisfy this rule.
