# 🧠 Cognitive Agent Runtime — Full Architecture

This project implements a **cognitive runtime for autonomous agents** in Go, featuring:

- filesystem manipulation tools  
- a hybrid planner (rules + LLM)  
- a contextual tool grounder  
- metacognition  
- continuous self‑optimization (AOC)  
- persistent cognitive memory  
- an execution monitor  
- specialized agents (analysis, refactor, docs)  

The goal is to provide a **safe, extensible, autonomous environment** for agents capable of analyzing, modifying, and documenting entire projects.

---

## 🧩 Overall Architecture

Main cognitive flow:

```text
Goal → Planner → Grounder → Runtime → Tools
          ↑         ↓
     Metacognition ← ExecutionMonitor
          ↓
     StrategyEngine
          ↓
           AOC
```

Each component contributes a different cognitive layer and relies on project‑scoped persistent memory.

---

## 🧠 Cognitive Memory & Planner Memory

### 🧠 CognitiveMemory

Project‑scoped cognitive memory stored at:

```
.opencode/<projectID>/cognitive_memory.json
```

Stores:

- last_tool  
- last_result  
- fail_count  
- success_count  
- detected patterns (dependency_cycles, long_functions_detected)  
- AOC snapshots  

API:

- Remember(key, value)  
- Recall(key)  
- Increment(key)  
- Reset(key)  

Helpers:

- RememberLastTool(name)  
- RememberLastResult(result)  
- RegisterFailure()  
- RegisterSuccess()  

---

### 🧠 PlannerMemory

Planner‑specific memory stored at:

```
.opencode/<projectID>/planner_memory.json
```

Stores:

- SuccessfulSteps  
- FailedSteps  
- GoalPatterns  

API:

- RecordSuccess  
- RecordFailure  
- RecordGoal  
- IsFrequentlyFailing  
- Score  

Used to:

- filter steps that fail too often  
- prioritize steps with better cognitive scores  

---

## 🔍 ExecutionMonitor

Tracks execution step‑by‑step:

- RepeatCount → repeated step  
- FailureCount → consecutive errors  
- StallCount → identical result hash  
- LastStepTime → temporal stalling  

Flags:

- loop_detected  
- repeated_failures  
- stalled  

These feed into Metacognition and StrategyEngine.

---

## 🧠 Metacognition

Produces:

```go
type MetaEvaluation struct {
    Confidence float64
    Flags      []string
    Advice     string
}
```

Detects:

- loops  
- repeated failures  
- stalling  
- ambiguous goals  
- lack of progress  
- cognitive memory patterns  
- incorrect planner/grounder modes  

---

## 🎛️ StrategyEngine

Dynamically adjusts:

- PlannerMode: coarse, fine, aggressive, conservative  
- GroundingMode: strict, flexible  
- ToolBias: analysis, refactor, docs, none  
- agent switching  

Based on:

- metacognition  
- cognitive memory  
- execution monitor  
- goal type  
- project size  

---

## 🔁 AOC — Autonomous Optimization Cycle

Learns persistent preferences:

```go
type AOCUpdate struct {
    PreferredPlannerMode   string
    PreferredGroundingMode string
    PreferredAgentBias     string
    UpdatedAt              time.Time
}
```

Stores:

- aoc_last_update  
- aoc_failures_total  
- aoc_loops_total  
- aoc_stalls_total  
- aoc_last_confidence  

---

## 🧠 HybridPlanner

Hybrid planner combining:

- deterministic rules  
- LLM reasoning  
- planner memory  

Modes:

- coarse  
- fine  
- aggressive  
- conservative  

Filters failing steps and sorts by cognitive score.

---

## 🧭 ContextualToolGrounder

Maps planner steps to actual tools.

Modes:

- strict → literal grounding  
- flexible → semantic grounding  

Uses:

- project statistics  
- cognitive memory  
- tool selection heuristics  

---

## 🤖 Specialized Agents

Interface:

```go
type SpecializedAgent interface {
    Name() string
    CanHandle(goal string) bool
    Run(goal string, ctx *AgentContext) AgentContext
}
```

### AnalysisAgent
- dependencies  
- metrics  
- structure  
- diagnostics  

### RefactorAgent
- cleanup  
- optimization  
- simplification  

### DocsAgent
- documentation  
- explanations  
- summaries  

All delegate to:

```go
return a.runtime.Run(goal)
```

---

## 🧠 AgentRuntime

Cognitive execution core:

- runs the plan  
- invokes tools  
- updates memory  
- queries metacognition  
- applies StrategyEngine  
- triggers AOC  
- monitors loops and stalling  

Signature:

```go
func (rt *AgentRuntime) Run(goal string) AgentContext
```

---

# 🛠️ Toolset: Filesystem & Patch Engine

All tools operate inside:

```
workspace/
```

---

## 📁 Basic Filesystem Tools

- read_file  
- write_file  
- create_file  
- delete_file  
- rename_file  
- copy_file  
- move_file  
- file_exists  
- read_dir  
- stat_file  
- touch_file  

---

## ✏️ Editing Tools

- append_file  
- truncate_file  

---

## 🔍 Search Tools

- search_in_file  
- grep  
- search_replace  
- search_regex  

---

## 🧩 Patch & Refactor Tools

- apply_patch  
- apply_patch_fuzzy  
- apply_patch_auto  
- apply_patch_structured  

---

## 🧰 Utility Tools

- diff_files  
- zip_dir  
- unzip  
- format_code  

---

# 🖥️ Sandbox Commands (run_command)

### Basic Commands

- count_lines  
- file_size  
- validate_json  
- echo  
- word_count  
- char_count  
- sha256  
- list_dir  
- head / tail  
- search  
- now  

### Code Analysis

- count_funcs  
- count_imports  
- find_structs  
- find_interfaces  

### Project Analysis

- project_stats  
- largest_files  
- file_tree  

### Intelligent Commands

- detect_language  
- summarize_file  
- extract_comments  

---

# 🚀 Agent Capabilities

With this architecture, an agent can:

- navigate the filesystem like an IDE  
- safely modify code  
- apply semantic refactors  
- auto‑format code  
- analyze entire projects  
- detect structural patterns  
- self‑adjust strategy  
- learn from previous runs  
- avoid loops and stalling  
- delegate across specialized agents  

It is a **complete cognitive system**, designed for high‑level autonomous agents.