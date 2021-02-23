package fofa_cli

import (
	"fmt"
	"os"
)

//fofaErr
func fofaErr(err error) {
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err.Error())
		os.Exit(1)
	}
}
