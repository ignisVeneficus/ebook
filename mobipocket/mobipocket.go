package mobipocket

import (
	"encoding/binary"
	"io"

	"github.com/ignisVeneficus/ebook/eBookData"
	"github.com/ignisVeneficus/ebook/mobipocket/palmdb"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

const ENCODING_CP1252 = "CP1252"
const ENCODING_UTF8 = "utf8"

func readLongInteger(data []byte, start int) (int, int) {
	return int(binary.BigEndian.Uint32(data[start : start+4])), start + 4
}

func readShortInteger(data []byte, start int) (int, int) {
	return int(binary.BigEndian.Uint16(data[start : start+2])), start + 2
}
func readString(data []byte, start int, length int, encoding string) (string, int, error) {
	switch encoding {
	case ENCODING_CP1252:
		dec := charmap.Windows1251.NewDecoder()
		ret, err := dec.Bytes(data[start : start+length])
		return string(ret), start + length, err
	default:
		return string(data[start : start+length]), start + length, nil
	}
}
func readStringFull(data []byte, encoding string) (string, error) {
	switch encoding {
	case ENCODING_CP1252:
		dec := charmap.Windows1251.NewDecoder()
		ret, err := dec.Bytes(data)
		return string(ret), err
	default:
		return string(data), nil
	}
}

type Mobipocket struct {
	db           palmdb.Db
	cover        []byte
	compression  int
	bookType     int
	textEncoding string
	exthRecords  []exthRecord
	metadata     mobiMetadata
}

/*
	type huffmanData struct {
		firstRecord int
		qtyRecord   int
		tableOffset int
		tableLength int
	}
*/
type exthRecord struct {
	recordType int
	length     int
	content    []byte
}

type mobiMetadata struct {
	title          string
	author         []string
	contributor    []string
	isbn           string
	publisher      string
	publishingDate string
}

func (m mobiMetadata) Author() []string {
	return m.author
}
func (m mobiMetadata) Title() string {
	return m.title
}
func (m mobiMetadata) Publisher() string {
	return m.publisher
}
func (m mobiMetadata) PubDate() string {
	return m.publishingDate
}
func (m mobiMetadata) ISBN() string {
	return m.isbn
}
func (m mobiMetadata) Contributor() []string {
	return m.contributor
}

func ReadMobi(f io.Reader) (*Mobipocket, error) {
	log.Logger.Debug().Msg("Start read mobipocket file")
	defer log.Logger.Debug().Msg("End read mobipocket file")
	mobi := Mobipocket{}
	var (
		err              error
		firstImageRecord int
	)
	if mobi.db, err = palmdb.ReadDb(f); err != nil {
		return nil, err
	}

	if len(mobi.db.Records) == 0 {
		return &mobi, nil
	}

	header := (mobi.db.Records[0].Data())
	headerRecordSize := len(header)

	mobiHeaderLength, _ := readLongInteger(header, 20)
	exhtStart := mobiHeaderLength + 16
	compression, _ := readShortInteger(header, 0)
	mobi.compression = compression
	bookType, _ := readLongInteger(header, 24)
	mobi.bookType = bookType
	encoding, _ := readLongInteger(header, 28)
	switch encoding {
	case 1252:
		mobi.textEncoding = ENCODING_CP1252
	case 65001:
		mobi.textEncoding = ENCODING_UTF8
	}
	titlePos, _ := readLongInteger(header, 84)
	titleLength, _ := readLongInteger(header, 88)

	firstImageRecord, _ = readLongInteger(header, 108)

	exth, _ := readLongInteger(header, 128)

	if headerRecordSize < exhtStart || (exth&0x40) == 0 {
		log.Logger.Debug().Msg("No metadata in the file")
		return &mobi, nil
	}
	//EXTH headers
	mobi.exthRecords = make([]exthRecord, 0)
	pos := exhtStart

	// "EXTH" text + 4byte of EXTH length
	pos += 8
	exthCount, pos := readLongInteger(header, pos)
	for i := 0; i < exthCount; i++ {
		var exthRecord exthRecord
		exthRecord, pos, _ = readExthRecord(header, pos)
		mobi.exthRecords = append(mobi.exthRecords, exthRecord)
	}
	metadata := mobiMetadata{
		author:      make([]string, 0),
		contributor: make([]string, 0),
	}

	//parse to metadata
	title, _, _ := readString(header, titlePos, titleLength, mobi.textEncoding)
	metadata.title = title

	var coverImage int
	for _, exthRecord := range mobi.exthRecords {
		// https://github.com/kevinhendricks/KindleUnpack/blob/master/lib/mobi_header.py
		// https://wiki.mobileread.com/wiki/MOBI
		switch exthRecord.recordType {
		case 100:
			if author, err := readStringFull(exthRecord.content, mobi.textEncoding); err == nil {
				metadata.author = append(metadata.author, author)
			}
		case 101:
			if publisher, err := readStringFull(exthRecord.content, mobi.textEncoding); err == nil {
				metadata.publisher = publisher
			}
		case 104:
			if isbn, err := readStringFull(exthRecord.content, mobi.textEncoding); err == nil {
				metadata.isbn = isbn
			}
		case 106:
			if publishingDate, err := readStringFull(exthRecord.content, mobi.textEncoding); err == nil {
				metadata.publishingDate = publishingDate
			}
		case 108:
			if contributor, err := readStringFull(exthRecord.content, mobi.textEncoding); err == nil {
				metadata.contributor = append(metadata.contributor, contributor)
			}
		case 121:
			// kf8start = int(binary.BigEndian.Uint32(exthRecord.content))
		case 201:
			coverImage = int(binary.BigEndian.Uint32(exthRecord.content))
		}
	}
	mobi.metadata = metadata

	// get cover
	// cover image is the nr.th image from the firstImageRecord
	coverImage = coverImage + firstImageRecord
	if len(mobi.db.Records) > coverImage {
		mobi.cover = mobi.db.Records[coverImage].Data()
	}
	return &mobi, nil
}
func readExthRecord(data []byte, pos int) (exthRecord, int, error) {
	rType, pos := readLongInteger(data, pos)
	rLength, pos := readLongInteger(data, pos)
	rData := data[pos : rLength+pos-8]
	ret := exthRecord{recordType: rType, length: rLength, content: rData}

	return ret, rLength + pos - 8, nil
}
func (mobi Mobipocket) Metadata() eBookData.Metadata {
	return mobi.metadata
}
func (mobi Mobipocket) Cover() []byte {
	return mobi.cover
}
