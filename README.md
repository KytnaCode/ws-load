# ws-load

*ws-load* is a websocket load generator written in pure Go.

## Usage

### Shoot

ws-load allows to _shoot_ a bunch of messages, that means, create a certain amount of connection, send a certain amount
of messages and gracefully close the connection.

```sh
ws-load shoot [options] <urls...>
```

`shoot` messages are just binary messages a random data, the generator have a fixed seed, so the data will be the same
between executions, all messages of all connections have the same data, the size of the generated buffer can be
specified with `--bufsize <size_in_bytes>` (defaults to 512).

*Example*:

```sh
ws-load shoot --amount 1000 --messages 50 ws://localhost:8080/ws
```

Where _amount_ is the amount of connections and _messages_ is the amout of messages send per connection.

*help*:

```
Starts a bunch of short-lived connections, send a certain
number of messages and then closes the connection, command
ends automatically when all messages are sent.

Usage:
  ws-load shoot [flags]

Aliases:
  shoot, s

Examples:
ws-load shoot --amount 1000 --messages 50 ws://localhost:3000
ws-load shoot -A 500 -M 250 -S 4096 ws://localhost:8080/ws ws://locahost:8081/ws

Flags:
  -A, --amount int     Amount of connections. (default 1)
  -S, --bufsize int    Size of each message content. (default 2048)
  -h, --help           help for shoot
  -M, --messages int   Amount of messages to send per connection. (default 1)```
```

### Help

You can display the help message with `ws-load help`:

```text
ws-load is a dead-simple websocket load generator.

Usage:
  ws-load [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  shoot        Start a connections, and shoot messages.

Flags:
  -h, --help   help for ws-load

Use "ws-load [command] --help" for more information about a command.
```

## Versioning

This projects follows [Semanting Versioning 2.0](https://semver.org/spec/v2.0.0.html).

## License

MIT. (see LICENSE.txt)
