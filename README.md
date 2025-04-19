# 📚 ebook

A Go module for parsing EPUB and Mobipocket (MOBI/AZW) ebook formats.  
It extracts metadata and cover images from ebook files with ease.

## ✨ Features

- 📖 Supports EPUB and MOBI formats
- 🏷️ Reads metadata (title, author, language, etc.)
- 🖼️ Extracts embedded cover images
- ⚡ Clean and simple API

## 🔧 Usage

```go
import "github.com/ignisVeneficus/ebook"

meta, cover, err := ebook.Parse("path/to/book.epub")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Title:", meta.Title)
fmt.Println("Author:", meta.Author)

// cover is a []byte (e.g. JPEG or PNG)
```

## 📦 Installation
```bash
go get github.com/ignisVeneficus/ebook
```
## 📁 Supported Formats
- .epub (EPUB 2 / EPUB 3)
- .mobi / .azw (Mobipocket)

## ✅ Planned Features
- ✅ EPUB metadata & cover support
- ✅ MOBI/AZW basic metadata support

## 📃 License
MIT License
Feel free to use, modify, and share.

Made with ☕ and Go 🐹