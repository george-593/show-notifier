# Show Notifier

A Go application that tracks your TV shows and sends Telegram notifications when new episodes are released. Shows and episode data are sourced from the [TVmaze API](https://www.tvmaze.com/api).

## Features

- Telegram notifications when new episodes air
- Telegram bot interface — ⁠/add, ⁠/remove, ⁠/list, ⁠/upcoming, ⁠/check
- Local CLI for managing shows when running locally
- Persistent JSON store with duplicate notification prevention
- Scheduled background checks via goroutines

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | Your Telegram bot token | Yes |
| `TELEGRAM_CHAT_ID` | Your Telegram chat ID | Yes |
| `MODE` | Set to `headless` to disable the CLI menu (recommended for server deployment) | No |
| `STORE_PATH` | The path that the json store will be at (include.json) | Yes |

## Telegram Commands

| Command | Description |
|---|---|
| `/add <show name>` | Search for and add a show to your watchlist |
| `/list` | List all tracked shows |
| `/remove <show name>` | Remove a show from your watchlist |
| `/upcoming` | Show episodes airing in the next 7 days |
| `/check` | Manually trigger a check for new episodes |

## Deployment

### Install Go
```bash
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Clone and Build
```bash
git clone https://github.com/george-593/show-notifier.git
cd show-notifier
nano .env
go build -o show-notifier .
./show-notifier
```

### Set up as a systemd service
```bash
sudo nano /etc/systemd/system/show-notifier.service
```

```ini
[Unit]
Description=Show Notifier
After=network.target

[Service]
Type=simple
User=show-notifier
WorkingDirectory=/home/show-notifier/show-notifier
EnvironmentFile=/home/show-notifier/show-notifier/.env
ExecStart=/home/show-notifier/show-notifier/show-notifier
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable show-notifier
sudo systemctl start show-notifier
```

### Check logs
```bash
journalctl -u show-notifier -f
```

# Future Improvements 
- Fetch from updates endpoint every 24hr to get any show updates, if detected not ran in last 24 hr revert to week, then month, then full rebuild