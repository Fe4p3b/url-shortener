package noosexit

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func Example() {
	// for single analizer
	singlechecker.Main(Analyzer())

	// for multiple analizers
	multichecker.Main(
		Analyzer(),
	)
}
