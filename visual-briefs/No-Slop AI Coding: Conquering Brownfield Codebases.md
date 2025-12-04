Title Text: No-Slop AI Coding: Conquering Brownfield Codebases
Youtube URL: https://www.youtube.com/watch?v=rmvDxxNubIg
Featured Speaker: Dex Horthy (Draw a simple caricature or a profile icon with "HumanLayer" beneath)
Central Metaphor: A well-engineered pipeline or a meticulously planned construction project, ensuring a clean and efficient flow.
Visual Hierarchy:
Header (Big Text): The Future of AI Software Engineering: Context Management for No-Slop Coding

Sub-Headers (Medium Text): The Core Problem: "Slop" in Brownfield Codebases
Details (Small Text):
*   AI struggles with complex, legacy codebases ("brownfield").
*   AI-generated code in these environments often leads to rework, churn, and technical debt ("slop").
*   Goal: Achieve "no-slop" coding and maintain "mental alignment."

Sub-Headers (Medium Text): The Constraint: The "Dumb Zone"
Details (Small Text):
*   LLMs are stateless; performance depends on the context window.
*   "The Dumb Zone": Model reasoning degrades rapidly once ~40% of the context window is filled.
*   **Memorable Quote:** "Once you fill more than roughly **40% of a context window**...the model's ability to reason, retrieve correct information, and follow complex instructions degrades rapidly."
*   Errors lead to context filling with "noise," trapping the model in a failure loop.

Sub-Headers (Medium Text): The Solution: Intentional Compaction
Details (Small Text):
*   Keep the context window small and relevant.
*   **Strategy 1: Start Fresh**
    *   Agents frequently restart with clean context.
    *   Carry over only compressed summaries (e.g., Markdown files) of learned information.
*   **Strategy 2: Sub-Agents for Retrieval**
    *   Use sub-agents to isolate context, not anthropomorphize roles.
    *   Sub-agents read massive code, find specific answers, and return *only* the answer.

Sub-Headers (Medium Text): The Workflow: Research, Plan, Implement
Details (Small Text):
*   **Phase 1: Research (Compression of Truth)**
    *   Goal: Understand the system before coding.
    *   Method: Agent explores codebase for relevant files/dependencies.
    *   Output: `research.md` (300-1000 lines) – an objective snapshot, compressing massive repo "truth."
*   **Phase 2: Plan (Compression of Intent)**
    *   Goal: Outline exact changes (filenames, line numbers, code snippets).
    *   **Crucial Step: Human Review!** Do not outsource thinking. A bad plan leads to bad code. Catching misunderstandings here prevents hundreds of useless lines.
*   **Phase 3: Implement (Reliable Execution)**
    *   Goal: Execute the approved plan.
    *   Result: High-reliability execution due to specific plan and research context.

Sub-Headers (Medium Text): Key Conclusions & Philosophy
Details (Small Text):
*   "Spec-Driven Development" (SDD) is actually "Context Engineering" – managing what AI sees.
*   Mental Alignment: Humans review the high-leverage *Plan*, not massive AI-generated code.
*   **Memorable Quote:** "Humans should focus on the high-leverage parts of the pipeline (Research and Planning). If the research is wrong, the code will be wrong."
*   Actionable: Pick one tool and get "reps" (practice) to build intuition.

Iconography Suggestions:
*   "Slop": A tangled ball of yarn or messy spaghetti code with a red 'X'.
*   "Brownfield": Old, crumbling brick wall or tangled, overgrown vines.
*   "Greenfield": A fresh, clean sprout or a new, simple building block.
*   "Dumb Zone": A brain with smoke coming out or a confused face, a context window filling with garbled text/noise, a digital "DANGER" sign.
*   "Context Window": A rectangular digital screen or a whiteboard.
*   "Compaction": A hydraulic press, a 'zip' file icon, or a small, neatly folded document.
*   "Start Fresh": A 'reset' button, a clean slate, or an empty notepad.
*   "Sub-Agents": Small, focused robots with magnifying glasses, delivering a single, precise answer in a scroll.
*   "Research": A magnifying glass over code, an open book/laptop, a camera taking a 'snapshot'.
*   "Plan": A blueprint, a checklist, a detailed map, a human hand with a pen over a document.
*   "Implement": A robotic arm building something, gears turning, a conveyor belt.
*   "Mental Alignment": Interlocking gears, two heads with a connecting thought bubble, puzzle pieces fitting.
*   "Human Leverage": A human hand subtly guiding a robotic arm or an architect at a drawing board.
*   "Reps": A dumbbell, a tally counter, or a growing stack of successful tasks.

Mood: Informative, Strategic, Empowering, with a touch of urgency regarding the "Dumb Zone."