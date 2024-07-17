<div align="center">
  <img class="logo" src="https://github.com/yalexaner/nag-meetings-go/raw/master/images/nag-meetings.png" width="200px" alt="Nag Meetings"/>
</div>

# Nag Meetings

This project runs a Telegram bot that manages subscriptions for meeting notifications and sends meeting links to subscribers.

## Features

- Telegram bot for user interaction
- User subscription management (subscribe/unsubscribe)
- Scheduled parsing of calendar events
- Automated meeting link distribution to subscribers

## How it works

1. The bot allows users to subscribe or unsubscribe using commands:
   - `/subscribe`: Adds the user's Telegram ID to the database
   - `/unsubscribe`: Removes the user's Telegram ID from the database

2. A cron job runs every weekday at 10:20 AM (UTC+5) to:
   - Parse the specified calendar URL for events
   - Fetch the meeting URL for the current day

3. If a meeting URL is found, the bot sends it to all subscribed users

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
