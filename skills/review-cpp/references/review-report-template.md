# C++ review report template

Use this structure for review responses; omit sections that do not apply. Never output the middle dot (U+00B7) or em dash (U+2014); separate fields with colons, commas, or ASCII hyphen-minus. Quote rule excerpts free of both characters. Never edit or normalize text presented as a verbatim quote.

## Review strategy

- **Mode and scope:** <mode>; <commit range, files, or warnings>
- **Guidelines source:** <absolute path or fetched URL>, or "unavailable: <reason>"
- **Build and environment:** <build system and affected target>; host <OS/architecture> to target <OS/architecture/ABI> (<class>); evidence: <paths>
- **Rule-affecting build facts:** <facts and effects>, or "none found"
- **Guideline scope:** <relevant lenses> to <sections actually consulted> to <rules read in full>
- **Checks run:** `<command>` to <result summary>, or "none"
- **Checks not run:** <check>; <reason>; <where it belongs>

## Findings

### 1. <Severity>: <Kind>: <short title> - `<path:line>`

- **Rule:** [<ID>: <title>](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#<anchor>) - "<verbatim quote>"; source lines `<start>-<end>` read
- **Evidence:** <code facts, path, sequence, interleaving, or tool output>
- **Consequence:** <what goes wrong>
- **Fix:** <one specific change or check, including affected callers>

## Questions

- `<path:line>`: <question>; <what the answer decides>

## Limitations

- <coverage gaps and unverified assumptions>

## Reporting notes

- **Kind:** Defect: demonstrated failure path. Risk: concrete interface misuse or departure from a project-enforced rule. Question: unresolved contract.
- **Severity:** Follow repository conventions. Otherwise:
  - High: reachable undefined behavior, races, deadlocks, common-path crashes, or a failing CI gate.
  - Medium: wrong results, error-path leaks, or contract violations.
  - Low: demonstrated maintainability, portability, or performance costs.
- **No findings:** State "No findings." and list the lenses and rules reviewed.
- **Analyzer triage:** Per warning: location, rule, verdict (true positive, false positive, or needs context), reason, and action.
- **Adoption plan:** Replace findings with ordered phases. Per phase: rules or profile, checks to enable, configuration change, expected impact, and exit condition.

## Example

### 1. High: Defect: Reference into a vector used after `push_back` - `src/registry.cpp:88`

- **Rule:** [ES.65: Don't dereference an invalid pointer](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#res-deref) - "Dereferencing an invalid pointer, such as `nullptr`, is undefined behavior"; source lines `<start>-<end>` read. A Note extends "pointer" to iterators and views.
- **Evidence:** `auto& first = entries_[0];` at line 84, then `entries_.push_back(e);` at line 87, then `first.count++;` at line 88. `push_back` reallocates when `size() == capacity()`, invalidating `first`.
- **Consequence:** Write through a dangling reference whenever the vector grows, causing heap corruption.
- **Fix:** Keep the index and re-read `entries_[0]` after `push_back`.
