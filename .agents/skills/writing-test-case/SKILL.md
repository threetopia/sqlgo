---
name: writing-test-case
description: Use when creating or revising Test Case Documents (TCDs), defining positive and negative test coverage, or documenting acceptance and AI agent evaluation scenarios. Complements writing-specs by defining verifiable test cases and expected results.
---

# Test Case Design

Create and maintain Test Case Documents (TCDs) that define verifiable expected behavior for software features, applications, integrations, and AI agents.

TCD defines what must be tested and the expected outcome. Automated tests implement executable checks where practical. Keep documented expectations, requirements, and automated tests consistent.

Use the user's language for generated TCD documents.

## Scope

This skill applies to creating or revising TCD documents in the repository's `docs/test-cases/` directory.

Use this skill when:
- Creating test cases for a new feature or application.
- Updating test cases when requirements or behavior change.
- Defining positive, negative, boundary, and failure-path coverage.
- Designing acceptance tests or AI agent evaluation scenarios.
- Reviewing test coverage or documenting regression scenarios.

Do not use this skill solely because an existing test must be executed or because implementation code needs to be changed. However, when a change affects expected behavior, assess whether the TCD must also be updated.

## Relationship to Specifications

Use `writing-specs` for creating or revising FSD, TSD, DRD, ERD, DFD, and related specification documents.

TCD complements these documents; it does not replace them.

| Source | Relationship to TCD |
|---|---|
| FSD | Primary source for functional requirements and expected business behavior. |
| TSD | Source for technical flows, interfaces, integrations, and technical failure behavior. |
| DRD | Source for data requirements, constraints, and integrity expectations. |
| ERD / DFD / UML | Supporting references for relationships, data movement, and system interactions. |
| TCD | Defines test scenarios, test data, execution steps, and expected results. |
| Automated tests | Executable verification of the relevant requirements and scenarios. |

Before creating a TCD, inspect relevant specifications and existing test cases.

Reference existing requirement IDs, such as `FR-001`, instead of inventing duplicate requirements. Use applicable identifiers from the source documents.

Do not invent business rules, technical contracts, schemas, or acceptance criteria. Record unresolved decisions as open questions and clearly label assumptions.

If a required specification is missing, continue with test cases that can be defined reliably and identify the missing information. Do not block all progress unnecessarily.

## Location and Filenames

Save TCD documents only in `docs/test-cases/` at the repository root. Create the directory if it does not exist.

Use this exact filename pattern:

`{yyyy-mm-dd}_TCD_{name-of-the-feature-or-app}.md`

Rules:
- Use uppercase `TCD`.
- Use a lowercase, hyphen-separated feature or application slug.
- Separate the date, type, and slug with one literal underscore.
- Use the user's local creation date or the explicitly requested date.
- When a related specification set already exists, reuse its exact date and slug, even if the TCD is created later.
- When creating a new documentation set, use the same date and slug for all related specification and test-case documents.
- When revising an existing TCD, preserve its filename date and slug.
- Record the actual creation or revision date in the document content.
- Do not overwrite an unrelated document with a matching filename.

Example:

```text
docs/specs/2026-10-27_FSD_customer-management.md
docs/specs/2026-10-27_TSD_customer-management.md
docs/specs/2026-10-27_DRD_customer-management.md
docs/test-cases/2026-10-27_TCD_customer-management.md
```

If the TCD is added later, it must retain the existing specification set's date and slug.

## Required Document Structure

Every TCD must include the following sections.

### 1. Document Information

Include:
- Title.
- Document type: TCD.
- Creation date.
- Last updated date.
- Status, such as Draft, In Review, or Approved.
- Scope and exclusions.
- Relative links to related specifications.
- Relevant feature, application, or module identifier.

### 2. Test Objectives and Scope

Describe:
- What behavior is being verified.
- Which features and requirements are covered.
- What is explicitly excluded.
- Relevant dependencies, preconditions, and environment assumptions.

### 3. Requirements Traceability

Map each test case to one or more existing requirement IDs, acceptance criteria, or documented defect IDs.

Identify requirements without sufficient test coverage when practical.

Do not create fictitious requirement IDs. If no formal requirement exists, reference the relevant behavior or record the missing requirement as an open question.

### 4. Test Cases

Each test case must have:

