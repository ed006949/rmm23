---
apply: off
---

- Parsing and validation
	- Parse defensively; validate at each step and return early on error.
	- Use switch-with-conditions for guard checks instead of deep if-nesting.
	- Parse into an interim value; only assign to the receiver after full validation.

- API shape
	- Implement standard text interfaces: UnmarshalText, MarshalText, String, Parse.
	- Keep methods small and single-purpose.

- Dependencies and cohesion
	- Prefer internal utility packages (mod_strings, mod_slices, mod_errors) over ad-hoc helpers.
	- Centralize constants (separators, sizes/limits) in shared modules.

- Data structures and allocation
	- Pre-size slices with exact len and cap to avoid reallocations.
	- Use simple structs (e.g., KV) for intermediate representations.

- Naming and readability
	- Use terse, descriptive names reflecting transformation stages (e.g., interim, interimElement).
	- Keep control flow linear and compact.

- Errors
	- Use shared, typed errors (e.g., mod_errors.EParse).
	- Provide helper functions (e.g., StripErr1) where ergonomic.

- String handling
	- Use strings.Join and split helpers from utilities; avoid manual loops for concatenation.
	- Trim and filter empties during tokenization via flags.

- Returns
	- Prefer named result parameters for concise returns.
	- Use zero-value returns; avoid redundant variable initializations.

- Consistency
	- Enforce separators and field counts via constants (e.g., DNPathSeparator, DNKVSeparator, KVElements).
	- Keep methods concise; avoid unnecessary branching and nesting.
