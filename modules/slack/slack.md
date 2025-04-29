# slack

The `slack` module supports easily creating and sending messages to Slack channels, managing conversations, and interacting with users and user groups.

## Module

```go copy filename="Function signature"
slack.new_client(token string) slack.client
```

Initialize a new Slack client with an API token.

```go copy filename="Example"
>>> client := slack.new_client("xoxb-your-token-here")
>>> client.post_message("#general", "Hello from Risor!")
{channel: "#general", timestamp: "1710000000.000000"}
```

## Client

The Slack client provides methods for interacting with the Slack API.

### post_message

```go filename="Method signature"
post_message(channel_id string, message string|map) map
```

Sends a message to a Slack channel.

```go filename="Example"
>>> client.post_message("#general", "Hello world!")
{channel: "#general", timestamp: "1710000000.000000"}

>>> client.post_message("#general", {
...   text: "Message with attachment",
...   attachments: [
...     {
...       title: "Attachment Title",
...       text: "Attachment text here",
...       color: "#36a64f"
...     }
...   ]
... })
{channel: "#general", timestamp: "1710000000.000000"}
```

### post_ephemeral_message

```go filename="Method signature"
post_ephemeral_message(channel_id string, user_id string, message string|map) map
```

Sends an ephemeral message to a Slack channel that is only visible to a specific user.

```go filename="Example"
>>> client.post_ephemeral_message("#general", "U123456", "Only you can see this message")
{channel: "#general", timestamp: "1710000000.000000", user: "U123456"}
```

### update_message

```go filename="Method signature"
update_message(channel_id string, timestamp string, message string|map) map
```

Updates an existing message in a Slack channel.

```go filename="Example"
>>> client.update_message("#general", "1710000000.000000", "Updated message content")
{channel: "#general", timestamp: "1710000000.000001"}
```

### delete_message

```go filename="Method signature"
delete_message(channel_id string, timestamp string)
```

Deletes a message from a Slack channel.

```go filename="Example"
>>> client.delete_message("#general", "1710000000.000000")
```

### add_reaction

```go filename="Method signature"
add_reaction(emoji string, item_ref map) map
```

Adds an emoji reaction to a message or other item.

```go filename="Example"
>>> client.add_reaction("thumbsup", {
...   channel: "#general",
...   timestamp: "1710000000.000000"
... })
{emoji: "thumbsup", added: true, channel: "#general", timestamp: "1710000000.000000"}
```

### remove_reaction

```go filename="Method signature"
remove_reaction(emoji string, item_ref map)
```

Removes an emoji reaction from a message or other item.

```go filename="Example"
>>> client.remove_reaction("thumbsup", {
...   channel: "#general",
...   timestamp: "1710000000.000000"
... })
```

### upload_file

```go filename="Method signature"
upload_file(channel_id string, options map) map
```

Uploads a file to a Slack channel.

```go filename="Example"
>>> client.upload_file("#general", {
...   content: "File content here",
...   filename: "example.txt",
...   title: "Example File"
... })
{id: "F123456", title: "Example File"}
```

### get_users

```go filename="Method signature"
get_users(options map={}) [slack.user]
```

Gets a list of all users in the workspace.

```go filename="Example"
>>> users := client.get_users()
>>> len(users)
42
>>> users[0].name
"johndoe"
```

### get_user_info

```go filename="Method signature"
get_user_info(user_id string) slack.user
```

Gets information about a specific user.

```go filename="Example"
>>> user := client.get_user_info("U123456")
>>> user.name
"johndoe"
>>> user.profile.email
"john@example.com"
```

### get_user_groups

```go filename="Method signature"
get_user_groups(options map={}) [map]
```

Gets all user groups in the workspace.

```go filename="Example"
>>> groups := client.get_user_groups({include_users: true})
>>> groups[0].name
"Engineering"
>>> groups[0].users
["U123456", "U234567"]
```

### get_conversation_info

```go filename="Method signature"
get_conversation_info(channel_id string) slack.channel
```

Gets information about a specific conversation.

```go filename="Example"
>>> channel := client.get_conversation_info("C123456")
>>> channel.name
"general"
>>> channel.is_private
false
```

