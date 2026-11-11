# Build, host, and target discovery

Use this file to establish which guideline rules apply, which checks can produce evidence, and where each check must run. Record each fact and its source (`path:line` or command output). Mark a fact "unknown" when evidence does not settle it.

Before classifying the target, inspect README build instructions, enumerate all available project presets (including `CMakeUserPresets.json` when present), inspect toolchain/cross files, and inspect every CI job that builds or tests an affected target. Record inspected inputs. A missing file does not prove the project lacks target-specific configuration.

## 1. Facts to collect

| Question | Where to look |
|---|---|
| Build system and entry points | `CMakeLists.txt`, `CMakePresets.json`, `CMakeUserPresets.json`; `*.sln`, `*.vcxproj`, `Directory.Build.props`; `MODULE.bazel`, `WORKSPACE`, `BUILD`, `.bazelrc`; `meson.build` and cross files; `Makefile`, `configure.ac`; `conanfile.*`, `vcpkg.json`; build scripts and README build steps |
| Target that contains changed files | `add_library`/`add_executable` sources, `.vcxproj` item lists, Bazel `srcs`, Meson targets |
| Target OS, architecture, and ABI | CMake toolchain files (`CMAKE_SYSTEM_NAME`, `CMAKE_SYSTEM_PROCESSOR`, `CMAKE_SYSROOT`, `CMAKE_CROSSCOMPILING_EMULATOR`), preset `toolchainFile` or `architecture`, `ANDROID_ABI`, `CMAKE_OSX_ARCHITECTURES`, MSBuild `Platform`, Conan host profile, vcpkg triplet, Bazel `--platforms`, Meson `[host_machine]` |
| Compiler, language standard, and flags | `CMAKE_CXX_COMPILER`, `CMAKE_CXX_STANDARD` or `cxx_std_*`, `PlatformToolset`, `LanguageStandard`, Bazel `copts`, `-std=` in `compile_commands.json` |
| Configured analysis | `.clang-tidy`, `CMAKE_CXX_CLANG_TIDY`, `CMAKE_CXX_CPPCHECK`, `EnableCppCoreCheck`, `RunCodeAnalysis`, `CodeAnalysisRuleSet`, `*.ruleset`, `sonar-project.properties`, CodeQL workflows, analyzer suppression files, warning flags (`-Wall`, `-Werror`, `/W4`, `/WX`, `-Wthread-safety`) |
| Compilation database | `compile_commands.json` in existing build trees such as `build/`, `out/build/<preset>/`, `cmake-build-*` |
| Runtime checks | sanitizer presets or options (`-fsanitize=`, `/fsanitize=address`), CTest or test framework targets, fuzz targets (`LLVMFuzzerTestOneInput`), Valgrind or `ctest -T memcheck` steps |
| CI coverage | `.github/workflows/*`, `azure-pipelines.yml`, `.gitlab-ci.yml`, `Jenkinsfile`: OS, compiler, configuration, sanitizer, and analysis jobs |
| Host | Commands: OS and architecture (`uname -sm`, or `$env:PROCESSOR_ARCHITECTURE` on Windows); tools on `PATH` (`clang-tidy --version`, `cppcheck --version`); existing build trees. On Windows, absence from `PATH` does not prove a tool unavailable: `"%ProgramFiles(x86)%\Microsoft Visual Studio\Installer\vswhere.exe" -all -format json` lists Visual Studio installations; inspect each installation's MSVC and bundled `VC\Tools\Llvm\x64\bin\clang-tidy.exe` tools. They may need the developer environment (`VC\Auxiliary\Build\vcvars64.bat`) |

## 2. Classify host and target

The host runs review tools. The target runs built code. Classify each affected target from positive evidence for its OS, architecture, and ABI. Host OS alone or missing toolchain file does not prove a native target. A CMake generator such as Ninja names the build tool, not the target ABI. Without an existing configured tree, use target compiler/configuration evidence, presets, build commands, or CI runner configuration. If these do not establish the target, mark it unknown.

| Class | Evidence | Consequence |
|---|---|---|
| Native | Host OS, architecture, and ABI match the target | Static and runtime checks can run on the host when a matching configuration exists |
| Cross, runnable | A configured emulator, simulator, or device runner, such as `CMAKE_CROSSCOMPILING_EMULATOR` | Runtime checks can run only through that runner |
| Cross, static only | Different OS, architecture, or ABI, and no configured runner (for example bare-metal firmware) | Static analysis runs only on the host; runtime evidence belongs to target CI or hardware tests |
| Unknown | Evidence is missing or conflicting | State the gap; plan only checks that do not depend on the answer |

