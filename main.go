package main

import (
	"context"
	"fmt"

	"github.com/carbonable-labs/indexer/internal/dispatcher"
	"github.com/carbonable-labs/indexer/internal/indexer"
	"github.com/carbonable-labs/indexer/internal/starknet"
	"github.com/carbonable-labs/indexer/internal/storage"
	"github.com/carbonable-labs/indexer/internal/synchronizer"
	"github.com/charmbracelet/log"
)

const welcomeMessage = "Sheshat ... Indexing"

// starting block
// indexing configuration

// get all contracts declared at in the genesis dataget_class_by_hash
// -> each time config is changed, know where to start indexing from
// -> keep all indexing configurations to enable fast retrieval / per contract

// we may want to ignore what is before starting block as it is not required to have data
// for each contract -> index all events in a event driven way
// store txs and state updates as customer may want to retrieve data based on this.

// First step
// fetch all data related to contracts
// store data

// Second step
// stream data into message broker

// Every reload or reindex is based off a last_event_id (ulid based on time)
// replayed from database to get faster indexing

func main() {
	log.SetLevel(log.DebugLevel)
	fmt.Println(welcomeMessage)

	syncErrCh := make(chan error)
	indexerErrCh := make(chan error)
	ctx := context.Background()

	client := starknet.NewSepoliaFeederGatewayClient()
	storage := storage.NewPebbleStorage()
	bus := dispatcher.NewInMemoryBus()

	go synchronizer.Run(ctx, client, storage, syncErrCh)
	go indexer.Run(ctx, client, storage, bus, indexerErrCh)

	for {
		select {
		case syncErr := <-syncErrCh:
			log.Error("error while syncing network", "error", syncErr)
		case indexerErr := <-indexerErrCh:
			log.Error("error while indexing network", "error", indexerErr)
		}
	}
}
