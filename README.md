Here is a professional, portfolio-ready `README.md`. It highlights the architectural complexity and the "cool" features you requested, framing the project as a resilient enterprise-grade system.

***

# Postificus: Centralized Content Distribution Engine

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) ![Echo Framework](https://img.shields.io/badge/Framework-Echo_v4-25A162?style=flat) ![Redis](https://img.shields.io/badge/Queue-Asynq%20%2F%20Redis-DC382D?style=flat&logo=redis) ![Status](https://img.shields.io/badge/Status-Active_Development-success)

**Postificus** is a high-performance, event-driven backend system designed to solve the fragmentation of technical blogging. It allows authors to write content once (in Markdown) and orchestrate its distribution across a disparate ecosystem of platforms (Medium, Dev.to, Hashnode, LinkedIn, X/Twitter) while strictly enforcing SEO authority via Canonical URLs.

## 🚀 How It Works

The system moves beyond simple API scripting by implementing a **Multi-Stage Resilience Pipeline**.

1.  **Ingestion:** The user submits a raw Markdown payload via the **Echo** REST API.
2.  **Orchestration:** The payload is validated and pushed to a **Redis** backing queue (via `Asynq`), immediately freeing the HTTP worker.
3.  **Processing:** Background workers pick up the job and execute parallel publishing tasks based on platform-specific adapters.
4.  **AI Transformation:** An embedded LLM agent analyzes the long-form blog to generate platform-native social media summaries for LinkedIn and Twitter.

## 🛡️ The "Unstoppable" Delivery System

Unlike standard cross-posters that fail when an API goes down, Postificus implements a **Hierarchical Fallback Strategy**. It attempts to deliver content using the most efficient method first, degrading gracefully to more heavy-handed approaches only when necessary.

| Priority | Method | Description |
| :--- | :--- | :--- |
| **1. Primary** | **Direct API** | Uses official REST/GraphQL APIs (e.g., Dev.to, Hashnode). Fastest and most reliable. |
| **2. Fallback** | **Email-to-Post** | If APIs fail or are rate-limited, the system constructs a MIME payload and posts via the platform's "Post via Email" feature. |
| **3. Last Resort** | **Stealth Automation** | The "Nuclear Option." Uses **Go-Rod** to launch a headless browser, bypass bot detection, log in (handling 2FA), and physically type the content into the editor. |

## ✨ Key Features

### 🧠 Core Backend
*   **Event-Driven Architecture:** Decoupled ingestion and processing layers allow the system to handle traffic spikes without blocking.
*   **SEO Guardrails:** Automatically manages `rel=canonical` tags. The system posts to the "Primary" domain first, retrieves the URL, and forces all secondary platforms to point back to the original source.
*   **Concurrency Control:** Worker pools are rate-limited per domain to prevent IP bans.

### 🤖 AI & Social Agents
*   **Context-Aware Repurposing:** Integrates with OpenAI/Anthropic APIs to read the full blog post and generate:
    *   A professional, thread-optimized summary for **LinkedIn**.
    *   A punchy, hashtag-optimized thread for **X (Twitter)**.
*   **Auto-Scheduling:** Social posts are queued to go out *after* the blog is confirmed live.

### 🕵️ Browser Automation
*   **Stealth Mode:** Uses `rod-stealth` to strip `navigator.webdriver` flags, allowing the bot to pass as a human user on Single Page Applications (SPAs).
*   **TOTP Automation:** Built-in 2FA handler. If a login screen requests a code, the bot generates a Time-based One-Time Password (TOTP) on the fly using the server's secret keys.
*   **Rich Text Handling:** Simulates human keyboard events to handle `contenteditable` divs (like Medium's editor) where standard value injection fails.

## 🛠️ Tech Stack

*   **Language:** Golang
*   **Web Framework:** Echo v4
*   **Task Queue:** Asynq (Redis)
*   **Browser Automation:** Go-Rod (Stealth + CDP)
*   **Database:** PostgreSQL (GORM/Sqlc)
*   **Auth:** TOTP (2FA generation), JWT
*   **AI:** OpenAI API / Local LLM bindings

## 🔮 Future Roadmap
*   [ ] Webhook listeners for "Post Published" events to trigger newsletter dispatch.
*   [ ] Image processing pipeline (auto-resize and watermark for different platforms).
*   [ ] Analytics aggregation (scraping view counts from all platforms into a single dashboard).
