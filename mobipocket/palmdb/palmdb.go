package palmdb

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

/** difference in seconds from Jan 01, 1904 to Jan 01, 1970 */
const TIME_DIFF = 2082844800
const HEADER_SIZE = 78

type palmString struct {
	length int
	value  string
}

func (st *palmString) String() string {
	return st.value
}

func readLongInteger(f io.Reader) (int, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(f, buf); err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(buf)), nil
}

func readShortInteger(f io.Reader) (int, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(f, buf); err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(buf)), nil
}
func readDate(f io.Reader) (time.Time, error) {
	// If the time has the top bit set, it's an unsigned 32-bit number counting from 1st Jan 1904
	// If the time has the top bit clear, it's a signed 32-bit number counting from 1st Jan 1970.
	buf := make([]byte, 4)
	if _, err := io.ReadFull(f, buf); err != nil {
		return time.Time{}, err
	}
	value := 0
	if buf[0] > 0x80 {
		value = int(binary.BigEndian.Uint32(buf)) - TIME_DIFF
	} else {
		value = int(int32(binary.BigEndian.Uint32(buf)))
	}
	return time.Unix(int64(value), 0), nil
}
func ReadByte(f io.Reader) (byte, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(f, buf); err != nil {
		return 0, err
	}
	return buf[0], nil

}
func readPalmString(f io.Reader, length int) (palmString, error) {
	buf := make([]byte, length)
	if _, err := io.ReadFull(f, buf); err != nil {
		return palmString{}, err
	}
	buf = bytes.Trim(buf, "\x00")
	value := string(buf)
	return palmString{length: len(buf), value: value}, nil
}

type Db struct {
	Records          []Record
	Name             palmString
	Attributes       int
	Version          int
	CreateDate       time.Time
	ModDate          time.Time
	BackupDate       time.Time
	ModNumber        int
	OffAppInfo       int
	OffSortInfo      int
	DbType           palmString
	Creator          palmString
	UniqueIdSeed     int
	NextRecordListId int
	NrOfRecords      int
}

func ReadDb(f io.Reader) (Db, error) {
	db := Db{}
	var err error
	if db.Name, err = readPalmString(f, 32); err != nil {
		return db, err
	}
	if db.Attributes, err = readShortInteger(f); err != nil {
		return db, err
	}
	if db.Version, err = readShortInteger(f); err != nil {
		return db, err
	}
	if db.CreateDate, err = readDate(f); err != nil {
		return db, err
	}
	if db.ModDate, err = readDate(f); err != nil {
		return db, err
	}
	if db.BackupDate, err = readDate(f); err != nil {
		return db, err
	}
	if db.ModNumber, err = readLongInteger(f); err != nil {
		return db, err
	}
	if db.OffAppInfo, err = readLongInteger(f); err != nil {
		return db, err
	}
	if db.OffSortInfo, err = readLongInteger(f); err != nil {
		return db, err
	}
	if db.DbType, err = readPalmString(f, 4); err != nil {
		return db, err
	}
	if db.Creator, err = readPalmString(f, 4); err != nil {
		return db, err
	}
	if db.UniqueIdSeed, err = readLongInteger(f); err != nil {
		return db, err
	}
	if db.NextRecordListId, err = readLongInteger(f); err != nil {
		return db, err
	}
	if db.NrOfRecords, err = readShortInteger(f); err != nil {
		return db, err
	}
	db.Records = make([]Record, db.NrOfRecords)
	for i := 0; i < db.NrOfRecords; i++ {
		rec := Record{length: -1}
		if err = rec.readRecordHeader(f); err != nil {
			return db, err
		}
		if i > 0 {
			from := db.Records[i-1].offset
			db.Records[i-1].length = rec.offset - from
		}
		db.Records[i] = rec
	}
	pos := HEADER_SIZE + 8*db.NrOfRecords
	skipOfBytes := 0
	if db.NrOfRecords > 0 {
		skipOfBytes = db.Records[0].offset - pos
	}
	if skipOfBytes > 0 {
		skipBytes := make([]byte, skipOfBytes)
		if _, err = io.ReadFull(f, skipBytes); err != nil {
			return db, err
		}
	}
	for i := 0; i < db.NrOfRecords; i++ {
		record := &db.Records[i]
		if err = record.readRecordValue(f); err != nil {
			return db, err
		}
		//fmt.Printf("Record read after write: len: %d, cap %d\n", len(record.data), cap(record.data))
		//fmt.Printf("%p, %d\n", &record, len(record.data))
	}
	return db, nil
}

type Record struct {
	offset     int
	attributes byte
	recordId   int
	length     int
	data       []byte
}

func (r *Record) readRecordHeader(f io.Reader) error {
	var (
		err error
		b   byte
		i   int
	)
	if r.offset, err = readLongInteger(f); err != nil {
		return err
	}
	if r.attributes, err = ReadByte(f); err != nil {
		return err
	}
	if b, err = ReadByte(f); err != nil {
		return err
	}
	if i, err = readShortInteger(f); err != nil {
		return err
	}
	buf := make([]byte, 4)
	buf[1] = b
	r.recordId = i + int(binary.BigEndian.Uint32(buf))

	return nil
}

func (r *Record) readRecordValue(f io.Reader) error {
	if r.length != -1 {
		buf := make([]byte, r.length)
		if _, err := io.ReadFull(f, buf); err != nil {
			return err
		}
		r.data = append(r.data, buf...)
		return nil
	}
	buf := make([]byte, 1024)
	var result []byte
	for {
		n, err := io.ReadFull(f, buf)
		result = append(result, buf[:n]...)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
	}
	r.data = append(r.data, result...)
	return nil
}

func (r *Record) Data() []byte {
	return r.data
}
