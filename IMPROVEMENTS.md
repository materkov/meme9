# Code Improvement Suggestions

This document outlines various improvements that can be made to the codebase to enhance code quality, security, performance, and maintainability.

## 🔴 Critical Issues

### 1. **Security: Authentication Middleware Doesn't Validate Tokens**
**Location**: All services' `api/utils.go`

**Problem**: The `AuthMiddleware` simply reads `x-user-id` from headers without validating the token. This means any client can set this header and impersonate any user.

**Current Code**:
```go
func AuthMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), UserIDKey, r.Header.Get("x-user-id"))
        handler.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Recommendation**: 
- Validate the token with auth-service before trusting the `x-user-id` header
- Or implement proper token validation in a shared middleware
- Consider using JWT tokens with signature verification

### 2. **HTTP Clients Without Timeouts**
**Location**: `posts-service/api/service.go`, `likes-service/api/api.go`

**Problem**: HTTP clients are created without timeouts, which can lead to hanging requests and resource exhaustion.

**Current Code**:
```go
return usersapi.NewUsersProtobufClient(usersServiceURL, &http.Client{})
```

**Recommendation**:
```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
return usersapi.NewUsersProtobufClient(usersServiceURL, httpClient)
```

### 3. **No Graceful Shutdown**
**Location**: All `cmd/main.go` files

**Problem**: Services don't handle shutdown signals gracefully, which can lead to data loss and connection issues.

**Recommendation**: Implement graceful shutdown with signal handling:
```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer cancel()

srv := &http.Server{
    Addr:    addr,
    Handler: handler,
}

