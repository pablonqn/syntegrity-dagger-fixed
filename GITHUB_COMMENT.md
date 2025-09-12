## 🚨 Dagger SDK Compatibility Issue with Go 1.24+ and 1.25.1

We've identified a **critical compatibility issue** when using Dagger SDK with Go 1.24+ and Go 1.25.1.

### ❌ Problem
Build fails with telemetry-related errors:
```
# dagger.io/dagger/telemetry
../../../go/pkg/mod/dagger.io/dagger@v0.18.17/telemetry/transform.go:69:29: cannot use res (variable of type *resource.Resource) as resource.Resource value in argument to ResourceToPB
```

### 🧪 Testing Results
| Go Version | Dagger Version | Status |
|------------|----------------|---------|
| 1.23.x | v0.18.17 | ✅ **WORKS** |
| 1.24.0 | v0.18.17 | ❌ **FAILS** |
| 1.24.1 | v0.18.17 | ❌ **FAILS** |
| 1.24.2 | v0.18.17 | ❌ **FAILS** |
| 1.24.7 | v0.18.17 | ❌ **FAILS** |
| 1.25.1 | v0.18.17 | ❌ **FAILS** |

**Tested Dagger versions**: v0.12.0, v0.14.0, v0.15.0, v0.16.0, v0.17.0, v0.18.0, v0.18.16, v0.18.17
**All versions fail** with Go 1.24+

### 🔍 Root Cause
**OpenTelemetry compatibility issues** between:
- Dagger's telemetry package (older OpenTelemetry APIs)
- Go 1.24+ toolchain (newer OpenTelemetry versions)

### 🛠️ Workarounds Attempted
- ✅ Version downgrades (multiple Dagger versions)
- ❌ Build tags (`-tags=no_telemetry`) - not supported
- ❌ Replace directives in `go.mod`
- ✅ Removed conflicting dependencies (`go-kit-logger` → `log/slog`)

### 📋 Current Status
**BLOCKED**: Cannot use Dagger SDK with Go 1.24+ or Go 1.25.1

### 💡 Recommendations
1. **For Dagger Team**: Update telemetry package for Go 1.24+ compatibility
2. **For Users**: Use Go 1.23.x until Dagger fixes this issue
3. **Monitor**: Dagger releases for Go 1.24+ support

### 📝 Environment
- **OS**: macOS (darwin/arm64)
- **Project**: `github.com/getsyntegrity/syntegrity-dagger`
- **Issue**: Affects all projects using Dagger SDK with Go 1.24+

---
**This is a known issue in the Dagger community that needs to be addressed by the Dagger team.**
