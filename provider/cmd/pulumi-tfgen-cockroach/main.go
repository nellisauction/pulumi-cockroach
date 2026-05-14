package main

import (
	pftfgen "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfgen"
	provider "github.com/nellisauction/pulumi-cockroach/provider"
)

func main() {
	pftfgen.Main("cockroach", provider.Provider())
}