go func() {
    <-ctx.Done()
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()
    srv.Shutdown(shutdownCtx)
}()
```

## 🟡 High Priority Issues

### 4. **Code Duplication: Shared Utilities**
**Location**: Multiple `api/utils.go` files across services

**Problem**: `AuthMiddleware`, `GetUserIDFromContext`, and `ErrAuthRequired` are duplicated across all services.

**Recommendation**: 
- Create a shared `common` or `shared` package
- Move common utilities there
- Import from shared package in all services

### 5. **Service Clients Created Per Request**
**Location**: `posts-service/api/service.go` lines 40-54

**Problem**: New HTTP clients are created on every request, preventing connection reuse.

**Current Code**:
```go
func (s *Service) getUsersServiceClient() usersapi.Users {
    // Creates new client every time
}
```

**Recommendation**: 
- Store clients as fields in the Service struct
- Initialize once in `NewService` or use dependency injection
- Reuse the same client instance

### 6. **Inconsistent Function Naming**
**Location**: `likes-service/api/utils.go`

**Problem**: `getUserIDFromContext` is lowercase while other services use `GetUserIDFromContext` (exported).

**Recommendation**: Standardize on `GetUserIDFromContext` across all services.

### 7. **Missing Pagination Limits**
**Location**: `posts-service/api/service.go` - `GetFeed` method

**Problem**: `GetAll()` can return unlimited posts, causing memory and performance issues.

**Recommendation**:
- Add pagination with `limit` and `offset` or cursor-based pagination
- Set maximum limits (e.g., 100 posts per request)
- Implement proper pagination in adapters

### 8. **Incomplete Feed Implementation**
**Location**: `posts-service/api/service.go` line 220-223

**Problem**: The subscriptions feed uses `GetAll()` instead of filtering by followed users.

**Current Code**:
```go
// Note: We need a GetFollowing method in subscriptions service
// For now, we'll use GetAll
postsList, err = s.posts.GetAll(ctx)
```

**Recommendation**: 
- Implement `GetFollowing` in subscriptions service
- Filter posts by followed user IDs
- This is a TODO that should be completed

### 9. **No Request Context Timeouts**
**Location**: Service-to-service calls

**Problem**: Service calls don't propagate or respect request timeouts.

**Recommendation**: 
- Use request context with timeout for all service calls
- Set appropriate timeouts based on operation type
- Use `context.WithTimeout(ctx, 5*time.Second)` for service calls

### 10. **MongoDB Connection Pool Not Configured**
**Location**: All `cmd/main.go` files

**Problem**: MongoDB client uses default connection pool settings which may not be optimal.

**Recommendation**:
```go
clientOptions := options.Client().
    ApplyURI(mongoURI).
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetMaxConnIdleTime(30 * time.Second)
```

## 🟢 Medium Priority Issues

### 11. **Error Handling Inconsistencies**
**Location**: Various service files

**Problem**: Some errors are wrapped with `fmt.Errorf`, others are returned directly, and some use `twirp.NewError` inconsistently.

**Recommendation**: 
- Standardize error handling patterns
- Always wrap internal errors with context
- Use appropriate Twirp error codes consistently

### 12. **Hardcoded Default Values**
**Location**: All `cmd/main.go` files

**Problem**: Default MongoDB URIs, ports, and other values are hardcoded.

**Recommendation**: 
- Use configuration structs
- Consider using a config library (viper, envconfig)
- Document all environment variables

### 13. **Missing Input Validation**
**Location**: Various service methods

**Problem**: Some inputs lack proper validation (e.g., username length, password strength, URL validation).

**Recommendation**:
- Add comprehensive input validation
- Validate username format and length
- Validate password strength requirements
- Validate URLs before storing

### 14. **No Rate Limiting**
**Location**: All services

**Problem**: No protection against abuse or DoS attacks.

**Recommendation**: 
- Implement rate limiting middleware
- Use per-user and per-IP rate limits
- Consider using a library like `golang.org/x/time/rate`

### 15. **Missing Observability**
**Location**: All services

**Problem**: Limited logging, no metrics, no distributed tracing.

**Recommendation**:
- Add structured logging (zap, logrus)
- Add metrics (Prometheus)
- Add distributed tracing (OpenTelemetry, Jaeger)
- Add request ID tracking

### 16. **No Health Check Endpoints**
**Location**: All services

**Problem**: No way to check if services are healthy.

**Recommendation**: 
- Add `/health` and `/ready` endpoints
- Check database connectivity in readiness probe
- Return appropriate HTTP status codes

### 17. **Token Expiration Not Implemented**
**Location**: `auth-service/api/service.go`

**Problem**: Tokens never expire, creating security risk.

**Recommendation**:
- Add expiration time to tokens
- Implement token refresh mechanism
- Clean up expired tokens periodically

### 18. **No Database Transaction Support**
**Location**: Adapters

**Problem**: Multi-step operations aren't atomic.

**Recommendation**: 
- Use MongoDB transactions for operations that need atomicity
- Wrap related operations in transactions

### 19. **Missing Indexes Documentation**
**Location**: Adapters with `EnsureIndexes`

**Problem**: Not all necessary indexes may be created.

**Recommendation**: 
- Document all required indexes
- Add indexes for frequently queried fields
- Consider compound indexes for common query patterns

### 20. **No Request Size Limits**
**Location**: All services

**Problem**: No limits on request body size, allowing potential DoS.

**Recommendation**: 
- Set `MaxRequestSize` on HTTP servers
- Validate file upload sizes
- Set appropriate limits for different endpoints

## 🔵 Low Priority / Nice to Have

### 21. **Use Structured Logging**
Replace `log.Printf` with structured logging library.

### 22. **Add Request ID Middleware**
Track requests across services with unique IDs.

### 23. **Implement Caching**
Cache frequently accessed data (user info, post counts).

### 24. **Add Database Migrations**
Use a migration tool for schema changes.

### 25. **Improve Test Coverage**
Add more integration and end-to-end tests.

### 26. **Add API Versioning**
Plan for future API changes with versioning strategy.

### 27. **Implement Retry Logic**
Add retry logic for transient failures in service-to-service calls.

### 28. **Add Circuit Breakers**
Protect against cascading failures with circuit breakers.

### 29. **Documentation**
- Add API documentation (OpenAPI/Swagger)
- Document deployment procedures
- Add architecture diagrams

### 30. **Code Organization**
- Consider using a more structured project layout
- Separate business logic from transport layer more clearly

## Summary

**Priority Actions:**
1. Fix authentication security issue (Critical)
2. Add HTTP client timeouts (Critical)
3. Implement graceful shutdown (Critical)
4. Create shared utilities package (High)
5. Fix service client reuse (High)
6. Complete subscriptions feed implementation (High)
7. Add pagination limits (High)

**Estimated Impact:**
- **Security**: Fixing auth middleware will significantly improve security
- **Reliability**: Timeouts and graceful shutdown will improve stability
- **Performance**: Client reuse and connection pooling will improve latency
- **Maintainability**: Shared utilities will reduce duplication