### get_conversations

```go filename="Method signature"
get_conversations(options map={}) slack.conversation_iterator
```

Gets a list of all conversations in the workspace.

Options:
- `types`: An optional array of conversation types to include. Valid values are:
  - `public_channel`: Public channels that anyone in the workspace can join
  - `private_channel`: Private channels with restricted membership 
  - `mpim`: Multi-person direct messages (group DMs)
  - `im`: Direct messages between two users
- `exclude_archived`: Whether to exclude archived channels (default: false)

```go filename="Example"
>>> convs := client.get_conversations({
...   types: ["public_channel", "private_channel"],
...   exclude_archived: true
... })
>>> for i, conv in convs {
...   print(conv.name)
... }
general
random
team-project
```

### create_conversation

```go filename="Method signature"
create_conversation(name string, options map={}) slack.channel
```

Creates a new conversation (channel).

```go filename="Example"
>>> channel := client.create_conversation("new-project", {is_private: true})
>>> channel.name
"new-project"
>>> channel.is_private
true
```

### get_conversation_history

```go filename="Method signature"
get_conversation_history(channel_id string, options map={}) slack.message_iterator
```

Gets the message history of a conversation.

```go filename="Example"
>>> messages := client.get_conversation_history("#general", {limit: 10})
>>> for i, msg in messages {
...   print(msg.text)
... }
Hello world!
Important announcement
Check out this link: https://example.com
```

### get_conversation_members

```go filename="Method signature"
get_conversation_members(channel_id string, options map={}) slack.conversation_members_iterator
```

Gets the members of a conversation.

```go filename="Example"
>>> members := client.get_conversation_members("#general")
>>> for i, member_id in members {
...   user := client.get_user_info(member_id)
...   print(user.name)
... }
johndoe
janedoe
bobsmith
```

### message_builder

```go filename="Method signature"
message_builder() slack.message_builder
```

Creates a new message builder for constructing complex messages.

```go filename="Example"
>>> builder := client.message_builder()
>>> builder.add_text("Hello world!")
>>> builder.add_divider()
>>> builder.add_section("This is a section")
>>> msg := builder.build()
>>> client.post_message("#general", msg)
{channel: "#general", timestamp: "1710000000.000000"}
```

## Types

### slack.client

The Slack client represents a connection to the Slack API and provides methods for interacting with it.

### slack.channel

The channel object represents a Slack channel and provides access to its properties.

#### Properties

- `id` - The channel ID
- `name` - The channel name
- `is_channel` - Whether the channel is a standard channel
- `is_group` - Whether the channel is a group
- `is_im` - Whether the channel is a direct message
- `is_mpim` - Whether the channel is a multi-person direct message
- `created` - When the channel was created
- `creator` - The user ID of the channel creator
- `is_archived` - Whether the channel is archived
- `is_general` - Whether the channel is the general channel
- `unlinked` - When the channel was unlinked
- `name_normalized` - The normalized channel name
- `is_shared` - Whether the channel is shared
- `is_ext_shared` - Whether the channel is externally shared
- `is_org_shared` - Whether the channel is org shared
- `is_pending_ext_shared` - Whether the channel is pending external sharing
- `is_member` - Whether the current user is a member
- `is_private` - Whether the channel is private
- `is_open` - Whether the channel is open
- `topic` - The channel topic
- `purpose` - The channel purpose
- `members` - The channel members
- `num_members` - The number of members in the channel

#### Methods

##### json

```go filename="Method signature"
json() map
```

Returns a map representation of the channel object.

```go filename="Example"
>>> channel := client.get_conversation_info("C123456")
>>> channel_map := channel.json()
>>> channel_map.id
"C123456"
```

### slack.message

The message object represents a Slack message and provides access to its properties.

#### Properties

- `text` - The message text
- `channel` - The channel ID
- `timestamp` - The message timestamp
- `thread_timestamp` - The thread timestamp
- `user` - The user ID of the message sender
- `type` - The message type
- `subtype` - The message subtype
- `team` - The team ID
- `bot_id` - The bot ID if sent by a bot
- `username` - The username if sent by a bot
- `reactions` - The reactions to the message
- `is_bot_message` - Whether the message was sent by a bot
- `reply_count` - The number of replies to the message
- `latest_reply` - The timestamp of the latest reply

