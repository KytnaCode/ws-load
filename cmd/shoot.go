package cmd

import (
	"math/rand/v2"
	"sync"
	"ws-load/pkg/load"
	random "ws-load/pkg/random/pure"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var (
	amount   int
	messages int
	size     int

	// shootCmd represents the shoot command
	shootCmd = &cobra.Command{
		Use:   "shoot",
		Short: "Start a connections, and shoot messages.",
		Long: `Starts a bunch of short-lived connections, send a certain
number of messages and then closes the connection, command
ends automatically when all messages are sent.`,
		Aliases:    []string{"s"},
		SuggestFor: []string{"load", "gen", "start"},
		Example: `ws-load --amount 1000 --messages 50 ws://localhost:3000
ws-load -A 500 -M 250 -S 4096 ws://localhost:8080/ws ws://locahost:8081/ws`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup

			wg.Add(len(args))

			msg := random.Bytes(rand.New(rand.NewPCG(8, 8)), size) //nolint:gosec

			loader := load.NewWSLoader(msg, websocket.BinaryMessage, messages, amount)

			for _, target := range args {
				go func() {
					defer wg.Done()

					loader.Shoot(cmd.Context(), target)
				}()
			}

			wg.Wait()
		},
	}
)

func init() {
	rootCmd.AddCommand(shootCmd)

	shootCmd.PersistentFlags().IntVarP(&amount, "amount", "A", 1, "Amount of connections.")
	shootCmd.PersistentFlags().IntVarP(&messages, "messages", "M", 1, "Amount of messages to send per connection.")
	shootCmd.PersistentFlags().IntVarP(&size, "bufsize", "S", 2*1024, "Size of each message content.")
}
