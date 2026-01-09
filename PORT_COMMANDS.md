# Commands to Find Processes Listening on TCP Ports

## 1. Using `lsof` (List Open Files)
```bash
# Find process listening on specific port
lsof -i :8080

# Find process with PID, user, and command
lsof -i :8080 -P -n

# Find all listening TCP ports
lsof -iTCP -sTCP:LISTEN -n -P
```

## 2. Using `ss` (Socket Statistics) - Modern & Recommended
```bash
# Find process listening on specific port
ss -tulpn | grep :8080

# Show all listening TCP ports with process info
ss -tulpn

# Only TCP listening ports
ss -tlnp
```

## 3. Using `netstat` (Older, but widely available)
```bash
# Find process listening on specific port
netstat -tulpn | grep :8080

# Show all listening TCP ports with process info
netstat -tulpn

# Only TCP listening ports
netstat -tlnp
```

## 4. Using `fuser`
```bash
# Find process using port (requires root)
sudo fuser 8080/tcp

# Kill process using port
sudo fuser -k 8080/tcp
```

## 5. Using `sockstat` (FreeBSD/OpenBSD)
```bash
sockstat -4 -l | grep 8080
```

## Quick Reference

**Most common (works on macOS and Linux):**
```bash
lsof -i :8080
```

**Modern Linux (most detailed):**
```bash
ss -tulpn | grep :8080
```

**Kill process on port:**
```bash
# Find PID first
lsof -ti :8080
# Then kill
kill $(lsof -ti :8080)
# Or force kill
kill -9 $(lsof -ti :8080)
```

## Example Output

```bash
$ lsof -i :8080
COMMAND   PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
node    12345 user   23u  IPv4  12345      0t0  TCP *:8080 (LISTEN)
```

```bash
$ ss -tulpn | grep :8080
tcp   LISTEN 0  128  *:8080  *:*  users:(("node",pid=12345,fd=23))
```

