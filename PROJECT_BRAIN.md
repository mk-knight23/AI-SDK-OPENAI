# Project Brain: AI-SDK-OPENAI

## Purpose

Operational OpenAI agents with mission routing and safety gates.

## Current State

- Shared Agents Army routing and support-agent selection.
- OpenAI Agent object creation path.
- FastAPI health/run service and CLI runner.
- Skill-aware mission plan and verification notes.
- Dockerfile, CI workflow, pytest contract tests, and portfolio metadata.

## Upgrade Direction

- Wire Runner execution with configured model and tools.
- Add safety gates for tool calls and external actions.
- Add evals for mission quality, refusal behavior, and cost.

## Quality Bar

- Keep the repository runnable from a fresh clone.
- Keep generated caches and local secrets out of git.
- Keep README, skill matrix, tests, and CI aligned with actual behavior.
