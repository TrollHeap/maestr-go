
# Je vais créer une structure complète pour l'outil Go
# Cela m'aidera à designer correctement avant de créer le code final

tool_structure = {
    "project_name": "maestro",
    "tagline": "Ultra-Learning Practice Tool for ADHD",
    "core_concepts": [
        "Micro-learning chunks",
        "Spaced repetition tracking",
        "Visual progress feedback",
        "Anti-procrastination timers",
        "Code kata execution",
        "System architecture visualization"
    ],
    "key_features": [
        {
            "name": "Quick Start Sessions",
            "description": "15-30 minute focused practice blocks",
            "adhd_fit": "Prevents overwhelm with short bursts"
        },
        {
            "name": "Visual Progress Trees",
            "description": "Tree-based visualization of learning paths",
            "adhd_fit": "Shows immediate progress and context"
        },
        {
            "name": "Streak Tracking",
            "description": "Daily practice streaks with visual indicators",
            "adhd_fit": "Dopamine rewards for consistency"
        },
        {
            "name": "Code Kata System",
            "description": "Run, test, and practice code challenges",
            "adhd_fit": "Hands-on learning reinforces concepts"
        },
        {
            "name": "Spaced Repetition Scheduler",
            "description": "SM-2 algorithm for optimal review timing",
            "adhd_fit": "Prevents both cramming and forgetting"
        },
        {
            "name": "System Context Diagrams",
            "description": "ASCII art mental models of architecture",
            "adhd_fit": "Visual analogies for complex concepts"
        }
    ],
    "data_structure": {
        "Exercise": {
            "id": "string",
            "title": "string",
            "description": "string",
            "difficulty": "1-3",
            "domain": "golang|linux|architecture",
            "content": "string (code template)",
            "tests": "[]TestCase",
            "created_at": "timestamp",
            "completed": "bool",
            "last_reviewed": "timestamp",
            "ease_factor": "2.5",
            "interval_days": "int",
            "repetitions": "int"
        }
    }
}

import json
print(json.dumps(tool_structure, indent=2))
