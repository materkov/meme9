# Options to Refactor Proto Schemas and Remove `replace` Directives

## Current Problem
Services have `replace github.com/materkov/meme9/web7 => ../web7` in their `go.mod` files because:
- Protobuf files are generated in `web7/pb/`
- Services import from `github.com/materkov/meme9/web7/pb/...`
- `web7` is a local module, not published
- Docker builds fail because `../web7` isn't available in build context

## Option 1: Create Dedicated `api` Module (RECOMMENDED) ⭐

**Structure:**
```
api/
  go.mod                    # module github.com/materkov/meme9/api
  pb/
    github.com/materkov/meme9/api/
      auth/
      users/
      posts/
      ...
```

**Changes:**
1. Create new `api/` directory with its own `go.mod`
2. Update `proto_gen.sh` to generate into `api/pb/`
3. Update `go_package` in `.proto` files to `github.com/materkov/meme9/api`
4. Services import from `github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/...`
5. Publish `api` module or use it as a local module that can be copied into Docker builds

**Pros:**
- Clean separation of concerns
- `api` module can be versioned independently
- No `replace` directives needed
- Can be published to a private Go module registry
- Docker builds work (can copy `api/` into build context)

**Cons:**
- Requires refactoring imports across all services
- Need to maintain another module

**Implementation:**
```bash
# 1. Create api module
mkdir -p api/pb
cd api
go mod init github.com/materkov/meme9/api

# 2. Update proto_gen.sh to generate into api/pb/
# 3. Update all .proto files go_package option
# 4. Update all service imports
# 5. Remove replace directives
```

---

## Option 2: Generate Protobuf in Each Service

**Structure:**
```
schema/
  *.proto
posts-service/
  pb/              # Generated per service
    github.com/materkov/meme9/api/...
  go.mod
```

**Changes:**
1. Each service generates its own protobuf files
2. Update `proto_gen.sh` to generate into each service's `pb/` directory
3. Services import from their own `pb/` directory
4. No cross-service dependency on `web7`

**Pros:**
- No shared dependency
- Each service is self-contained
- No `replace` directives needed
- Docker builds work (everything is in service directory)

**Cons:**
- Duplicated generated code across services
- Need to regenerate in each service
- More maintenance (multiple generation steps)
- Type mismatches if services use different versions

**Implementation:**
```bash
# Update proto_gen.sh to generate into each service:
protoc -I schema --go_out=posts-service/pb --twirp_out=posts-service/pb schema/*.proto
protoc -I schema --go_out=likes-service/pb --twirp_out=likes-service/pb schema/*.proto
# etc.
```

---

## Option 3: Use Go Workspaces

**Structure:**
```
go.work
web7/
  go.mod
posts-service/
  go.mod
likes-service/
  go.mod
```

**Changes:**
1. Create `go.work` file at project root
2. Add all modules to workspace
3. Remove `replace` directives (workspace handles it)
4. Services can import from `web7` directly

**Pros:**
- Native Go solution (Go 1.18+)
- No `replace` directives needed
- Works for local development

**Cons:**
- **Doesn't work in Docker builds** (workspaces aren't supported in Docker)
- Still need `replace` for Docker or copy `web7` into build
- Not a complete solution

**Implementation:**
```bash
# Create go.work
go work init web7 posts-service likes-service ...

# Remove replace directives from go.mod files
# Works for local dev, but Docker still needs web7
```

---

## Option 4: Publish `web7` as a Go Module

**Changes:**
1. Tag `web7` with version (e.g., `v0.1.0`)
2. Push to Git repository
3. Services import from published module
4. Remove `replace` directives

**Pros:**
- Standard Go module approach
- Versioned dependencies
- Works in Docker builds

**Cons:**
- Requires Git repository access
- Need to tag versions
- Need to push changes to access in Docker
- More complex CI/CD

**Implementation:**
```bash
# In web7/
git tag v0.1.0
git push origin v0.1.0

# In services, update go.mod:
require github.com/materkov/meme9/web7 v0.1.0
```

---

## Option 5: Vendor Protobuf Files

**Structure:**
```
posts-service/
  vendor/
    github.com/materkov/meme9/web7/pb/...
  go.mod
```

**Changes:**
1. Vendor `web7/pb/` into each service
2. Use `go mod vendor`
3. Services import from vendor directory

**Pros:**
- No external dependencies
- Works in Docker builds
- Reproducible builds

**Cons:**
- Large vendor directories
- Need to update vendor when protos change
- Not ideal for microservices

---

## Option 6: Change Build Context to Project Root

**Changes:**
1. Update `docker-compose.yml` to use project root as build context
2. Update Dockerfiles to copy `web7/` into build
3. Keep `replace` directives

**Pros:**
- Minimal code changes
- Works with current structure

**Cons:**
- Still uses `replace` directives
- Larger Docker build contexts
- Need to update deploy scripts

**Implementation:**
```dockerfile
# Dockerfile
FROM golang:1.24.9-alpine3.22
WORKDIR /build
COPY web7 /web7
COPY posts-service/go.mod posts-service/go.sum ./
COPY posts-service ./
RUN go mod download
RUN go build -o /app ./cmd
```

---

## Recommendation

**Option 1 (Dedicated `api` Module)** is the best long-term solution:
- Clean architecture
- Proper module boundaries
- Works in Docker (can copy `api/` directory)
- Can be published later if needed
- No `replace` directives

**Quick Fix (Option 6)** for immediate Docker build fix:
- Update build contexts and Dockerfiles
- Keep `replace` for now
- Plan migration to Option 1

---

## Migration Path (Option 1)

1. **Create `api` module:**
   ```bash
   mkdir -p api/pb
   cd api
   go mod init github.com/materkov/meme9/api
   ```

2. **Update `.proto` files:**
   ```protobuf
   option go_package = "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api";
   ```

3. **Update `proto_gen.sh`:**
   ```bash
   protoc -I schema --go_out=api/pb --twirp_out=api/pb schema/*.proto
   ```

4. **Update service imports:**
   ```go
   // Old:
   import "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/posts"
   
   // New:
   import "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/posts"
   ```

5. **Update service `go.mod`:**
   ```go
   require github.com/materkov/meme9/api v0.0.0
   replace github.com/materkov/meme9/api => ../api
   ```

6. **For Docker builds:**
   - Copy `api/` into build context, or
   - Publish `api` module to a registry