| Field | Description |
|---|---|
| Test Case ID | Unique and stable identifier, e.g. `TC-CUSTOMER-001`. |
| Requirement ID | Related FSD requirement, acceptance criterion, or defect ID. |
| Title | Concise description of the behavior under test. |
| Category | Positive, Negative, Boundary, Security, etc. |
| Priority | Critical, High, Medium, or Low. |
| Preconditions | Initial state, permissions, dependencies, and required setup. |
| Test Data | Inputs, fixtures, and relevant data conditions. |
| Test Steps | Ordered actions required to execute the case. |
| Expected Result | Specific, observable, and verifiable outcome. |
| Actual Result | Recorded outcome after execution. |
| Status | Not Run, Pass, Fail, or Blocked. |
| Evidence / Defect ID | Relevant logs, artifacts, or defect reference. |

Use `Not Run` as the default status for test cases that have not been executed.

Add optional fields when relevant, including environment, expected API status, database assertions, expected tool calls, model version, prompt version, and evaluation criteria.

## Test Coverage Requirements

Consider both positive and negative behavior for every applicable feature.

Do not mechanically create every category for every feature. Include categories according to the requirements, risks, and possible failure modes.

### Positive Test Cases

Verify valid inputs, authorized operations, successful workflows, and expected outputs.

Examples:
- Valid data is accepted.
- Authorized users can perform permitted operations.
- Valid API requests return the documented response.
- Successful transactions persist the correct data.

### Negative Test Cases

Verify that invalid, unauthorized, missing, or unsupported conditions produce the correct behavior.

Examples:
- Required fields are missing.
- Input formats or values are invalid.
- Referenced records do not exist.
- The user lacks the required permission.
- A dependency returns an error.
- Duplicate operations violate a business rule.

Negative test cases must verify both the appropriate error response and any required side-effect behavior, such as ensuring invalid operations do not modify data.

### Boundary and Validation Test Cases

Consider:
- Empty, minimum, maximum, and out-of-range values.
- Length and quantity limits.
- Invalid formats and malformed payloads.
- Nullability and optional fields.
- Duplicate values and uniqueness constraints.
- Pagination, sorting, and filtering limits.

### Error Handling and Resilience

When applicable, verify:
- Database and external API failures.
- Network errors, timeouts, and rate limits.
- Retry and recovery behavior.
- Partial failures and transaction rollback.
- Duplicate requests and idempotency.
- Unexpected or malformed dependency responses.

### Security Test Cases

When applicable, verify:
- Authentication and authorization.
- Access isolation between users or tenants.
- Input injection and untrusted data handling.
- Sensitive data exposure.
- Permission boundaries and prohibited operations.

### Regression Test Cases

For bug fixes and behavior changes, include cases that protect existing behavior and reproduce previously fixed defects where practical.

## AI Agent Testing

AI agents require testing beyond conventional API and CRUD behavior.

For an AI agent, assess the following categories as applicable:

1. Intent interpretation and instruction compliance.
2. Tool selection and tool-call arguments.
3. Tool execution, error handling, and recovery.
4. Response correctness and grounding in available evidence.
5. Missing information and hallucination prevention.
6. Output format and schema compliance.
7. Multi-step workflow completion.
8. Prompt injection and untrusted input handling.
9. Authentication, authorization, and data isolation.
10. Token limits, timeouts, retries, and rate limits.
11. Regression across model, prompt, tool, or retrieval changes.

### AI Agent Expected Results

Define expected behavior independently of the model's generated response.

Where appropriate:
- Assert exact tool names and arguments when they are deterministic requirements.
- Verify that required tool calls occur and prohibited tool calls do not occur.
- Verify that responses contain correct facts supported by the available data.
- Verify structured outputs against the required schema.
- Define explicit criteria for relevance, completeness, grounding, and refusal behavior.
- Verify that missing information does not result in fabricated facts.
- Check for unintended side effects and unauthorized data access.

Do not require exact natural-language equality for probabilistic responses unless the requirement explicitly demands an exact output.

Use deterministic assertions wherever possible. Use semantic evaluation or model-based evaluators when deterministic checks are insufficient. Define scoring criteria and acceptance thresholds before evaluating results.

Record the model identifier, prompt version, dataset version, and relevant configuration when they materially affect reproducibility.

For model-based evaluation, treat scores as evidence against defined criteria, not as definitive proof of correctness.

### AI Agent Test Data

