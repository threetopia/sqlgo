# Implementation Plan Template

## General Information
- Title: <feature / task name>
- Date: 2026-08-30
- Status: Draft / In Progress / Completed
- Owner: <name or team>

## Objective
Describe the main goal of this task in one or two sentences.

Examples:
- Complete the UI redesign for the student chat feature to match the mockup.
- Ensure the new template can be used in the student app without breaking routing or rendering.

## Scope of Work
- [ ] Areas to be changed
- [ ] Areas that will not be changed
- [ ] Related dependencies / modules

## Current Condition
Describe the current problem or the requirement that needs to be met.

Examples:
- The student templates still use the legacy Metronic layout.
- The UI does not match the ChatSiswa-v3 mockup.
- The page must remain compatible with the existing Go template engine and funcMap.

## Implementation Plan

### 1. Analysis and Mapping
- Identify the templates involved.
- Check base templates, layouts, and assets used.
- Note the differences between the reference design and the current state.

### 2. Preparation and Backup
- Back up the old template if necessary.
- Make sure changes are made in the correct location.
- Verify that asset / template paths match the application that uses them.

### 3. Create / Update UI
- Update the main page according to the reference design.
- Ensure the HTML structure stays compatible with Go templates.
- Reuse existing components / blocks whenever possible.

### 4. Integration and Validation
- Build or test the relevant application.
- Confirm that the template renders without errors.
- Verify that the layout on the main page matches the target design.

## Detailed Steps
1. <Step 1>
2. <Step 2>
3. <Step 3>
4. <Step 4>

## Risks / Notes
- Risk 1: Asset path does not match the static root.
- Risk 2: Template layout conflicts with existing CSS.
- Risk 3: HTMX components do not work correctly when the page is re-rendered.

## Verification
Record evidence that the work is complete and correct.

Examples:
- `go build ./cmd/api_dummy/http`
- Check the rendered template output in the browser
- Confirm the layout matches the mockup

## Final Outcome
Describe the expected result once implementation is complete.

Examples:
- All student pages use the modern ChatSiswa-v3 layout.
- No layout issues remain.
- The base template and routing continue to work as expected.
