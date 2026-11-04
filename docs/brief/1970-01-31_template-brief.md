# Feature Brief Template

## General Information
- Title: <feature / initiative name>
- Date: 2026-08-30
- Status: Draft / In Progress / Ready for Plan
- Owner: <name or team>
- Related Plan: <optional link or filename>

## Objective
Describe the main goal of this feature in one or two sentences.

Examples:
- Protect student privacy by encrypting student names while keeping other user roles unaffected.
- Add a new data export flow for teacher dashboard reporting without changing the existing api_user contract.

## Problem / Background
Explain the current problem, pain point, or requirement that triggered this feature.

Examples:
- Student names are stored in plain text in api_user.users.
- The dashboard needs to show data without exposing personally identifiable information.
- Existing architecture expects one source of truth for user data, but the current implementation duplicates logic.

## Scope of Work
- [ ] In scope: <areas to be changed>
- [ ] Out of scope: <areas that will not be changed>
- [ ] Related modules / apps: <api_user, api_core, entity, dao, adapter, etc.>

## Requirements
List what must be true for this feature to be considered successful.

Examples:
- Only users with role STUDENT are anonymized/encrypted.
- Non-STUDENT roles remain unchanged.
- Existing business flows continue to work without breaking the api_user contract.
- The solution must respect the repo architecture and dependency boundaries.

## Functional Details
- Describe the user-facing or system-facing behavior.
- Clarify the business rule, data rule, or security rule.
- Specify any edge cases or constraints.

Examples:
- Encrypt student full_name at persistence time only when role = STUDENT.
- Keep teacher/admin names readable and unchanged.
- Use the existing encryption helper or entity wrapper if already available.
- Do not introduce cross-app direct service dependencies.

## Acceptance Criteria
List the expected outcomes in clear, testable form.

Examples:
- Student names are stored encrypted in the database for STUDENT users only.
- Teacher and admin names remain plain text.
- Read and write flows for student data still work as expected.
- No unrelated app modules are changed.

## Relevant Files / Modules
- api_user
- api_core
- internal/app/.../service
- pkg/entity
- outbound adapter / model contract
- migration / SQL / schema
- other related components

## Risks / Unknowns
- Risk 1: <uncertain dependency or data contract>
- Risk 2: <migration or legacy data issue>
- Risk 3: <hidden side effect on another app or API>

## Verification
Describe how this feature will be validated.

Examples:
- Run relevant Go build / tests for the affected app.
- Verify student role data is encrypted only in the target field.
- Confirm teacher/admin records remain unaffected.
- Validate API behavior and UI output with sample data.

## Final Outcome
Describe the expected result once the plan is implemented.

Examples:
- Student personal names are protected while preserving the app workflow.
- The feature follows the repo architecture and keeps user ownership clear.
- The feature is ready to be converted into a concrete implementation plan in docs/plan.

## Prompt for AI Planning
Use this brief to generate a concrete implementation plan in docs/plan.

Required output:
- Objective
- Scope
- Current condition
- Detailed implementation steps
- Risks / notes
- Verification checklist
- Final expected result
