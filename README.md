# CLI2UI

CLI2UI is a simple tool that converts a CLI command to a fully functional web UI.

## Usage

First one has to compile the frontend:

```bash
cd ui/
yarn
yarn build
```

Then compile CLI2UI:

```bash
go build -o cli2ui cmd/main.go
```

Simply run the binary with a well-defined config file (samples are found under `/samples`, JSON schema is also provided `cli.schema.json`):

```bash
./cli2ui <path-to-config>
```

CLI2UI supports both JSON and YAML as the config file; feel free to test out the samples under `/samples` to get a feel of how it works.

## Development

```
.
├── cmd // Entry of the app.
├── patch // Auto-generated patch using `sunmao-ui`.
├── pkg
│   ├── config // Config file parser and script constructor. If the generated script is found wrong or malformed, this is the place to look at.
│   ├── executor // Script executor. It is now stable for the most part, but in case of any bugs with GoRoutine management, read this.
│   └── ui // Implementations of different styles of UIs. One can add new style, change how events are handled, or define new components here.
│       ├── flat
│       └── naive
├── samples // As the name suggests.
├── test
│   └── config
└── ui // Frontend implementation.
    └── src
        ├── application
        ├── assets
        ├── editor
        ├── modules
        └── sunmao
            └── components
```
