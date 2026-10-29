package cockroach

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	pf "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	tfbridgetokens "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
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
		Resources: map[string]*tfbridge.ResourceInfo{
			"cockroach_sql_user": {
				Fields: map[string]*tfbridge.SchemaInfo{
					"password": {
						DeprecationMessage: "The `password` attribute persists the clear text password in Terraform state. " +
							"Migrate to `passwordWo` (with `passwordWoVersion` for rotations): Pulumi keeps `passwordWo` " +
							"only as an encrypted secret input and never in resource outputs or Terraform state. " +
							"`password` will be removed in a future major release of the provider.",
					},
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

	clusterDS := prov.DataSources["cockroach_cluster"]
	clusterDS.Fields = markComputedObjectsOptional(
		prov.P.DataSourcesMap().Get("cockroach_cluster").Schema(), clusterDS.Fields, "backup_config")
	markOptional(prov.P.DataSourcesMap().Get("cockroach_cluster").Schema(), clusterDS.Fields,
		"account_id", "plan", "edition", "parent_id",
		"serverless.spend_limit", "serverless.with_empty_ip_allowlist",
		"serverless.usage_limits.request_unit_limit",
		"serverless.usage_limits.storage_mib_limit",
		"serverless.usage_limits.provisioned_virtual_cpus",
		"dedicated.supports_cluster_virtualization",
		"regions.machine_type", "regions.num_virtual_cpus")

	clusterRes := prov.Resources["cockroach_cluster"]
	if clusterRes.Fields == nil {
		clusterRes.Fields = map[string]*tfbridge.SchemaInfo{}
	}
	markOptional(prov.P.ResourcesMap().Get("cockroach_cluster").Schema(), clusterRes.Fields,
		"account_id", "plan", "edition", "parent_id",
		"serverless.with_empty_ip_allowlist",
		"regions.machine_type", "regions.num_virtual_cpus")
	prov.SetAutonaming(255, "-")

	return prov
}

// markOptional types computed scalars as optional outputs. The upstream read
// leaves these null whenever the API omits them (plan vs edition, no parent
// folder, homogeneous regions) or never sets them outside a plan (data source
// and import), but the bridge types every computed output as always present.
// Paths are dot-separated TF attribute names.
func markOptional(sch shim.SchemaMap, fields map[string]*tfbridge.SchemaInfo, paths ...string) {
	for _, path := range paths {
		keys := strings.Split(path, ".")
		cur, curFields := sch, fields
		for i, key := range keys {
			s, ok := cur.GetOk(key)
			if !ok {
				panic(fmt.Sprintf("markOptional: unknown attribute %q in %q", key, path))
			}
			info := curFields[key]
			if info == nil {
				info = &tfbridge.SchemaInfo{}
				curFields[key] = info
			}
			if i == len(keys)-1 {
				info.MarkAsOptional = tfbridge.True()
				break
			}
			if elem, ok := s.Elem().(shim.Schema); ok {
				if info.Elem == nil {
					info.Elem = &tfbridge.SchemaInfo{}
				}
				info, s = info.Elem, elem
			}
			obj, ok := s.Elem().(shim.Resource)
			if !ok {
				panic(fmt.Sprintf("markOptional: %q in %q is not an object", key, path))
			}
			if info.Elem == nil {
				info.Elem = &tfbridge.SchemaInfo{}
			}
			if info.Elem.Fields == nil {
				info.Elem.Fields = map[string]*tfbridge.SchemaInfo{}
			}
			cur, curFields = obj.Schema(), info.Elem.Fields
		}
	}
}

// The cluster data source only populates the nested objects that apply to the
// cluster (one of serverless/dedicated/host, one cloud account, etc.), but the
// bridge types computed outputs as always present. alwaysPresent lists keys
// the upstream read always populates (backup_config is fetched separately and
// the read fails without it).
func markComputedObjectsOptional(
	sch shim.SchemaMap, fields map[string]*tfbridge.SchemaInfo, alwaysPresent ...string,
) map[string]*tfbridge.SchemaInfo {
	if fields == nil {
		fields = map[string]*tfbridge.SchemaInfo{}
	}
	sch.Range(func(key string, s shim.Schema) bool {
		obj, ok := s.Elem().(shim.Resource)
		if !ok || s.Type() != shim.TypeMap || !s.Computed() || s.Optional() ||
			slices.Contains(alwaysPresent, key) {
			return true
		}
		info := fields[key]
		if info == nil {
			info = &tfbridge.SchemaInfo{}
			fields[key] = info
		}
		info.MarkAsOptional = tfbridge.True()
		if info.Elem == nil {
			info.Elem = &tfbridge.SchemaInfo{}
		}
		info.Elem.Fields = markComputedObjectsOptional(obj.Schema(), info.Elem.Fields)
		return true
	})
	return fields
}
