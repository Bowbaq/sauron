package flagx

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

// MustParse attempts to parse the flags defined by the struct tags on v. If an error occurs, the error
// will be logged and the program will exit.
func MustParse(v interface{}) {
	parser := flags.NewParser(v, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		if ferr, ok := err.(*flags.Error); ok && ferr.Type == flags.ErrHelp {
			fmt.Println(err)
			os.Exit(0)
		}

		log.Fatalf("flagx: Couldn't parse flags: %v", err)
	}
}
