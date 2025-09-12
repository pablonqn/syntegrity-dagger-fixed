# GitHub Actions Workflows

Este directorio contiene los workflows de GitHub Actions para el proyecto Syntegrity Dagger.

## Workflows Disponibles

### 1. CI/CD Pipeline (`ci.yml`)
**Trigger:** Push y Pull Requests a `main` y `develop`

Ejecuta:
- **Lint**: Verificación de código con golangci-lint
- **Test**: Ejecución de tests unitarios
- **Build**: Compilación del binario
- **Security**: Escaneo de seguridad con Gosec

### 2. Test Suite (`test.yml`)
**Trigger:** Push y Pull Requests a `main` y `develop`

Ejecuta:
- Tests con coverage
- Validación de umbral de coverage (90%)
- Generación de reportes de coverage
- Upload de reportes a Codecov

### 3. Release (`release.yml`)
**Trigger:** 
- Push de tags (formato `v*`)
- Manual dispatch

Ejecuta:
- Creación de releases con GoReleaser
- Verificaciones pre-release
- Creación automática de tags (manual)

### 4. Auto Tag (`auto-tag.yml`)
**Trigger:**
- Push a `main`
- Manual dispatch

Ejecuta:
- Creación automática de tags semánticos
- Creación de releases de GitHub
- Build y upload de binarios multi-plataforma

### 5. Dependabot (`dependabot.yml`)
**Trigger:** Pull Requests de Dependabot

Ejecuta:
- Auto-merge de dependencias menores y patches
- Auto-aprobación de PRs de Dependabot

## Configuración

### Secrets Requeridos
- `GITHUB_TOKEN`: Token automático de GitHub (ya configurado)

### Variables de Entorno
- `GO_VERSION`: Versión de Go (1.25.1)
- `COVERAGE_THRESHOLD`: Umbral de coverage (90%)

## Uso

### Crear un Release Manual
1. Ve a Actions → Release
2. Click en "Run workflow"
3. Ingresa el tag (ej: `v1.0.0`)
4. El workflow creará el tag y release

### Crear Tag Automático
1. Haz merge a `main`
2. El workflow `auto-tag.yml` se ejecutará automáticamente
3. Creará el siguiente tag semántico (patch por defecto)

### Crear Tag Manual con Versión Específica
1. Ve a Actions → Auto Tag
2. Click en "Run workflow"
3. Selecciona el tipo de versión (patch/minor/major)
4. El workflow creará el tag correspondiente

## Estructura de Tags

Los tags siguen [Semantic Versioning](https://semver.org/):
- `v1.0.0` - Major release
- `v1.1.0` - Minor release  
- `v1.1.1` - Patch release

## Binarios Generados

Cada release incluye binarios para:
- Linux AMD64
- macOS AMD64
- macOS ARM64
- Windows AMD64

## Coverage

- **Umbral mínimo**: 90%
- **Reportes**: Generados en formato HTML y texto
- **Upload**: Automático a Codecov
- **Artifacts**: Disponibles por 30 días

## Dependencias

### Automatización
- **Dependabot**: Actualiza dependencias semanalmente
- **Auto-merge**: Patches y minor updates se auto-mergean
- **Labels**: Automáticos para organización

### Herramientas
- **golangci-lint**: Linting de código
- **Gosec**: Escaneo de seguridad
- **GoReleaser**: Generación de releases
- **Codecov**: Reportes de coverage

## Troubleshooting

### Coverage Falla
- Verifica que el coverage esté por encima del 90%
- Revisa los reportes en la sección de artifacts

### Release Falla
- Verifica que el tag no exista previamente
- Asegúrate de que el formato del tag sea correcto (`v1.0.0`)

### Build Falla
- Verifica que todas las dependencias estén actualizadas
- Revisa los logs de linting y tests
