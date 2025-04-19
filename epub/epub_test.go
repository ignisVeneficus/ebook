package epub

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

// Helper: create an XML document from string
func parseXML(t *testing.T, xml string) *xmlquery.Node {
	doc, err := xmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}
	return doc
}

func TestParseTitle(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:title>Test Book Title</dc:title>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	title, err := parseTitle(doc, nsMap)
	if err != nil || title != "Test Book Title" {
		t.Errorf("Expected title 'Test Book Title', got '%s' (err: %v)", title, err)
	}
}

func TestParseAuthor(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
			<dc:creator opf:role="aut">Jane Smith</dc:creator>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	authors, err := parseAuthor(doc, nsMap, "normal")
	if err != nil || len(authors) != 1 || authors[0] != "Jane Smith" {
		t.Errorf("Expected author 'Jane Smith', got %v (err: %v)", authors, err)
	}
}

func TestParseContributor(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:contributor>Editor One</dc:contributor>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	contributors, err := parseContributor(doc, nsMap)
	if err != nil || len(contributors) != 1 || contributors[0] != "Editor One" {
		t.Errorf("Expected contributor 'Editor One', got %v (err: %v)", contributors, err)
	}
}

func TestParseIsbn(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:identifier id="ISBN">1234567890</dc:identifier>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	isbn, err := parseIsbn(doc, nsMap)
	if err != nil || isbn != "1234567890" {
		t.Errorf("Expected ISBN '1234567890', got '%s' (err: %v)", isbn, err)
	}
}

func TestParsePublisher(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:publisher>GoLang Books</dc:publisher>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	publisher, err := parsePublisher(doc, nsMap)
	if err != nil || publisher != "GoLang Books" {
		t.Errorf("Expected publisher 'GoLang Books', got '%s' (err: %v)", publisher, err)
	}
}

func TestParsePubDate(t *testing.T) {
	xml := `<package xmlns="http://www.idpf.org/2007/opf">
		<metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:date>2022-07-15</dc:date>
		</metadata>
	</package>`
	doc := parseXML(t, xml)
	nsMap := map[string]string{"dc": "http://purl.org/dc/elements/1.1/", "opf": "http://www.idpf.org/2007/opf"}

	date, err := parsePubDate(doc, nsMap)
	if err != nil || date != "2022" {
		t.Errorf("Expected pub date '2022', got '%s' (err: %v)", date, err)
	}
}

func TestReadEpubMinimal(t *testing.T) {
	// Create in-memory ZIP EPUB
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// META-INF/container.xml
	w, _ := zw.Create("META-INF/container.xml")
	io.WriteString(w, `<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
	<rootfiles>
		<rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
	</rootfiles>
</container>`)

	// content.opf
	w, _ = zw.Create("OEBPS/content.opf")
	io.WriteString(w, `<?xml version="1.0"?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="BookId">
	<metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
		<dc:title>Test EPUB</dc:title>
		<dc:creator opf:role="aut">Test Author</dc:creator>
	</metadata>
</package>`)

	zw.Close()

	epub, err := ReadEpub(bytes.NewReader(buf.Bytes()), "normal")
	if err != nil {
		t.Fatalf("ReadEpub failed: %v", err)
	}

	meta := epub.Metadata()
	if meta.Title() != "Test EPUB" || len(meta.Author()) == 0 || meta.Author()[0] != "Test Author" {
		t.Errorf("Unexpected metadata: %+v", meta)
	}
}

func TestLoadCover(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// META-INF/container.xml
	w, _ := zw.Create("META-INF/container.xml")
	io.WriteString(w, `<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
	<rootfiles>
		<rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
	</rootfiles>
</container>`)

	// content.opf
	w, _ = zw.Create("OEBPS/content.opf")
	io.WriteString(w, `<?xml version="1.0"?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="BookId">
	<metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
		<meta name="cover" content="cover-image"/>
	</metadata>
	<manifest>
		<item id="cover-image" href="images/cover.jpg" media-type="image/jpeg"/>
	</manifest>
</package>`)

	// images/cover.jpg (dummy JPEG data)
	w, _ = zw.Create("OEBPS/images/cover.jpg")
	jpgData := []byte{0xFF, 0xD8, 0xFF, 0xDB, 0x00, 0x43, 0x00, 0xFF, 0xD9} // JPEG start + dummy + end
	w.Write(jpgData)

	zw.Close()

	epub, err := ReadEpub(bytes.NewReader(buf.Bytes()), "normal")
	if err != nil {
		t.Fatalf("ReadEpub failed: %v", err)
	}
	if !bytes.HasPrefix(epub.cover, []byte{0xFF, 0xD8}) {
		t.Errorf("Cover does not start with JPEG header: %x", epub.cover[:4])
	}
}
