<div align="center">
  <img class="logo" src="https://github.com/yalexaner/nag-meetings-go/raw/master/images/nag-meetings.png" width="200px" alt="Nag Meetings"/>
</div>

# Nag Meetings

A configurable Telegram bot that manages subscriptions for meeting notifications and automatically sends meeting links to subscribers by parsing calendar HTML pages.

## Features

- 🤖 Telegram bot for user interaction
- 👥 User authorization and subscription management
- 🔍 Configurable HTML parser for any calendar system
- ⏰ Scheduled parsing of calendar events via cron
- 📤 Automated meeting link distribution to subscribers
- 🛡️ Admin controls for user authorization

## How it works

1. **User Commands**:
   - `/start`: Register with the bot (requires admin authorization)
   - `/subscribe`: Subscribe to meeting notifications
   - `/unsubscribe`: Unsubscribe from notifications
   - `/admin`: (Admin only) Authorize or reject new users

2. **Automated Notifications**:
   - A cron job runs every weekday at 10:20 AM (UTC+5)
   - Parses the specified calendar URL for events
   - Extracts meeting links based on your configuration
   - Sends links to all subscribed users

3. **Authorization Flow**:
   - New users must use `/start` to register
   - Admin receives authorization request
   - Admin approves or rejects via inline buttons
   - Only authorized users can subscribe

## Setup

### Prerequisites

- Go 1.22 or higher
- A Telegram Bot Token (get one from [@BotFather](https://t.me/botfather))
- A publicly accessible calendar URL with HTML content

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yalexaner/nag-meetings-go.git
cd nag-meetings-go
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file (see Configuration section below)

4. Build and run:
```bash
go build -o nag-meetings
./nag-meetings
```

## Configuration

Create a `.env` file in the project root with the following variables:

### Required Configuration

```bash
# Telegram Bot Token from @BotFather
TELEGRAM_BOT_TOKEN=your_bot_token_here

# URL of the calendar page to parse
CALENDAR_URL=https://example.com/calendar

# Directory where the database will be stored (include trailing slash)
WORKING_DIRECTORY=/path/to/directory/
```

### Optional Parser Configuration

Customize these to match your calendar's HTML structure:

```bash
# CSS selector for event container elements (default: .b-content-event)
PARSER_EVENT_SELECTOR=.calendar-event

# CSS selector for event title within the event container (default: h1)
PARSER_TITLE_SELECTOR=.event-title

# Meeting name to match (default: STB Daily Meeting)
PARSER_MEETING_NAME=Daily Standup

# CSS selector for links within events (default: .e-description a)
PARSER_LINK_SELECTOR=.event-link a

# HTML attribute to extract the meeting URL from (default: href)
PARSER_LINK_ATTRIBUTE=href
```

### Optional Settings

```bash
# Set to "debug" to use local index.html instead of fetching URL
ENVIRONMENT=debug
```

## Configuration Examples

### Example 1: Default Configuration (Original Use Case)

If your calendar HTML looks like this:
```html
<div class="b-content-event">
  <h1>STB Daily Meeting</h1>
  <div class="e-description">
    <a href="https://meet.google.com/abc-defg">Join Meeting</a>
  </div>
</div>
```

No parser configuration needed! The defaults will work.

### Example 2: Custom Calendar System

If your calendar HTML looks like this:
```html
<article class="calendar-event">
  <h2 class="event-title">Team Standup</h2>
  <div class="event-details">
    <a class="meeting-link" href="https://zoom.us/j/123456">Join</a>
  </div>
</article>
```

Configure your `.env`:
```bash
PARSER_EVENT_SELECTOR=.calendar-event
PARSER_TITLE_SELECTOR=.event-title
PARSER_MEETING_NAME=Team Standup
PARSER_LINK_SELECTOR=.event-details .meeting-link
PARSER_LINK_ATTRIBUTE=href
```

### Example 3: Data Attributes

If the meeting link is in a data attribute:
```html
<div class="event">
  <span class="title">Daily Meeting</span>
  <button data-meeting-url="https://teams.microsoft.com/l/meetup/...">Join</button>
</div>
```

Configure your `.env`:
```bash
PARSER_EVENT_SELECTOR=.event
PARSER_TITLE_SELECTOR=.title
PARSER_MEETING_NAME=Daily Meeting
PARSER_LINK_SELECTOR=button
PARSER_LINK_ATTRIBUTE=data-meeting-url
```

## Debugging

To test your parser configuration:

1. Download the calendar HTML page:
```bash
curl "YOUR_CALENDAR_URL" > index.html
```

2. Set debug mode in `.env`:
```bash
ENVIRONMENT=debug
```

3. Run the bot - it will parse the local `index.html` file instead of fetching the URL

4. Check logs for parsing results

## Setting Up Admin User

1. Run the bot for the first time
2. Find the `main.db` file in your `WORKING_DIRECTORY`
3. Manually set the first user as admin:

```bash
sqlite3 /path/to/main.db
UPDATE users SET authorized = 1, admin = 1 WHERE user_id = YOUR_TELEGRAM_ID;
.exit
```

Replace `YOUR_TELEGRAM_ID` with your actual Telegram user ID (you can get this from [@userinfobot](https://t.me/userinfobot))

## Cron Schedule

The default schedule is **Monday-Friday at 10:20 AM (UTC+5)**.

To change it, edit `main.go` line 30:
```go
_, err = c.AddFunc("20 10 * * 1-5", func() {
```

Cron format: `minute hour day month weekday`
- `20 10 * * 1-5`: 10:20 AM, Monday-Friday
- `0 9 * * *`: 9:00 AM, every day
- `30 14 * * 1,3,5`: 2:30 PM, Monday/Wednesday/Friday

## Technologies Used

- **Go 1.22** - Primary language
- **telegram-bot-api** - Telegram bot integration
- **goquery** - HTML parsing (jQuery-like)
- **robfig/cron** - Scheduled job execution
- **SQLite3** - Local database for user management
- **godotenv** - Environment configuration

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
