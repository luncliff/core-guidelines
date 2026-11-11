# Build, host, and target discovery

Establish which rules apply, which checks produce evidence, and where each check runs. Record each fact with its source (`path:line` or command output). Mark unsettled facts "unknown".

Before classifying the target, inspect README build steps, all project presets (including `CMakeUserPresets.json`), toolchain and cross files, and every CI job that builds or tests an affected target. Record inputs inspected. A missing file does not prove the project lacks target-specific configuration.

## 1. Facts to collect

| Fact | Where to look |
|---|---|
| Build system and entry points | `CMakeLists.txt`, `CMakePresets.json`, `CMakeUserPresets.json`; `*.sln`, `*.vcxproj`, `Directory.Build.props`; `MODULE.bazel`, `WORKSPACE`, `BUILD`, `.bazelrc`; `meson.build` and cross files; `Makefile`, `configure.ac`; `conanfile.*`, `vcpkg.json`; build scripts and README build steps |
| Target containing changed files | `add_library`/`add_executable` sources, `.vcxproj` item lists, Bazel `srcs`, Meson targets |
| Target OS, architecture, ABI | CMake toolchain files (`CMAKE_SYSTEM_NAME`, `CMAKE_SYSTEM_PROCESSOR`, `CMAKE_SYSROOT`, `CMAKE_CROSSCOMPILING_EMULATOR`), preset `toolchainFile` or `architecture`, `ANDROID_ABI`, `CMAKE_OSX_ARCHITECTURES`, MSBuild `Platform`, Conan host profile, vcpkg triplet, Bazel `--platforms`, Meson `[host_machine]` |
| Compiler, language standard, flags | `CMAKE_CXX_COMPILER`, `CMAKE_CXX_STANDARD` or `cxx_std_*`, `PlatformToolset`, `LanguageStandard`, Bazel `copts`, `-std=` in `compile_commands.json` |
| Configured analysis | `.clang-tidy`, `CMAKE_CXX_CLANG_TIDY`, `CMAKE_CXX_CPPCHECK`, `EnableCppCoreCheck`, `RunCodeAnalysis`, `CodeAnalysisRuleSet`, `*.ruleset`, `sonar-project.properties`, CodeQL workflows, analyzer suppression files, warning flags (`-Wall`, `-Werror`, `/W4`, `/WX`, `-Wthread-safety`) |
| Compilation database | `compile_commands.json` in existing build trees: `build/`, `out/build/<preset>/`, `cmake-build-*` |
| Runtime checks | Sanitizer presets or options (`-fsanitize=`, `/fsanitize=address`), CTest or test framework targets, fuzz targets (`LLVMFuzzerTestOneInput`), Valgrind or `ctest -T memcheck` steps |
| CI coverage | `.github/workflows/*`, `azure-pipelines.yml`, `.gitlab-ci.yml`, `Jenkinsfile`: OS, compiler, configuration, sanitizer, and analysis jobs |
| Host | OS and architecture: `uname -sm`, or `$env:PROCESSOR_ARCHITECTURE` on Windows. Tools on `PATH`: `clang-tidy --version`, `cppcheck --version`. Existing build trees. On Windows, absence from `PATH` does not prove a tool unavailable: `"%ProgramFiles(x86)%\Microsoft Visual Studio\Installer\vswhere.exe" -all -format json` lists Visual Studio installations. Check each for MSVC and bundled `VC\Tools\Llvm\x64\bin\clang-tidy.exe`; they may need the developer environment (`VC\Auxiliary\Build\vcvars64.bat`) |

## 2. Classify host and target

Host runs review tools; target runs built code. Classify each affected target from positive evidence of its OS, architecture, and ABI.

- Host OS or a missing toolchain file alone does not prove a native target.
- A CMake generator such as Ninja names the build tool, not the target ABI.
- No configured tree: use target compiler or configuration evidence, presets, build commands, or CI runner configuration. Still unsettled: mark unknown.

| Class | Evidence | Consequence |
|---|---|---|
| Native | Host OS, architecture, and ABI match the target | Static and runtime checks run on the host when a matching configuration exists |
| Cross, runnable | Configured emulator, simulator, or device runner, such as `CMAKE_CROSSCOMPILING_EMULATOR` | Runtime checks run only through that runner |
| Cross, static only | Different OS, architecture, or ABI; no configured runner (for example bare-metal firmware) | Static analysis only, on the host; runtime evidence belongs to target CI or hardware tests |
| Unknown | Missing or conflicting evidence | State the gap; plan only checks that do not depend on the answer |

