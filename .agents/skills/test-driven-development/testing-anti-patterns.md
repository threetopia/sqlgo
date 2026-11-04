# Testing Anti-Patterns

**Load this reference when:** writing or changing tests, adding mocks, or tempted to add test-only methods to production code.

## Overview

Tests must verify real behavior, not mock behavior. Mocks are a means to isolate, not the thing being tested.

**Core principle:** Test what the code does, not what the mocks do.

**Following strict TDD prevents these anti-patterns.**

## The Iron Laws

```
1. NEVER test mock behavior
2. NEVER add test-only methods to production types/structs
3. NEVER mock without understanding dependencies
```

## Anti-Pattern 1: Testing Mock Behavior

**The violation:**
```go
// ❌ BAD: Test hanya memverifikasi marker dari mock
package page_test

import (
  "strings"
  "testing"
)

type Sidebar interface{ Render() string }

type Page struct{ Sidebar Sidebar }

func (p Page) Render() string {
  return "<main>" + p.Sidebar.Render() + "</main>"
}

type sidebarMock struct{}

func (sidebarMock) Render() string {
  return `<div data-testid="sidebar-mock"></div>`
}

func TestPageRender_Bad(t *testing.T) {
  html := Page{Sidebar: sidebarMock{}}.Render()
  if !strings.Contains(html, "sidebar-mock") {
    t.Fatal("expected sidebar mock marker")
  }
}
```

**Why this is wrong:**
- You're verifying the mock works, not that the component works
- Test passes when mock is present, fails when it's not
- Tells you nothing about real behavior

**your human partner's correction:** "Are we testing the behavior of a mock?"

**The fix:**
```go
// ✅ GOOD: Uji perilaku nyata komponen
package page_test

import (
  "strings"
  "testing"
)

type Sidebar interface{ Render() string }

type Page struct{ Sidebar Sidebar }

func (p Page) Render() string {
  return "<main>" + p.Sidebar.Render() + "</main>"
}

type realSidebar struct{}

func (realSidebar) Render() string {
  return `<nav aria-label="Main Navigation"></nav>`
}

func TestPageRender_Good(t *testing.T) {
  html := Page{Sidebar: realSidebar{}}.Render()
  if !strings.Contains(html, "<nav") {
    t.Fatal("expected navigation element")
  }
}

// Jika mock wajib untuk isolasi, jangan assert marker mock.
// Tetap uji perilaku Page yang relevan bagi user.
```

### Gate Function

```
BEFORE asserting on any mock element:
  Ask: "Am I testing real component behavior or just mock existence?"

  IF testing mock existence:
    STOP - Delete the assertion or unmock the component

  Test real behavior instead
```

## Anti-Pattern 2: Test-Only Methods in Production

**The violation:**
```go
// ❌ BAD: Destroy hanya dipakai test tapi ditaruh di production type
package session

type WorkspaceManager interface {
	DestroyWorkspace(id string) error
}

type Session struct {
	id string
	wm WorkspaceManager
}

// Terlihat seperti API production, padahal hanya kebutuhan test.
func (s Session) Destroy() error {
	return s.wm.DestroyWorkspace(s.id)
}

// In tests:
// t.Cleanup(func() { _ = session.Destroy() })
```

**Why this is wrong:**
- Production type polluted with test-only code
- Dangerous if accidentally called in production
- Violates YAGNI and separation of concerns
- Confuses object lifecycle with entity lifecycle

**The fix:**
```go
// ✅ GOOD: Cleanup khusus test dipindah ke test utility
package session_test

import "testing"

type Session struct {
  WorkspaceID string
}

type WorkspaceManager interface {
  DestroyWorkspace(id string) error
}

func cleanupSession(t *testing.T, wm WorkspaceManager, s Session) {
  t.Helper()
  if s.WorkspaceID == "" {
    return
  }
  if err := wm.DestroyWorkspace(s.WorkspaceID); err != nil {
    t.Fatalf("cleanup failed: %v", err)
  }
}

// In tests:
// t.Cleanup(func() { cleanupSession(t, workspaceManager, session) })
```

### Gate Function

```
BEFORE adding any method to a production type:
  Ask: "Is this only used by tests?"

  IF yes:
    STOP - Don't add it
    Put it in test utilities instead

  Ask: "Does this type own this resource's lifecycle?"

  IF no:
    STOP - Wrong type for this method
```

## Anti-Pattern 3: Mocking Without Understanding

**The violation:**
```go
// ❌ BAD: Stub level tinggi menghilangkan side effect penting
package server_test

import (
  "context"
  "testing"
)

type Config struct{ Name string }

type ToolCatalog interface {
  DiscoverAndCacheTools(ctx context.Context, cfg Config) error
}

type catalogStub struct{}

func (catalogStub) DiscoverAndCacheTools(context.Context, Config) error {
  // Tidak menulis config/cached state yang dibutuhkan deteksi duplicate
  return nil
}

func TestAddServerDetectsDuplicate_Bad(t *testing.T) {
  ctx := context.Background()
  cfg := Config{Name: "alpha"}
  cat := catalogStub{}

  _ = addServer(ctx, cat, cfg)
  err := addServer(ctx, cat, cfg)
  if err == nil {
    t.Fatal("expected duplicate error")
  }
}
```

**Why this is wrong:**
- Mocked method had side effect test depended on (writing config)
- Over-mocking to "be safe" breaks actual behavior
- Test passes for wrong reason or fails mysteriously

