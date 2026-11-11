# Review lenses

A lens maps changed code to candidate Guidelines rules. Each lens lists activating signals, defect-exposing questions, candidate rules, guards against reporting correct code, and corroborating evidence.

- Find sections and rules with `toc.md` and the selected document. Read each applied rule in full.
- Rule IDs are a starting set. Follow "See also" links and cross-references in rule text when they bear on the code.
- Checks are optional evidence: `tools.md` describes each tool; `build-discovery.md` decides which can run.

## Router

| The change touches | Lenses |
|---|---|
| `new`/`delete`, `malloc`/`free`, handles, owning raw members, factories | 1, 2, 9 |
| destructors, copy/move operations, constructors, class invariants | 2, 1 |
| `virtual`, inheritance, casts between class pointers | 3 |
| function signatures, parameters, return types, public headers | 4, 5 |
| references, views, iterators, lambdas, callbacks, coroutines | 5 |
| raw arrays, pointer arithmetic, indexes, buffers, `memcpy` | 6, 8 |
| casts, unions, varargs, type punning, initialization | 7 |
| integer math, mixed signedness, shifts, size computations | 8 |
| `throw`/`catch`, `noexcept`, error codes, cleanup paths | 9 |
| threads, atomics, mutexes, condition variables, `async`, coroutines | 10, 5 |
| templates, concepts, forwarding, overloads, customization points | 11 |
| headers, includes, namespaces, macros, globals | 12 |
| containers, strings, streams, C library calls, C interop | 13 |
| `const`, `constexpr`, enumerations | 14 |
| performance-motivated changes or claims | 15 |
| new modules, abstractions, or public APIs | 16 |
| analyzer suppressions (`NOLINT`, `[[gsl::suppress]]`, `#pragma warning(disable: ...)`) or analyzer configuration | 0, plus the lens of each suppressed rule |

Lens 0 applies to every review.

## 0. Calibration

- Rules: NR.1, NR.2, NR.3, NR.4 (demands the Guidelines reject); NR.5, NR.6, NR.7 (patterns to avoid).
- Sections: `s-not`, `s-modernizing`, `s-naming`, `ss-force`, `s-profile`.
- Guards:
  - Acceptable (NR.1–NR.4): declarations away from the top, multiple `return` statements, exceptions, several classes per file.
  - Exception and Note subsections override a rule's title; read them before reporting.
  - Naming and layout (NL): findings only when the project adopts them or the user asks.
  - Legacy code: apply rules to changed lines and their direct dependencies (Appendix B, `s-modernizing`).
- Suppressions: In.force (`ss-force`) defines `[[gsl::suppress("tag", justification: "message")]]`. The tag is a rule anchor (such as `rh-public`), a profile (`type`, `bounds`, `lifetime`), or a profile rule (such as `type.4`). Read the suppressed rule.
  - Finding against that rule: a new suppression without justification, or one that hides a reachable defect.
  - Acceptable: a justified suppression at a necessary boundary.

## 1. Ownership and resource release

- Signals: `new`, `delete`, `malloc`, `free`, `fopen`, OS handles, owning raw pointer members, `release()`, manual cleanup, factory functions.
- Ask: Who releases each resource? Release guaranteed on every path, including early return and exceptions? Ownership visible in the type? Can two owners release one object?
- Rules: P.8, R.1, R.3, R.5, R.10, R.11, R.12, R.13, R.20, R.21, R.23, I.11, C.31, C.33, E.6, E.13, ES.60, ES.61.
- Sections: `s-resource`, `ss-alloc`, `ss-smart`, `ss-ownership`.
- Guards:
  - A non-owning raw `T*` or `T&` is correct (R.3, R.4, F.7, F.22).
  - A class that acquires a resource and releases it in its destructor is a resource handle. Its raw allocation alone is not an R.1 finding; check ownership and copy/move semantics separately.
  - ES.60 allows `new`/`delete` inside resource-management functions.
- Evidence: clang-tidy `cppcoreguidelines-owning-memory`, `cppcoreguidelines-no-malloc`, `clang-analyzer-cplusplus.NewDelete`, `clang-analyzer-cplusplus.NewDeleteLeaks`, `bugprone-shared-ptr-array-mismatch`; MSVC C26400–C26409, C26865; ASan, LSan, Valgrind Memcheck.

## 2. Class design and special members

