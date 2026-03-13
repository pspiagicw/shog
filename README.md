# shog

`shog` is a personal portfolio served as an SSH application.

Instead of sending someone to a traditional website, `shog` opens a small terminal-native experience built with [Wish](https://github.com/charmbracelet/wish) and [Bubble Tea](https://github.com/charmbracelet/bubbletea). The current version focuses on a simple, readable flow: an introduction page, project showcase pages, and a foundation that can later grow into blog and contact views.

## Why

Personal sites tend to collapse into the same visual template. `shog` takes a different route: lightweight, interactive, and intentionally text-first.

The goal is not novelty for its own sake. It is to create a portfolio that:

- feels personal instead of generic
- works naturally in a terminal
- keeps the content front and center
- can expand into projects, writing, and contact information over time

## Current Features

- SSH server powered by Wish
- Terminal UI built with Bubble Tea and Lip Gloss
- Multi-page navigation for intro and project highlights
- Simple keyboard controls for moving between pages
- Minimal project structure that is easy to customize

## Getting Started

### Prerequisites

- Go 1.25+

### Run locally

```bash
go run .
```

By default, the SSH server listens on `0.0.0.0:23234`.

In another terminal:

```bash
ssh localhost -p 23234
```

The app will generate a local SSH host key on first run. That key is ignored by git and should remain local to the machine running the server.

## Configuration

`shog` currently supports two environment variables:

- `SHOG_HOST` sets the listen host
- `SHOG_PORT` sets the listen port

Example:

```bash
SHOG_HOST=127.0.0.1 SHOG_PORT=23234 go run .
```

## Controls

- `left`, `h`, `p`: previous page
- `right`, `l`, `n`: next page
- `q` or `ctrl+c`: quit

## Project Structure

```text
.
├── internal/app/    Bubble Tea model and page rendering
├── main.go          Wish server entrypoint
├── TODO.md          Planned roadmap and next steps
├── makefile         Basic build and formatting commands
└── build-release.sh Release packaging helper
```

## Roadmap

Near-term priorities:

- replace placeholder content with real bio and project details
- improve styling for a wider range of terminal sizes
- add a contact-oriented page or footer
- document deployment for serving the app publicly over SSH

Planned future direction:

- blog index and post navigation
- richer project metadata and layout
- tests for navigation and rendering
- CI and release automation

## Customization

Most of the experience lives in `internal/app/model.go`. That is where page content, navigation behavior, and presentation are currently defined.

For a personal deployment, the usual next steps are:

1. Replace the placeholder page content with your actual introduction.
2. Add real projects with sharper summaries and outcomes.
3. Adjust styling to match your tone.
4. Deploy the app on a host with port `22` or another SSH-accessible port.

## License

This project is available under the terms of the [MIT License](./LICENSE).
