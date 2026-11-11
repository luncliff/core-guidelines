# C++ build and analysis tools

Tool map for C++ correctness, safety, security, and maintainability. Most tools complement analyzers or tests; they do not implement the Guidelines. No single tool checks every rule or proves a program correct.

## Compiler diagnostics and guideline checkers

| Tool | Use | Guidelines coverage |
|---|---|---|
| [GCC](https://gcc.gnu.org/onlinedocs/gcc/Warning-Options.html), [Clang](https://clang.llvm.org/docs/DiagnosticsReference.html), and [MSVC](https://learn.microsoft.com/en-us/cpp/build/reference/compiler-options?view=msvc-170) warnings | Catch language errors and suspicious constructs at compile time. Enable a deliberate warning set; fix warnings in CI. Flags and warning quality differ by compiler. | General diagnostics, not a Core Guidelines ruleset. |
| [clang-tidy](https://clang.llvm.org/extra/clang-tidy/checks/list.html) | Extensible AST-based linting and static checks. `cppcoreguidelines-*` checks selected guideline patterns; other groups cover bug-prone constructs, readability, modernization, and more. Works best with actual compile commands. | Direct but partial and configurable. `.clang-tidy` selects checks and options. Fix-its: `-export-fixes` saves them, `-fix` applies them. Review automated edits like any patch. [Usage and configuration](https://clang.llvm.org/extra/clang-tidy/) |
| [MSVC C++ Code Analysis / C++ Core Check](https://learn.microsoft.com/en-us/cpp/code-quality/using-the-cpp-core-guidelines-checkers?view=msvc-170) | Static analysis integrated with Visual Studio and MSBuild. C++ Core Check covers Microsoft's supported subset of guideline rules. | Direct but partial; best for MSVC and Windows projects. The default Microsoft Native Recommended ruleset enables only a subset; project rule sets can select more. Diagnostics such as C26494 and C26485 name the related guideline; `[[gsl::suppress]]` documents intentional exceptions. |

## General-purpose static analysis

These tools inspect code without running analyzed paths. Findings depend on configuration, language support, and analysis depth.

| Tool | Use | Notes |
|---|---|---|
| [Clang Static Analyzer](https://clang.llvm.org/docs/ClangStaticAnalyzer.html) | Path-sensitive, interprocedural symbolic analysis of C and C++; some workflows support cross-translation-unit analysis. | Reports show the code path to each finding. `scan-build` emits HTML reports; [CodeChecker](https://codechecker.readthedocs.io/) adds storage, comparison, filtering, and triage. Not a Core Guidelines checker. [Driver comparison and output](https://clang.llvm.org/docs/analyzer/user-docs/CommandLineUsage.html) |
| [GCC `-fanalyzer`](https://gcc.gnu.org/onlinedocs/gcc/Static-Analyzer-Options.html) | Costlier coverage-guided symbolic analysis of program paths; complements GCC warnings. | GCC documentation describes it as suitable for C only. Do not count it as C++ coverage. |
| [Cppcheck](https://cppcheck.sourceforge.io/manual.html) | Cross-platform C/C++ checks for errors, undefined behavior, and risky patterns. | A compilation database or project file supplies compiler options and preprocessor configuration. Its build directory caches results for repeat runs; diagnostics and suppressions fit CI. Not a Core Guidelines ruleset. [Project and build options](https://github.com/cppcheck-opensource/cppcheck/blob/main/man/manual.md) |
| [SonarQube CFamily](https://docs.sonarsource.com/sonarqube-server/analyzing-source-code/languages/c-family/) | C, C++, and Objective-C analysis for bugs, vulnerabilities, and maintainability. Some rules use symbolic execution or cross-translation-unit taint analysis from untrusted data to sensitive operations. | Pull request analysis reports new issues, decorates the PR, and applies a quality gate; runs in CI with build context. MISRA-related rule tags are not full MISRA or Core Guidelines conformance. Features and language support vary by edition and version. [Analysis details](https://docs.sonarsource.com/sonarqube-server/2026.2/analyzing-source-code/languages/c-family/understanding-the-analysis); [PR analysis](https://docs.sonarsource.com/sonarqube-server/discovering/code-analysis/pull-request-analysis) |
| [PVS-Studio](https://pvs-studio.com/en/docs/) | Commercial C/C++ analysis for defects and security issues; build-system integrations; standards diagnostics such as MISRA, AUTOSAR, and OWASP. | IDE and CI workflows, result baselining, incremental analysis. Modified-files PR analysis focuses on changed code; vendor docs mark it Enterprise-license-only. Not a Core Guidelines implementation. [Modified-file analysis](https://pvs-studio.com/en/docs/manual/0055/) |
| [Coverity Static Analysis](https://www.blackduck.com/static-analysis-tools-sast/coverity.html) | Commercial interprocedural analysis for complex quality defects and security issues in large codebases. | Issues show source events and path context; triage classifies, prioritizes, assigns, and comments. Supports MISRA and CERT reporting. Findings depend on Coverity checkers and supported compiler and build capture. [Triage workflow](https://docs.blackduck.com/r/coverity/2026.6/coverity-documentation/triaging-issues.html) |
| [CodeQL](https://docs.github.com/en/code-security/code-scanning/introduction-to-code-scanning/about-code-scanning-with-codeql) | Query-based security and correctness analysis for repeatable CI code scanning and custom queries. | C/C++ analysis builds a database from the project. GitHub code scanning shows repository alerts and data-flow paths for supported findings. Default and security-extended suites trade breadth for runtime. Custom queries can encode project checks; built-in suites are not Core Guidelines checks. [C/C++ query suites](https://docs.github.com/en/code-security/reference/code-scanning/codeql/codeql-queries/c-cpp-built-in-queries) |

## Runtime instrumentation and memory diagnostics

These tools find errors only on executed paths; they do not review diffs. Each reports a failure, stack, or allocation context when a test or workload reaches faulty behavior. Run them in CI with broad test coverage; keep a reproducing input or test case when possible.

| Tool | Detects | Notes |
|---|---|---|
| [AddressSanitizer (ASan)](https://clang.llvm.org/docs/AddressSanitizer.html) | Out-of-bounds access, use-after-free, invalid frees, related memory errors. | Clang and GCC on supported platforms; MSVC on Windows. Instrument the code, then run tests. |
| [UndefinedBehaviorSanitizer (UBSan)](https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html) | Selected undefined behavior: invalid shifts, misaligned access, some integer and type errors. | Checks are selectable; it does not detect every undefined behavior. |
| [ThreadSanitizer (TSan)](https://clang.llvm.org/docs/ThreadSanitizer.html) | Data races in multithreaded execution. | Platform support varies. Instrumented runs can be much slower and use more memory. |
| [MemorySanitizer (MSan)](https://clang.llvm.org/docs/MemorySanitizer.html) | Use of uninitialized memory. | Needs compatible instrumentation across code and many dependencies; harder to adopt. |
| [LeakSanitizer (LSan)](https://clang.llvm.org/docs/LeakSanitizer.html) | Heap memory leaks. | Often available with ASan on supported platforms. |
| [Valgrind Memcheck](https://valgrind.org/docs/manual/mc-manual.html) | Invalid reads and writes, uninitialized values, invalid frees, leaks. | Runs programs under instrumentation at higher runtime cost; availability depends on target platform. |

## Fuzz testing

Fuzzers generate inputs to explore code paths and expose crashes, hangs, and sanitizer-detected defects. They need a fuzz target at an input boundary and work best with sanitizers.

| Tool | Characteristics |
|---|---|
| [libFuzzer](https://llvm.org/docs/LibFuzzer.html) | In-process, coverage-guided. Integrates with Clang instrumentation and sanitizers. Best for libraries and parsers with a callable fuzz-target function. Keeps coverage-expanding inputs; minimizes failures to reproducers. |
| [AFL++](https://aflplus.plus/) | Coverage-guided, with multiple instrumentation modes for source-instrumented and some binary-only targets. Suits long campaigns and varied build environments. Saves crashing inputs for reproduction and regression tests. |

## Build-system and CI integration

Build systems make analysis reproducible: same defines, include paths, and language options as the real build. They orchestrate tools; they do not replace them.

- **CMake** can run [clang-tidy](https://cmake.org/cmake/help/latest/prop_tgt/LANG_CLANG_TIDY.html) and [Cppcheck](https://cmake.org/cmake/help/latest/prop_tgt/LANG_CPPCHECK.html) in supported Makefile and Ninja builds. `CMAKE_EXPORT_COMPILE_COMMANDS` generates `compile_commands.json`; CMake documents it for Makefile and Ninja generators.
- **Visual Studio/MSBuild** integrates MSVC Code Analysis and C++ Core Check into Windows project builds.
- **SonarQube Build Wrapper** captures compiler invocations for CFamily analysis. Supported scanner workflows also accept a compilation database.
- **Bazel** can run clang-tidy through external tooling such as [bazel_clang_tidy](https://github.com/erenon/bazel_clang_tidy); check compatibility with the repository's Bazel version first.
- **CI** runs deterministic compiler and static checks on changes, plus sanitizer tests and fuzzing when time and platform allow. Keep slow or platform-specific jobs off the fast feedback path.

## Using findings in code review

- Prefer PR annotations or a quality-gate summary for changed code; keep a full-project report for broader cleanup. Some products have changed-file or new-code modes; others need a full scan compared to a stored baseline.
- Version the rule set so local and CI results agree. Supply compiler flags through a compilation database, build wrapper, or supported project integration.
- Triage each finding against its trace, source, and runtime conditions before calling it a defect. Static analyzers have false positives and false negatives; runtime tools cover executed paths only.
- Suppress only understood findings, with a reason where supported. Baseline legacy issues apart from new findings so the gate stays useful and regressions stay visible.
- Before choosing a tool, check edition, language-standard support, compiler compatibility, and CI integration; commercial features vary by license and release.

## Choosing a combination

- Direct guideline feedback: selected clang-tidy `cppcoreguidelines-*` checks; add MSVC C++ Core Check with MSVC.
- Broader static findings: compiler warnings plus one or more of Clang Static Analyzer, Cppcheck, SonarQube, PVS-Studio, Coverity, or CodeQL. Choose by language and build support, reporting workflow, and licensing.
- Runtime-only defects: tests under the matching sanitizer; Memcheck where appropriate.
- Parsers and other input-heavy components: fuzz, preferably with ASan or UBSan. Keep minimized failing inputs as regression tests.
- Adopt incrementally: establish a baseline, tune false-positive handling, then block new high-confidence issues.