- Signals: user-declared destructor, copy or move operations; pointer or handle members; constructors that can fail; single-argument constructors; member order; `memset`/`memcpy` on objects.
- Ask: Do copy, move, and destruction preserve the invariant? Rule of zero possible? Every member initialized? Self-assignment safe? Moved-from object still valid?
- Rules: C.20, C.21, C.22, C.30, C.31, C.33, C.35, C.36, C.37, C.41, C.42, C.46, C.47, C.49, C.62, C.64, C.65, C.66, C.82, C.90, C.12.
- Sections: `ss-defop`, `ss-dtor`, `ss-ctor`, `ss-copy`, `ss-concrete`.
- Guards:
  - C.20 prefers no user-declared default operations when members manage themselves.
  - C.21: account for all five operations. A defaulted move of an owning raw pointer copies the address without clearing the source. Check that ownership transfers safely and affected callers compile.
- Evidence: clang-tidy `cppcoreguidelines-special-member-functions`, `cppcoreguidelines-pro-type-member-init`, `cppcoreguidelines-noexcept-destructor`, `cppcoreguidelines-noexcept-move-operations`, `cppcoreguidelines-noexcept-swap`, `cppcoreguidelines-prefer-member-initializer`, `cppcoreguidelines-avoid-const-or-ref-data-members`, `bugprone-unhandled-self-assignment`, `bugprone-copy-constructor-init`; MSVC C26432, C26439, C26495; ASan for double free after copy.

## 3. Class hierarchies and polymorphism

- Signals: `virtual`, `override`, `final`, base classes, casts between class pointers, polymorphic objects passed or stored by value, virtual calls in constructors or destructors, default arguments on virtual functions.
- Ask: Can a derived object be destroyed through a base pointer, or sliced? Does each override actually override? Every downcast checked?
- Rules: C.35, C.67, C.82, C.121, C.127, C.128, C.130, C.133, C.140, C.145, C.146, C.147, C.148, C.152, C.153, ES.63, Type.2.
- Sections: `ss-hier`.
- Guards: A protected non-virtual destructor satisfies C.35 when deletion through the base is not intended.
- Evidence: clang-tidy `cppcoreguidelines-virtual-class-destructor`, `cppcoreguidelines-slicing`, `cppcoreguidelines-explicit-virtual-functions`, `cppcoreguidelines-pro-type-static-cast-downcast`, `bugprone-virtual-near-miss`, `bugprone-parent-virtual-call`, `bugprone-pointer-arithmetic-on-polymorphic-object`; MSVC C26433–C26437, C26443, C26466, C26491.

## 4. Interfaces, parameters, and contracts

- Signals: new or changed signatures, public headers, parameter passing modes, output parameters, smart pointer parameters, preconditions, null handling, global state, exported or ABI-visible types.
- Ask: Does the signature state ownership, nullability, size, and direction? Can callers swap adjacent arguments? Preconditions stated or checked? Source or binary compatibility broken for existing callers?
- Rules: I.1, I.2, I.4, I.5, I.7, I.10, I.11, I.12, I.13, I.23, I.24, I.26, I.27, F.7, F.15, F.16, F.17, F.18, F.19, F.20, F.21, F.22, F.23, F.24, F.25, F.26, F.60, R.30, R.32, R.33, R.34, R.36, R.37.
- Sections: `s-interfaces`, `ss-call`, `ss-smart`.
- Guards:
  - F.16 size thresholds are heuristics. Pass-by-value versus `const&` is Low unless it changes behavior or has measured cost (Per.6).
  - Raw pointers and references are the default for non-owning parameters (F.7).
- Evidence: clang-tidy `cppcoreguidelines-missing-std-forward`, `cppcoreguidelines-rvalue-reference-param-not-moved`, `performance-unnecessary-value-param`, `cppcoreguidelines-avoid-non-const-global-variables`, `cppcoreguidelines-interfaces-global-init`; MSVC C26410, C26411, C26415–C26418, C26426, C26427, C26429–C26431, C26460, C26461.

## 5. Lifetime, dangling, and aliasing

