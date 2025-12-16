# Feature: Dual Player UI & Streaming Improvements

## Overview
This PR introduces the "Dual Player" layout, allowing users to watch the YouTube video alongside the generated sketchnote. It also includes critical fixes for the streaming infrastructure and security compliance for embedded playback.

## Changes
### 🎨 UI & UX
- **Dual Player Layout**: Side-by-side grid layout (`web/index.html`, `web/style.css`) showing:
  - Left: YouTube Video Player
  - Right: Sketchnote Result / Progress
- **Responsive Design**: Stacks vertically on smaller screens.
- **Harmonized Design**: Unified background colors using `--player-bg` (#000) for a seamless look.
- **Smart Visibility**: The player stage is hidden (`hidden` class) until a video is submitted (`web/app.js`).

### 🛠️ Functionality & Fixes
- **Streaming Fix**: Updated frontend event parser (`web/app.js`) to correctly handle ADK's nested `content.parts[].functionCall` structure.
- **Payload Fix**: Added `streaming: true` to the request API payload to ensure proper SSE behavior.
- **Security Compliance**:
  - Added `referrerpolicy="origin"` to the YouTube iframe.
  - Appended `origin=${window.location.origin}` to the embed URL to satisfy YouTube's security checks.

## Verification
- **E2E Tests Certified**:
  - `e2e/api.spec.ts`: Passed (Verified streaming payload and image generation).
  - `e2e/webui.spec.ts`: Passed (Verified UI interactions on localhost).
- **Manual Review**: Verified "Me at the zoo" flow visually.

## Screenshots
*(See Walkthrough artifact for visual verification)*

## Tasks Closed
- Implemented Dual Player UI
- Fixed Streaming Events
- Hardened Embed Security
