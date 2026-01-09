# Go Workspace Setup (Option 3)

This project uses Go workspaces to manage local module dependencies without `replace` directives in `go.mod` files.

## Local Development

The workspace is configured in `go.work` at the project root. All services can import from `web7` without needing `replace` directives in their `go.mod` files.

### Workspace Modules

- `web7` - Contains protobuf generated code
- `auth-service`
- `users-service`
- `subscriptions-service`
- `likes-service`
- `posts-service`
- `photos`

### Using the Workspace

1. **Build services:**
   ```bash
   cd likes-service
   go build ./cmd/main.go
   ```

2. **Run tests:**
   ```bash
   cd posts-service
   go test ./...
   ```

3. **Add new modules to workspace:**
   ```bash
   go work use ./new-service
   ```

4. **Sync workspace:**
   ```bash
   go work sync
   ```

## Docker Builds

**Important:** Go workspaces don't work in Docker builds. The solution:

1. Keep `replace` directive in `go.mod` files (for Docker builds)
2. Workspace overrides `replace` for local development
3. Docker build context must include `web7` directory (use project root as context)

**Note:** The `replace` directive in `go.mod` is ignored when using workspaces locally, but is used in Docker builds.

### Docker Build Process

```dockerfile
# Copy web7 first (needed for replace directive in Docker)
COPY web7 /web7
# Copy service files
COPY posts-service/go.mod .
COPY posts-service/go.sum .
# Add replace directive for Docker build (workspaces don't work in Docker)
RUN echo "replace github.com/materkov/meme9/web7 => ../web7" >> go.mod
```

## Benefits

✅ No `replace` directives in `go.mod` files (for local development)
✅ Clean module dependencies
✅ Native Go solution (Go 1.18+)
✅ All services use same workspace configuration

## Limitations

⚠️ Workspaces don't work in Docker builds
⚠️ Need to add `replace` directive in Dockerfiles
⚠️ Need to copy `web7/` into Docker build context

## Migration Notes

- Removed `replace` directives from:
  - `likes-service/go.mod`
  - `posts-service/go.mod`
- Other services don't have `replace` directives (they don't use web7 directly)
- `go.work` file should be committed to the repository

