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
- Three tabs: About, Projects, and Blog
- Project browsing within the Projects tab
- Keyboard controls for switching tabs, browsing projects, and scrolling
- Minimal project structure that is easy to customize

## Getting Started

### Prerequisites

- Go 1.25+

### Run locally

```bash
./scripts/fetch-project-readmes.sh
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

Application content is file-backed:

- `content/app.toml` controls the app title and blog placeholder
- `content/about.md` provides the About tab body
- `content/projects.urls` lists GitHub repositories to vendor into the Projects tab
- `content/blogs/` holds markdown blog posts, including Hugo-style content files

## Controls

- `left`, `h`, `tab`: previous tab
- `right`, `l`, `shift+tab`: next tab
- `1`, `2`, `3`: jump directly to About, Projects, or Blog
- `n`, `p`: next and previous project or blog post inside the active content tab
- `up`, `k`: scroll up
- `down`, `j`: scroll down
- `space`, `f`: page down
- `b`: page up
- `g`: top
- `G`: bottom
- `q` or `ctrl+c`: quit

## Project Structure

```text
.
├── app/                 Bubble Tea model and content loading
├── content/             Embedded app content and vendored project READMEs
├── main.go              Wish server entrypoint
├── scripts/             Content fetch helpers
├── TODO.md              Planned roadmap and next steps
└── makefile             Basic build and formatting commands
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

Most of the behavior lives in `app/model.go`, while the `content/` directory drives the visible content.

For a personal deployment, the usual next steps are:

1. Replace `content/about.md` with your actual introduction.
2. Edit `content/app.toml` to set the app title and blog placeholder text.
3. Add repository URLs to `content/projects.urls`, then run `./scripts/fetch-project-readmes.sh`.
4. Drop Hugo markdown posts into `content/blogs/`.
5. Adjust styling in `app/model.go` to match your tone.
6. Deploy the app on a host with port `22` or another SSH-accessible port.

## License

This project is available under the terms of the [MIT License](./LICENSE).
