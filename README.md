# Todo
- Scheduler
- Update from API when episodes run out
- Docker deployment on server
- Full interaction via telegram instead of CLI
- Check updates endpoint for any updates every 24hr (contingency if we detect a missed update?)

# AI Generated Todo
Quality of Life
	•	Duplicate notification guard across restarts — ⁠NotifiedIDs grows forever and never gets cleaned up. Old episode IDs from years ago are being stored pointlessly. Prune any IDs older than 30 days on startup
    •	Show status awareness — if a show's ⁠Ended field is set, stop checking it and notify you via Telegram that the show has ended so you can decide whether to remove it
    •	Next episode display — in your ⁠/list Telegram command, show the next upcoming episode and its air date, not just the show name
Robustness
    •	Logging — replace your ⁠fmt.Printf error prints with proper structured logging using Go's ⁠log/slog package (added in Go 1.21). Useful when running headless in Docker
    •	Graceful shutdown — listen for OS signals (⁠SIGTERM, ⁠SIGINT) and save the store cleanly before exiting rather than just dying mid-write
Features
    •	Telegram inline keyboards — when adding a show via Telegram, show search results as buttons rather than typing y/n	•	Episode countdown — notify you 24 hours before an episode airs as well as when it drops
    •	Season premiere alerts — flag when it's a season premiere specifically, not just any episode