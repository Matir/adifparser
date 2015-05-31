package adifparser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type ADIFReader interface {
	ReadRecord() (ADIFRecord, error)
}

type baseADIFReader struct {
	// Underlying io Reader
	rdr io.Reader
	// Whether or not the header has been read
	headerRead bool
	// Version of the adif file
	version float64
	// Excess read data
	excess []byte
}

func (ardr *baseADIFReader) ReadRecord() (*ADIFRecord, error) {
	if !ardr.headerRead {
		ardr.readHeader()
	}
	return nil, nil
}

func NewADIFReader(r io.Reader) *baseADIFReader {
	reader := new(baseADIFReader)
	reader.rdr = bufio.NewReader(r)
	reader.headerRead = false
	// Assumption
	reader.version = 2
	return reader
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

func (ardr *baseADIFReader) readChunk() ([]byte, error) {
	chunk := make([]byte, 1024)
	n, err := ardr.rdr.Read(chunk)
	if err != nil {
		// TODO: Log the error somewhere
		return nil, err
	}
	return chunk[:n], nil
}

func (ardr *baseADIFReader) readRecord() ([]byte, error) {
	eor := []byte("<eor>")
	buf := ardr.excess
	ardr.excess = nil
	for !bytes.Contains(buf, eor) {
		newchunk, err := ardr.readChunk()
		buf = bytes.TrimSpace(buf)
		if err != nil {
			if err == io.EOF {
				return buf, nil
			}
			//TODO: Log the error somewhere
			return buf, err
		}
		buf = append(buf, newchunk...)
	}
	record_end := bytes.Index(buf, eor)
	ardr.excess = buf[record_end+len(eor):]
	return bytes.TrimSpace(buf[:record_end]), nil
}
