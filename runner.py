"""CLI runner for OpenAI Agents adapter."""

import argparse

try:
    from .app import run_openai_agents_mission
except ImportError:
    from app import run_openai_agents_mission


def demo(mission: str) -> None:
    out = run_openai_agents_mission(mission)
    print("[OpenAI Agents] primary:", out.get("primary"))
    print("[OpenAI Agents] support:", out.get("support"))
    print("[OpenAI Agents] result:", out.get("result"))
    print("[OpenAI Agents] verification:", out.get("verification"))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--mission", default="implement ai workflow and add evals")
    args = parser.parse_args()
    demo(args.mission)
