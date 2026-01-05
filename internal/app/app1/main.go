package main

import (
	"github.com/pivaldi/go-cleanstack/internal/app/app1/cmd"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/clierr"
)

func main() {
	if err := cmd.GetRootCmd().Execute(); err != nil {
		clierr.ExitOnError(err, true)
	}

}
