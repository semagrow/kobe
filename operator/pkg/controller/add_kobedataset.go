package controller

import (
	"github.com/semagrow/kobe/operator/pkg/controller/kobedataset"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kobedataset.Add)
}
