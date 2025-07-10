[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/firewut/go-json-map) 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/firewut/go-json-map/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/firewut/go-json-map)](https://goreportcard.com/report/github.com/firewut/go-json-map)

# go-json-map

> Navigate and manipulate Go maps with simple dot notation paths, just like working with JSON!

## Table of Contents

- [Installation](#installation)
- [Why go-json-map?](#why-go-json-map)
- [Key Features](#key-features)
- [Usage](#usage)
  - [Get Property](#get-property)
  - [Create Property](#create-property)
  - [Update Property](#update-property)
  - [Delete Property](#delete-property)
- [Custom Separators](#custom-separators)
  - [Why Use Custom Separators?](#why-use-custom-separators)
  - [Working with Email Addresses or URLs](#working-with-email-addresses-or-urls)
  - [Empty Separator Behavior](#empty-separator-behavior)
- [Common Use Cases](#common-use-cases)
- [Important Notes](#important-notes)
- [API Reference](#api-reference)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/firewut/go-json-map
```

## Why go-json-map?

Working with deeply nested map structures in Go can be painful. This library makes it simple:

❌ **Without go-json-map:**
```go
// Accessing nested values
if val1, ok := data["one"].(map[string]interface{}); ok {
    if val2, ok := val1["two"].(map[string]interface{}); ok {
        if arr, ok := val2["three"].([]interface{}); ok {
            if len(arr) > 0 {
                value := arr[0]
            }
        }
    }
}
```

✅ **With go-json-map:**
```go
value, err := gjm.GetProperty(data, "one.two.three[0]")
```

## Key Features

- 🎯 **Simple dot notation** for nested access (`"user.profile.name"`)
- 🔢 **Array indexing** support (`"items[2].price"`)
- 🔧 **Custom separators** - use `/`, `::`, or any delimiter you prefer
- 📦 **Zero dependencies** - pure Go standard library
- 🛡️ **Type-safe operations** with proper error handling
- 🚀 **Fast and lightweight** - minimal overhead

## Usage

### Import the package

```go
import gjm "github.com/firewut/go-json-map"
```

### Sample data structure

```go
document := map[string]interface{}{
    "user": map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "John Doe",
            "scores": []int{100, 200, 300},
        },
        "settings": map[string]interface{}{
            "theme": "dark",
            "notifications": true,
        },
    },
}
```

### Get Property

Retrieve values from nested structures:

```go
import (
    "fmt"
    "log"
    gjm "github.com/firewut/go-json-map"
)

// Get a simple value
name, err := gjm.GetProperty(document, "user.profile.name")
if err != nil {
    log.Fatal(err)
}
fmt.Println(name) // Output: John Doe

// Get array element
score, err := gjm.GetProperty(document, "user.profile.scores[1]")
fmt.Println(score) // Output: 200

// Get entire array
scores, err := gjm.GetProperty(document, "user.profile.scores")
fmt.Println(scores) // Output: [100 200 300]

// Use custom separator
theme, err := gjm.GetProperty(document, "user/settings/theme", "/")
fmt.Println(theme) // Output: dark
```

### Create Property

Create new properties (fails if property already exists):

```go
import (
    "log"
    gjm "github.com/firewut/go-json-map"
)

// Create a new property
err := gjm.CreateProperty(document, "user.profile.email", "john@example.com")
if err != nil {
    log.Fatal(err)
}

// Create in array (creates array if needed)
err = gjm.CreateProperty(document, "user.profile.tags[0]", "premium")

// Create nested structure
err = gjm.CreateProperty(document, "user.profile.address.city", "New York")
```

### Update Property

Update existing properties or create new ones:

```go
import (
    "time"
    gjm "github.com/firewut/go-json-map"
)

// Update existing value
err := gjm.UpdateProperty(document, "user.settings.theme", "light")

// Update array element
err = gjm.UpdateProperty(document, "user.profile.scores[0]", 150)

// Create or update - won't fail if property doesn't exist
err = gjm.UpdateProperty(document, "user.lastLogin", time.Now())
```

### Delete Property

Remove properties from the map:

```go
// Delete a property
err := gjm.DeleteProperty(document, "user.settings.theme")

// Delete array element (shifts remaining elements)
err = gjm.DeleteProperty(document, "user.profile.scores[1]")

// Delete entire nested structure
err = gjm.DeleteProperty(document, "user.profile")
```

## Custom Separators

### Why Use Custom Separators?

The default separator is `.` (dot), but you can use any string as a separator. This is particularly useful when your property names contain dots or other special characters.

```go
import (
    "fmt"
    gjm "github.com/firewut/go-json-map"
)

// Data with dots in property names
config := map[string]interface{}{
    "servers": map[string]interface{}{
        "api.example.com": map[string]interface{}{
            "host": "192.168.1.100",
            "port": 8080,
        },
        "db.example.com": map[string]interface{}{
            "host": "192.168.1.101",
            "port": 5432,
        },
    },
}

// Use "/" separator to access properties with dots
host, err := gjm.GetProperty(config, "servers/api.example.com/host", "/")
fmt.Println(host) // Output: 192.168.1.100

// Or use "::" for clearer nesting
port, err := gjm.GetProperty(config, "servers::db.example.com::port", "::")
fmt.Println(port) // Output: 5432
```

### Working with Email Addresses or URLs

```go
import (
    "fmt"
    gjm "github.com/firewut/go-json-map"
)

users := map[string]interface{}{
    "accounts": map[string]interface{}{
        "john.doe@example.com": map[string]interface{}{
            "name": "John Doe",
            "role": "admin",
        },
    },
}

// Use a custom separator to handle email addresses
name, err := gjm.GetProperty(users, "accounts|john.doe@example.com|name", "|")
fmt.Println(name) // Output: John Doe
```

### Empty Separator Behavior

**Important**: Passing an empty string `""` as separator will **not** treat the entire path as a single property name. Instead, it defaults to using `.` as the separator:

```go
import gjm "github.com/firewut/go-json-map"

// These two calls are equivalent:
val1, _ := gjm.GetProperty(data, "one.two.three", "")    // empty separator
val2, _ := gjm.GetProperty(data, "one.two.three", ".")   // explicit dot
// val1 == val2
```

If you need to access a property that literally contains dots in its name, use a different separator that doesn't appear in your property names.

## Common Use Cases

### Configuration Management

```go
import (
    "os"
    gjm "github.com/firewut/go-json-map"
)

config := make(map[string]interface{})

// Build configuration dynamically
gjm.UpdateProperty(config, "database.host", "localhost")
gjm.UpdateProperty(config, "database.port", 5432)
gjm.UpdateProperty(config, "database.credentials.username", os.Getenv("DB_USER"))
gjm.UpdateProperty(config, "features.authentication.enabled", true)
gjm.UpdateProperty(config, "features.authentication.providers[0]", "google")
gjm.UpdateProperty(config, "features.authentication.providers[1]", "github")
```

### API Response Manipulation

```go
import (
    "time"
    gjm "github.com/firewut/go-json-map"
)

// Modify API response before sending
response := getAPIResponse()

// Add metadata
gjm.UpdateProperty(response, "meta.version", "1.0")
gjm.UpdateProperty(response, "meta.timestamp", time.Now().Unix())

// Remove sensitive data
gjm.DeleteProperty(response, "data.users[0].password")
gjm.DeleteProperty(response, "data.users[0].ssn")
```

### Dynamic Form Processing

```go
import (
    "fmt"
    gjm "github.com/firewut/go-json-map"
)

formData := make(map[string]interface{})

// Process form fields with dynamic paths
for fieldPath, value := range request.Form {
    gjm.UpdateProperty(formData, fieldPath, value)
}

// Validate required fields
required := []string{"user.email", "user.name", "preferences.newsletter"}
for _, path := range required {
    if _, err := gjm.GetProperty(formData, path); err != nil {
        return fmt.Errorf("required field missing: %s", path)
    }
}
```

## Important Notes

### Property Name Limitations

Due to the regex-based parsing, property names must follow these rules:
- Allowed characters: letters, numbers, underscores, and hyphens
- Pattern: `(\w+[\_]?[\-]?)+`
- ✅ Valid: `user_name`, `item-1`, `score_2`
- ❌ Invalid: `user.name` (dots are separators), `item@price`, `name with spaces`

### Error Handling

Always check for errors, especially when:
- Property doesn't exist
- Type mismatches occur (e.g., treating a string as an array)
- Array index is out of bounds

```go
import gjm "github.com/firewut/go-json-map"

// Safe property access
if value, err := gjm.GetProperty(data, "might.not.exist"); err != nil {
    // Handle missing property
    value = "default"
}
```

### Type Assertions

Retrieved values are `interface{}` - use type assertions as needed:

```go
import (
    "fmt"
    gjm "github.com/firewut/go-json-map"
)

value, err := gjm.GetProperty(data, "user.age")
if err == nil {
    if age, ok := value.(int); ok {
        fmt.Printf("User is %d years old\n", age)
    }
}
```

## API Reference

### CRUD Operations

The library follows CRUD naming conventions:

- **Create**: `CreateProperty()` - Creates a new property (fails if exists)
- **Read**: `GetProperty()` - Retrieves a property value
- **Update**: `UpdateProperty()` - Creates or updates a property
- **Delete**: `DeleteProperty()` - Removes a property

### Deprecated Functions

- `AddProperty()` - Deprecated alias for `CreateProperty()`. Use `CreateProperty()` for new code.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.