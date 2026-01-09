# Deployment Notes

## Docker Build Context Requirement

The `posts-service` requires `web7` module for protobuf generated code. The Docker build context **must** be the project root to include the `web7` directory.

### Required docker-compose.yml Configuration

**Update your deployment `docker-compose.yml`** (e.g., `~/mypage/docker-compose.yml`):

```yaml
posts-service:
  build:
    context: /Users/m.materkov/projects/meme9  # Full path to project root
    dockerfile: ./posts-service/Dockerfile
```

Or if the docker-compose.yml is in the project root:
```yaml
posts-service:
  build:
    context: .
    dockerfile: ./posts-service/Dockerfile
```

### Why?

The `go.mod` has:
```
replace github.com/materkov/meme9/web7 => ../web7
```

This requires `web7` to be accessible at `../web7` relative to the service directory, which means the build context must be the project root.

### Error if Wrong

If using `context: ./posts-service`, you'll get:
```
go: reading /web7/go.mod: open /web7/go.mod: no such file or directory
```

