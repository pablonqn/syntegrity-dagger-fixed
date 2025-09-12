# Dagger SDK Compatibility Issue with Go 1.24+ and 1.25+

## Problem Description

We encountered a critical compatibility issue when trying to use Dagger SDK with Go 1.24+ and Go 1.25.1. The build fails with telemetry-related errors in the Dagger SDK.

## Error Details

### Go 1.25.1 with Dagger v0.18.17
```
# dagger.io/dagger/telemetry
../../../go/pkg/mod/dagger.io/dagger@v0.18.17/telemetry/transform.go:69:29: cannot use res (variable of type *resource.Resource) as resource.Resource value in argument to ResourceToPB
```

### Go 1.24.x with Dagger v0.18.17
```
# dagger.io/dagger/telemetry
../../../go/pkg/mod/dagger.io/dagger@v0.18.17/telemetry/transform.go:69:29: cannot use res (variable of type *resource.Resource) as resource.Resource value in argument to ResourceToPB
```

### Go 1.24.x with Dagger v0.15.0
```
# dagger.io/dagger/telemetry
../../../go/pkg/mod/dagger.io/dagger@v0.15.0/telemetry/transform.go:69:29: cannot use res (variable of type *resource.Resource) as resource.Resource value in argument to ResourceToPB
../../../go/pkg/mod/dagger.io/dagger@v0.15.0/telemetry/transform.go:957:25: cannot use processor (variable of type *collectLogProcessor) as "go.opentelemetry.io/otel/sdk/log".Processor value in argument to sdklog.WithProcessor: *collectLogProcessor does not implement "go.opentelemetry.io/otel/sdk/log".Processor (wrong type for method OnEmit)
                have OnEmit(context.Context, "go.opentelemetry.io/otel/sdk/log".Record) error
                want OnEmit(context.Context, *"go.opentelemetry.io/otel/sdk/log".Record) error
```

## Versions Tested

### Go Versions Tested:
- ✅ Go 1.23.x - **WORKS** (compatible with Dagger)
- ❌ Go 1.24.0 - **FAILS** (telemetry compatibility issues)
- ❌ Go 1.24.1 - **FAILS** (telemetry compatibility issues)
- ❌ Go 1.24.2 - **FAILS** (telemetry compatibility issues)
- ❌ Go 1.24.7 - **FAILS** (telemetry compatibility issues)
- ❌ Go 1.25.1 - **FAILS** (telemetry compatibility issues)

### Dagger Versions Tested:
- ❌ v0.18.17 - **FAILS** with Go 1.24+
- ❌ v0.18.16 - **FAILS** with Go 1.24+
- ❌ v0.18.0 - **FAILS** with Go 1.24+
- ❌ v0.17.0 - **FAILS** with Go 1.24+
- ❌ v0.16.0 - **FAILS** with Go 1.24+
- ❌ v0.15.0 - **FAILS** with Go 1.24+
- ❌ v0.14.0 - **FAILS** with Go 1.24+
- ❌ v0.12.0 - **FAILS** with Go 1.24+

## Root Cause Analysis

The issue appears to be related to **OpenTelemetry compatibility** between:
1. **Dagger's telemetry package** - which uses older OpenTelemetry APIs
2. **Go 1.24+ toolchain** - which includes newer OpenTelemetry versions

### Specific Issues:
1. **Resource type mismatch**: `*resource.Resource` vs `resource.Resource` in `ResourceToPB` function
2. **OpenTelemetry SDK log processor interface changes**: Method signature changes in `OnEmit` method
3. **Pointer vs value type mismatches** in OpenTelemetry APIs

## Workarounds Attempted

### 1. Version Downgrades
- Tried multiple Dagger versions (v0.12.0 to v0.18.17)
- Tried Go version downgrades (1.25.1 → 1.24.x → 1.23.x)

### 2. Build Tags
- Attempted to disable telemetry with `-tags=no_telemetry`
- **Result**: Dagger doesn't support telemetry disabling via build tags

### 3. Replace Directives
- Used `replace` directives in `go.mod` to force specific versions
- **Result**: Issue persists across all tested versions

### 4. Dependency Management
- Removed conflicting dependencies (`go-kit-logger`)
- Updated to use standard `log/slog` package
- **Result**: Build still fails due to Dagger telemetry issues

## Current Status

**BLOCKED**: Cannot use Dagger SDK with Go 1.24+ or Go 1.25.1 due to telemetry compatibility issues.

## Recommendations

### For Dagger Team:
1. **Update telemetry package** to use compatible OpenTelemetry APIs
2. **Add build tag support** to disable telemetry (`//go:build !no_telemetry`)
3. **Test compatibility** with Go 1.24+ and 1.25+ in CI/CD
4. **Update OpenTelemetry dependencies** to latest compatible versions

### For Users:
1. **Use Go 1.23.x** until Dagger fixes compatibility
2. **Monitor Dagger releases** for Go 1.24+ compatibility
3. **Consider alternative solutions** if telemetry is not required

## Environment Details

- **OS**: macOS (darwin/arm64)
- **Go Module**: `github.com/getsyntegrity/syntegrity-dagger`
- **Dependencies**: Standard Go modules + Dagger SDK
- **Build System**: Make + GitHub Actions

## Related Issues

This appears to be a known issue in the Dagger community. The telemetry package needs to be updated to support newer Go versions and OpenTelemetry APIs.

---

**Note**: This issue affects all projects using Dagger SDK with Go 1.24+ or Go 1.25.1. The problem is in the Dagger SDK itself, not in user code.
