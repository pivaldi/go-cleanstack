package clierr

import (
	"fmt"
	"os"

	"github.com/pivaldi/go-cleanstack/internal/common/platform/apperr"
)

func ExitOnError(err error, debug bool) {
	if err == nil {
		return
	}

	ae := apperr.As(err)
	if ae == nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Always print stable info
	fmt.Fprintf(os.Stderr, "error [%s]: %s\n", ae.Code, ae.Message)

	if ae.IsPrivate() || debug {
		if ae.Cause != nil {
			fmt.Fprintln(os.Stderr, "cause:", ae.Cause)
		}

		if ae.Stack != "" {
			fmt.Fprintln(os.Stderr, "\nstack:\n"+ae.Stack)
		}

		if ae.Fields != nil {
			fmt.Fprintln(os.Stderr, "\nfields:", ae.Fields)
		}

		if ae.Req != nil {
			fmt.Fprintln(os.Stderr, "\nreq:", ae.Req)
		}
	}

	os.Exit(1)
}
