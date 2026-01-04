# csp

Minimal CSP violation report collector. Parses both legacy `report-uri` and modern `report-to` formats.

## Install

Download from [releases](https://github.com/coalaura/csp/releases) or build:

```bash
go build -o csp .
```

## Usage

```bash
./csp
```

Reports are logged to stdout:

```
2026-01-04T02:30:00 CSP [script-src] https://example.com/ blocked https://evil.com/bad.js
```

## nginx

```nginx
more_set_headers "Content-Security-Policy: default-src 'self'; ...; report-uri https://your-domain.com/report;";
```

Proxy to the collector:

```nginx
location = /report {
    proxy_pass http://127.0.0.1:9393;
}
```

## Supported Formats

| Format | Browsers | Header |
|--------|----------|--------|
| Legacy | All | `report-uri` |
| Modern | Chrome/Edge | `report-to` + `Reporting-Endpoints` |

Use `report-uri` alone for universal support. Chrome falls back to it when `report-to` is absent.

## License

MIT
