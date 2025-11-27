# Session Summary: Miniflux on rss.pinescore.com

Date: 27/11/2025 17:10 GMT

## State after this session
- Miniflux is deployed at `/home/loopnova/domains/rss.pinescore.com/public_html` with Go 1.25.x, PostgreSQL 15, and `miniflux.service` behind Apache on `https://rss.pinescore.com`.
- Feed polling is configured via `/etc/miniflux.conf` with `POLLING_FREQUENCY=10` (seconds) and `SCHEDULER_ROUND_ROBIN_MIN_INTERVAL=1` (minute); scheduler ticks every ~10s and feeds are effectively refreshed about once per minute.
- The `/unread` page uses a 10-second JSON snapshot poller so unread counts and “Last fetch” update live without a full page reload.

## Next steps (front page behaviour)
- On `/unread`, after new items are fetched, the header can show e.g. `Unread (2)`, `0 unread entry`, `Last fetch: <timestamp>`, plus “There are no unread entries.” until the page is manually reloaded; counts update, but the entry list does not.
- Next session should extend the existing `/unread` poller to refresh the visible entry list (re-render the `.items` container from a JSON or HTML endpoint) while preserving scroll position and keyboard focus, so new articles appear automatically as they arrive.
