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

## Component Specifications

### 1. Main Layout (`.centered-layout`)
-   **Width**: Max `900px`, centered.
-   **Padding**: `2rem`.
-   **Gap**: `2.5rem` between elements.

### 2. Header (`.minimal-header`)
-   **Logo**: Material Icon `gesture` (Size `2.5rem`, Color `#FF0000`).
-   **Title**: "Sketchnote Artist" (Size `1.75rem`, Bold 700, Tracking `-0.5px`).

### 3. Input Section (`.input-section`)
-   **Width**: Max `640px`.
-   **Field** (`.search-box`):
    -   Background: `#FFFFFF`.
    -   Border Radius: `24px` (Pill).
    -   Shadow: `0 1px 4px rgba(0, 0, 0, 0.05)`.
    -   Focus State: Blue border `#1C62B9`.
-   **Button** (`#generateBtn`):
    -   Size: `48x48px` Circle.
    -   Background: `#F0F0F0` (Hover: `#E5E5E5`).

### 4. Display Stage (`.display-stage`)
-   **Structure**: Side-by-side Grid (Dual Pane).
-   **Columns**: 1fr 1fr (Equal width).
-   **Gap**: `2rem`.
-   **Responsive**: Stacks vertically on smaller screens (< 900px).

#### Left Pane: Video Player
-   **Content**: Embedded YouTube Player (iframe).
-   **Aspect Ratio**: 16:9.
-   **Styling**: Matches Sketchnote container (Rounded corners, Shadow).

#### Right Pane: Sketchnote Result (`.player-container`)
-   **Content**: Generated Image.
-   **Aspect Ratio**: 16:9.
-   **Background**: Black `#000`.
-   **Border Radius**: `12px`.
-   **Shadow**: `0 4px 12px rgba(0, 0, 0, 0.15)`.

### 5. States
1.  **Idle** (`#idleState`): Shows large `smart_display` icon.
2.  **Progress** (`#progressSection`):
    -   Spinner: `48px`, Top border Red `#FF0000`.
    -   Text: "Watching video...", "Summarizing...", "Sketching...".
3.  **Result** (`#resultSection`):
    -   Image: Contains `img` tag (object-fit: contain).
    -   **Actions Overlay** (`.result-actions`):
        -   Position: Bottom Right (`24px`).
        -   Buttons: Download, Close, Share (White circle `44px`, Backdrop Blur).
        -   Interaction: Fade in on hover.

## Technical Implementation
-   **Icons**: Google Material Icons.
-   **Fonts**: Google Fonts (Roboto).
-   **CSS Variables**: Defined in `:root` for consistency.

