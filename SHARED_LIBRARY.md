# Syntegrity Dagger - Shared Library

Syntegrity Dagger es una librería compartida que proporciona pipelines de CI/CD unificados para proyectos Go. Esta librería se distribuye como un binario pre-compilado que otros servicios pueden descargar y usar en sus propios pipelines.

## 🚀 Instalación Rápida

### Opción 1: Script de Instalación Automática

```bash
# Instalar la última versión
curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash

# Instalar una versión específica
curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash -s -- -v v1.0.0

# Instalar en un directorio personalizado
curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash -s -- -d ~/bin
```

### Opción 2: Descarga Manual

```bash
# Detectar tu plataforma y descargar el binario apropiado
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/')
VERSION="v1.0.0"  # o "latest" para la última versión

# Descargar
curl -L "https://github.com/getsyntegrity/syntegrity-dagger/releases/download/${VERSION}/syntegrity-dagger-${PLATFORM}" -o syntegrity-dagger
chmod +x syntegrity-dagger
```

## 📦 Uso en CI/CD

### GitHub Actions

```yaml
name: CI/CD Pipeline
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v5
    
    # Instalar Syntegrity Dagger
    - name: Install Syntegrity Dagger
      run: |
        curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash
    
    # Usar en el pipeline
    - name: Run Pipeline
      run: |
        syntegrity-dagger --pipeline go-kit --env dev --coverage 90
```

### GitLab CI

```yaml
stages:
  - build
  - test
  - deploy

variables:
  SYNTERGRITY_VERSION: "v1.0.0"

before_script:
  # Instalar Syntegrity Dagger
  - curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash -s -- -v $SYNTERGRITY_VERSION

build:
  stage: build
  script:
    - syntegrity-dagger --pipeline go-kit --only-build
```

### Jenkins

```groovy
pipeline {
    agent any
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash
                '''
            }
        }
        
        stage('Build') {
            steps {
                sh 'syntegrity-dagger --pipeline go-kit --only-build'
            }
        }
        
        stage('Test') {
            steps {
                sh 'syntegrity-dagger --pipeline go-kit --only-test --coverage 90'
            }
        }
    }
}
```

## 🔧 Comandos Disponibles

### Comandos Básicos

```bash
# Mostrar ayuda
syntegrity-dagger --help

# Mostrar versión
syntegrity-dagger --version

# Listar pipelines disponibles
syntegrity-dagger --list-pipelines

# Listar pasos de un pipeline
syntegrity-dagger --list-steps --pipeline go-kit
```

### Ejecutar Pipelines

```bash
# Pipeline completo
syntegrity-dagger --pipeline go-kit --env dev --coverage 90

# Solo build
syntegrity-dagger --pipeline go-kit --only-build

# Solo tests
syntegrity-dagger --pipeline go-kit --only-test --coverage 85

# Ejecutar paso específico
syntegrity-dagger --pipeline go-kit --step build

# Pipeline local (sin Docker)
syntegrity-dagger --pipeline go-kit --local
```

### Configuración Avanzada

```bash
# Con archivo de configuración YAML
syntegrity-dagger --config .syntegrity-dagger.yml

# Con variables de entorno
export SYNTERGRITY_ENV=staging
export SYNTERGRITY_COVERAGE=95
syntegrity-dagger --pipeline go-kit
```

## 📋 Pipelines Disponibles

### go-kit
Pipeline optimizado para servicios Go con arquitectura go-kit.

```bash
syntegrity-dagger --pipeline go-kit --env dev
```

### docker-go
Pipeline para aplicaciones Go con contenedores Docker.

```bash
syntegrity-dagger --pipeline docker-go --env prod
```

### infra
Pipeline para infraestructura y deployments.

```bash
syntegrity-dagger --pipeline infra --env staging
```

## 🔄 Actualizaciones

### Verificar Versión Actual

```bash
syntegrity-dagger --version
```

### Actualizar a la Última Versión

```bash
curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash
```

### Actualizar a Versión Específica

```bash
curl -fsSL https://raw.githubusercontent.com/getsyntegrity/syntegrity-dagger/main/install.sh | bash -s -- -v v1.1.0
```

## 🛠️ Configuración

### Variables de Entorno

```bash
# Configuración básica
export SYNTERGRITY_ENV=dev
export SYNTERGRITY_COVERAGE=90
export SYNTERGRITY_BRANCH=develop

# Configuración de Git
export SYNTERGRITY_GIT_AUTH=ssh  # o https
export SYNTERGRITY_GIT_REF=main

# Configuración de Docker
export SYNTERGRITY_SKIP_PUSH=false
export SYNTERGRITY_REGISTRY=registry.example.com
```

### Archivo de Configuración YAML

Crea un archivo `.syntegrity-dagger.yml`:

```yaml
pipeline:
  name: go-kit
  coverage: 90
  skip_push: false
  only_build: false
  only_test: false
  verbose: true

environment: dev

git:
  ref: main
  protocol: ssh

steps:
  - name: setup
    required: true
    timeout: 5m
  - name: build
    required: true
    timeout: 10m
  - name: test
    required: true
    timeout: 15m
```

## 🔗 Integración con Otros Servicios

### Webhook para Notificaciones

```bash
# Configurar webhook para notificaciones de release
curl -X POST https://your-service.com/webhook/syntegrity-update \
  -H "Content-Type: application/json" \
  -d '{"version": "v1.0.0", "download_url": "https://github.com/getsyntegrity/syntegrity-dagger/releases/download/v1.0.0/syntegrity-dagger-linux-amd64"}'
```

### API para Verificar Versiones

```bash
# Verificar última versión disponible
curl -s https://api.github.com/repos/getsyntegrity/syntegrity-dagger/releases/latest | jq '.tag_name'

# Verificar si hay actualizaciones
CURRENT_VERSION=$(syntegrity-dagger --version)
LATEST_VERSION=$(curl -s https://api.github.com/repos/getsyntegrity/syntegrity-dagger/releases/latest | jq -r '.tag_name')

if [[ "$CURRENT_VERSION" != "$LATEST_VERSION" ]]; then
    echo "Nueva versión disponible: $LATEST_VERSION"
fi
```

## 🐛 Troubleshooting

### Problemas Comunes

1. **Error de permisos**: Asegúrate de que el binario tenga permisos de ejecución
   ```bash
   chmod +x syntegrity-dagger
   ```

2. **Error de conectividad**: Verifica que tengas acceso a GitHub
   ```bash
   curl -I https://github.com/getsyntegrity/syntegrity-dagger/releases/latest
   ```

3. **Versión incorrecta**: Verifica la versión instalada
   ```bash
   syntegrity-dagger --version
   ```

### Logs y Debug

```bash
# Ejecutar con verbose para más información
syntegrity-dagger --pipeline go-kit --verbose

# Ver logs detallados
syntegrity-dagger --pipeline go-kit --verbose 2>&1 | tee pipeline.log
```

## 📞 Soporte

- **Issues**: [GitHub Issues](https://github.com/getsyntegrity/syntegrity-dagger/issues)
- **Documentación**: [Wiki](https://github.com/getsyntegrity/syntegrity-dagger/wiki)
- **Releases**: [GitHub Releases](https://github.com/getsyntegrity/syntegrity-dagger/releases)

## 📄 Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.
