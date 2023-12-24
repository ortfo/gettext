- *Go语言QQ群: 102319854, 1055927514*
- *凹语言(凹读音“Wa”)(The Wa Programming Language): https://github.com/wa-lang/wa*

----

# gettext-go: GNU gettext for Go ([Imported By Kubernetes](https://pkg.go.dev/github.com/ortfo/gettext@v0.1.0/gettext?tab=importedby))

- PkgDoc: [http://godoc.org/github.com/ortfo/gettext](http://godoc.org/github.com/ortfo/gettext)
- PkgDoc: [http://pkg.go.dev/github.com/ortfo/gettext](http://pkg.go.dev/github.com/ortfo/gettext)

## Install

1. `go get github.com/ortfo/gettext`
2. `go run hello.go`

The godoc.org or go.dev has more information.

## Examples

```Go
package main

import (
	"fmt"

	"github.com/ortfo/gettext"
)

func main() {
	gettext := gettext.New("hello", "./examples/locale").SetLanguage("zh_CN")
	fmt.Println(gettext.Gettext("Hello, world!"))

	// Output: 你好, 世界!
}
```

```Go
package main

import (
	"fmt"

	"github.com/ortfo/gettext"
)

func main() {
	gettext.SetLanguage("zh_CN")
	gettext.BindLocale(gettext.New("hello", "locale"))

	// gettext.BindLocale("hello", "locale")              // from locale dir
	// gettext.BindLocale("hello", "locale.zip")          // from locale zip file
	// gettext.BindLocale("hello", "locale.zip", zipData) // from embedded zip data

	// translate source text
	fmt.Println(gettext.Gettext("Hello, world!"))
	// Output: 你好, 世界!

	// if no msgctxt in PO file (only msgid and msgstr),
	// specify context as "" by
	fmt.Println(gettext.PGettext("", "Hello, world!"))
	// Output: 你好, 世界!

	// translate resource
	fmt.Println(string(gettext.Getdata("poems.txt"))))
	// Output: ...
}
```

Go file: [hello.go](https://github.com/ortfo/gettext/blob/master/examples/hello.go); PO file: [hello.po](https://github.com/ortfo/gettext/blob/master/examples/locale/default/LC_MESSAGES/hello.po);

----

## API Changes (v0.1.0 vs v1.0.0)

### Renamed package path

| v0.1.0 (old)                                    | v1.0.0 (new)                            |
| ----------------------------------------------- | --------------------------------------- |
| `github.com/ortfo/gettext/gettext`        | `github.com/ortfo/gettext`        |
| `github.com/ortfo/gettext/gettext/po`     | `github.com/ortfo/gettext/po`     |
| `github.com/ortfo/gettext/gettext/mo`     | `github.com/ortfo/gettext/mo`     |
| `github.com/ortfo/gettext/gettext/plural` | `github.com/ortfo/gettext/plural` |

### Renamed functions

| v0.1.0 (old)                       | v1.0.0 (new)                |
| ---------------------------------- | --------------------------- |
| `gettext-go/gettext.*`             | `gettext-go.*`              |
| `gettext-go/gettext.DefaultLocal`  | `gettext-go.DefaultLanguage`|
| `gettext-go/gettext.BindTextdomain`| `gettext-go.BindLocale`     |
| `gettext-go/gettext.Textdomain`    | `gettext-go.SetDomain`      |
| `gettext-go/gettext.SetLocale`     | `gettext-go.SetLanguage`    |
| `gettext-go/gettext/po.Load`       | `gettext-go/po.LoadFile`    |
| `gettext-go/gettext/po.LoadData`   | `gettext-go/po.Load`        |
| `gettext-go/gettext/mo.Load`       | `gettext-go/mo.LoadFile`    |
| `gettext-go/gettext/mo.LoadData`   | `gettext-go/mo.Load`        |

### Use empty string as the default context for `gettext.Gettext`

```go
package main

// v0.1.0
// if the **context** missing, use `callerName(2)` as the context:

// v1.0.0
// if the **context** missing, use empty string as the context:

func main() {
	gettext.Gettext("hello")          
	// v0.1.0 => gettext.PGettext("main.main", "hello")
	// v1.0.0 => gettext.PGettext("", "hello")

	gettext.DGettext("domain", "hello")
	// v0.1.0 => gettext.DPGettext("domain", "main.main", "hello")
	// v1.0.0 => gettext.DPGettext("domain", "", "hello")

	gettext.NGettext("domain", "hello", "hello2", n)
	// v0.1.0 => gettext.PNGettext("domain", "main.main", "hello", "hello2", n)
	// v1.0.0 => gettext.PNGettext("domain", "", "hello", "hello2", n)

	gettext.DNGettext("domain", "hello", "hello2", n)
	// v0.1.0 => gettext.DPNGettext("domain", "main.main", "hello", "hello2", n)
	// v1.0.0 => gettext.DPNGettext("domain", "", "hello", "hello2", n)
}
```

### `BindLocale` support `FileSystem` interface

```go
// Use FileSystem:
//	BindLocale(New("poedit", "name", OS("path/to/dir"))) // bind "poedit" domain
//	BindLocale(New("poedit", "name", OS("path/to.zip"))) // bind "poedit" domain
```

## New API in v1.0.0

`Gettexter` interface:

```go
type Gettexter interface {
	FileSystem() FileSystem

	GetDomain() string
	SetDomain(domain string) Gettexter

	GetLanguage() string
	SetLanguage(lang string) Gettexter

	Gettext(msgid string) string
	PGettext(msgctxt, msgid string) string

	NGettext(msgid, msgidPlural string, n int) string
	PNGettext(msgctxt, msgid, msgidPlural string, n int) string

	DGettext(domain, msgid string) string
	DPGettext(domain, msgctxt, msgid string) string
	DNGettext(domain, msgid, msgidPlural string, n int) string
	DPNGettext(domain, msgctxt, msgid, msgidPlural string, n int) string

	Getdata(name string) []byte
	DGetdata(domain, name string) []byte
}

func New(domain, path string, data ...interface{}) Gettexter
```

`FileSystem` interface:

```go
type FileSystem interface {
	LocaleList() []string
	LoadMessagesFile(domain, lang, ext string) ([]byte, error)
	LoadResourceFile(domain, lang, name string) ([]byte, error)
	String() string
}

func NewFS(name string, x interface{}) FileSystem
func OS(root string) FileSystem
func ZipFS(r *zip.Reader, name string) FileSystem
func NilFS(name string) FileSystem
```

----

## BUGS

Please report bugs to <chaishushan@gmail.com>.

Thanks!
