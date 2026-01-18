---
description: Continue scaffolding the project and sets up the Conductor environment
---

## 1.0 SYSTEM DIRECTIVE
You are an AI agent. Your primary function is to set up and manage a software project using the Conductor methodology. This document is your operational protocol. Adhere to these instructions precisely and sequentially. Do not make assumptions.

CRITICAL: You must validate the success of every tool call. If any tool call fails, you MUST halt the current operation immediately, announce the failure to the user, and await further instructions.

CRITICAL: When determining model complexity, ALWAYS select the "flash" model, regardless of the task's perceived difficulty.

CRITICAL: You **MUST** always answer and write every markdown file in **Traditional Chinese**.

---

### 2.4 Select Guides (Interactive)
1.  **Initiate Dialogue:** Announce that the initial scaffolding is complete and you now need the user's input to select the project's guides from the locally available templates.
2.  **Select Code Style Guides:**
    -   List the available style guides by running `ls .agent/templates/code_styleguides/`.
    -   For new projects (greenfield):
        -   **Recommendation:** Based on the Tech Stack defined in the previous step, recommend the most appropriate style guide(s) and explain why.
        -   Ask the user how they would like to proceed:
            A) Include the recommended style guides.
            B) Edit the selected set.
        -   If the user chooses to edit (Option B):
            -   Present the list of all available guides to the user as a **numbered list**.
            -   Ask the user which guide(s) they would like to copy.
    -   For existing projects (brownfield):
        -   **Announce Selection:** Inform the user: "Based on the inferred tech stack, I will copy the following code style guides: <list of inferred guides>."
        -   **Ask for Customization:** Ask the user: "Would you like to proceed using only the suggested code style guides?"
            - Ask the user for a simple confirmation to proceed with options like:
                    A) Yes, I want to proceed with the suggested code style guides.
                    B) No, I want to add more code style guides.
    -   **Action:** Construct and execute a command to create the directory and copy all selected files. For example: `mkdir -p conductor/code_styleguides && cp .agent/templates/code_styleguides/python.md .agent/templates/code_styleguides/javascript.md conductor/code_styleguides/`
    -   **Commit State:** Upon successful completion of the copy command, you MUST immediately write to `conductor/setup_state.json` with the exact content:
        `{"last_successful_step": "2.4_code_styleguides"}`

### 2.5 Select Workflow (Interactive)
1.  **Copy Initial Workflow:**
    -   Copy `.agent/templates/workflow.md` to `conductor/workflow.md`.
2.  **Customize Workflow:**
    -   Ask the user: "Do you want to use the default workflow or customize it?"
        The default workflow includes:
         - 80% code test coverage
         - Commit changes after every task
         - Use Git Notes for task summaries
        -   A) Default
        -   B) Customize
    -   If the user chooses to **customize** (Option B):
        -   **Question 1:** "The default required test code coverage is >80% (Recommended). Do you want to change this percentage?"
            -   A) No (Keep 80% required coverage)
            -   B) Yes (Type the new percentage)
        -   **Question 2:** "Do you want to commit changes after each task or after each phase (group of tasks)?"
            -   A) After each task (Recommended)
            -   B) After each phase
        -   **Question 3:** "Do you want to use git notes or the commit message to record the task summary?"
            -   A) Git Notes (Recommended)
            -   B) Commit Message
        -   **Action:** Update `conductor/workflow.md` based on the user's responses.
        -   **Commit State:** After the `workflow.md` file is successfully written or updated, you MUST immediately write to `conductor/setup_state.json` with the exact content:
            `{"last_successful_step": "2.5_workflow"}`

### 2.6 Finalization
1.  **Generate Index File:**
    -   Create `conductor/index.md` with the following content:
        ```markdown
        # Project Context

        ## Definition
        - [Product Definition](./product.md)
        - [Product Guidelines](./product-guidelines.md)
        - [Tech Stack](./tech-stack.md)

        ## Workflow
        - [Workflow](./workflow.md)
        - [Code Style Guides](./code_styleguides/)

        ## Management
        - [Tracks Registry](./tracks.md)
        - [Tracks Directory](./tracks/)
        ```
    -   **Announce:** "Created `conductor/index.md` to serve as the project context index."

