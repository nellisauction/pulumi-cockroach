package cockroach

import (
	"fmt"
	"path/filepath"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	tfbridgetokens "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	cockroachshim "github.com/cockroachdb/terraform-provider-cockroach/shim"

	"github.com/nellisauction/pulumi-cockroach/provider/pkg/version"
)

const (
	mainPkg = "cockroach"
	mainMod = "index"
)

//go:embed cmd/pulumi-resource-cockroach/bridge-metadata.json
var metadata []byte

func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		P:                 pf.ShimProvider(cockroachshim.NewProvider()),
		Name:              "cockroach",
		Version:           version.Version,
		DisplayName:       "CockroachDB",
		Publisher:         "NellisAuction",
		PluginDownloadURL: "github://api.github.com/nellisauction/pulumi-cockroach",
		Description:       "A Pulumi package for managing CockroachDB Cloud resources.",
		Keywords:          []string{"pulumi", "cockroach", "cockroachdb", "category/database"},
		License:           "Apache-2.0",
		Homepage:          "https://github.com/nellisauction/pulumi-cockroach",
		Repository:        "https://github.com/nellisauction/pulumi-cockroach",
		GitHubOrg:         "cockroachdb",
		Config: map[string]*tfbridge.SchemaInfo{
			"apikey": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"COCKROACH_API_KEY"},
				},
			},
			"apijwt": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"COCKROACH_API_JWT"},
				},
			},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName:          "@nellisauction/pulumi-cockroach",
			RespectSchemaVersion: true,
		},
		Python: (func() *tfbridge.PythonInfo {
			i := &tfbridge.PythonInfo{RespectSchemaVersion: true}
			i.PyProject.Enabled = true
			return i
		})(),
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/nellisauction/pulumi-%[1]s/sdk/", mainPkg),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RespectSchemaVersion: true,
			PackageReferences:    map[string]string{"Pulumi": "3.*"},
			Namespaces:           map[string]string{mainPkg: "CockroachDB"},
		},
		MetadataInfo:                   tfbridge.NewProviderMetadata(metadata),
		EnableZeroDefaultSchemaVersion: true,
		EnableAccurateBridgePreview:    true,
	}

	prov.MustComputeTokens(tfbridgetokens.SingleModule("cockroach_", mainMod,
		tfbridgetokens.MakeStandard(mainPkg)))
	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}
