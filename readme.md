# gochat

just another messaging application supporting text, sticker and binary data.

**Note: currently under framework shift from gin to gorrilla mux, May have breaking dependency**

## features

- [x] text/sticker/binary messages retrival
- [x] system messages
- [x] group channels
- [ ] realtime websocket messaging
- [ ] call / meetings

## Usage

using docker compose:

```bash
docker compose up
```

endpoints:
|method  | uri | context |
|--------|-----|---------|
| POST   | <http://localhost:1212/user>             | get access token                      |
| GET    | <http://localhost:1212/user/messages>    | get user channel with recent messages |
| DELETE | <http://localhost:1212/user>             | delete user                           |
| POST   | <http://localhost:1212/channel>          | create a new channel                  |
| GET    | <http://localhost:1212/channel/:id>      | get messages from channel             |
| POST   | <http://localhost:1212/channel/:id>      | post message in channel               |
| POST   | <http://localhost:1212/channel/:id/join> | join a channel                        |
| DELETE | <http://localhost:1212/channel/:id>      | delete a channel                      |
