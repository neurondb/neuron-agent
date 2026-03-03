# Contributing to NeuronAgent

Thank you for your interest in contributing.

1. Fork the repository and create a feature branch.
2. Build and test: `make build` and `make test` (run from repo root; Go code is under `src/`).
3. Commit with clear messages (see [Commit Message Guidelines](#commit-message-guidelines)).
4. Submit a pull request with a clear description.

For questions, open an issue or contact the maintainers.

## Commit Message Guidelines

Commit messages must contain relevant information and follow these rules:

**General Rules**

- The subject line must end with a period and should be concise and clear, typically not exceeding a single line.
- Write commit messages in paragraph form rather than as bullet points or lists, making sure to clearly communicate the content of the change.
- Do not reference specific file names or locations within the commit message text.
- Omit any mention of merge or cherry-pick actions, such as "Cherry-picked commit..." or "Merge commit...".
- Exclude any references to code cleanup, coding standards, or style violations (e.g., "C90 violation", "cleanup", or "coding standard changes").
- Leave out references to version compatibility or APIs, such as "PostgreSQL compatibility" or "API changes".
- Avoid including statements about the status of the codebase, such as whether it compiles correctly or if all errors have been resolved.
- Focus exclusively on what has changed in the commit. Do not explain why it was done, and do not comment on compliance with standards or practices.
- The body of the message should be written in clear paragraphs, providing a concise narrative that describes the change, breaking details into additional paragraphs as necessary for clarity.
- Each line contains max 80–90 characters.

**Module Prefix**

Prefix the first line of the commit message with `NeuronAgent:` for this repository.

**Example**

```
NeuronAgent: Improve embedding vector normalization logic.

This commit adjusts the normalization routine to
use a more numerically stable approach, addressing
issues with denormalized input data. Additional
refactoring ensures consistent vector sizing across all
embedding interfaces, providing clearer behavior for
callers and simplifying future maintenance.
```
