# C++ Build and Analysis Tools

Practical tool map for better C++ correctness, safety, security, and maintainability. Scope exceeds C++ Core Guidelines checking: most tools complement analyzers or tests; they do not implement the Guidelines. No single tool checks every rule or proves a program correct.

## Compiler diagnostics and guideline checkers

| Tool | What it is useful for | Guidelines coverage |
|---|---|---|
| [GCC](https://gcc.gnu.org/onlinedocs/gcc/Warning-Options.html), [Clang](https://clang.llvm.org/docs/DiagnosticsReference.html), and [MSVC](https://learn.microsoft.com/en-us/cpp/build/reference/compiler-options?view=msvc-170) warnings | Catch language errors and suspicious constructs at compile time. Enable a deliberate warning set and address warnings in CI; flags and warning quality differ by compiler. | General compiler diagnostics, not a Core Guidelines ruleset. |
| [clang-tidy](https://clang.llvm.org/extra/clang-tidy/checks/list.html) | Extensible, AST-based C++ linting and static checks. The `cppcoreguidelines-*` group checks selected guideline patterns; other groups cover bug-prone constructs, readability, modernization, and more. Usually works best with actual compile commands. | Direct, but partial and configurable. Use `.clang-tidy` to select checks and options. Diagnostics can include fix-its; `-export-fixes` saves suggested edits and `-fix` applies them. Review automated changes like any other patch. [Usage and configuration](https://clang.llvm.org/extra/clang-tidy/) |
| [MSVC C++ Code Analysis / C++ Core Check](https://learn.microsoft.com/en-us/cpp/code-quality/using-the-cpp-core-guidelines-checkers?view=msvc-170) | Visual Studio/MSBuild-integrated static analysis. C++ Core Check includes Microsoft's supported subset of guideline rules. | Direct, but partial; especially handy for MSVC and Windows projects. The default Microsoft Native Recommended ruleset enables only a subset; project rule sets can select more. Diagnostics such as C26494 and C26485 include the related guideline identifier, and `[[gsl::suppress]]` can document intentional exceptions. |

## General-purpose static analysis

These tools inspect code without running analyzed paths. Findings depend on configuration, supported language features, and analysis depth.

| Tool | What it is useful for | Notes |
|---|---|---|
| [Clang Static Analyzer](https://clang.llvm.org/docs/ClangStaticAnalyzer.html) | Path-sensitive, interprocedural symbolic analysis for C and C++ bugs; selected workflows support cross-translation-unit analysis. | Reports include a code path showing how a condition leads to a finding. `scan-build` can emit HTML reports; [CodeChecker](https://codechecker.readthedocs.io/) adds report storage, comparison, filtering, and triage. Not a Core Guidelines checker. [Driver comparison and output](https://clang.llvm.org/docs/analyzer/user-docs/CommandLineUsage.html) |
| [GCC `-fanalyzer`](https://gcc.gnu.org/onlinedocs/gcc/Static-Analyzer-Options.html) | More expensive, coverage-guided symbolic analysis of program paths; complements GCC warnings. | GCC documentation currently describes this analyzer as suitable for C only. Do not count it as C++ analyzer coverage. |
| [Cppcheck](https://cppcheck.sourceforge.io/manual.html) | Cross-platform C/C++ checks for errors, undefined behavior, and risky coding patterns; can analyze projects with build information. | Importing a compilation database or project file supplies compiler options and preprocessor configuration. Its build directory caches analysis results to speed repeat runs; diagnostics and suppressions can go into CI. Not a Core Guidelines ruleset. [Project and build options](https://github.com/cppcheck-opensource/cppcheck/blob/main/man/manual.md) |
| [SonarQube CFamily](https://docs.sonarsource.com/sonarqube-server/analyzing-source-code/languages/c-family/) | C, C++, and Objective-C analysis for bugs, vulnerabilities, and maintainability findings. Some CFamily rules use symbolic execution or cross-translation-unit taint analysis to trace untrusted data to sensitive operations. | Pull request analysis can report new issues, decorate the PR, and apply a quality gate; it runs in CI and requires build context. It can tag some rules by MISRA relation, but that is not full MISRA or Core Guidelines conformance. Product features and language support vary by edition/version. [Analysis details](https://docs.sonarsource.com/sonarqube-server/2026.2/analyzing-source-code/languages/c-family/understanding-the-analysis); [PR analysis](https://docs.sonarsource.com/sonarqube-server/discovering/code-analysis/pull-request-analysis) |
| [PVS-Studio](https://pvs-studio.com/en/docs/) | Commercial C/C++ static analysis for defects and security-relevant issues; supports build-system integrations and standards-oriented diagnostics such as MISRA, AUTOSAR, and OWASP. | Offers IDE and CI workflows, result baselining, and incremental analysis. Modified-files PR/merge analysis can focus reports on changed code, but vendor documentation marks this workflow Enterprise-license-only. Not a Core Guidelines implementation. [Modified-file analysis](https://pvs-studio.com/en/docs/manual/0055/) |
| [Coverity Static Analysis](https://www.blackduck.com/static-analysis-tools-sast/coverity.html) | Commercial, interprocedural static analysis for complex code-quality defects and security issues across large codebases. | Issue details include source events/path context; triage workflows classify, prioritize, assign, and comment on findings. Supports coding-standard reporting such as MISRA and CERT. Findings depend on Coverity checkers and supported compiler/build capture. [Triage workflow](https://docs.blackduck.com/r/coverity/2026.6/coverity-documentation/triaging-issues.html) |
| [CodeQL](https://docs.github.com/en/code-security/code-scanning/introduction-to-code-scanning/about-code-scanning-with-codeql) | Query-based security and correctness analysis, especially useful for repeatable CI/code-scanning workflows and custom queries. | C/C++ analysis uses a database built from the project. GitHub code scanning presents repository alerts and can show data-flow paths for supported findings; default and security-extended query suites trade scan breadth for runtime. Custom queries can encode project-specific checks, but built-in suites are not Core Guidelines checks. [C/C++ query suites](https://docs.github.com/en/code-security/reference/code-scanning/codeql/codeql-queries/c-cpp-built-in-queries) |

## Runtime instrumentation and memory diagnostics

These tools detect errors only when tests or workloads execute relevant code. They complement static analysis and benefit from broad test coverage.

Sanitizers and Memcheck do not review diffs. They report a runtime failure, stack, or allocation context only when a test/workload reaches faulty behavior. Run their tests in CI and preserve a reproducing input or test case when possible.

| Tool | Detects / helps investigate | Notes |
|---|---|---|
| [AddressSanitizer (ASan)](https://clang.llvm.org/docs/AddressSanitizer.html) | Out-of-bounds memory access, use-after-free, invalid frees, and related memory errors. | Available with Clang/GCC on supported platforms; MSVC also provides ASan on Windows. Instrument the code and run tests. |
| [UndefinedBehaviorSanitizer (UBSan)](https://clang.llvm.org/docs/UndefinedBehaviorSanitizer.html) | Selected kinds of undefined behavior, such as invalid shifts, misaligned access, and some integer or type errors. | Checks are selectable; it does not detect every form of undefined behavior. |
| [ThreadSanitizer (TSan)](https://clang.llvm.org/docs/ThreadSanitizer.html) | Data races in multithreaded execution. | Platform support varies. Instrumented execution can be substantially slower and use more memory. |
| [MemorySanitizer (MSan)](https://clang.llvm.org/docs/MemorySanitizer.html) | Use of uninitialized memory. | Requires compatible instrumentation across code and many dependencies, which can make adoption harder. |
| [LeakSanitizer (LSan)](https://clang.llvm.org/docs/LeakSanitizer.html) | Heap memory leaks. | Often available alongside ASan on supported platforms. |
| [Valgrind Memcheck](https://valgrind.org/docs/manual/mc-manual.html) | Invalid reads/writes, use of uninitialized values, invalid frees, and leaks. | Runs programs under instrumentation and generally costs more runtime; availability depends on target platform. |

## Fuzz testing

Fuzzers generate many inputs to explore code paths and expose crashes, hangs, and sanitizer-detected defects. They need a suitable input boundary (a fuzz target) and work best with sanitizers.

| Tool | Characteristics |
|---|---|
| [libFuzzer](https://llvm.org/docs/LibFuzzer.html) | In-process, coverage-guided fuzzing. Integrates with Clang instrumentation and can run with sanitizers. Best for libraries and parsers with a callable fuzz-target function. Keeps inputs that expand coverage and can minimize failures to reproducers. |
| [AFL++](https://aflplus.plus/) | Coverage-guided fuzzing with multiple instrumentation modes and support for source-instrumented and some binary-only targets. Useful for longer-running campaigns and varied build environments; saves crashing inputs for reproducing and regression testing. |

## Build-system and CI integration

Build systems make analysis reproducible: they supply the same defines, include paths, and language options as the real build. They orchestrate tools; they do not replace them.

- **CMake** can invoke [clang-tidy](https://cmake.org/cmake/help/latest/prop_tgt/LANG_CLANG_TIDY.html) and [Cppcheck](https://cmake.org/cmake/help/latest/prop_tgt/LANG_CPPCHECK.html) alongside supported Makefile and Ninja builds. `CMAKE_EXPORT_COMPILE_COMMANDS` can generate `compile_commands.json` for tools that consume a compilation database; CMake documents this option for Makefile and Ninja generators.
- **Visual Studio/MSBuild** integrates MSVC Code Analysis and C++ Core Check into Windows project builds.
- **SonarQube Build Wrapper** captures compiler invocations for CFamily analysis. SonarQube also documents a compilation-database mode for supported scanner workflows.
- **Bazel** users can integrate clang-tidy through external tooling such as [bazel_clang_tidy](https://github.com/erenon/bazel_clang_tidy); check compatibility with the repository's Bazel version before adoption.
- **CI** should run deterministic compiler/static-analysis checks on changes, plus sanitizer tests and fuzzing jobs when time and platform allow. Keep slower or platform-specific jobs separate from the fast feedback path.

## Using findings in code review

- Prefer pull-request annotations or a quality-gate summary for findings on changed code; keep a full-project report for broader cleanup. Some products offer changed-file/new-code modes; others need a full scan and comparison to a stored baseline.
- Configure and version the rule set so local and CI results agree. For tools that need compiler flags, use a compilation database, build wrapper, or supported project integration.
- Triage findings against their trace, relevant source, and runtime conditions before calling them defects. Static analyzers can have false positives and false negatives; runtime tools cover executed paths only.
- Suppress only understood findings, with a reason where supported. Baseline legacy issues separately from new findings so the review gate stays useful without hiding regressions.
- Check product edition, language-standard support, compiler compatibility, and CI integration before choosing a tool; commercial feature availability can differ by license and release.

## Choosing a useful combination

- For direct guideline feedback, start with selected clang-tidy `cppcoreguidelines-*` checks; add MSVC C++ Core Check when using MSVC.
- For broader static findings, combine compiler warnings with one or more analyzers: Clang Static Analyzer, Cppcheck, SonarQube, PVS-Studio, Coverity, or CodeQL. Choose by language/build support, reporting workflow, and licensing.
- For runtime-only defects, run tests with the matching sanitizer and use Memcheck where appropriate.
- Fuzz parsers and other input-heavy components, preferably with ASan or UBSan. Preserve minimized failing inputs as regression tests.
- Introduce findings incrementally: establish a baseline, tune false-positive handling, then block new high-confidence issues from entering the codebase.
