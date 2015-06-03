package adifparser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

// Interface for ADIFReader
type ADIFReader interface {
	ReadRecord() (ADIFRecord, error)
	RecordCount() int
}

// Real implementation of ADIFReader
type baseADIFReader struct {
	// Underlying io Reader
	rdr io.Reader
	// Whether or not the header has been read
	headerRead bool
	// Version of the adif file
	version float64
	// Excess read data
	excess []byte
	// Record count
	records int
}

type dedupeADIFReader struct {
	baseADIFReader
	// Store seen entities
	seen map[string]bool
}

func (ardr *baseADIFReader) ReadRecord() (ADIFRecord, error) {
	if !ardr.headerRead {
		ardr.readHeader()
	}
	buf, err := ardr.readRecord()
	if err != nil {
		if err != io.EOF {
			adiflog.Printf("readRecord: %v", err)
		}
		return nil, err
	}
	record, err := ParseADIFRecord(buf)
	if err == nil {
		ardr.records += 1
		return record, nil
	}
	return record, err
}

func (ardr *dedupeADIFReader) ReadRecord() (ADIFRecord, error) {
	for true {
		record, err := ardr.baseADIFReader.ReadRecord()
		if err != nil {
			return nil, err
		}
		fp := record.Fingerprint()
		if _, ok := ardr.seen[fp]; !ok {
			ardr.seen[fp] = true
			return record, nil
		}
	}
	return nil, nil
}

func NewADIFReader(r io.Reader) *baseADIFReader {
	reader := &baseADIFReader{}
	reader.init(r)
	return reader
}

func NewDedupeADIFReader(r io.Reader) *dedupeADIFReader {
	reader := &dedupeADIFReader{}
	reader.init(r)
	reader.seen = make(map[string]bool)
	return reader
}

func (ardr *baseADIFReader) init(r io.Reader) {
	ardr.rdr = bufio.NewReader(r)
	ardr.headerRead = false
	// Assumption
	ardr.version = 2
	ardr.records = 0
}

func (ardr *baseADIFReader) readHeader() {
	ardr.headerRead = true
	eoh := []byte("<eoh>")
	adif_version := []byte("<adif_ver:")
	chunk, err := ardr.readChunk()
	if err != nil {
		// TODO: Log the error somewhere
		return
	}
	if bytes.HasPrefix(chunk, []byte("<")) {
		if bytes.HasPrefix(bytes.ToLower(chunk), adif_version) {
			ver_len_str_end := bytes.Index(chunk, []byte(">"))
			ver_len_str := string(chunk[len(adif_version):ver_len_str_end])
			ver_len, err := strconv.ParseInt(ver_len_str, 0, 0)
			if err != nil {
				// TODO: Log the error somewhere
			}
			ver_len_end := int64(ver_len_str_end) + 1 + ver_len
			ardr.version, err = strconv.ParseFloat(
				string(chunk[ver_len_str_end+1:ver_len_end]), 0)
			excess := chunk[ver_len_end:]
			excess = excess[bytes.Index(bytes.ToLower(excess), eoh)+len(eoh):]
			ardr.excess = excess[bytes.Index(excess, []byte("<")):]
		} else {
			ardr.excess = chunk
		}
		return
	}
	for !bytes.Contains(bytes.ToLower(chunk), eoh) {
		newchunk, _ := ardr.readChunk()
		chunk = append(chunk, newchunk...)
	}
	offset := bytes.Index(bytes.ToLower(chunk), eoh) + len(eoh)
	chunk = chunk[offset:]
	ardr.excess = chunk[bytes.Index(chunk, []byte("<")):]
}

func (ardr *baseADIFReader) readChunk() ([]byte, error) {
	chunk := make([]byte, 1024)
	n, err := ardr.rdr.Read(chunk)
	if err != nil {
		return nil, err
	}
	return chunk[:n], nil
}

func (ardr *baseADIFReader) readRecord() ([]byte, error) {
	eor := []byte("<eor>")
	buf := ardr.excess
	ardr.excess = nil
	for !bytes.Contains(bytes.ToLower(buf), eor) {
		newchunk, err := ardr.readChunk()
		buf = bytes.TrimSpace(buf)
		if err != nil {
			ardr.excess = nil
			if err == io.EOF {
				// Expected, pass it up the chain
				if len(buf) > 0 {
					return buf, nil
				}
				return nil, err
			}
			adiflog.Println(err)
			return buf, err
		}
		buf = append(buf, newchunk...)
	}
	record_end := bytes.Index(bytes.ToLower(buf), eor)
	ardr.excess = buf[record_end+len(eor):]
	return bytes.TrimSpace(buf[:record_end]), nil
}

func (ardr *baseADIFReader) RecordCount() int {
	return ardr.records
}
