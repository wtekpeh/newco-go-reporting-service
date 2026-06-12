# newco-go-reporting-service

# AI Operational Intelligence Assistant

## Overview

The NewCo AI Assistant combines internal operational intelligence with external market intelligence to help management make better decisions.

The assistant does not access the database directly and does not generate SQL. All data access is controlled through approved Go tools and existing reporting services.

Architecture:

User
→ OpenAI Intent Router
→ Approved Tool Selection
→ Go Service Execution
→ Structured Result
→ OpenAI Explanation Layer
→ React Chart Rendering

---

## Security Model

### Executive Users

Executive users (Boss and Managing Director) have access to all branches and can ask questions across the entire organization.

Examples:

- Which site is overloaded?
- Compare Accra and Kumasi.
- What should management focus on?

### Branch Managers

Branch Managers can access the AI assistant but are automatically restricted to their assigned branch.

Examples:

- How is my site performing?
- What risks should I monitor?
- Show ingredient variance.

Branch managers cannot access data from branches outside their assigned scope.

---

## Internal Intelligence Tools

### executive_summary

Provides high-level operational summaries.

Examples:

- Summarize operations for management.
- Give me an executive overview.

### branch_summary

Provides site performance information.

Examples:

- Show performance for Accra Site.
- Which site is the most active?

### site_staff_load

Provides staffing and workload analysis.

Examples:

- Which site is overloaded?
- How many staff are in Kumasi Site?
- What is the workload per staff member?

### ingredient_variance_risk

Provides ingredient variance and usage analysis.

Examples:

- Show ingredient variance for Accra Site.
- Which ingredients have the highest variance?

### planning_risk_summary

Provides planning readiness and execution risk analysis.

Examples:

- What planning risks exist?
- Are there any incomplete plans?

### management_action_summary

Combines operational intelligence into management recommendations.

Examples:

- What should management focus on?
- What are the biggest operational risks?
- What actions should management take?

---

## External Intelligence Tools

### internet_search

Uses Tavily search to retrieve current external information.

Examples:

- What is the current price of cooking oil in Ghana?
- What are food inflation trends?
- What is happening in the rice market?

Data Source:

- Tavily Search API

---

## Strategic Intelligence Tools

### operational_market_intelligence

Combines internal operational intelligence with external market intelligence.

Inputs:

- Site performance
- Staffing load
- Planning risks
- Ingredient variance
- Market intelligence

Examples:

- Considering our current operations, what external risks should management monitor?
- What market conditions could affect our operations?
- How do external market conditions affect our current operations?

Outputs:

- Operational risks
- Staffing risks
- Ingredient risks
- Market risks
- Supply chain risks
- Recommended management actions

---

## Conversational AI

The assistant supports conversational follow-up questions.

Examples:

- Why so?
- What should management do?
- What about Kumasi Site?
- Thank you.

The assistant maintains conversation context and automatically reuses previous operational context when appropriate.

---

## Chart Support

The assistant can generate chart recommendations and chart data for supported tools.

Supported chart types:

- Bar
- Line
- Pie

Charts are rendered by the React frontend using structured datasets returned by the backend.

---

## AI Design Principles

1. OpenAI selects the most appropriate approved tool.
2. Go services execute all data access and business logic.
3. OpenAI explains results using natural language.
4. No direct database access from the LLM.
5. No SQL generation by the LLM.
6. All branch access is enforced by backend authorization and access scoping.
7. Internal and external intelligence can be combined for strategic recommendations.
