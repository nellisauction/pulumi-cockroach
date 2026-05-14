package shim

import (
	tfpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/cockroachdb/terraform-provider-cockroach/internal/provider"
)

func NewProvider() tfpf.Provider {
	return provider.New("dev")()
}
