# MatchMaker Game Server

A generic websocket server for hosting matches of two-player, turn-based games. Has support for round-robin tournament hosting.https://github.com/jar2333/MatchMaker 

Extensible using a provided game interface. To be provided is an example implementation of Rock Paper Scissors.

## HTTP API

All requests must specify an authorized username in the `Authorization` header, using the `Basic` authentication scheme. 

<!-- 

Question: Should /game and /game/ws be different endpoints? 
The motivation is the following formulation:

The GET /game HTTP endpoint starts a matchmaking process to find a game,
then it returns a token which can be used to connect to the game through 
the GET /game/ws WebSocket endpoint (HTTP + Connection: Upgrade header).

Alternatively, GET /game is a WebSocket endpoint, and it establishes a 
connection to a game which the server implicitly matchmakes for. 

Going with the second one, for now.

-->

<!-- 
   Add documentation for sent payloads/parameters and response payloads/parameters.
-->

### HTTP endpoints (WIP)

1. `GET /tournament`

2. `POST /tournament`

3. `DELETE /tournament/<id>`

4. `GET /tournament/<id>`

5. `GET /tournament/<id>/registered/`

6. `POST /tournament/<id>/registered/`

7. `DELETE /tournament/<id>/registered/<rid>`

8. `GET /tournament/<id>/registered/<rid>`

### WebSocket endpoints

1. `GET /game`

This request initiates a websocket connection with the server. Once established, the client must utilize the WebSocket Game API to communicate with the server.
The client must wait to receive the message indicating a game has started. Matchmaking is handled by the server, and the specifics are an implementation detail.
The Game API is documented more below.

2. `GET /tournament/<id>/game`

This request initiates a websocket connection with the server, allowing the player to participate in the tournament. 
The `Authorization` header should be provided with the proper credentials for the tournament (username of registered user + tournament password), using the `Basic` authentication scheme. In this connection, multiple games may be played, corresponding to the tournament schedule. If the tournament is already taking place, the connection will fail. The client must attempt to connect before the tournament is scheduled to begin, and maintain that connection throughout.


<!-- This command registers the client as a player. Closing the WebSocket connection before the torunament starts will unregister the player, requiring that this command be run again upon reconnection. The `<id>` field must be a unique identifier for a tournament. -->

## WebSocket Game API

### Game agnostic API:

Before a game, a WebSocket connection is established. The Game API is designed around events sent to the player through this connection 
which denote the beginning and end of certain scopes during a game. The nested structure of a tournament (or generally, a sequence of games) is as follows, expressed using XML:

```xml
<tournament>
   <game>
      <turn>
         <move/>
         ...
      </turn>
      ...
   </game>
   ...
</tournament>
```

A move is an atomic unit. An invalid move will cause no state change in the game. A valid move, conversely, will usually cause some change in game state. A move does not necessarily end a turn. A move which doesn't end the turn is termed a "free action". Moves are the way a client can send arbitrary JSON payloads to the server, corresponding to a specific game's API.

There can be multiple moves in a turn, multiple turns in a game, and multiple games in a tournament. It is the responsibility of the client in a game to adequately respond to the server's messages which delimit the beginning of a game, turn, or move. Conversely,
the server should respond to messages sent by the player that specify a move.

The following JSON API is employed to send/receive these messages.

 * Client messages:
    1. `{"type": "move", "move": <move>}`

    Sent to specify a move to be played, using a JSON payload `<move>`. The list of available move payloads is specified by the game's specific API.

 * Server messages:
     1. `{"type": "game_start"}`

     Received when a game involving this client has started. If this client is player 1, a message indicating that their turn has started will be sent shortly afterwards (3.).

     2. `{"type": "game_ended"}`

     Received when a game involving this client has ended.

     3. `{"type": "turn_started"}`

     Received when the turn for the player associated with this client has started. This will be followed by a message indicating that input is being read (5.)

     4. `{"type": "turn_ended"}`

     Received when the turn for the player associated with this client has ended. If game has ended, it will be followed by a game ended message (2.). Otherwise, a turn started message (3.) will be received when it is the player's turn again.

     5. `{"type": "reading_move"}`

     Received when server is reading the websocket for commands. Note, that the websocket will buffer all newline-separated commands sent at any time, but will not read them until after this message has been sent to the client. This message also indicates that the timer has been started.

     6. `{"type": "invalid_move"}`

     Received when message sent by client could not be parsed. It can be assumed that no side-effects occurred. An invalid move does not reset the timer. It is recommended to never have to rely on this message, it is provided for debugging purposes. Followed by a state message (8.).

     7. `{"type": "valid_move"}`

     Received when message sent by client was successfully parsed, leading to a game action. This resets the timer. Followed by a state message (8.).

     8. `{"type": "state", "state": <state>, "timer": <time>}`

     Received after a game start message (1.), and after a valid or invalid move message (6., 7.). Contains a JSON payload `<state>` encoding the state of the game, and a float `<time>` denoting how many seconds are left in the timer. Note, a move may not end the player's turn. An explicit turn ended message (4.) will be sent shortly afterwards if the turn was ended. Otherwise, another reading move message (5.) will be sent instead, indicating that another player move is being read. The shape of `<state>` is determined by the specific game being played.
