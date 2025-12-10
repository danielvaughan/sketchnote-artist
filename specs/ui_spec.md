# UI Specification: Sketchnote Artist (YouTube Inspired Minimalist)

## Overview
This document specifies the UI design for Sketchnote Artist. The design philosophy is **Minimalist & Premium**, drawing inspiration from YouTube's "Dark Mode" aesthetic (colors, typography, rounded corners) without attempting to clone YouTube's functionality or layout complexity. The focus is strictly on the core workflow: Input Video -> View Sketchnote.

## Design System

### Color Palette (YouTube Light Mode)
- **Background**: `#F9F9F9` (Main Background - Very light gray)
- **Surface**: `#FFFFFF` (Input fields, Cards)
- **Text Primary**: `#0F0F0F` (Main text)
- **Text Secondary**: `#606060` (Secondary text, icons)
- **Accent**: `#FF0000` (YouTube Red)
- **Border/Separator**: `rgba(0, 0, 0, 0.1)`

### Typography
- **Font**: 'Roboto', Arial, sans-serif.
- **Weights**: Regular (400), Medium (500), Bold (700).
- **Title Styling**: Ensure the app title is visually distinct but not overwhelming. Tidy alignment with the search bar.

### Visual Language
- **Rounded Corners**: Generous border-radius (e.g., 24px or 30px) for inputs/buttons.
- **Shadows**: Subtle drop shadows (e.g., `box-shadow: 0 1px 5px rgba(0,0,0,0.1)`) to separate white surfaces from the light gray background.

## Layout Structure

The layout is a clean, single-column view centered on the screen. No sidebars, no complex headers.

### 1. Header (Minimal)
- **Content**: App Logo/Title ("Sketchnote Artist").
- **Style**: Centered or Top-Left. Clean, unobtrusive.

### 2. Main Action Area (Centered)
- **Input**: A prominent, aesthetically pleasing input field for the YouTube URL.
    -   Style: Dark gray background (`#272727`), rounded corners (pill-shaped), subtle border.
    -   Placeholder: "Paste YouTube Link..."
- **Action Button**:
    -   Style: Simple icon or distinct button next to or inside the input.
    -   Color: Accent Red or White on dark.

### 3. Display Stage
- **Container**: A dedicated area to display the result.
- **State: Idle**: Hidden or shows a subtle tagline/illustration.
- **State: Processing**:
    -   Visual feedback indicating activity (e.g., a pulsating glow, a simple progress bar, or a loading skeleton).
    -   Text feedback: "Watching video...", "Sketching..."
- **State: Result**:
    -   The generated Sketchnote image displayed beautifully.
    -   Rounded corners on the image.
    -   **Actions**: Simple "Download" (Icon) and "Close/Reset" (Icon) buttons floating or positioned neatly below.

## Functional Requirements
-   **No Menus**: No hamburger menus, no user profile icons, no notification bells.
-   **No "Fake" Buttons**: Do not include "Subscribe", "Like", "Dislike" unless they perform a real function within *this* app (which they don't).
-   **Focus**: All visual weight should be on the Input (start) and the Image (end).

