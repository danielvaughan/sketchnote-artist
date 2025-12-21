
I have always found visual learning to be the most effective way to retain information, but manually sketching technical talks is time-consuming. To bridge this gap, I built an AI agent to automate the process.

**Sketchnote Artist** is an open-source project designed to "watch" YouTube videos and translate them into illustrated, hand-drawn style sketchnotes using a multi-agent workflow.

Developed as part of the Google #AISprintH2, the project uses three technologies:

* **Google ADK for Go**: Used for agent orchestration, allowing me to build the logic in my preferred language.
* **Gemini 3 Flash**: Efficiently processes long-form video content to generate a structured "Visual Brief."
* **Gemini 3 Pro Image**: Interprets those briefs to produce high-quality visual summaries that maintain a human, sketched aesthetic.

Getting these tools to work together to solve a personal learning bottleneck has been a fun technical challenge.

I have shared the full development story in a blog post and made the source code available on GitHub. Links to both are in the comments.

# GoogleCloud #GeminiAI #Golang #Sketchnoting #GenAI #AgenticWorkflows #AISprintH2

---

📖 **Read the Blog Post**: [blog-post.md](blog-post.md)
💻 **Git Repo**: [https://github.com/danielvaughan/sketchnote-artist](https://github.com/danielvaughan/sketchnote-artist)
🌐 **Live Demo**: [https://sketchnotes.danielvaughan.com](https://sketchnotes.danielvaughan.com)

Check out this example sketchnote generated from a video! 👇

![Example Sketchnote](../sketchnotes/The_Future_of_AI_Beyond_Scaling__Towards_Super-Intelligent_Learning_in_large_bold_hand-lettered_bloc.png)
