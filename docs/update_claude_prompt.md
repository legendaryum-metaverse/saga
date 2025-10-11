# Update CLAUDE.md Prompt

**Purpose**: This document provides instructions for updating `CLAUDE.md` - the technical documentation that provides context to Claude Code in future sessions.

---

## Core Principle

CLAUDE.md should be **concise yet comprehensive** - a technical reference document, not a tutorial or detailed implementation guide.

**Target**: 350-400 lines (measured with `wc -l CLAUDE.md`)

---

## What to Include

### 1. Essential Context

- **What**: Project purpose and key capabilities
- **Why**: Architectural decisions and design rationale
- **How**: Quick reference patterns for common tasks

### 2. Critical Information

- Current metadata (date, Go version, dependencies)
- System architecture (high-level flow diagrams)
- Core components with file locations
- Key design patterns and their tradeoffs
- Developer guidelines for common modifications

### 3. Quick Reference Section

Place at the top for immediate access:

- Initialize connection
- Publish events
- Consume events
- Handle saga steps
- Error handling patterns

---

## What NOT to Include

### ❌ Avoid These

- **Verbose code walkthroughs**: Show pattern, not every line
- **Implementation details**: Those belong in code comments
- **Duplicate examples**: One clear example per pattern
- **Line-by-line explanations**: Trust the reader's engineering background
- **Complete file contents**: Reference key functions/structs only
- **Historical information**: Focus on current state
- **Granular metrics**: High-level stats only

### ❌ Red Flags

- Sections longer than 100 lines (except for comprehensive tables)
- Repeated code examples with minor variations
- Step-by-step tutorials (use Quick Reference instead)
- Detailed change logs (use git history)
- Implementation pseudo-code (reference actual files)

---

## Document Structure

```
1. Quick Reference (40-50 lines)
   - Essential patterns for immediate use
   - No explanations, just working code snippets

2. Project Overview (20-30 lines)
   - What the library does
   - Key features (bullet points)
   - Tech stack

3. Architecture (50-60 lines)
   - System flow diagram (ASCII art)
   - Core components table
   - Directory structure overview

4. Feature Sections (50-80 lines each)
   - One section per major feature
   - Overview → Key Concepts → Usage Example
   - Tables for structured data

5. Developer Guidelines (50-60 lines)
   - Adding new events
   - Adding new microservices
   - Common modification patterns

6. Design Decisions (20-30 lines)
   - Why X over Y (comparison tables)
   - Architectural rationale

7. Summary (20-30 lines)
   - Production-ready checklist
   - Key statistics
   - Critical concepts to remember
```

---

## Update Process

### Step 1: Analyze Current State

```bash
# Check current line count
wc -l CLAUDE.md

# Identify recent changes
git log --oneline -10

# Check Go version
go version

# List core files
ls -1 *.go event/*.go micro/*.go
```

### Step 2: Verify Accuracy

- Cross-reference mentioned files with actual codebase
- Validate code examples compile
- Check constants/types exist as documented
- Remove references to non-existent code

### Step 3: Update Content

Focus updates on:

1. **Metadata**: Date, Go version, dependency versions
2. **New features**: Add section if >100 LOC changed
3. **Removed features**: Delete outdated sections
4. **Architecture changes**: Update flow diagrams
5. **API changes**: Update Quick Reference

### Step 4: Condense if Needed

If over 400 lines:

- Merge similar examples
- Convert verbose text to tables
- Remove redundant explanations
- Simplify code examples (show concept, not production code)
- Delete low-value sections

---

## Writing Style Guidelines

### Code Examples

```go
// ✅ Good: Shows pattern clearly
emitter.On(event.SomeEvent, func(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.SomePayload{})
    // Process...
    handler.Channel.AckMessage()
})

// ❌ Bad: Too much context, not reusable
func handleAuthNewUserForSocialMicroservice(handler saga.EventHandler) {
    payload := saga.ParsePayload(handler.Payload, &event.AuthNewUserPayload{})

    // Check if user exists in database
    existingUser, err := db.FindUserByEmail(payload.Email)
    if err != nil {
        log.Error("Database error:", err)
        handler.Channel.NackWithFibonacciStrategy(19, 3)
        return
    }
    // ... 20 more lines
}
```

### Explanations

```markdown
✅ Good: Concise with key insight
**Direct Exchange** chosen for audit because:

- Single consumer (audit microservice)
- No broadcast needed
- More efficient than headers exchange

❌ Bad: Over-explained
The audit feature uses a direct exchange instead of a headers exchange.
This is because the audit events only need to be consumed by a single
microservice (the audit microservice), and we don't need the broadcast
functionality that comes with a headers exchange. A direct exchange is
more efficient in this use case because it routes messages directly to
the queue using a routing key, whereas a headers exchange has to evaluate
multiple header attributes to determine routing...
```

