package main

import (
	"context"
	"log"

	"github.com/Disassembly-AI/terraform-provider-disassembly/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/Disassembly-AI/disassembly",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