#### Methods

##### conversation

```go filename="Method signature"
conversation() [slack.message]
```

Gets the conversation thread for this message.

```go filename="Example"
>>> msg := messages[0]
>>> replies := msg.conversation()
>>> for reply in replies {
...   print(reply.text)
... }
This is a reply
Another reply
```

##### json

```go filename="Method signature"
json() map
```

Returns a map representation of the message object.

```go filename="Example"
>>> msg := messages[0]
>>> msg_map := msg.json()
>>> msg_map.text
"Hello world!"
```

### slack.user

The user object represents a Slack user and provides access to their properties.

#### Properties

- `id` - The user ID
- `team_id` - The team ID
- `name` - The username
- `deleted` - Whether the user has been deleted
- `color` - The user's color preference
- `real_name` - The user's real name
- `tz` - The user's timezone
- `tz_label` - The user's timezone label
- `tz_offset` - The user's timezone offset
- `profile` - The user's profile
- `is_admin` - Whether the user is an admin
- `is_owner` - Whether the user is an owner
- `is_primary_owner` - Whether the user is the primary owner
- `is_restricted` - Whether the user is restricted
- `is_ultra_restricted` - Whether the user is ultra restricted
- `is_bot` - Whether the user is a bot
- `is_app_user` - Whether the user is an app user
- `updated` - When the user was last updated
- `has_2fa` - Whether the user has two-factor authentication enabled
- `two_factor_type` - The user's two-factor authentication type
- `has_files` - Whether the user has files
- `presence` - The user's presence status
- `locale` - The user's locale setting
- `is_stranger` - Whether the user is a stranger
- `is_invited_user` - Whether the user was invited

#### Methods

##### json

```go filename="Method signature"
json() map
```

Returns a map representation of the user object.

```go filename="Example"
>>> user := client.get_user_info("U123456")
>>> user_map := user.json()
>>> user_map.id
"U123456"
```

### slack.user_profile

The user profile object represents a Slack user's profile and provides access to its properties.

#### Properties

- `real_name` - The user's real name
- `real_name_normalized` - The user's normalized real name
- `display_name` - The user's display name
- `display_name_normalized` - The user's normalized display name
- `email` - The user's email address
- `first_name` - The user's first name
- `last_name` - The user's last name
- `phone` - The user's phone number
- `skype` - The user's Skype handle
- `title` - The user's title
- `team` - The user's team
- `status_text` - The user's status text
- `status_emoji` - The user's status emoji
- `bot_id` - The user's bot ID
- `image_24` - URL to a 24x24 image for the user
- `image_32` - URL to a 32x32 image for the user
- `image_48` - URL to a 48x48 image for the user
- `image_72` - URL to a 72x72 image for the user
- `image_192` - URL to a 192x192 image for the user
- `image_512` - URL to a 512x512 image for the user
- `image_original` - URL to the original image for the user

#### Methods

##### json

```go filename="Method signature"
json() map
```

Returns a map representation of the user profile object.

```go filename="Example"
>>> user := client.get_user_info("U123456")
>>> profile_map := user.profile.json()
>>> profile_map.email
"john@example.com"
```

### slack.conversation_iterator

An iterator for paging through conversations in a workspace.

```go filename="Example"
>>> convs := client.get_conversations()
>>> for i, conv in convs {
...   print(conv.name)
... }
general
random
team-project
```

Note: Iterators must be consumed through iteration and don't support `len()` operations directly.

### slack.message_iterator

An iterator for paging through messages in a conversation.

```go filename="Example"
>>> messages := client.get_conversation_history("#general")
>>> for i, msg in messages {
...   print(msg.text)
... }
Hello world!
Important announcement
Check out this link: https://example.com
```

### slack.conversation_members_iterator

An iterator for paging through members of a conversation.

```go filename="Example"
>>> members := client.get_conversation_members("#general")
>>> for i, member_id in members {
...   user := client.get_user_info(member_id)
...   print(user.name)
... }
johndoe
janedoe
bobsmith
```
