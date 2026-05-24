# croncheck

> Parses and audits cron expressions in config files to detect conflicts and misconfigurations.

---

## Installation

```bash
go install github.com/yourusername/croncheck@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/croncheck.git
cd croncheck && go build -o croncheck .
```

---

## Usage

Run `croncheck` against a config file containing cron expressions:

```bash
croncheck --file /etc/cron.d/myjobs
```

**Example output:**

```
[WARN]  Line 4:  "*/5 * * * *"  overlaps with line 7: "0,5,10,... * * * *"
[ERROR] Line 11: "60 * * * *"   invalid value for minute field (0-59)
[OK]    Line 15: "0 2 * * 0"    no issues detected
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the config file to audit | _(required)_ |
| `--format` | Output format: `text`, `json` | `text` |
| `--strict` | Treat warnings as errors | `false` |

---

## What It Detects

- Invalid field values (e.g., minute `> 59`)
- Overlapping or duplicate schedules
- Overly aggressive intervals (e.g., `* * * * *`)
- Unreachable expressions (e.g., `0 0 31 2 *`)

---

## License

MIT © 2024 yourusername