Use controlled datasets containing representative cases, including:
- Normal user requests.
- Ambiguous requests.
- Missing or conflicting information.
- Invalid tool arguments.
- Tool errors and unavailable dependencies.
- Adversarial or prompt-injection inputs.
- Requests outside the agent's authorized scope.

Keep sensitive information out of test fixtures unless its use is explicitly authorized and necessary.

## BDD and Automated Testing Strategy

Use Given–When–Then when it improves clarity and communication across development, QA, and business stakeholders.

Example:

```gherkin
Feature: Customer lookup by AI Agent

  Scenario: Customer exists
    Given customer "CUST-001" exists
    When the agent is asked to find the customer
    Then the agent returns the correct customer details

  Scenario: Customer does not exist
    Given customer "INVALID-001" does not exist
    When the agent is asked to find the customer
    Then the agent reports that the customer was not found
    And the agent does not invent customer details

  Scenario: Customer service fails
    Given the customer service is unavailable
    When the agent is asked to find a customer
    Then the agent handles the tool error gracefully
```

Choose the appropriate testing approach for each layer:

| Layer | Purpose | Preferred approach |
|---|---|---|
| Unit | Verify isolated business rules and functions. | Existing language-native test framework. |
| Integration | Verify database, API, adapter, and tool interactions. | Test framework with controlled dependencies and fixtures. |
| Acceptance / BDD | Verify externally observable behavior against requirements. | Given–When–Then; Godog or Cucumber when justified. |
| LLM evaluation | Evaluate response quality and agent behavior against a dataset. | Repeatable evaluation harness with explicit assertions and criteria. |
| Security | Verify access boundaries and resistance to adversarial inputs. | Targeted automated and manual security tests. |

For Go projects:
- Prefer the standard `testing` package for unit and integration tests.
- Use table-driven tests for multiple inputs and expected outcomes.
- Use `httptest` for HTTP handlers and API behavior where appropriate.
- Use Godog or another compatible BDD runner when Gherkin scenarios provide sufficient value.
- Use a repeatable test harness and controlled fixtures for LLM evaluation.

Cucumber/Gherkin is optional. Do not introduce a new testing framework solely to satisfy this skill. Follow existing project conventions unless there is a clear reason to change them.

TCD defines the expected behavior; it does not require every documented case to have a separate automated test when automation is impractical. Document relevant limitations and manual verification requirements.

## Workflow

When creating or updating test cases:

1. Inspect relevant specifications, existing TCDs, implementation contracts, and existing tests.
2. Identify requirements, business rules, assumptions, dependencies, and open questions.
3. Define applicable positive, negative, boundary, validation, and failure scenarios.
4. Add security, regression, concurrency, performance, and AI-specific scenarios when relevant.
5. Define preconditions, test data, execution steps, and observable expected results.
6. Map cases to existing requirements and acceptance criteria.
7. Update the TCD using the existing filename and document conventions.
8. Add or update automated tests where practical and consistent with the project's testing strategy.
9. Execute relevant tests when requested or when the implementation workflow calls for it.
10. Record actual results accurately and identify failures, blocked tests, and uncovered requirements.
11. Report unresolved assumptions, missing requirements, and material coverage gaps.

When modifying an existing TCD, preserve stable test case IDs where their meaning remains unchanged. Assign new IDs to new cases and do not silently reuse an ID for a different behavior.

Do not claim that a test passed unless it was actually executed and its result verified.

## Definition of Done

A TCD is ready for review when:
- Scope and related specifications are identified.
- Applicable positive and negative cases are defined.
- Relevant boundary, error, and security scenarios have been considered.
- Expected results are observable and verifiable.
- Requirement traceability is provided where possible.
- AI-specific behavior is covered when an AI agent is involved.
- Automated versus manual verification needs are identified.
- Open questions, assumptions, and known coverage gaps are documented.

Testing work is complete only when the applicable verification has been performed and its actual outcome reported. A drafted TCD alone does not prove that the software passes its tests.

## Checks Before Completion

Before finishing:
- Confirm that the TCD is in `docs/test-cases/`.
- Confirm that the filename follows the required pattern.
- Confirm that related specification documents use the matching date and slug.
- Validate relative links to related documents.
- Check requirement IDs and test case ID consistency.
- Check for applicable positive and negative coverage.
- Check that expected results are specific and testable.
- Check consistency between the TCD, requirements, and automated tests.
- Report documents created or updated, tests executed, results, and remaining gaps.
