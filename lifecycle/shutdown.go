package lifecycle

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var shutdownRegistry = map[string]func() error{}

// RegisterShutdownCallback will register a shutdown handler function.
func RegisterShutdownCallback(name string, fn func() error) {
	shutdownRegistry[name] = fn
}

// WaitForShutdown is a blocking function that will run shutdown callbacks
// before the process is allowed to finish.
func WaitForShutdown() {
	// Watch for process lifecycle signals, then run callbacks.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	// Wait until we receive a shutdown message.
	s := <-sigCh
	if s == nil {
		return
	}

	log.Printf("Received shutdown signal ‘%s’", s)
	log.Println("Running shutdown callbacks")

	// Start another goroutine to immediately shutdown if required.
	go func() {
		sig := <-sigCh
		log.Printf("Second signal received ‘%s’\n", sig)
		log.Println("Forcing Exit")
		os.Exit(1)
	}()

	// Run all shutdown functions in goroutines.
	wg := sync.WaitGroup{}
	for name, fn := range shutdownRegistry {
		wg.Add(1)
		go func(fn func() error, name string) {
			log.Printf("Running shutdown callback ‘%s’\n", name)
			err := fn()
			if err != nil {
				log.Printf("Failed shutdown callback ‘%s’: %s\n", name, err)
			} else {
				log.Printf("Finished shutdown callback ‘%s’\n", name)
			}
			wg.Done()
		}(fn, name)
	}
	wg.Wait()
}