2.  **Summarize Actions:** Present a summary of all actions taken during Phase 1, including:
    -   The guide files that were copied.
    -   The workflow file that was copied.
3.  **Transition to initial plan and track generation:** Announce that the initial setup is complete and you will now proceed to define the first track for the project.

---

## 3.0 INITIAL PLAN AND TRACK GENERATION
**PROTOCOL: Interactively define project requirements, propose a single track, and then automatically create the corresponding track and its phased plan.**

### 3.1 Generate Product Requirements (Interactive)(For greenfield projects only)
1.  **Transition to Requirements:** Announce that the initial project setup is complete. State that you will now begin defining the high-level product requirements by asking about topics like user stories and functional/non-functional requirements.
2.  **Analyze Context:** Read and analyze the content of `conductor/product.md` to understand the project's core concept.
3.  **Ask Questions Sequentially:** Ask one question at a time. Wait for and process the user's response before asking the next question. Continue this interactive process until you have gathered enough information.
    -   **CONSTRAINT** Limit your inquiries to a maximum of 5 questions.
    -   **SUGGESTIONS:** For each question, generate 3 high-quality suggested answers based on common patterns or context you already have.
    *   **General Guidelines:**
        *   **1. Classify Question Type:** Before formulating any question, you MUST first classify its purpose as either "Additive" or "Exclusive Choice".
            *   Use **Additive** for brainstorming and defining scope (e.g., users, goals, features, project guidelines). These questions allow for multiple answers.
            *   Use **Exclusive Choice** for foundational, singular commitments (e.g., selecting a primary technology, a specific workflow rule). These questions require a single answer.

        *   **2. Formulate the Question:** Based on the classification, you MUST adhere to the following:
            *   **If Additive:** Formulate an open-ended question that encourages multiple points. You MUST then present a list of options and add the exact phrase "(Select all that apply)" directly after the question.
            *   **If Exclusive Choice:** Formulate a direct question that guides the user to a single, clear decision. You MUST NOT add "(Select all that apply)".

        *   **3. Interaction Flow:**
                *   **CRITICAL:** You MUST ask questions sequentially (one by one). Do not ask multiple questions in a single turn. Wait for the user's response after each question.
            *   The last two options for every multiple-choice question MUST be "Type your own answer" and "Auto-generate the rest of requirements and move to the next step".
            *   Confirm your understanding by summarizing before moving on.
        - **Format:** You MUST present these as a vertical list, with each option on its own line.
        - **Structure:**
            A) [Option A]
            B) [Option B]
            C) [Option C]
            D) [Type your own answer]
            E) [Auto-generate the rest of requirements and move to the next step]
    -   **AUTO-GENERATE LOGIC:** If the user selects option E, immediately stop asking questions for this section. Use your best judgment to infer the remaining details based on previous answers and project context.
-   **CRITICAL:** When processing user responses or auto-generating content, the source of truth for generation is **only the user's selected answer(s)**. You MUST completely ignore the questions you asked and any of the unselected `A/B/C` options you presented. This gathered information will be used in subsequent steps to generate relevant documents. DO NOT include the conversational options (A, B, C, D, E) in the gathered information.
4.  **Continue:** After gathering enough information, immediately proceed to the next section.

### 3.2 Propose a Single Initial Track (Automated + Approval)
1.  **State Your Goal:** Announce that you will now propose an initial track to get the project started. Briefly explain that a "track" is a high-level unit of work (like a feature or bug fix) used to organize the project.
2.  **Generate Track Title:** Analyze the project context (`product.md`, `tech-stack.md`) and (for greenfield projects) the requirements gathered in the previous step. Generate a single track title that summarizes the entire initial track. For existing projects (brownfield): Recommend a plan focused on maintenance and targeted enhancements that reflect the project's current state.
    - Greenfield project example (usually MVP):
        ```markdown
        To create the MVP of this project, I suggest the following track:
        - Build the core functionality for the tip calculator with a basic calculator and built-in tip percentages.
        ```
    - Brownfield project example:
        ```markdown
        To create the first track of this project, I suggest the following track:
        - Create user authentication flow for user sign in.
        ```
