# slack

The `slack` module supports easily creating and saving Slack clients.

Wraps the [slack-go](https://github.com/slack-go/slack) library.

## Usage

```risor
// Create a new Slack client
client = slack.new("xoxb-your-token-here")

// Send a simple message
client.post_message("C12345678", "Hello, world!")

// Send a message with options
client.post_message("C12345678", {
    "text": "Hello with attachments",
    "thread_ts": "1234567890.123456",
    "attachments": [
        {
            "title": "Attachment Title",
            "text": "Attachment text",
            "color": "#36a64f"
        }
    ]
})

// Build rich messages with blocks
builder = slack.message_builder(client)
builder.add_header("Important Announcement")
builder.add_section("Hello everyone! This is an important message.")
builder.add_divider()
builder.add_section("Please read carefully.")
builder.send("C12345678")

// Create and use a message
msg = slack.message(client, "C12345678")
msg.reply("This is a threaded reply")

// Get user info
user = client.get_user_info("U12345678")
print(user["real_name"])

// Get channels
channels = client.get_channels()
```

## Functions

### `slack.new(token)`

Creates a new Slack client with the given token.

### `slack.message(client, [channel])`

Creates a new Slack message object associated with the given client and optional channel.

### `slack.message_builder(client)`

Creates a new Slack message builder for creating rich messages with blocks.

## Slack Client Methods

- `get_user_groups([include_users])`: Get all user groups for the team
- `get_user_info(user_id)`: Get information about a user
- `get_billable_info([user_id])`: Get billable info for the team
- `post_message(channel_id, text, [options])`: Send a message to a channel
- `upload_file(content, channels, [options])`: Upload a file to Slack
- `get_channels([options])`: Get conversations for a user

## Slack Message Methods

- `reply(text)`: Reply to the message in a thread

## Slack Message Builder Methods

- `add_section(text)`: Add a section block with text
- `add_divider()`: Add a divider block
- `add_header(text)`: Add a header block
- `build()`: Build the message as a map
- `send(channel, [text])`: Send the message to a channel
