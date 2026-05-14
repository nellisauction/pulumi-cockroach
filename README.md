# CockroachDB Cloud Resource Provider

The CockroachDB Cloud Resource Provider lets you manage [CockroachDB Cloud](https://www.cockroachlabs.com/product/cloud/) resources.

## Installing

This package is available for several languages/platforms:

### Node.js (JavaScript/TypeScript)

To use from JavaScript or TypeScript in Node.js, install using either `npm`:

```bash
npm install @nellisauction/pulumi-cockroach
```

or `yarn`:

```bash
yarn add @nellisauction/pulumi-cockroach
```

### Python

To use from Python, install using `pip`:

```bash
pip install pulumi_cockroach
```

### Go

To use from Go, use `go get` to grab the latest version of the library:

```bash
go get github.com/nellisauction/pulumi-cockroach/sdk/go/...
```

### .NET

To use from .NET, install using `dotnet add package`:

```bash
dotnet add package Pulumi.CockroachDB
```

## Configuration

The following configuration points are available for the `cockroach` provider:

- `cockroach:apikey` (environment: `COCKROACH_API_KEY`) - a CockroachDB Cloud API key
- `cockroach:apijwt` (environment: `COCKROACH_API_JWT`) - a JWT from a JWT Issuer configured for the CockroachDB Cloud Organization

## Reference

For detailed reference documentation, please visit [the Pulumi registry](https://www.pulumi.com/registry/packages/cockroach/api-docs/).