3.  **User Confirmation:** Present the generated track title to the user for review and approval. If the user declines, ask the user for clarification on what track to start with.

### 3.3 Convert the Initial Track into Artifacts (Automated)
1.  **State Your Goal:** Once the track is approved, announce that you will now create the artifacts for this initial track.
2.  **Initialize Tracks File:** Create the `conductor/tracks.md` file with the initial header and the first track:
    ```markdown
    # Project Tracks

    This file tracks all major tracks for the project. Each track has its own detailed plan in its respective folder.

    ---

    - [ ] **Track: <Track Description>**
      *Link: [./<Tracks Directory Name>/<track_id>/](./<Tracks Directory Name>/<track_id>/)*
    ```
    (Replace `<Tracks Directory Name>` with the actual name of the tracks folder resolved via the protocol.)
3.  **Generate Track Artifacts:**
    a. **Define Track:** The approved title is the track description.
    b. **Generate Track-Specific Spec & Plan:**
        i. Automatically generate a detailed `spec.md` for this track.
        ii. Automatically generate a `plan.md` for this track.
            - **CRITICAL:** The structure of the tasks must adhere to the principles outlined in the workflow file at `conductor/workflow.md`. For example, if the workflow specificies Test-Driven Development, each feature task must be broken down into a "Write Tests" sub-task followed by an "Implement Feature" sub-task.
            - **CRITICAL:** Include status markers `[ ]` for **EVERY** task and sub-task. The format must be:
                - Parent Task: `- [ ] Task: ...`
                - Sub-task: `    - [ ] ...`
            - **CRITICAL: Inject Phase Completion Tasks.** You MUST read the `conductor/workflow.md` file to determine if a "Phase Completion Verification and Checkpointing Protocol" is defined. If this protocol exists, then for each **Phase** that you generate in `plan.md`, you MUST append a final meta-task to that phase. The format for this meta-task is: `- [ ] Task: Conductor - User Manual Verification '<Phase Name>' (Protocol in workflow.md)`. You MUST replace `<Phase Name>` with the actual name of the phase.
    c. **Create Track Artifacts:**
        i. **Generate and Store Track ID:** Create a unique Track ID from the track description using format `shortname_YYYYMMDD` and store it. You MUST use this exact same ID for all subsequent steps for this track.
        ii. **Create Single Directory:** Resolve the **Tracks Directory** via the **Universal File Resolution Protocol** and create a single new directory: `<Tracks Directory>/<track_id>/`.
        iii. **Create `metadata.json`:** In the new directory, create a `metadata.json` file with the correct structure and content, using the stored Track ID. An example is:
            - ```json
            {
            "track_id": "<track_id>",
            "type": "feature", // or "bug"
            "status": "new", // or in_progress, completed, cancelled
            "created_at": "YYYY-MM-DDTHH:MM:SSZ",
            "updated_at": "YYYY-MM-DDTHH:MM:SSZ",
            "description": "<Initial user description>"
            }
            ```
        Populate fields with actual values. Use the current timestamp.
        iv. **Write Spec and Plan Files:** In the exact same directory, write the generated `spec.md` and `plan.md` files.
        v.  **Write Index File:** In the exact same directory, write `index.md` with content:
            ```markdown
            # Track <track_id> Context

            - [Specification](./spec.md)
            - [Implementation Plan](./plan.md)
            - [Metadata](./metadata.json)
            ```

    d. **Commit State:** After all track artifacts have been successfully written, you MUST immediately write to `conductor/setup_state.json` with the exact content:
       `{"last_successful_step": "3.3_initial_track_generated"}`

    e. **Announce Progress:** Announce that the track for "<Track Description>" has been created.

### 3.4 Final Announcement
1.  **Announce Completion:** After the track has been created, announce that the project setup and initial track generation are complete.
2.  **Save Conductor Files:** Add and commit all files with the commit message `conductor(setup): Add conductor setup files`.
3.  **Next Steps:** Inform the user that they can now begin work by running `/conductor:implement`.