Static analyzers run on the host but need the target's flags, defines, include paths, and sysroot. Analyze target-specific `#if` branches only with the target's configuration; host-only analysis can miss or misread them.

## 3. Build facts that change rule selection

| Build fact | Effect |
|---|---|
| Language standard | C++20 provides concepts (T.10–T.13), `std::span` (F.24), and `std::jthread` (CP.25). Earlier standard: recommend GSL equivalents or T.48 |
| Exceptions disabled (`-fno-exceptions`, or MSVC without an `/EH` model) | Apply E.25–E.27 instead of E.2 and E.3 |
| RTTI disabled (`-fno-rtti`, MSVC `/GR-`) | `dynamic_cast` (C.146–C.148) unavailable; review the alternative dispatch (C.153) |
| No heap, real-time, or freestanding target | Weight R.5 and Per.15; allocation on a critical path is a defect |
| GSL dependency present | Suggest `owner`, `not_null`, `span`, `narrow`, `finally`, and `Expects` only when available or adoptable |
| Target threading model | Single-threaded targets lower CP severity; interrupts, RTOS tasks, and multicore targets raise it |
| Data model (LP64, LLP64, 32-bit) | Changes integer widths for ES.46 and ES.100–ES.107; `long` is 32-bit on LLP64 Windows |
| Warnings as errors or a blocking analyzer gate | Report CI-failing findings first |

## 4. Choose checks

Prefer the project's configured command over ad hoc invocation: it carries target flags and the team's rule selection. Choose the smallest set that confirms or refutes the change's risks.

| Situation | Plan |
|---|---|
| `compile_commands.json` exists, `clang-tidy` available | `clang-tidy -p <build-dir> <changed files>`; reads the project's `.clang-tidy`. No `.clang-tidy`: a supplemental `--checks=-*,cppcoreguidelines-*,bugprone-*` run is allowed; label it ad hoc |
| No compilation database, native target, `clang-tidy` available | `clang-tidy <changed files> -- -std=<standard> -I<include dirs> -D<defines>`, flags read from build files. Label it ad hoc: missing flags can hide or invent warnings |
| Cppcheck available, compilation database exists | `cppcheck --project=<build-dir>/compile_commands.json --file-filter=<changed path>` |
| MSVC project with Code Analysis configured | Use the configured build; Core Check needs the MSVC toolchain on Windows. Building falls under section 5 |
| Native target with existing sanitizer configuration and tests | Run affected tests from that configuration |
| Cross or unknown target | Run static analysis with target flags when possible; name the CI job or device for runtime checks |
| No usable tool or build data | Review manually against the rules; list checks that would confirm each finding |

Header diagnostics appear only for headers matching `HeaderFilterRegex` (or `--header-filter`). On Windows, clang-tidy can report mixed separators such as `include\lib/buffer.h`, so a filter written with `/` can silently drop header warnings. Pass a separator-agnostic `--header-filter` (for example `.*lib.*` or `include[\\/]lib[\\/].*`) and state this in the strategy.

Sanitizer support varies by toolchain and platform. MSVC provides `/fsanitize=address`, `kernel-address`, and `fuzzer`. Clang documents ThreadSanitizer for 64-bit platforms and MemorySanitizer for Linux, NetBSD, and FreeBSD. Check `tools.md` before planning a runtime check for a target.

## 5. Decision boundaries

| Action | Rule |
|---|---|
| Read files, inspect configuration, query tool versions | Proceed |
| Run a read-only analyzer against an existing build tree, or with flags read from build files, with output to the console or a temporary directory | Proceed; record the exact command and a result summary |
| Run existing tests from a build tree present before the review began, when the command is documented and needs no network, credentials, or devices | Proceed; record the command and results |
| Configure or populate a new build tree (including `cmake --preset`, `cmake -S ... -B ...`, and later build or test commands), run a full or long build, install tools, contact servers, use devices, or apply fixes (`-fix`) | Ask the user first. Non-interactive review: do not run it; list it under "Checks not run" with the reason and the CI job or environment where it belongs. Temporary and git-ignored directories still count as new build trees |
| Edit tracked files, commit, or push | Outside this skill's scope |

## 6. Strategy statement

Before findings, state the strategy in the "Review strategy" section of `review-report-template.md`: affected targets and build system, host and target class, rule-affecting build facts, and checks run or not run, each with evidence.
