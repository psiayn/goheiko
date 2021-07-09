package daemon

import (
	"log"
	"os"

	"github.com/psiayn/heiko/internal/scheduler"
	"github.com/sevlyar/go-daemon"
)

func stopHandler(sig os.Signal) error {
	log.Println("Terminating heiko daemon...")

	log.Println("Waiting for tasks to finish.")
	for _, stop := range scheduler.Stops {
		stop <- struct{}{}
	}

	log.Println("Waiting for scheduler to stop.")
	for _, done := range scheduler.Dones {
		<-done
	}

	return daemon.ErrStop
}
