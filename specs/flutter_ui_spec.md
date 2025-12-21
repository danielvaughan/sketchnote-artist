# UI Specification: Sketchnote Artist (Flutter Edition)

## Overview

This document specifies the technical design for a Flutter-based frontend for the Sketchnote Artist application. The goal is to provide a "pixel-perfect" alternative to the web UI, utilizing **Flutter 3.38.3+** and the `Material 3` design system to maintain the **Minimalist & Premium** aesthetic.

## 🎨 Design System (ThemeData)

| Property | Value | Description |
| :--- | :--- | :--- |
| **Brightness** | `Brightness.light` | Matches YouTube Light Mode |
| **ColorPrimary** | `#FF0000` | YouTube Red |
| **ScaffoldBackground** | `#F9F9F9` | Main Background |
| **CardColor / Surface** | `#FFFFFF` | Input fields, Result cards |
| **FontFamily** | `Roboto` | Via `google_fonts` package |
| **BorderRadius** | `BorderRadius.circular(24)` | Pill shape for inputs |

## 📐 Layout Structure

The app uses a `SingleChildScrollView` to handle varying screen heights, with content centered using `ConstrainedBox` to limit maximum width to `900px`.

### 📱 Responsive Design (Breakpoints)

- **Mobile (< 900px)**: Vertical stack (Input -> Video -> Sketchnote).
- **Desktop (>= 900px)**: Single column with side-by-side video and image in a `Row`.

---

## 🏗️ Component Mapping

### 1. Header (`AppBar` inspired)

- **Leading**: `Icon(Icons.gesture, color: Colors.red, size: 32)`
- **Title**: `Text("Sketchnote Artist", style: TextStyle(fontWeight: FontWeight.w700))`

### 2. Search Section (`SketchnoteInput`)

- **Widget**: `TextField` inside a `Container` with `BoxShadow`.
- **Decoration**: `InputDecoration` with `OutlineInputBorder` (Radius 24), no border side (shadow provided by parent).
- **Suffix**: `IconButton` with `Icons.play_arrow_rounded` (Red on Hover).

### 3. Display Pane (`dual_pane.dart`)

- **Video Player**: Uses `youtube_player_flutter` or `webview_flutter` for the iframe embed.
- **Sketchnote Card**: `Container` with `ClipRRect` (Radius 12) and `Image.network` for the generated result.

### 4. Actions Overlay (`ResultActions`)

- **Positioned**: Bottom Right of the Sketchnote Card.
- **Widgets**: `FloatingActionButton.small` (Transparent white background, backdrop blur filter).
- **Icons**: `Icons.download`, `Icons.share`, `Icons.refresh`.

---

## ⚙️ Logic & State Management

### 1. Identity & Auth

- At startup (`initState`), call `GET /me` (via `http` package).
- Store `userId` in a global/controller state.

### 2. Session & Running

- `POST /apps/sketchnote-artist/users/$userId/sessions` to initialize.
- Use `fetch_client` or a custom stream wrapper for **SSE (Server-Sent Events)** streaming from `POST /run_sse`.

### 3. Streaming Events

Map SSE data lines to a `SketchnoteController`:

- `event: status` -> Update a `statusMessage` string.
- `data: { content: { parts: [...] } }` -> Parse for `.png` filename regex.
- `event: done` -> Finalize UI state.

---

## 📦 Dependencies

- `google_fonts`: For Roboto typography.
- `http`: For standard REST calls.
- `flutter_client_sse` (or similar): For handling server-sent events.
- `youtube_player_flutter`: For native-like YouTube playback.
- `path_provider` & `dio`: For image download functionality.
- `share_plus`: For native sharing capabilities.
