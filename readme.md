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

### Flags

#### Change port

```bash
plain -p 8080
```

Default: `5000`

#### Change host

```bash
plain -H 192.168.0.103
```

Default: `localhost`

#### Start watcher

```bash
plain -w
```

This will tell plain to watch changes in `pages` directory. If you create or delete file, plain will update routes accordingly.
