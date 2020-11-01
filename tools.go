//+build tools

package main

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/rjeczalik/interfaces/cmd/interfacer"
	_ "github.com/gostaticanalysis/vetgen"
)
