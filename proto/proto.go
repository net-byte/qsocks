package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func Encode(data []byte) ([]byte, error) {
	length := int32(len(data))
	pkg := new(bytes.Buffer)
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func Decode(reader *bufio.Reader) ([]byte, int32, error) {
	len, _ := reader.Peek(4)
	blen := bytes.NewBuffer(len)
	var dlen int32
	err := binary.Read(blen, binary.LittleEndian, &dlen)
	if err != nil {
		return nil, 0, err
	}
	if int32(reader.Buffered()) < dlen+4 {
		return nil, 0, err
	}
	pack := make([]byte, 4+dlen)
	_, err = reader.Read(pack)
	if err != nil {
		return nil, 0, err
	}
	return pack[4:], dlen, nil
}
