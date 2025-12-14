Here is a formal summary of the legal and compliance standing for your application. You can save this text as a **Compliance Reference Document** for your internal records.

***

# Legal & Compliance Reference: AI Video Summarization & Sketchnoting

**Date:** December 14, 2025
**Application Workflow:** YouTube API Transcript Fetch $\rightarrow$ AI Summarization (RAM only) $\rightarrow$ Visual Sketchnote Output.

## 1. Executive Summary
The application workflow operates in a **Low-to-Moderate Risk** zone. By adhering to a strict "No Storage" policy for transcript data and using the Official YouTube API, the application avoids the most common legal pitfalls (scraping and caching violations). The output (Sketchnote) constitutes a strong case for "Fair Use" under copyright law due to its transformative nature and lack of market substitution.

---

## 2. Platform Compliance (YouTube API Services)
*Governed by the YouTube API Services Terms of Service (ToS) and Developer Policies.*

### A. Data Access (The "Scraping" Prohibition)
* **Rule:** YouTube Terms strictly prohibit the collection of content via automated means (scraping, robots, botnets) outside of the official API.
* **Compliance Status:** **PASS.**
* **Reasoning:** The application uses the Official YouTube API to retrieve captions, complying with the authorized method of access.

### B. Data Retention (The "Stale Data" Policy)
* **Rule:** API Clients must not store or cache API Data (specifically captions) for longer than 30 days, and in many cases, storage is discouraged entirely.
* **Compliance Status:** **PASS.**
* **Reasoning:** The application processes transcripts in **volatile memory (RAM)** only. The raw transcript data is used strictly for the duration of the summarization session and is immediately discarded/overwritten. No permanent copy is created in any database or file storage.

### C. Separation of Components (The "Separation" Clause)
* **Rule:** Developers must not "separate, isolate, or modify the audio or video components" (e.g., ripping audio for a music player).
* **Compliance Status:** **Acceptable Interpretation.**
* **Reasoning:** While the application isolates the text track for processing, it does not **display** or **distribute** this isolated component to the end user. The isolated text is an internal data dependency used to generate a wholly new asset (the visual brief), not a consumer-facing product that replaces the video experience.

---

## 3. Intellectual Property (Copyright Law)
*Governed by US Copyright Act (Title 17) and international equivalents (Berne Convention).*

### A. Derivative Works
* **Legal Concept:** A "derivative work" is any work based on one or more preexisting works (e.g., a translation, summary, or adaptation).
* **Assessment:** The Sketchnote is technically a **derivative work** of the original video content.
* **Defense:** The creation of this derivative work is defensible under the doctrine of **Fair Use**.

### B. Fair Use Analysis (The Four Factors)
1.  **Purpose and Character of the Use:**
    * *Argument:* The use is **Transformative**. It shifts the medium from a time-based audiovisual format to a static, spatial visual format (Sketchnote). It adds new aesthetics and organization rather than just copying.
2.  **Nature of the Copyrighted Work:**
    * *Argument:* The original works are factual/educational videos (stronger Fair Use case) rather than purely creative fiction (weaker Fair Use case).
3.  **Amount and Substantiality:**
    * *Argument:* The application extracts only "Main Points" (a small fraction of the total expression), not the full script or detailed plot.
4.  **Effect on the Potential Market:**
    * *Argument:* **(Critical Factor)** The Sketchnote does not serve as a "Market Substitute." Viewing the Sketchnote does not replace the experience of watching the video; it acts as a companion or teaser. It does not cannibalize the creator's ad revenue or view counts.

---

## 4. Risk Mitigation Strategy

To maintain this legal standing, the application adheres to the following protocols:

* **Ephemeral Processing:** No database storage of transcripts.
* **Attribution & Link:** Every Sketchnote includes a clear citation: *"Based on [Video Title] by [Channel Name]"* and includes the original URL/QR Code.
* **Traffic Driver:** The output is designed to drive interest *back* to the source video, aligning the application's incentives with the content creator's.

***

*Disclaimer: This document is a summary of legal concepts based on general principles of technology and intellectual property law. It does not constitute formal legal advice or create an attorney-client relationship. If commercial disputes arise, consult a qualified intellectual property attorney.*