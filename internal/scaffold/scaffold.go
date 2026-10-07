package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuyacarlo/gitforge-toolkit/internal/config"
)

func PreCommitTemplate(lang config.StackType) string {
	base := `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.6.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
`
	switch lang {
	case config.StackGo:
		return base + `  - repo: https://github.com/golangci/golangci-lint
    rev: v1.59.1
    hooks:
      - id: golangci-lint
`
	case config.StackPython:
		return base + `  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.5.0
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format
`
	case config.StackRust:
		return base + `  - repo: https://github.com/doublify/pre-commit-rust
    rev: v1.0
    hooks:
      - id: fmt
      - id: cargo-check
      - id: clippy
`
	case config.StackC:
		return base + `  - repo: https://github.com/pocc/pre-commit-hooks
    rev: v1.3.5
    hooks:
      - id: clang-format
        args: [-i]
`
	case config.StackJava:
		return base + `  - repo: https://github.com/macisamuele/language-formatters-pre-commit-hooks
    rev: v2.13.0
    hooks:
      - id: pretty-format-java
        args: [--autofix]
`
	default:
		return base
	}
}

func WorkflowTemplate(name string, lang config.StackType) string {
	var step string
	switch lang {
	case config.StackGo:
		step = `      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Test
        run: go test -v ./...`
	case config.StackPython:
		step = `      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'
      - name: Install uv & test
        run: |
          pip install uv
          uv run pytest`
	case config.StackRust:
		step = `      - name: Set up Rust
        uses: dtolnay/rust-toolchain@stable
      - name: Test
        run: cargo test --all`
	case config.StackC:
		step = `      - name: Build and Test
        run: |
          mkdir -p build && cd build
          cmake .. && make && ctest`
	case config.StackJava:
		step = `      - name: Set up JDK
        uses: actions/setup-java@v4
        with:
          java-version: '21'
          distribution: 'temurin'
      - name: Test with Gradle
        run: ./gradlew test`
	default:
		step = `      - name: Step
        run: echo "CI initialized"`
	}

	return fmt.Sprintf(`name: CI

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  build:
    name: %s
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
%s
`, name, step)
}

func DockerfileTemplate(lang config.StackType) string {
	switch lang {
	case config.StackGo:
		return `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server /app/server
ENTRYPOINT ["/app/server"]
`
	case config.StackPython:
		return `FROM python:3.12-slim
WORKDIR /app
COPY requirements*.txt pyproject.toml ./
RUN pip install --no-cache-dir -r requirements.txt || true
COPY . .
CMD ["python", "main.py"]
`
	case config.StackRust:
		return `FROM rust:1.80-alpine AS builder
WORKDIR /app
COPY Cargo.* ./
COPY src ./src
RUN cargo build --release

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/target/release/app /app/app
ENTRYPOINT ["/app/app"]
`
	case config.StackC:
		return `FROM gcc:latest AS builder
WORKDIR /app
COPY . .
RUN gcc -O2 -o app main.c

FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/app /app/app
ENTRYPOINT ["/app/app"]
`
	case config.StackJava:
		return `FROM gradle:jdk21-alpine AS builder
WORKDIR /app
COPY . .
RUN gradle build --no-daemon

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
`
	default:
		return `FROM alpine:latest
WORKDIR /app
COPY . .
`
	}
}

func RenderProject(targetDir string, cfg *config.ProjectConfig) error {
	if cfg.PreCommit {
		prePath := filepath.Join(targetDir, ".pre-commit-config.yaml")
		if err := os.WriteFile(prePath, []byte(PreCommitTemplate(cfg.Language)), 0644); err != nil {
			return fmt.Errorf("writing pre-commit config: %w", err)
		}
	}

	if cfg.Workflows {
		wfDir := filepath.Join(targetDir, ".github", "workflows")
		if err := os.MkdirAll(wfDir, 0755); err != nil {
			return fmt.Errorf("creating workflow dir: %w", err)
		}
		wfPath := filepath.Join(wfDir, "ci.yml")
		content := WorkflowTemplate(cfg.Name, cfg.Language)
		if err := os.WriteFile(wfPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing workflow file: %w", err)
		}
	}

	if cfg.Container {
		dfPath := filepath.Join(targetDir, "Dockerfile")
		if err := os.WriteFile(dfPath, []byte(DockerfileTemplate(cfg.Language)), 0644); err != nil {
			return fmt.Errorf("writing Dockerfile: %w", err)
		}
	}

	cfgPath := filepath.Join(targetDir, "gitforge.yaml")
	if err := cfg.Save(cfgPath); err != nil {
		return fmt.Errorf("saving gitforge.yaml: %w", err)
	}

	return nil
}
