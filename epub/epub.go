package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/ignisVeneficus/ebook/eBookData"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/rs/zerolog/log"
)

type EpubError struct {
	msg  string
	root error
}

func (m *EpubError) Error() string {
	if m.root == nil {
		return m.msg
	}
	return fmt.Sprintf("%s : %s", m.msg, m.root.Error())
}

func createEpubFormatError(root error) *EpubError {
	return &EpubError{msg: "Epub format error", root: root}
}
func createCustomEpubFormatError(msg string) *EpubError {
	lastError := EpubError{msg: msg}
	return createEpubFormatError(&lastError)
}

func (m *EpubError) Unwrap() error {
	// Return the inner error.
	return m.root
}

type epubMetadata struct {
	title          string
	author         []string
	contributor    []string
	isbn           string
	publisher      string
	publishingDate string
}

func (e epubMetadata) Author() []string {
	return e.author
}
func (e epubMetadata) Title() string {
	return e.title
}
func (e epubMetadata) Publisher() string {
	return e.publisher
}
func (e epubMetadata) PubDate() string {
	return e.publishingDate
}
func (e epubMetadata) ISBN() string {
	return e.isbn
}
func (e epubMetadata) Contributor() []string {
	return e.contributor
}

type Epub struct {
	metadata epubMetadata
	cover    []byte
}

func parseTitle(doc *xmlquery.Node, nsMap map[string]string) (string, error) {
	expr, _ := xpath.CompileWithNS("/opf:package/opf:metadata/dc:title/text()", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		title := node.InnerText()
		log.Logger.Trace().Str("Title", title).Msg("Title parsed")
		return title, nil
	}
	log.Logger.Trace().Msg("No title found")
	return "", nil
}
func parseIsbn(doc *xmlquery.Node, nsMap map[string]string) (string, error) {
	expr, _ := xpath.CompileWithNS("/opf:package/opf:metadata/dc:identifier[@id='ISBN']/text()", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		isbn := node.InnerText()
		log.Logger.Trace().Str("ISBN", isbn).Msg("ISBN parsed")
		return isbn, nil
	}
	log.Logger.Trace().Msg("No ISBN found")
	return "", nil
}

func parsePublisher(doc *xmlquery.Node, nsMap map[string]string) (string, error) {
	expr, _ := xpath.CompileWithNS("/opf:package/opf:metadata/dc:publisher/text()", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		publisher := node.InnerText()
		log.Logger.Trace().Str("Publisher", publisher).Msg("Publisher parsed")
		return publisher, nil
	}
	log.Logger.Trace().Msg("No publisher found")
	return "", nil
}
func parsePubDate(doc *xmlquery.Node, nsMap map[string]string) (string, error) {
	expr, _ := xpath.CompileWithNS("/opf:package/opf:metadata/dc:date/text()", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		pubDate := node.InnerText()
		if len(pubDate) > 4 {
			pubDate = pubDate[:4]
		}
		log.Logger.Trace().Str("PubDate", pubDate).Msg("PubDate parsed")
		return pubDate, nil
	}
	log.Logger.Trace().Msg("No puDate found")
	return "", nil
}
func parseCoverFile(doc *xmlquery.Node, nsMap map[string]string) (string, error) {
	expr, _ := xpath.CompileWithNS("/opf:package/opf:metadata/meta[@name='cover']/@content", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	var coverId string
	for _, node := range nodes {
		coverId = node.InnerText()
		log.Logger.Trace().Str("CoverId", coverId).Msg("CoverId parsed")
		break
	}
	if coverId == "" {
		log.Logger.Trace().Msg("No CoverId found")
		return "", nil
	}
	expr, _ = xpath.CompileWithNS("/opf:package/opf:manifest/item[@id='"+coverId+"']/@href", nsMap)
	nodes = xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		cover := node.InnerText()
		log.Logger.Trace().Str("Cover", cover).Msg("Cover parsed")
		return cover, nil
	}
	log.Logger.Trace().Msg("No Cover found")
	return "", nil
}

