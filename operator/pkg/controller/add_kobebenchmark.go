package controller

import (
	"github.com/semagrow/kobe/operator/pkg/controller/kobebenchmark"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kobebenchmark.Add)
}
