# Current Feature Inventory

## Repository

- Name: `AI-SDK-OPENAI`
- SDK: OpenAI Agents SDK
- Positioning: Operational OpenAI agents with mission routing and safety gates.

## Implemented Today

- Shared Agents Army routing and support-agent selection.
- OpenAI Agent object creation path.
- FastAPI health/run service and CLI runner.
- Skill-aware mission plan and verification notes.
- Dockerfile, CI workflow, pytest contract tests, and portfolio metadata.

## Not Yet Implemented

- Wire Runner execution with configured model and tools.
- Add safety gates for tool calls and external actions.
- Add evals for mission quality, refusal behavior, and cost.

## Verification Contract

- The local runner must complete without crashing when optional SDK credentials are missing.
- The API contract must return routing and verification fields.
- Tests must prove mission routing and a security-focused SENTINEL route.
