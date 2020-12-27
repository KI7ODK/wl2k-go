// Copyright 2020 Martin Hebnes Pedersen (LA5NTA). All rights reserved.
// Use of this source code is governed by the MIT-license that can be
// found in the LICENSE file.

// Package agwpe provides means of connecting to a TNC using AGWPE
package agwpe

const (
	DefaultAddr    = "localhost:8000" // The default address AGWPE listens on
	AX25MaxLen     = 255              // Max length of 255 bytes
	AGWHeaderLen   = 36               // Header length in bytes
	AGWMaxCallLen  = 9                // Max callsign-ssid length in bytes
	MinAGWMajorVer = 2000             // Minimum supported major version
	MinAGWMinorVer = 78               // Minimum supported minior version
)

const (
	ErrConnectTimeout = "Connect timeout"
	ErrInvalidAddr    = "Invalid address format"
)

// AGWPE DataKind specifies the frame type being sent/received
type datakind uint8

// AGWPE DataKind codes
const (
	// Setup codes
	AGWLogin     datakind = 0x50 // ASCII "P" - Login to AGWPE
	AGWRegCall   datakind = 0x58 // ASCII "X" - Register callsign/answer success or failure
	AGWUnregCall datakind = 0x78 // ASCII "x" - Unregister callsign

	// Query codes
	AGWVersion    datakind = 0x52 // ASCII "R" - Query AGWPE version/answer AGWPE version
	AGWPortInfo   datakind = 0x47 // ASCII "G" - Query AGWPE ports/answer AGWPE ports
	AGWPortCap    datakind = 0x67 // ASCII "g" - Query AGWPE port capabilities/Answer port capabilities
	AGWFramesPort datakind = 0x79 // ASCII "y" - Query outstanding frames at port/answer outstanding frames
	AGWFramesConn datakind = 0x59 // ASCII "Y" - Query outstanding frames for connection/answer outstanding frames
	AGWHeard      datakind = 0x48 // ASCII "H" - Query heard stations on port/answer heard stations

	// Connection and data codes
	AGWConnect    datakind = 0x43 // ASCII "C" - Start an AX.25 connection/answer success or failure
	AGWConnectVia datakind = 0x76 // ASCII "v" - Start an AX.25 connection using digipeaters/answered with "C"
	AGWDisconnect datakind = 0x64 // ASCII "d" - Disconnect AX.25 connection/answer success or failure
	AGWConnectNS  datakind = 0x63 // ASCII "c" - Start a non-standard AX.25 connection/answered with "C"
	AGWData       datakind = 0x44 // ASCII "D" - Send data over connection/received data over connection
	AGWUnproto    datakind = 0x4D // ASCII "M" - Send unproto info
	AGWUnprotoVia datakind = 0x56 // ASCII "V" - Send unproto info VIA
	AGWRawFrame   datakind = 0x4B // ASCII "K" - Send raw frame/received raw frame

	// Monitor codes
	AGWMonitor    datakind = 0x6D // ASCII "m" - Start or stop monitoring data
	AGWMonitorRaw datakind = 0x6B // ASCII "k" - Start or stop monitoring raw frames
	AGWUnnumFrame datakind = 0x55 // ASCII "U" - Received unnumbered info frame for a registered application
	AGWInfoFrame  datakind = 0x49 // ASCII "I" - Monitored info frame
	AGWSpvsrFrame datakind = 0x53 // ASCII "S" - Monitored supervisory frame
	AGWMntrFrame  datakind = 0x54 // ASCII "T" - Monitored frame sent by this application
)
