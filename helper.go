package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MAGIC_FOUR_BYTES_V1 = uint32(0xFFFFFF01)

func writeHeader(w io.Writer, seqId int64) error {
	if w == nil {
		return errors.New("write to nil writer")
	}

	err := binary.Write(w, binary.BigEndian, uint32(MAGIC_FOUR_BYTES_V1))
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, uint32(seqId))
}

func WriteMessage(w io.Writer, input []byte, seqId int64) error {
	err := writeHeader(w, seqId)
	if err != nil {
		return err
	}

	contentLen := uint32(len(input))
	err = binary.Write(w, binary.BigEndian, contentLen)
	if err != nil {
		return err
	}

	_, err = w.Write(input)
	return err
}

/*
seqId < 0 means don't check it
*/
func readHeader(r io.Reader, seqId int64) (uint32, error) {
	buf := [4]byte{}
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	if binary.BigEndian.Uint32(buf[:]) != MAGIC_FOUR_BYTES_V1 {
		return 0, fmt.Errorf("got %v expected header v1", buf)
	}

	_, err = io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	realSeqId := binary.BigEndian.Uint32(buf[:])
	if seqId >= 0 && uint32(seqId) != realSeqId {
		return 0, fmt.Errorf("seqId invalid %v, current %v", buf, seqId)
	}
	return realSeqId, nil
}

func ReadMessage(r io.Reader, seqId int64) ([]byte, uint32, error) {
	realSeqId, err := readHeader(r, seqId)
	if err != nil {
		return nil, realSeqId, err
	}

	lengthBuf := [4]byte{}
	_, err = io.ReadFull(r, lengthBuf[:])
	if err != nil {
		return nil, realSeqId, err
	}
	length := int(binary.BigEndian.Uint32(lengthBuf[:]))

	if length > 10*1024*1024 {
		return nil, realSeqId, fmt.Errorf("message too long %v", length)
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, realSeqId, err
	}
	return buf, realSeqId, nil
}
