# SSH Config Parser

A Go library for parsing, modifying, and writing SSH config files while preserving comments and formatting.

## Features

- Parse SSH config files with full preservation of formatting, whitespace, and comments
- Modify existing SSH config entries
- Add new host configurations
- Delete host entries
- Write modified configs back to files

## Installation

```go
go get github.com/patrikkj/sshconf
```

## Usage

### Parsing Config Files

```go
// Parse from a file
config, err := sshconf.ParseConfigFile("~/.ssh/config")
if err != nil {
log.Fatal(err)
}

// Parse from a string
configStr := `Host example
    HostName example.com
    User admin`
config := sshconf.ParseConfig(configStr)
```

### Modifying Configs

```go
// Patch (update or append) a host entry
err := config.Patch("Host example", `Host example
    HostName example.com
    User admin
    Port 2222`)

// Delete a host entry
err := config.Delete("Host example")
```

### Writing Config Files

```go
// Write to a file
err := config.WriteFile("~/.ssh/config")

// Write to any io.Writer
err := config.Write(os.Stdout)

// Get config as string
configStr := config.Render()
```

### Raw Parsing Mode

If you need to parse without hierarchical organization (preserving exact file structure):

```go
config := sshconf.ParseConfigRaw(content)
```

## Example

```go
package main

import (
"fmt"
"github.com/patrikkj/sshconf"
)

func main() {
// Parse existing config
config, err := sshconf.ParseConfigFile("~/.ssh/config")
if err != nil {
panic(err)
}

    // Add or update a host entry
    err = config.Patch("Host example", `Host example
    HostName example.com
    User admin
    Port 2222
    IdentityFile ~/.ssh/example_key`)
    if err != nil {
        panic(err)
    }

    // Write the modified config back
    err = config.WriteFile("~/.ssh/config")
    if err != nil {
        panic(err)
    }

}
```

## Features

- Preserves comments and formatting when parsing and writing
- Maintains indentation styles
- Supports all SSH config directives
- Handles multi-level configurations (e.g., Match blocks)

## License

[Your chosen license]
