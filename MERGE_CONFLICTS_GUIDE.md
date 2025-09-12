# Guía para Resolver Conflictos de Merge

Esta guía te ayudará a resolver los conflictos de merge que pueden aparecer al integrar los cambios del pipeline de release.

## 🚨 Archivos que Pueden Tener Conflictos

### 1. `.github/workflows/release.yml`
**Conflicto esperado**: Diferencias entre la versión antigua y la nueva del pipeline de release.

**Resolución recomendada**:
- Mantener la versión nueva (Release Pipeline) que incluye:
  - Auto-versioning
  - Multi-platform builds
  - Changelog generation
  - GitHub release creation

**Pasos**:
```bash
# Mantener nuestra versión
git checkout --ours .github/workflows/release.yml
git add .github/workflows/release.yml
```

### 2. `go.mod`
**Conflicto esperado**: Diferencias en versiones de dependencias.

**Resolución recomendada**:
- Mantener Go 1.24.2
- Mantener Dagger v0.9.11 (versión estable)
- Limpiar dependencias innecesarias

**Pasos**:
```bash
# Mantener nuestra versión optimizada
git checkout --ours go.mod
go mod tidy
git add go.mod
```

### 3. `go.sum`
**Conflicto esperado**: Checksums de dependencias diferentes.

**Resolución recomendada**:
- Regenerar go.sum después de resolver go.mod

**Pasos**:
```bash
# Regenerar go.sum
go mod tidy
git add go.sum
```

## 🔧 Resolución Automática

### Opción 1: Script Automático
```bash
./resolve-conflicts.sh
```

### Opción 2: Resolución Manual
```bash
# 1. Resolver conflictos en cada archivo
git checkout --ours .github/workflows/release.yml
git checkout --ours go.mod

# 2. Limpiar dependencias
go mod tidy

# 3. Agregar archivos resueltos
git add .github/workflows/release.yml go.mod go.sum

# 4. Completar merge
git commit -m "Resolve merge conflicts: keep optimized release pipeline"
```

## 📋 Verificación Post-Resolución

Después de resolver los conflictos, verifica que todo funciona:

```bash
# 1. Verificar que compila
go build ./...

# 2. Verificar que los tests pasan
go test ./...

# 3. Verificar que el linter funciona
make lint

# 4. Verificar que el build de release funciona
make build-release
```

## 🎯 Configuración Final Recomendada

### go.mod
```go
module github.com/getsyntegrity/syntegrity-dagger

go 1.25.1

require (
    dagger.io/dagger v0.18.17
    github.com/getsyntegrity/go-kit-logger v0.0.0-20250828114729-566d9913c10b
    // ... otras dependencias esenciales
)
```

### .github/workflows/release.yml
- Mantener el pipeline completo con:
  - Auto-versioning
  - Multi-platform builds
  - Changelog generation
  - GitHub release creation

## 🚀 Próximos Pasos

1. **Resolver conflictos** usando el método preferido
2. **Verificar** que todo compila y funciona
3. **Hacer commit** de la resolución
4. **Push** a la rama
5. **Crear Pull Request** si es necesario
6. **Merge** a main para activar el primer release

## 🆘 Troubleshooting

### Error: "go mod tidy failed"
```bash
# Limpiar cache de Go
go clean -modcache
go mod download
go mod tidy
```

### Error: "Dagger compatibility issues"
```bash
# Asegurar versión correcta de Dagger
go get dagger.io/dagger@v0.9.11
go mod tidy
```

### Error: "Workflow syntax error"
- Verificar sintaxis YAML del workflow
- Usar un validador YAML online
- Comparar con la versión que funciona

## 📞 Soporte

Si encuentras problemas:
1. Revisa los logs de GitHub Actions
2. Verifica la sintaxis de los archivos
3. Consulta la documentación de GitHub Actions
4. Abre un issue en el repositorio