Static analyzers run on the host but need the target's flags, defines, include paths, and sysroot. Analyze code under target-specific `#if` branches only with the target's configuration; host-only analysis can miss or misread it.

## 3. Build facts that change which rules apply

| Build fact | Effect on the review |
|---|---|
| Language standard | C++20 provides concepts (T.10–T.13), `std::span` (F.24), and `std::jthread` (CP.25); with an earlier standard, recommend GSL equivalents or T.48 |
| Exceptions disabled (`-fno-exceptions`, or MSVC without an `/EH` model) | Apply E.25–E.27 instead of E.2 and E.3 |
| RTTI disabled (`-fno-rtti`, MSVC `/GR-`) | `dynamic_cast` (C.146–C.148) is unavailable; review the alternative dispatch (C.153) |
| No heap, real-time, or freestanding target | Weight R.5 and Per.15; allocation on a critical path is a defect |
| GSL dependency present | Suggest `owner`, `not_null`, `span`, `narrow`, `finally`, and `Expects` only when available or adoptable |
| Target threading model | Single-threaded targets lower CP severity; interrupts, RTOS tasks, and multicore targets raise it |
| Data model (LP64, LLP64, 32-bit) | Changes integer widths for ES.46 and ES.100–ES.107; `long` is 32-bit on LLP64 Windows |
| Warnings as errors or a blocking analyzer gate | Report findings that will fail CI first |

## 4. Choose checks

Prefer the project's configured command over an ad hoc invocation: it carries target flags and the team's rule selection. Choose the smallest set that can confirm or refute risks found in the change.

| Situation | Plan |
|---|---|
| `compile_commands.json` exists and `clang-tidy` is available | `clang-tidy -p <build-dir> <changed files>`; it reads the project's `.clang-tidy`. Without one, a supplemental `--checks=-*,cppcoreguidelines-*,bugprone-*` run is allowed; label it ad hoc |
| No compilation database, native target, and `clang-tidy` is available | `clang-tidy <changed files> -- -std=<standard> -I<include dirs> -D<defines>`, with flags read from the build files; label it ad hoc, because missing flags can hide or invent warnings |
| Cppcheck is available and a compilation database exists | `cppcheck --project=<build-dir>/compile_commands.json --file-filter=<changed path>` |
| MSVC project with Code Analysis configured | Use the configured build; Core Check needs the MSVC toolchain on Windows. A build is subject to section 5 |
| Native target with an existing sanitizer configuration and tests | Run the affected tests from that configuration |
| Cross or unknown target | Run static analysis with target flags when possible; name the CI job or device where runtime checks belong |
| No usable tool or build data | Review manually against the rules; list checks that would confirm each finding |

Header diagnostics appear only for headers matching `HeaderFilterRegex` (or `--header-filter`). On Windows, clang-tidy can report paths with mixed separators such as `include\lib/buffer.h`. A filter written with `/` can silently drop header warnings. Pass a separator-agnostic `--header-filter` (for example `.*lib.*` or `include[\\/]lib[\\/].*`) and state this in the strategy.

Sanitizer support depends on toolchain and platform: MSVC provides `/fsanitize=address`, `kernel-address`, and `fuzzer`; Clang documents ThreadSanitizer for 64-bit platforms and MemorySanitizer for Linux, NetBSD, and FreeBSD. Check the tool's documentation (`tools.md`) before planning a runtime check for a target.

## 5. Decision boundaries

| Action | Rule |
|---|---|
| Read files; inspect configuration; query tool versions | Proceed |
| Run a read-only analyzer against an existing build tree, or with flags read from the build files, with output to the console or a temporary directory | Proceed; record the exact command and a summary of results |
| Run existing tests from a build tree that existed before the review began, when the command is documented and needs no network, credentials, or devices | Proceed; record the command and results |
| Configure or populate a new build tree (including `cmake --preset`, `cmake -S ... -B ...`, and subsequent build or test commands), run a full or long build, install tools, contact servers, use devices, or apply fixes (`-fix`) | Ask the user first. Otherwise, in non-interactive reviews, do not run it. List it under "Checks not run" with the reason and the CI job or environment where it belongs. Temporary and git-ignored directories are still new build trees |
| Edit tracked files, commit, or push | Outside this skill's scope |

## 6. Strategy statement

State the strategy in the review before the findings:

- Affected target(s) and build system, with evidence.
- Host versus target classification, with evidence.
- Build facts that changed rule selection.
- Checks run (command and result summary) and checks not run (reason, and where they belong).
