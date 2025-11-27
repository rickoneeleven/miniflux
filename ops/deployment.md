# Miniflux deployment: rss.pinescore.com

DATETIME of last agent review: 27/11/2025 14:55 GMT

## Host and runtime
- OS: Debian 12 (bookworm) with Virtualmin/Apache managing `rss.pinescore.com`.
- Application code lives in `/home/loopnova/domains/rss.pinescore.com/public_html` and tracks `main` from `git@github.com:rickoneeleven/miniflux.git`.
- Go toolchain: `/usr/local/go` (go1.25.x) installed from upstream tarball; Debian `golang-go` (go1.19) is present but not used for builds.

## PostgreSQL
- Version: PostgreSQL 15 (`postgresql`/`postgresql-client` from Debian packages).
- Service: managed by `systemd` as `postgresql.service`.
- Database: `miniflux` (UTF8), owner `miniflux`.
- Role: `miniflux` with a strong password; the current value is stored in `/etc/miniflux.conf` via `DATABASE_URL` and should not be committed to version control.
- Connection string pattern: `postgres://miniflux:<password>@localhost/miniflux?sslmode=disable`.

## Application binary
- Source directory: `/home/loopnova/domains/rss.pinescore.com/public_html`.
- Build command (run with `/usr/local/go/bin` first on `PATH`): `make miniflux`.
- Installed binary: `/usr/local/bin/miniflux` (copied from the repo root after build).
- To rebuild and redeploy on this host:
  - `cd /home/loopnova/domains/rss.pinescore.com/public_html`
  - `git pull`
  - `PATH=/usr/local/go/bin:$PATH make miniflux`
  - `sudo install -o root -g root -m 0755 miniflux /usr/local/bin/miniflux`
  - `sudo systemctl restart miniflux`

## Miniflux configuration
- Config file: `/etc/miniflux.conf` (mode `600`, owned by `root`).
- Key options currently set:
  - `DATABASE_URL=postgres://miniflux:<password>@localhost/miniflux?sslmode=disable`
  - `RUN_MIGRATIONS=1`
  - `LISTEN_ADDR=127.0.0.1:8180`
  - `CREATE_ADMIN=1`
  - `ADMIN_USERNAME=admin`
  - `ADMIN_PASSWORD=<initial admin password>`
  - `ADMIN_EMAIL=admin@rss.pinescore.com`
  - `BASE_URL=https://rss.pinescore.com`
- Recommended hardening after first login:
  - Change the admin password in the UI (already done).
  - Edit `/etc/miniflux.conf` to remove `ADMIN_PASSWORD` and set `CREATE_ADMIN=0`.
  - `sudo systemctl restart miniflux`.

## Systemd service
- Unit file: `/etc/systemd/system/miniflux.service` (copied from `packaging/systemd/miniflux.service` and adjusted).
- Key settings:
  - `ExecStart=/usr/local/bin/miniflux`
  - `User=miniflux` (created via `useradd --system --home /var/lib/miniflux --shell /usr/sbin/nologin miniflux`)
  - `EnvironmentFile=/etc/miniflux.conf`
  - `Type=notify`, `Restart=always`, `RuntimeDirectory=miniflux`
- Management commands:
  - `sudo systemctl status miniflux`
  - `sudo systemctl restart miniflux`
  - `sudo journalctl -u miniflux`

## Apache / Virtualmin integration
- Vhost file: `/etc/apache2/sites-available/rss.pinescore.com.conf` (backed up as `rss.pinescore.com.conf.bak-<timestamp>` before changes).
- Both HTTP (`VirtualHost 192.168.1.206:80`) and HTTPS (`VirtualHost 192.168.1.206:443`) blocks include:
  - `ProxyPreserveHost On`
  - `ProxyPass / http://127.0.0.1:8180/ retry=1 acquire=3000 timeout=600 Keepalive=On`
  - `ProxyPassReverse / http://127.0.0.1:8180/`
  - `RequestHeader set X-Forwarded-Proto "https" env=HTTPS`
  - `ProxyPass /.well-known !` (existing Let’s Encrypt/Virtualmin handling, retained)
- Apache modules required and enabled:
  - `proxy`
  - `proxy_http`
  - `headers`
- After editing the vhost:
  - Validate config: `sudo apachectl configtest`
  - Reload Apache: `sudo systemctl reload apache2`

## Update / recovery notes
- If Apache vhost changes break the site:
  - Restore from backup: copy the latest `rss.pinescore.com.conf.bak-<timestamp>` over `rss.pinescore.com.conf`.
  - Run `sudo apachectl configtest` and `sudo systemctl reload apache2`.
- If Miniflux stops responding but the service is running:
  - Check logs: `sudo journalctl -u miniflux`.
  - Verify listener: `ss -tlnp | grep 8180`.
  - Ensure PostgreSQL is healthy: `sudo systemctl status postgresql`.

