[![Build Status](https://travis-ci.org/firewut/go-json-map.svg)](https://travis-ci.org/firewut/go-json-map)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/firewut/go-json-map) 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/firewut/go-json-map/master/LICENSE)


# Go json map

Using this package you can [Get/Create/Update/Delete] your map's nested properties as easy as json.

For example:
    
```go
document := map[string]interface{}{
	"one": map[string]interface{}{
		"two": map[string]interface{}{
			"three": []int{
				1, 2, 3,
			},
		},
		"four": map[string]interface{}{
			"five": []int{
				11, 22, 33,
			},
		},
	},
}
```

Get a property

```go
property, err := GetProperty(document, "one.two.three[0]")
fmt.Println(property)
// property => 1

property, err = GetProperty(document, "one.two.three", ".")
fmt.Println(property)
// property => 1, 2, 3
```

Create

```go
err := CreateProperty(document, "one.three", "third value")
```

Update

```go
err := UpdateProperty(document, "one.two.three[0]", "updated value")
err := UpdateProperty(document, "one/two/three[4]", []int{1,2,3,4], "/")
```

Delete

```go
err := DeleteProperty(document, "one.four")
err := DeleteProperty(document, "one.two.three[3]")
```
