# C++ review report template

Use this structure for review responses. Omit sections that do not apply. Do not use the Unicode middle dot (U+00B7) or em dash (U+2014); separate fields with colons, commas, or the ASCII hyphen-minus.

When quoting a rule, choose a verbatim excerpt that does not contain either prohibited character. Never edit or normalize text presented as a verbatim quote.

## Review strategy

- **Mode and scope:** <mode>; <commit range, files, or warnings>
- **Guidelines source:** <absolute selected path, SHA-256, refresh or offline status, completeness unverified or confirmed with evidence, and timestamp when available>
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

- **Kind:** Defect for a demonstrated failure path; Risk for concrete interface misuse or a project-enforced rule departure; Question for an unresolved contract.
- **Severity:** Follow repository conventions. Otherwise use High for reachable undefined behavior, races, deadlocks, common-path crashes, or a failing CI gate; Medium for wrong results, error-path leaks, or contract violations; Low for demonstrated maintainability, portability, or performance costs.
- **No findings:** State "No findings." and list the lenses and rules reviewed.
- **Analyzer triage:** For each warning, report its location, rule, verdict (true positive, false positive, or needs context), reason, and action.
- **Adoption plan:** Replace findings with ordered phases. For each phase, name the rules or profile, checks to enable, configuration change, expected impact, and exit condition.

## Example

### 1. High: Defect: Reference into a vector used after `push_back` - `src/registry.cpp:88`

- **Rule:** [ES.65: Don't dereference an invalid pointer](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#res-deref) - "Dereferencing an invalid pointer, such as `nullptr`, is undefined behavior"; source lines `<start>-<end>` read. A Note extends "pointer" to iterators and views.
- **Evidence:** `auto& first = entries_[0];` at line 84, then `entries_.push_back(e);` at line 87, then `first.count++;` at line 88. `push_back` reallocates when `size() == capacity()`, invalidating `first`.
- **Consequence:** A write through a dangling reference whenever the vector grows, causing heap corruption.
- **Fix:** Keep the index and re-read `entries_[0]` after `push_back`.