### Tables vs Prose

```markdown
✅ Good: Scannable table
| Feature | Events | Saga Commands |
|--------------|------------------|------------------|
| Pattern | Pub/Sub | Point-to-point |
| Exchange | matching_exchange| commands_exchange|
| Consumers | Multiple | Single |

❌ Bad: Dense paragraph
Events use a publish/subscribe pattern with the matching_exchange which
is a headers exchange, whereas saga commands use a point-to-point pattern
with the commands_exchange which is a direct exchange. Events can have
multiple consumers across different microservices, but saga commands...
```

---

## Validation Checklist

Before committing updates:

- [ ] Line count: 350-400 lines (`wc -l CLAUDE.md`)
- [ ] All file references exist (`ls -1 <files_mentioned>`)
- [ ] Code examples are minimal and reusable
- [ ] No section exceeds 100 lines
- [ ] Quick Reference is at the top
- [ ] Metadata is current (date, versions)
- [ ] Tables used for structured data
- [ ] No duplicate examples
- [ ] Design decisions are explained
- [ ] Summary includes critical concepts

---

## Examples of Good Updates

### ✅ Adding a New Feature (50-80 lines)

```markdown
## New Feature Name

### Overview

Brief description (2-3 sentences).

### Key Concept

One critical insight or design decision.

### Usage

// Minimal code example (10-15 lines)

### Infrastructure

| Resource | Type | Purpose |
| -------- | ---- | ------- |
| ...      | ...  | ...     |
```

### ✅ Updating Quick Reference

```markdown
## Quick Reference

// Add new pattern without explanation
// 5. New Pattern
result := saga.NewMethod(&saga.NewPayload{...})
```

### ❌ Adding Tutorial Content

```markdown
## Step-by-Step Guide to Implementing X

In this section, we'll walk through how to implement feature X.

First, you'll need to create a new file...
Then, you'll need to define a struct...
After that, you'll implement the interface...
Let's look at each part in detail...
[continues for 150 lines]
```

---

## Common Scenarios

### Scenario 1: Major Feature Added (>200 LOC)

**Action**: Add new feature section (50-80 lines)

- Overview paragraph
- Key design decision
- Minimal usage example
- Infrastructure table if applicable
- Update Quick Reference

### Scenario 2: Minor Bug Fix (<50 LOC)

**Action**: No CLAUDE.md update needed

- Bug fixes don't change architecture or usage

### Scenario 3: API Changed

**Action**: Update affected sections (10-30 lines)

- Update Quick Reference
- Update relevant usage examples
- Add migration note in Developer Guidelines if breaking

### Scenario 4: Dependency Upgraded

**Action**: Update metadata (1-3 lines)

- Update Tech Stack section
- Update metadata header

### Scenario 5: Refactoring (file renamed, reorganized)

**Action**: Update file references (5-15 lines)

- Update Directory Structure
- Update Core Components table
- Verify all file:line references

---

## Anti-Patterns to Avoid

### 🚫 Documentation Drift

Don't document planned features or outdated code:

```markdown
❌ "In the future, we plan to add..."
❌ "The old ConnectToAudit() method (deprecated)..."
✅ Document only what exists in the current codebase
```

### 🚫 Tutorial Creep

Don't turn CLAUDE.md into a tutorial:

```markdown
❌ "Let's walk through how to build a microservice from scratch..."
❌ "Step 1: First, open your editor and create..."
✅ "See Quick Reference for initialization pattern"
```

### 🚫 Over-Specification

Don't duplicate information available in code:

```markdown
❌ Listing all 42 event types with descriptions
❌ Showing every field of every struct
✅ "45 events total (42 business + 3 audit). See event/microserviceEvent.go"
```

### 🚫 Inconsistent Depth

Don't explain some concepts in detail while glossing over others:

```markdown
❌ 150 lines on audit feature, 10 lines on saga orchestration
✅ Balanced coverage (50-80 lines per major feature)
```

---

## Maintenance Frequency

- **After major feature**: Update immediately
- **After breaking change**: Update immediately
- **After bug fix**: Usually no update needed
- **After refactoring**: Update file references if changed
- **Monthly**: Verify accuracy and prune if >400 lines
- **Quarterly**: Refresh metadata and dependency versions

---

## Meta: Updating This Document

When updating `update_claude_prompt.md`:

- Keep under 300 lines
- Balance specificity with brevity
- Include concrete examples (good/bad comparisons)
- Focus on engineering judgment, not rigid rules
- Update when CLAUDE.md update patterns change

---

**Last Updated**: 2025-10-11
**CLAUDE.md Target**: 350-400 lines
**This Document Target**: <300 lines