- Signals: returned references, pointers, or views; `string_view` or `span` members; iterators kept across container mutation; reference-capturing lambdas that escape (stored, returned, or run on another thread); coroutine reference parameters; use after `std::move`; kept `get()`, `data()`, or `c_str()` results.
- Ask: For each reference-like value, which object does it refer to, and does that object outlive every use? Can reallocation, a move, or scope exit invalidate it first?
- Rules: F.24, F.42, F.43, F.45, F.52, F.53, F.54, ES.56, ES.65, Lifetime.1, R.37, C.64, CP.24, CP.26, CP.31, CP.51, CP.53, SL.str.2.
- Sections: `s-interfaces`, `ss-call`, `ss-lifetime`, `ss-views`.
- Guards:
  - Reference capture is preferred for locally used lambdas, including algorithm arguments (F.52).
  - `string_view` and `span` parameters are correct (SL.str.2, F.24); the risk is storing or returning them.
- Borrowed returns: for a function returning a pointer, reference, iterator, or view into a non-owning parameter, check whether a caller can pass a temporary owner.
  - In C++20, `span<const T>` can be constructed from a temporary container; the owner still dies at the end of the full-expression.
  - Read F.42 and F.43 for the escaped position and lifetime rule; F.24 establishes that the span does not own the elements.
  - Verify the API keeps the borrowed result from outliving that owner, or report the concrete dangling path.
- Evidence: clang-tidy `bugprone-dangling-handle`, `bugprone-use-after-move`, `clang-analyzer-cplusplus.Move`, `bugprone-return-const-ref-from-parameter`, `cppcoreguidelines-avoid-capturing-lambda-coroutines`, `cppcoreguidelines-avoid-reference-coroutine-parameters`; MSVC C26445, C26449, C26800, C26810, C26811, C26815, C26816; ASan.

## 6. Bounds and memory access

- Signals: C arrays, pointer arithmetic, computed indexes, pointer-plus-length pairs, `memcpy`, `strcpy`, `sprintf`, array parameters, hand-written loop bounds.
- Ask: What keeps each index or pointer in range? Does the size travel with the pointer?
- Rules: I.13, F.24, R.14, ES.27, ES.42, ES.55, ES.62, SL.con.1, SL.con.3, SL.con.4, Bounds.1, Bounds.2, Bounds.3, Bounds.4, C.152.
- Sections: `ss-bounds`, `ss-con`, `ss-views`.
- Guards: The Bounds profile is strict. Unless the project enforces it, report a defect only when an access can leave its range or an interface loses the size; otherwise report a Low guideline risk.
- Evidence: clang-tidy `cppcoreguidelines-pro-bounds-pointer-arithmetic`, `cppcoreguidelines-pro-bounds-constant-array-index`, `cppcoreguidelines-pro-bounds-array-to-pointer-decay`, `cppcoreguidelines-pro-bounds-avoid-unchecked-container-access`, `cppcoreguidelines-avoid-c-arrays`, `bugprone-sizeof-expression`, `bugprone-not-null-terminated-result`, `bugprone-unsafe-functions`; MSVC C26446, C26481–C26485; ASan; fuzzing for parsers.

## 7. Types, casts, and initialization

- Signals: `reinterpret_cast`, C-style or functional casts, `const_cast`, unions, `va_arg`, `void*`, uninitialized locals or members, type punning.
- Ask: Every object initialized before it is read? Does each cast preserve type and value? Can memory be accessed as a type it does not hold?
- Rules: ES.20, ES.23, ES.48, ES.49, ES.50, ES.34, F.55, C.41, C.181, C.183, Type.1, Type.3, Type.4, Type.5, Type.6, Type.7, Type.8, I.30.
- Sections: `ss-type`, `ss-union`.
- Guards:
  - `{}` initialization (ES.23) is a preference.
  - Hardware, serialization, and C API boundaries can need casts. Report a cast that is unchecked, loses information, or has a safer alternative; prefer encapsulated violations (I.30).
- Evidence: clang-tidy `cppcoreguidelines-pro-type-reinterpret-cast`, `cppcoreguidelines-pro-type-cstyle-cast`, `cppcoreguidelines-pro-type-const-cast`, `cppcoreguidelines-pro-type-union-access`, `cppcoreguidelines-pro-type-vararg`, `cppcoreguidelines-pro-type-member-init`, `cppcoreguidelines-init-variables`, `bugprone-casting-through-void`, `bugprone-bitwise-pointer-cast`, `bugprone-undefined-memory-manipulation`; MSVC C26465, C26466, C26471–C26475, C26490–C26495; UBSan, MSan, Valgrind Memcheck.

## 8. Arithmetic and evaluation order

