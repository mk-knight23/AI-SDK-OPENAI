"""Production-style OpenAI Agents runtime for Kazi's Agents Army."""

from pathlib import Path
import sys

sys.path.append(str(Path(__file__).resolve().parent / "core"))
from agents_army_core import MissionRequest, build_mission_plan, render_system_instructions


def run_openai_agents_mission(mission_text: str) -> dict:
    plan = build_mission_plan(MissionRequest(mission_text))
    system_msg = render_system_instructions(plan)

    try:
        from agents import Agent
    except Exception as exc:
        return {
            "primary": plan.primary,
            "support": plan.support,
            "result": None,
            "verification": f"OpenAI Agents dependency missing: {exc}",
        }

    # Keep this compatible with minimal SDK surface.
    _ = Agent(name=plan.primary, instructions=system_msg)

    return {
        "primary": plan.primary,
        "support": plan.support,
        "result": "Agent object instantiated. Connect Runner + model in environment.",
        "verification": "Routing + agent instantiation succeeded.",
    }
