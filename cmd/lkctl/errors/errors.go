package errors

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Check logs a fatal message and exits with ErrorGeneric if err is not nil
func Check(err error) {
	if err != nil {
		exitfunc := func() {
			os.Exit(1)
		}
		log.RegisterExitHandler(exitfunc)
		log.Fatal(err)
	}
}
