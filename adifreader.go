package adifparser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type AdifReader interface {
	ReadRecord() (AdifRecord, error)
}

type baseAdifReader struct {
	// Underlying io Reader
	rdr io.Reader
	// Whether or not the header has been read
	headerRead bool
	// Version of the adif file
	version float64
	// Excess read data
	excess []byte
}

func (ardr *baseAdifReader) ReadRecord() (*AdifRecord, error) {
	if !ardr.headerRead {
		ardr.readHeader()
	}
	return nil, nil
}

func NewAdifReader(r io.Reader) *baseAdifReader {
	reader := new(baseAdifReader)
	reader.rdr = bufio.NewReader(r)
	reader.headerRead = false
	// Assumption
	reader.version = 2
	return reader
}

func (ardr *baseAdifReader) readHeader() {
	ardr.headerRead = true
	eoh := []byte("<eoh>")
	adif_version := []byte("<adif_ver:")
	chunk, err := ardr.readChunk()
	if err != nil {
		// TODO: Log the error somewhere
		return
	}
	if bytes.HasPrefix(chunk, []byte("<")) {
		if bytes.HasPrefix(chunk, adif_version) {
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
			excess = excess[bytes.Index(excess, eoh)+len(eoh):]
			ardr.excess = excess[bytes.Index(excess, []byte("<")):]
		} else {
			ardr.excess = chunk
		}
		return
	}
	for !bytes.Contains(chunk, eoh) {
		newchunk, _ := ardr.readChunk()
		chunk = append(chunk, newchunk...)
	}
	offset := bytes.Index(chunk, eoh) + len(eoh)
	chunk = chunk[offset:]
	ardr.excess = chunk[bytes.Index(chunk, []byte("<")):]
}

func (ardr *baseAdifReader) readChunk() ([]byte, error) {
	chunk := make([]byte, 1024)
	_, err := ardr.rdr.Read(chunk)
	if err != nil {
		// TODO: Log the error somewhere
		return nil, err
	}
	return chunk, nil
}

func (ardr *baseAdifReader) readRecord() ([]byte, error) {
	return nil, nil
}