- Signals: mixed signed and unsigned values, narrowing, allocation size computations, shifts, division, accumulations that can overflow, unsigned countdown loops, expressions with side effects on shared operands.
- Ask: Can any operation overflow, underflow, divide by zero, shift out of range, or narrow? Does the result depend on evaluation order?
- Rules: ES.40, ES.41, ES.43, ES.44, ES.46, ES.100, ES.101, ES.102, ES.103, ES.104, ES.105, ES.106, ES.107.
- Sections: `ss-numbers`, `s-expr`.
- Guards: Magic constants (ES.45) and `gsl::index` (ES.107) are maintainability concerns unless they cause a wrong result.
- Evidence: clang-tidy `cppcoreguidelines-narrowing-conversions`, `bugprone-implicit-widening-of-multiplication-result`, `bugprone-misplaced-widening-cast`, `bugprone-integer-division`, `bugprone-signed-char-misuse`; MSVC C26472, C26831–C26833, C26838, C26839; compiler conversion warnings; UBSan.

## 9. Error handling and exception safety

- Signals: `throw`, `try`/`catch`, `noexcept`, destructors that do work, error codes, `errno`, cleanup on error paths, partial updates before a call that can fail, `catch (...)`.
- Ask: If a call fails or throws, are resources released and invariants kept? Can a `noexcept` function or destructor throw and terminate the program? Any error result ignored?
- Rules: E.2, E.3, E.5, E.6, E.12, E.13, E.14, E.15, E.16, E.17, E.18, E.19, E.25, E.26, E.27, E.28, E.31, C.36, C.37, C.42, F.6, I.10.
- Sections: `s-errors`.
- Guards: Build disables exceptions: apply E.25–E.27 instead of E.2 and E.3 (`build-discovery.md`).
- Evidence: clang-tidy `bugprone-exception-escape`, `misc-throw-by-value-catch-by-reference`, `bugprone-throw-keyword-missing`, `bugprone-unhandled-exception-at-new`, `cppcoreguidelines-noexcept-destructor`; MSVC C26439, C26440, C26447, C26448; tests that inject failures.

## 10. Concurrency, parallelism, and coroutines

- Signals: `std::thread`, `std::jthread`, `std::async`, executors, atomics, mutexes, condition variables, `volatile`, shared mutable state, function-local statics, callbacks invoked under a lock, `co_await`.
- Ask:
  - For each shared writable object: which threads access it, and what synchronizes every access?
  - Lock order consistent? Can a wait miss a notification or wake spuriously?
  - Does a thread or coroutine outlive what it references?
  - Does each lock guard live through the entire critical section? A standalone temporary `std::lock_guard{mutex}` dies at the end of that full-expression, before the next statement; check CP.44/ES.84 rather than assume it protects following code.
- Rules: CP.1, CP.2, CP.3, CP.8, CP.20, CP.21, CP.22, CP.23, CP.24, CP.25, CP.26, CP.31, CP.32, CP.42, CP.43, CP.44, CP.50, CP.51, CP.52, CP.53, CP.100, CP.110, CP.200, ES.84.
- Sections: `s-concurrency`, `sscp-con`, `sscp-coro`, `sscp-free`.
- Guards: CP.1 assumes code may run in a multi-threaded program. A data race finding still needs a concurrent path, such as two threads, a callback from another thread, or an executor.
- Evidence: clang-tidy `concurrency-mt-unsafe`, `bugprone-spuriously-wake-up-functions`, `bugprone-unused-raii`, `cppcoreguidelines-no-suspend-with-lock`, `cppcoreguidelines-avoid-capturing-lambda-coroutines`; Clang `-Wthread-safety` with capability annotations; MSVC C26100–C26167, C26441; TSan.

## 11. Templates and generic code

- Signals: new templates, concepts, `requires`, `enable_if`, forwarding references, variadic templates, specializations, unqualified calls inside templates.
- Ask: Do stated requirements reject misuse at the call site? Each forwarding reference forwarded exactly once? Can argument-dependent lookup select an unintended overload?
- Rules: T.10, T.11, T.41, T.47, T.48, T.61, T.62, T.68, T.69, T.83, T.84, T.120, T.143, T.144, T.150, F.19.
- Sections: `s-templates`, `ss-concepts`, `ss-temp-interface`, `ss-temp-def`.
- Guards: Concepts require C++20; with an earlier standard, apply T.48. Most template rules are design guidance and rarely exceed Medium.
- Evidence: clang-tidy `cppcoreguidelines-missing-std-forward`, `bugprone-move-forwarding-reference`; instantiation in tests.