**The fix:**
```go
// ✅ GOOD: Mock di level operasi lambat, bukan di level logika inti
package server_test

import (
  "context"
  "errors"
  "testing"
)

type Config struct{ Name string }

type Starter interface {
  Start(context.Context, Config) error
}

type fastStarterStub struct{}

func (fastStarterStub) Start(context.Context, Config) error { return nil }

func TestAddServerDetectsDuplicate_Good(t *testing.T) {
  ctx := context.Background()
  cfg := Config{Name: "alpha"}

  err := addServerWithStarter(ctx, fastStarterStub{}, cfg)
  if err != nil {
    t.Fatalf("first add failed: %v", err)
  }

  err = addServerWithStarter(ctx, fastStarterStub{}, cfg)
  if !errors.Is(err, ErrDuplicateServer) {
    t.Fatalf("expected ErrDuplicateServer, got %v", err)
  }
}
```

### Gate Function

```
BEFORE mocking any method:
  STOP - Don't mock yet

  1. Ask: "What side effects does the real method have?"
  2. Ask: "Does this test depend on any of those side effects?"
  3. Ask: "Do I fully understand what this test needs?"

  IF depends on side effects:
    Mock at lower level (the actual slow/external operation)
    OR use test doubles that preserve necessary behavior
    NOT the high-level method the test depends on

  IF unsure what test depends on:
    Run test with real implementation FIRST
    Observe what actually needs to happen
    THEN add minimal mocking at the right level

  Red flags:
    - "I'll mock this to be safe"
    - "This might be slow, better mock it"
    - Mocking without understanding the dependency chain
```

## Anti-Pattern 4: Incomplete Mocks

**The violation:**
```go
// ❌ BAD: Mock parsial - hanya field yang terlihat dipakai saat ini
package api_test

type Metadata struct {
  RequestID string
  Timestamp int64
}

type UserData struct {
  UserID string
  Name   string
}

type APIResponse struct {
  Status   string
  Data     UserData
  Metadata *Metadata
}

func badMockResponse() APIResponse {
  return APIResponse{
    Status: "success",
    Data:   UserData{UserID: "123", Name: "Alice"},
    // Missing: Metadata yang dipakai downstream
  }
}

// Nanti gagal saat code akses response.Metadata.RequestID
```

**Why this is wrong:**
- **Partial mocks hide structural assumptions** - You only mocked fields you know about
- **Downstream code may depend on fields you didn't include** - Silent failures
- **Tests pass but integration fails** - Mock incomplete, real API complete
- **False confidence** - Test proves nothing about real behavior

**The Iron Rule:** Mock the COMPLETE data structure as it exists in reality, not just fields your immediate test uses.

**The fix:**
```go
// ✅ GOOD: Mock mencerminkan struktur response asli secara lengkap
package api_test

type Metadata struct {
  RequestID string
  Timestamp int64
}

type UserData struct {
  UserID string
  Name   string
}

type APIResponse struct {
  Status   string
  Data     UserData
  Metadata *Metadata
}

func goodMockResponse() APIResponse {
  return APIResponse{
    Status: "success",
    Data:   UserData{UserID: "123", Name: "Alice"},
    Metadata: &Metadata{
      RequestID: "req-789",
      Timestamp: 1234567890,
    },
  }
}
```

### Gate Function

```
BEFORE creating mock responses:
  Check: "What fields does the real API response contain?"

  Actions:
    1. Examine actual API response from docs/examples
    2. Include ALL fields system might consume downstream
    3. Verify mock matches real response schema completely

  Critical:
    If you're creating a mock, you must understand the ENTIRE structure
    Partial mocks fail silently when code depends on omitted fields

  If uncertain: Include all documented fields
```

## Anti-Pattern 5: Integration Tests as Afterthought

**The violation:**
```
✅ Implementation complete
❌ No tests written
"Ready for testing"
```

**Why this is wrong:**
- Testing is part of implementation, not optional follow-up
- TDD would have caught this
- Can't claim complete without tests

**The fix:**
```
TDD cycle:
1. Write failing test
2. Implement to pass
3. Refactor
4. THEN claim complete
```

## When Mocks Become Too Complex

**Warning signs:**
- Mock setup longer than test logic
- Mocking everything to make test pass
- Mocks missing methods real components have
- Test breaks when mock changes

**your human partner's question:** "Do we need to be using a mock here?"

**Consider:** Integration tests with real components often simpler than complex mocks

## TDD Prevents These Anti-Patterns

**Why TDD helps:**
1. **Write test first** → Forces you to think about what you're actually testing
2. **Watch it fail** → Confirms test tests real behavior, not mocks
3. **Minimal implementation** → No test-only methods creep in
4. **Real dependencies** → You see what the test actually needs before mocking

**If you're testing mock behavior, you violated TDD** - you added mocks without watching test fail against real code first.

## Quick Reference

| Anti-Pattern | Fix |
|--------------|-----|
| Assert on mock elements | Test real component or unmock it |
| Test-only methods in production | Move to test utilities |
| Mock without understanding | Understand dependencies first, mock minimally |
| Incomplete mocks | Mirror real API completely |
| Tests as afterthought | TDD - tests first |
| Over-complex mocks | Consider integration tests |

## Red Flags

- Assertion checks for `*-mock` test IDs
- Methods on production types only called in test files
- Mock setup is >50% of test
- Test fails when you remove mock
- Can't explain why mock is needed
- Mocking "just to be safe"

## The Bottom Line

**Mocks are tools to isolate, not things to test.**

If TDD reveals you're testing mock behavior, you've gone wrong.

Fix: Test real behavior or question why you're mocking at all.
