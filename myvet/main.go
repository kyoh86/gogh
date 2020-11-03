// This file is generated by vetgen.
// Do NOT modify this file.
//
// You can run this tool with go vet such as:
//	go vet -vettool=$(which myvet) pkgname
package main

import (
	"github.com/gostaticanalysis/vetgen/analyzers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

var myAnayzers = []*analysis.Analyzer{}

func main() {
	unitchecker.Main(append(
		analyzers.Recommend(),
		myAnayzers...,
	)...)
}