func parseAuthor(doc *xmlquery.Node, nsMap map[string]string, mode string) ([]string, error) {
	ret := make([]string, 0)
	expr, err := xpath.CompileWithNS("/opf:package/opf:metadata/dc:creator[@opf:role='aut' or not(@opf:role)]", nsMap)
	if err != nil {
		panic(err)
	}
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		author := ""
		if mode == "file-as" {
			author = node.SelectAttr("opf:file-as")
			if author == "" {
				author = node.InnerText()
			}
		} else {
			author = node.InnerText()
		}
		log.Logger.Trace().Str("Author", author).Msg("Author parsed")
		ret = append(ret, author)
	}
	log.Logger.Trace().Int("Author qrt", len(ret)).Msg("All authors parsed")
	return ret, nil
}
func parseContributor(doc *xmlquery.Node, nsMap map[string]string) ([]string, error) {
	ret := make([]string, 0)
	expr, err := xpath.CompileWithNS("/opf:package/opf:metadata/dc:contributor", nsMap)
	if err != nil {
		panic(err)
	}
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		contributor := node.InnerText()
		log.Logger.Trace().Str("Contributor", contributor).Msg("Contributor parsed")
		ret = append(ret, contributor)
	}
	log.Logger.Trace().Int("Contributor qrt", len(ret)).Msg("All contributors parsed")
	return ret, nil
}
func parseMetadata(file *zip.File, mode string) (*epubMetadata, string, error) {
	metadata := &epubMetadata{}
	reader, err := file.Open()
	if err != nil {
		return metadata, "", createCustomEpubFormatError("Root file not readable")
	}
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	nsMap := map[string]string{
		"dc":      "http://purl.org/dc/elements/1.1/",
		"dcterms": "http://purl.org/dc/terms/",
		"opf":     "http://www.idpf.org/2007/opf",
	}
	if metadata.title, err = parseTitle(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	if metadata.author, err = parseAuthor(doc, nsMap, mode); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	if metadata.contributor, err = parseContributor(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	if metadata.isbn, err = parseIsbn(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	if metadata.publisher, err = parsePublisher(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	if metadata.publishingDate, err = parsePubDate(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	var cover string
	if cover, err = parseCoverFile(doc, nsMap); err != nil {
		return metadata, "", createEpubFormatError(err)
	}
	return metadata, cover, nil
}

func getRoot(file *zip.File) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", createCustomEpubFormatError("Root file not readable")
	}
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return "", createEpubFormatError(err)
	}
	nsMap := map[string]string{
		"a": "urn:oasis:names:tc:opendocument:xmlns:container",
	}
	//	expr, _ := xpath.CompileWithNS("/a:container/a:rootfiles/a:rootfile/@full-path", nsMap)
	expr, _ := xpath.CompileWithNS("/container/rootfiles/rootfile/@full-path", nsMap)
	nodes := xmlquery.QuerySelectorAll(doc, expr)
	for _, node := range nodes {
		return node.InnerText(), nil
	}
	return "", nil
}
func loadCover(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return make([]byte, 0), createCustomEpubFormatError("Cover file not readable")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return make([]byte, 0), createEpubFormatError(err)
	}
	return data, nil

}

func ReadEpub(f io.Reader, mode string) (*Epub, error) {
	log.Logger.Debug().Msg("Start reading epub file")
	defer log.Logger.Debug().Msg("End reading epub file")
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, f)
	if err != nil {
		return nil, createEpubFormatError(err)
	}

	reader := bytes.NewReader(buff.Bytes())

	// Open a zip archive for reading.
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, createEpubFormatError(err)
	}
	files := make(map[string]*zip.File)

	for _, file := range zipReader.File {
		files[file.Name] = file
	}
	fileList := slices.Collect(maps.Keys(files))
	log.Logger.Trace().Strs("files", fileList).Msg("Files in the zip")

	metainfFile := files["META-INF/container.xml"]
	if metainfFile == nil {
		return nil, createCustomEpubFormatError("No META-INF/container file")
	}
	rootFile, err := getRoot(metainfFile)
	if err != nil {
		return nil, createEpubFormatError(err)
	}
	if rootFile == "" {
		return nil, createCustomEpubFormatError("Invalid META-INF/container file")
	}
	contentFile := files[rootFile]
	if contentFile == nil {
		return nil, createCustomEpubFormatError("No content.opf file")
	}
	//metadata, coverFile, err := getMetadata(contentFile)
	metadata, cover, err := parseMetadata(contentFile, mode)
	if err != nil {
		return nil, createEpubFormatError(err)
	}
	coverData := make([]byte, 0)
	if cover != "" {
		// first try
		coverFile := files[cover]
		if coverFile == nil {
			for _, f := range fileList {
				if strings.HasSuffix(f, cover) {
					coverFile = files[f]
					break
				}
			}
		}
		if coverFile != nil {
			coverData, err = loadCover(coverFile)
			if err != nil {
				return nil, createEpubFormatError(err)
			}
		} else {
			log.Logger.Warn().Msg("No cover file, but defined")
		}
	}

	ret := Epub{metadata: *metadata, cover: coverData}
	return &ret, nil
}
func (epub Epub) Metadata() eBookData.Metadata {
	return epub.metadata
}
func (epub Epub) Cover() []byte {
	return epub.cover
}
