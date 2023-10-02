# plain

File-based router written in Go

## Install

```bash
go install github.com/tsivinsky/plain/cmd/plain@latest
```

## Usage

```bash
plain -p 5000
```

### Pages

Put html files inside `pages` directory and they will be available at corresponding routes.

### Static files

Create `public` directory and put static file there
