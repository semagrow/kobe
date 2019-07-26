package controller

import (
	"github.com/semagrow/kobe/operator/pkg/controller/kobefederator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kobefederator.Add)
}