## 12. Source files, headers, macros, and globals

- Signals: header edits, new includes, `using namespace`, unnamed namespaces, non-inline definitions in headers, include guards, macros, mutable globals, cross-file initialization order.
- Ask: Each header self-contained and safe under the one-definition rule? Can a macro change meaning at a use site? Does global initialization depend on another translation unit?
- Rules: SF.2, SF.7, SF.8, SF.9, SF.10, SF.11, SF.21, SF.22, ES.30, ES.31, ES.33, I.2, I.3, I.22, R.6.
- Sections: `s-source`.
- Guards: File suffixes (SF.1) follow the project's existing convention.
- Evidence: clang-tidy `misc-definitions-in-headers`, `google-build-using-namespace`, `misc-include-cleaner`, `cppcoreguidelines-macro-usage`, `bugprone-macro-parentheses`, `cppcoreguidelines-interfaces-global-init`, `cert-err58-cpp`; MSVC C26426, C26427; linking every affected configuration.

## 13. Standard library and C interop

- Signals: containers, iterators, strings, `string_view`, iostreams, `<regex>`, `<chrono>`, `strcpy`, `sprintf`, `memset`, `setjmp`, `extern "C"`, `errno`, C APIs that take buffers.
- Ask: Do C calls receive correct sizes and terminators? Iterators and references valid after container mutation? Would a standard facility express the intent more safely?
- Rules: SL.1, SL.2, SL.3, SL.4, SL.con.1, SL.con.2, SL.con.3, SL.con.4, SL.str.1, SL.str.2, SL.str.3, SL.str.4, SL.str.5, SL.io.2, SL.C.1, CPL.1, CPL.2, CPL.3, F.25, E.28, R.10.
- Sections: `s-stdlib`, `ss-con`, `ss-string`, `ss-io`, `ss-clib`, `s-cpl`.
- Guards: `endl` (SL.io.50) and `sync_with_stdio` (SL.io.10) are performance suggestions.
- Evidence: clang-tidy `bugprone-unsafe-functions`, `bugprone-not-null-terminated-result`, `bugprone-string-constructor`, `bugprone-stringview-nullptr`, `bugprone-suspicious-stringview-data-usage`, `bugprone-unchecked-optional-access`, `bugprone-signal-handler`; MSVC C26829, C26830, C26859, C26860, C26865; ASan; fuzzing.

## 14. Constants, immutability, and enumerations

- Signals: mutable values that never change, missing `const` on members or parameters, values computable at compile time, plain enums, macros used as constants.
- Ask: Can unintended mutation or implicit enum conversion produce a wrong result?
- Rules: Con.1, Con.2, Con.3, Con.4, Con.5, ES.25, Enum.1, Enum.2, Enum.3, Enum.6.
- Sections: `s-const`, `s-enum`.
- Guards: Usually Low maintainability findings. Raise them when a missing `const` permits an actual unintended mutation or a plain enum silently converts into a wrong value.
- Evidence: clang-tidy `cppcoreguidelines-use-enum-class`, `cppcoreguidelines-macro-to-enum`, `misc-misplaced-const`; MSVC C26460–C26464, C26496–C26498, C26812, C26814.

## 15. Performance

- Signals: changes justified by speed or memory, allocations in hot paths, copies of large objects, `std::move` usage, data layout changes, added indirection.
- Rules: Per.1, Per.2, Per.3, Per.4, Per.5, Per.6, Per.7, Per.14, Per.15, Per.19, F.16, F.48, ES.56, C.66.
- Sections: `s-performance`.
- Guards: Report a performance problem when a measurement, a complexity change, or an avoidable cost on a known hot path backs it (Per.6). Phrase other concerns as questions.
- Evidence: clang-tidy `performance-unnecessary-value-param`, `performance-move-const-arg`; MSVC C26817, C26820; the project's benchmarks or profiles.

## 16. Design and philosophy

- Signals: new modules, abstractions, public APIs, frameworks, large refactors.
- Rules: P.1, P.3, P.4, P.5, P.7, P.8, P.10, P.11, I.25, I.27, C.10, C.129.
- Sections: `s-philosophy`, `s-a`.
- Guards: Use these to frame design opinions. A blocking finding still needs a concrete rule from lenses 1–15 and a consequence.
