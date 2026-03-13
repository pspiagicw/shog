+++
title = "Bringing Projects into the Binary"
date = "2026-03-01"
tags = ["build", "content", "github"]
description = "Why vendoring project README files at build time is a good fit for this app."
+++

# Bringing Projects into the Binary

The Projects tab is strongest when it is built from real project material instead of parallel hand-maintained summaries.

Vendoring repository README files during the build keeps the content source close to the project itself while still producing a single binary that can be deployed cleanly.

## Benefits

- Content stays close to the source repo.
- The app deploys as one artifact.
- Updating the portfolio becomes a content refresh rather than a UI rewrite.

That same pattern works well for blog content too, especially when the writing already exists in a Hugo content tree.
