// Copyright 2020 Martin Hebnes Pedersen (LA5NTA). All rights reserved.
// Use of this source code is governed by the MIT-license that can be
// found in the LICENSE file.

package agwpe

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type AGWHeader struct {
	port uint8 // AGWPE port
	// reserved byte x 3
	fType datakind // AGWPE DataKind code
	// reserved bytes x 1
	pid uint8 // PID usually 0x00 or 0xF0 for standard AX.25 frames
	// reserved bytes x 1
	callFrom   string // Sender's callsign-SSID null terminated (0x00)
	callTo     string // Receiver's callsign-SSID null terminated (0x00)
	dataLength uint32 // Length of data in frame
	// user (reserved) byte x 4
}

type AGWFrame struct {
	header AGWHeader
	data   []byte
}

func makeFrame(port uint8, fType datakind, callFrom string, callTo string, data []byte) AGWFrame {
	// Set PID
	var pid uint8
	switch fType {
	case AGWData:
		pid = 0xF0
	default:
		pid = 0x00
	}

	// Put info into frame struct
	return AGWFrame{
		header: AGWHeader{
			port:       port,
			fType:      fType,
			pid:        pid,
			callFrom:   callFrom,
			callTo:     callTo,
			dataLength: uint32(len(data)),
		},
		data: data,
	}
}

func encodeFrame(f AGWFrame) ([]byte, error) {
	// Check data length (do we need/should we have this check? Does AGWPE handle longer packets?)
	if f.header.dataLength > AX25MaxLen && f.header.fType != AGWLogin {
		return nil, errors.New("Data length exceeds AX.25 frame cap")
	}

	// Make byte slice for frame
	fLen := AGWHeaderLen + f.header.dataLength
	frame := make([]byte, fLen)

	// Put port, DataKind (frame type), and PID into frame
	frame[0] = f.header.port
	frame[4] = byte(f.header.fType)
	frame[6] = f.header.pid

	// Copy callsigns into frame
	if n := copy(frame[8:], f.header.callFrom); n > AGWMaxCallLen {
		return nil, fmt.Errorf("Callsign %s exceeds allowed length", f.header.callFrom)
	}

	if n := copy(frame[18:], f.header.callTo); n > AGWMaxCallLen {
		return nil, fmt.Errorf("Callsign %s exceeds allowed length", f.header.callTo)
	}

	// Put data length into frame
	binary.LittleEndian.PutUint32(frame[28:32], f.header.dataLength)

	// Copy data into frame
	if f.header.dataLength > 0 {
		if n := copy(frame[AGWHeaderLen:], f.data); n != int(f.header.dataLength) {
			return nil, errors.New("Data length does not match length declared in header")
		}
	}

	return frame, nil
}

func decodeFrame(b []byte) (AGWFrame, error) {
	// Check to make sure that header is complete
	if len(b) < AGWHeaderLen {
		return AGWFrame{}, errors.New("Incomplete header")
	}

	// Fill in header
	header := AGWHeader{
		port:       b[0],
		fType:      datakind(b[4]),
		pid:        b[6],
		callFrom:   strings.TrimRight(string(b[8:17]), "\x00"),  // Callsigns are null-terminated and data
		callTo:     strings.TrimRight(string(b[18:27]), "\x00"), // after 0x00 should be considered garbage
		dataLength: binary.LittleEndian.Uint32(b[28:32]),
	}

	// Check full frame length
	fLength := AGWHeaderLen + int(header.dataLength)
	if len(b) < fLength {
		return AGWFrame{}, errors.New("Incomplete data")
	}

	// Assemble full frame
	return AGWFrame{
		header: header,
		data:   b[AGWHeaderLen:fLength],
	}, nil
}

func readFrame(reader *bufio.Reader) (AGWFrame, error) {
	// Read header
	header := make([]byte, AGWHeaderLen)
	var err error
	for i := 0; i < AGWHeaderLen && err == nil; i++ {
		header[i], err = reader.ReadByte()
	}
	if err != nil {
		return AGWFrame{}, fmt.Errorf("Error reading frame: %w", err)
	}
	fmt.Println(header)
	// Get data length from header
	length := int(binary.LittleEndian.Uint32(header[28:32]))

	if length > 0 {
		// Read data
		data := make([]byte, length)
		var n int
		for read := 0; read < int(length) && err == nil; {
			n, err = reader.Read(data[read:])
			read += n
		}
		if err != nil {
			return AGWFrame{}, fmt.Errorf("Error reading frame: %w", err)
		}

		// Assemble complete frame
		frame := make([]byte, AGWHeaderLen+length)
		copy(frame, header)
		copy(frame[AGWHeaderLen:], data)

		fmt.Println(frame)
		return decodeFrame(frame)
	}

	// No data, header-only frame
	return decodeFrame(header)
}
