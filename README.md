# MatchMaker Game Server

A generic websocket server for hosting matches of two-player, turn-based games. Has support for round-robin tournament hosting.https://github.com/jar2333/MatchMaker 

Extensible using a provided game interface. To be provided is an example implementation of Rock Paper Scissors.

## Resource API

When connecting to the server, a RESTful HTTP API is specified to request a game connection:

1. `GET /game`

This request asks the server to find a game with a matching player for the client to play against, returning a websocket connection to play the game with the Game API.

2. `GET /tournament`

2. `POST /tournament`

2. `DELETE /tournament/<id>`

2. `GET /tournament/<id>/registered/`

2. `POST /tournament/<id>/registered/`

<!-- This command registers the client as a player. Closing the WebSocket connection before the torunament starts will unregister the player, requiring that this command be run again upon reconnection. The `<id>` field must be a unique identifier for a tournament. -->

## Game API

### Game agnostic API:

During a game, a WebSocket connection is established, where the following JSON API is employed:

 * Requests:
    1. `{"type": "move", "move": <move>}`
    Sent to specify a move to be played. The list of available move payloads is specified by the game's specific API.
    
 * Responses:
     1. `{"type": "game_start"}`

     Received when a game involving this client has started. If this client is player 1, a message indicating that their turn has started will be sent shortly afterwards (see 3. below).

     2. `{"type": "game_ended"}`

     Received when a game involving this client has ended.

     3. `{"type": "turn_started"}`

     Received when the turn for the player associated with this client has started. This will be followed by a message indicating that input is being read (see 5. below)

     4. `{"type": "turn_ended"}`

     Received when the turn for the player associated with this client has ended. If game has ended, it will be followed by a game ended message (2. above). Otherwise, a turn started message (3. above) will be received when it is the player's turn again.

     5. `{"type": "reading_move"}`

     Received when server is reading the websocket for commands. Note, that the websocket will buffer all newline-separated commands sent at any time, but will not read them until after this message has been sent to the client. This message also indicates that the timer has been started.

     6. `{"type": "invalid_move"}`

     Received when message sent by client could not be parsed. It can be assumed that no side-effects occurred. An invalid move does not reset the timer. It is recommended to never have to rely on this message, it is provided for debugging purposes. Followed by a state message (8. below).

     7. `{"type": "valid_move"}`

     Received when message sent by client was successfully parsed, leading to a game action. This resets the timer. Followed by a state message (8. below).

     8. `{"type": "state", "state": <state_dict>, "timer": <time>}`

     Received after a game start message (see 1. above), and after a valid or invalid move message (6. and 7. above). Contains a dictionary `<state_dict>` encoding the state of the game, and a float `<time>` denoting how many seconds are left in the timer. Note, not all posible moves end the player's turn. An explicit turn ended message (see 4. above) will be sent shortly afterwards if the turn was ended. Otherwise, another reading move message (see 5. above) will be sent instead, indicating that another player move is being read.
