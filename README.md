# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <file1> <file2>
```

**Example:**

```bash
envdiff .env.development .env.production
```

**Sample output:**

```
MISSING IN .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISSING IN .env.development:
  - CDN_BASE_URL

MISMATCHED VALUES:
  - LOG_LEVEL: "debug" vs "error"
  - PORT: "3000" vs "8080"
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare keys only, ignore values |
| `--quiet` | Exit with non-zero code if differences found, no output |
| `--json` | Output results as JSON |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)