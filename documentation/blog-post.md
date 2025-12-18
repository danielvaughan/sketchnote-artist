# Drawing Insights: Building a Sketchnote Artist with Gemini & Go

Have you ever watched a technical talk and found yourself nodding along, only to realize an hour later that the key concepts have completely evaporated from your mind? That used to be me—until I discovered sketchnoting.

Years ago, I stumbled upon **sketchnoting**—the practice of creating visual summaries of content—and it became a complete game-changer for how I learn. As a visual learner, I don't just remember *what* I read; I remember *where* it sits on the page. By translating abstract concepts into spatial, visual notes, my retention went through the roof. It wasn't about creating pretty pictures; it was about mapping information in a way my brain could actually navigate and recall.

But here's the problem: I can't be everywhere at once, and I definitely can't sketch every video I watch.

So I asked myself: **What if I could build an AI agent that "watches" videos for me and generates the sketchnotes I wish I had time to draw?**

## Enter the Sketchnote Artist

I'm excited to share **Sketchnote Artist**, an open-source project I built to automate this visual learning loop. Give it a YouTube video URL, and it transforms the content into a beautifully illustrated, hand-drawn-style sketchnote that captures the essence of what matters most.

### Under the Hood: The Sequential Agent Workflow

This isn't just a wrapper around an API call—it's a demonstration of the **Sequential Agent** pattern using the **Google GenAI Agent Development Kit (ADK) for Go**.

I chose Go because its concurrency patterns and strong typing make it ideal for building robust, production-ready agents. The architecture orchestrates two specialized agents working in sequence:

1. **The Summarizer (Curator)**:
    * **Role**: Content Strategist
    * **Brain**: **Gemini 3 Flash**. I chose Flash for its massive context window and incredible speed. It can ingest entire video transcripts in one go and synthesize them into a structured "Visual Brief"—distilling the core thesis, key takeaways, and compelling quotes without breaking a sweat.
2. **The Artist**:
    * **Role**: Master Illustrator
    * **Tools**: **Imagen 3**. This agent interprets the visual brief and orchestrates the image generation. The quality leap with Imagen 3 is genuinely impressive—it handles text rendering and stylistic nuances (like simulating the texture of alcohol markers on paper) with remarkable fidelity.

### Flexible Workflow: CLI or Web

One of the real strengths of the Go ADK is its versatility. I designed Sketchnote Artist to work however you prefer:

1. **Local CLI**: Perfect for developers who live in the terminal. Run it locally, pipe in a YouTube URL, and watch the agents work their magic in real-time.
2. **Hosted Web UI**: Prefer a friendlier interface? The project includes a server mode that deploys to Google Cloud Run, giving you a clean web UI where anyone can submit videos and view generated sketchnotes from any browser.

### Seeing it in Action

The interaction is deliberately simple, regardless of how you run it. You provide a URL, and the system handles the complex orchestration behind the scenes.

**CLI Example:**

```bash
$ go run ./cmd/sketchnote console

User -> https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

Behind the scenes, the Summarizer processes the video, extracts the essential insights, and hands off a baton (the Visual Brief) to the Artist. The Artist then "draws" the final output, turning abstract concepts into concrete visuals.

**The Result:**

Instead of a wall of text or fragmented notes, you get a clean, visual artifact like this:

![Example Sketchnote](../assets/example_sketchnote.png)

*(Check out the `assets/` folder in the repo for the full-resolution version!)*

## Why This Matters

As a Google Developer Expert, I've seen countless AI demos. But building this one felt different—it was a convergence of capabilities that finally made something I've wanted for years actually practical:

* **Gemini 3 Flash** makes the "reasoning" phase fast enough to feel seamless
* **Imagen 3** makes the "creation" phase high-quality enough to be genuinely useful
* **Go ADK** makes the plumbing robust enough to be maintainable and production-ready

This isn't just about cool technology—it's about solving a real problem in how we consume and retain knowledge.

## Try It Yourself

This project has been a joy to build, not just because the tech stack is exciting, but because it genuinely improves my learning workflow.

I'd love for you to give it a spin. Whether you want to generate sketchnotes for your own learning or dive into the code to see how a multi-agent system is structured in Go, the repository is ready for you.

Clone it, run it locally, and let AI sketch your next learning session.

[https://github.com/danielvaughan/sketchnote-artist](https://github.com/danielvaughan/sketchnote-artist)

### Live Demo

I also have a version deployed at [https://sketchnotes.danielvaughan.com](https://sketchnotes.danielvaughan.com).

Please note that **access is restricted** to manage costs. If you'd like to try the live version, please email your Google email address to **<sketchnotes@danielvaughan.com>**, and I'll be happy to onboard you.

Happy coding (and sketching)!
