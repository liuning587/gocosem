package gocosem

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

const (
	HDLC_FRAME_DIRECTION_CLIENT_INBOUND  = 1
	HDLC_FRAME_DIRECTION_CLIENT_OUTBOUND = 2
	HDLC_FRAME_DIRECTION_SERVER_INBOUND  = 3
	HDLC_FRAME_DIRECTION_SERVER_OUTBOUND = 4
)

const (
	HDLC_ADDRESS_LENGTH_1 = 1
	HDLC_ADDRESS_LENGTH_2 = 2
	HDLC_ADDRESS_LENGTH_4 = 4
)

const (
	HDLC_CONTROL_I    = 1 // I frame
	HDLC_CONTROL_RR   = 2 // response ready
	HDLC_CONTROL_RNR  = 3 // response not ready
	HDLC_CONTROL_SNRM = 4 // set normal response mode
	HDLC_CONTROL_DISC = 5 // disconnect
	HDLC_CONTROL_UA   = 6 // unnumbered acknowledgement
	HDLC_CONTROL_DM   = 7 // disconnected mode
	HDLC_CONTROL_FRMR = 8 // frame reject
	HDLC_CONTROL_UI   = 9 // unnumbered information
)

type HdlcTransport struct {
	rw                         io.ReadWriter
	responseTimeout            int // in milliseconds
	connected                  bool
	windowSize                 uint8
	maxInfoFieldLengthReceive  int
	maxInfoFieldLengthTransmit int
	expectedServerAddrLength   int // HDLC_ADDRESS_BYTE_LENGTH_1, HDLC_ADDRESS_BYTE_LENGTH_2, HDLC_ADDRESS_BYTE_LENGTH_4
}

type HdlcClientConnection struct {
	htran *HdlcTransport
	vs    uint8 // V(S) - send sequence variable
	vr    uint8 // V(R) - receive sequence variable
}

type HdlcServerConnection struct {
	htran *HdlcTransport
	vs    uint8 // V(S) - send sequence variable
	vr    uint8 // V(R) - receive sequence variable
}

type HdlcFrame struct {
	direction             int
	formatType            uint8
	segmentation          bool
	length                int
	logicalDeviceId       uint16
	physicalDeviceId      *uint16 // may not be present
	clientId              uint8
	pf                    bool  // poll/final bit
	nr                    uint8 // N(R) - receive sequence number
	ns                    uint8 // N(S) - send sequence number
	control               int
	fcs16                 uint16 // current fcs16 checksum
	infoField             []byte // information
	callingPhysicalDevice bool
}

var fcstab = []uint16{0x0000, 0x1189, 0x2312, 0x329b, 0x4624, 0x57ad, 0x6536, 0x74bf,
	0x8c48, 0x9dc1, 0xaf5a, 0xbed3, 0xca6c, 0xdbe5, 0xe97e, 0xf8f7,
	0x1081, 0x0108, 0x3393, 0x221a, 0x56a5, 0x472c, 0x75b7, 0x643e,
	0x9cc9, 0x8d40, 0xbfdb, 0xae52, 0xdaed, 0xcb64, 0xf9ff, 0xe876,
	0x2102, 0x308b, 0x0210, 0x1399, 0x6726, 0x76af, 0x4434, 0x55bd,
	0xad4a, 0xbcc3, 0x8e58, 0x9fd1, 0xeb6e, 0xfae7, 0xc87c, 0xd9f5,
	0x3183, 0x200a, 0x1291, 0x0318, 0x77a7, 0x662e, 0x54b5, 0x453c,
	0xbdcb, 0xac42, 0x9ed9, 0x8f50, 0xfbef, 0xea66, 0xd8fd, 0xc974,
	0x4204, 0x538d, 0x6116, 0x709f, 0x0420, 0x15a9, 0x2732, 0x36bb,
	0xce4c, 0xdfc5, 0xed5e, 0xfcd7, 0x8868, 0x99e1, 0xab7a, 0xbaf3,
	0x5285, 0x430c, 0x7197, 0x601e, 0x14a1, 0x0528, 0x37b3, 0x263a,
	0xdecd, 0xcf44, 0xfddf, 0xec56, 0x98e9, 0x8960, 0xbbfb, 0xaa72,
	0x6306, 0x728f, 0x4014, 0x519d, 0x2522, 0x34ab, 0x0630, 0x17b9,
	0xef4e, 0xfec7, 0xcc5c, 0xddd5, 0xa96a, 0xb8e3, 0x8a78, 0x9bf1,
	0x7387, 0x620e, 0x5095, 0x411c, 0x35a3, 0x242a, 0x16b1, 0x0738,
	0xffcf, 0xee46, 0xdcdd, 0xcd54, 0xb9eb, 0xa862, 0x9af9, 0x8b70,
	0x8408, 0x9581, 0xa71a, 0xb693, 0xc22c, 0xd3a5, 0xe13e, 0xf0b7,
	0x0840, 0x19c9, 0x2b52, 0x3adb, 0x4e64, 0x5fed, 0x6d76, 0x7cff,
	0x9489, 0x8500, 0xb79b, 0xa612, 0xd2ad, 0xc324, 0xf1bf, 0xe036,
	0x18c1, 0x0948, 0x3bd3, 0x2a5a, 0x5ee5, 0x4f6c, 0x7df7, 0x6c7e,
	0xa50a, 0xb483, 0x8618, 0x9791, 0xe32e, 0xf2a7, 0xc03c, 0xd1b5,
	0x2942, 0x38cb, 0x0a50, 0x1bd9, 0x6f66, 0x7eef, 0x4c74, 0x5dfd,
	0xb58b, 0xa402, 0x9699, 0x8710, 0xf3af, 0xe226, 0xd0bd, 0xc134,
	0x39c3, 0x284a, 0x1ad1, 0x0b58, 0x7fe7, 0x6e6e, 0x5cf5, 0x4d7c,
	0xc60c, 0xd785, 0xe51e, 0xf497, 0x8028, 0x91a1, 0xa33a, 0xb2b3,
	0x4a44, 0x5bcd, 0x6956, 0x78df, 0x0c60, 0x1de9, 0x2f72, 0x3efb,
	0xd68d, 0xc704, 0xf59f, 0xe416, 0x90a9, 0x8120, 0xb3bb, 0xa232,
	0x5ac5, 0x4b4c, 0x79d7, 0x685e, 0x1ce1, 0x0d68, 0x3ff3, 0x2e7a,
	0xe70e, 0xf687, 0xc41c, 0xd595, 0xa12a, 0xb0a3, 0x8238, 0x93b1,
	0x6b46, 0x7acf, 0x4854, 0x59dd, 0x2d62, 0x3ceb, 0x0e70, 0x1ff9,
	0xf78f, 0xe606, 0xd49d, 0xc514, 0xb1ab, 0xa022, 0x92b9, 0x8330,
	0x7bc7, 0x6a4e, 0x58d5, 0x495c, 0x3de3, 0x2c6a, 0x1ef1, 0x0f78}

const PPPINITFCS16 = uint16(0xffff) // Initial FCS value
const PPPGOODFCS16 = uint16(0xf0b8) // Good final FCS value

/*
 * Calculate a new fcs given the current fcs and the new data.
 */
func pppfcs16(fcs16 uint16, p []byte) uint16 {
	for i := 0; i < len(p); i++ {
		// fcs = (fcs >> 8) ^ fcstab[(fcs ^ *cp++) & 0xff];
		fcs16 = (fcs16 >> 8) ^ fcstab[(fcs16^uint16(p[i]))&0x00ff]
	}
	return fcs16
}

/*
    // How to use the fcs

   tryfcs16(cp, len)
       register unsigned char *cp;
       register int len;
   {
       u16 trialfcs;

       // add on output
       trialfcs = pppfcs16( PPPINITFCS16, cp, len );
       trialfcs ^= 0xffff;                  // complement
       cp[len] = (trialfcs & 0x00ff);       // least significant byte first
       cp[len+1] = ((trialfcs >> 8) & 0x00ff);




       // check on input
       trialfcs = pppfcs16( PPPINITFCS16, cp, len + 2 );
       if ( trialfcs == PPPGOODFCS16 )
           printf("Good FCS\n");
   }
*/

//TODO: better error reporting
var HdlcErrorMalformedSegment = errors.New("malformed segment")
var HdlcErrorInvalidValue = errors.New("invalid value")
var HdlcErrorTimeout = errors.New("time out")
var HdlcErrorConnected = errors.New("time out")
var HdlcErrorProtocolError = errors.New("protocl error")

func NewHdlcTransport(rw io.ReadWriter, responseTimeoutMilliseconds int) *HdlcTransport {
	htran := new(HdlcTransport)
	htran.rw = rw
	htran.responseTimeout = responseTimeoutMilliseconds
	return htran
}

func (htran *HdlcTransport) decodeServerAddress(frame *HdlcFrame) (err error, n int) {
	var r io.Reader = htran.rw
	n = 0

	var b0, b1, b2, b3 byte
	p := make([]byte, 1)

	if !((HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength)) {
		panic("wrong expected server address length value")
	}

	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n += 1
	frame.fcs16 = pppfcs16(frame.fcs16, p)
	b0 = p[0]

	if b0&0x01 > 0 {
		if HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength {
			frame.logicalDeviceId = (uint16(b0) & 0x00FE) >> 1
			frame.physicalDeviceId = nil
		} else {
			errorLog("short server address")
			return HdlcErrorMalformedSegment, n
		}
	} else {
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		b1 = p[0]

		if b1&0x01 > 0 {
			upperMAC := (uint16(b0) & 0x00FE) >> 1
			lowerMAC := (uint16(b1) & 0x00FE) >> 1
			if HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength {
				frame.logicalDeviceId = upperMAC
				frame.physicalDeviceId = new(uint16)
				*frame.physicalDeviceId = lowerMAC
			} else if HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength {
				if 0x007F == lowerMAC {
					// all station broadcast
					frame.logicalDeviceId = lowerMAC
					frame.physicalDeviceId = nil
				} else {
					errorLog("long server address")
					return HdlcErrorMalformedSegment, n
				}
			} else if HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength {
				frame.logicalDeviceId = upperMAC
				frame.physicalDeviceId = new(uint16)
				*frame.physicalDeviceId = lowerMAC
			} else {
				panic("assertion failed")
			}
		} else {
			_, err = r.Read(p)
			if nil != err {
				errorLog("r.Read() failed: %v", err)
				return err, n
			}
			n += 1
			frame.fcs16 = pppfcs16(frame.fcs16, p)
			b2 = p[0]

			if b2&0x01 > 0 {
				errorLog("short server address")
				return HdlcErrorMalformedSegment, n
			}

			_, err = r.Read(p)
			if nil != err {
				errorLog("r.Read() failed: %v", err)
				return err, n
			}
			n += 1
			frame.fcs16 = pppfcs16(frame.fcs16, p)
			b3 = p[0]

			if b3&0x01 > 0 {
				upperMAC := ((uint16(b0)&0x00FE)>>1)<<7 + ((uint16(b1) & 0x00FE) >> 1)
				lowerMAC := ((uint16(b2)&0x00FE)>>1)<<7 + ((uint16(b3) & 0x00FE) >> 1)

				if HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength {

					frame.logicalDeviceId = upperMAC
					frame.physicalDeviceId = new(uint16)
					*frame.physicalDeviceId = lowerMAC

				} else if HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength {
					if (0x3FFF == upperMAC) && (0x3FFF == lowerMAC) {
						// all station broadcast 0x3FFF
						frame.logicalDeviceId = 0x3FFF
						frame.physicalDeviceId = new(uint16)
						*frame.physicalDeviceId = 0x3FFF
					} else {
						errorLog("long server address")
						return HdlcErrorMalformedSegment, n
					}
				} else if HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength {
					if (0x3FFF == upperMAC) && (0x3FFF == lowerMAC) {
						// all station broadcast 0x3FFF
						frame.logicalDeviceId = 0x3FFF
						frame.physicalDeviceId = new(uint16)
						*frame.physicalDeviceId = 0x3FFF
					} else if (upperMAC == 0x3FFF) && (0x0001 == lowerMAC) && frame.callingPhysicalDevice {
						// event reporting
						frame.logicalDeviceId = upperMAC
						frame.physicalDeviceId = new(uint16)
						*frame.physicalDeviceId = lowerMAC
					} else {
						errorLog("long server address")
						return HdlcErrorMalformedSegment, n
					}
				} else {
					panic("assertion failed")
				}
			} else {
				errorLog("long server address")
				return HdlcErrorMalformedSegment, n
			}
		}
	}
	return nil, n
}

func (htran *HdlcTransport) encodeServerAddress(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw

	var v16 uint16
	p := make([]byte, 1)

	if !((HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength)) {
		panic("wrong expected server address length value")
	}

	if HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength {
		p := make([]byte, 1)

		// logicalDeviceId

		logicalDeviceId := frame.logicalDeviceId
		if logicalDeviceId&0xFF80 > 0 {
			errorLog("logicalDeviceId exceeds limit")
			return HdlcErrorInvalidValue
		}

		v16 = (logicalDeviceId << 1) | 0x0001

		p[0] = byte(v16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// physicalDeviceId

		if nil != frame.physicalDeviceId {
			errorLog("physicalDeviceId specified (expected to be nil)")
			return HdlcErrorInvalidValue
		}

	} else if HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength {

		// logicalDeviceId

		logicalDeviceId := frame.logicalDeviceId
		if logicalDeviceId&0xFF80 > 0 {
			errorLog("logicalDeviceId exceeds limit")
			return HdlcErrorInvalidValue
		}

		v16 = (logicalDeviceId << 1) | 0x0001

		p[0] = byte(v16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// physicalDeviceId

		if nil == frame.physicalDeviceId {
			errorLog("physicalDeviceId not specified")
			return HdlcErrorInvalidValue
		}

		physicalDeviceId := *frame.physicalDeviceId
		if physicalDeviceId&0xFF80 > 0 {
			errorLog("physicalDeviceId exceeds limit")
			return HdlcErrorInvalidValue
		}

		v16 = (physicalDeviceId << 1) | 0x0001

		p[0] = byte(v16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength {

		// logicalDeviceId

		logicalDeviceId := frame.logicalDeviceId
		if logicalDeviceId&0x1000 > 0 {
			errorLog("logicalDeviceId exceeds limit")
			return HdlcErrorInvalidValue
		}

		v16 = (logicalDeviceId << 1) | 0x0001

		p[0] = byte((v16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		p[0] = byte(v16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// physicalDeviceId

		if nil == frame.physicalDeviceId {
			errorLog("physicalDeviceId not specified")
			return HdlcErrorInvalidValue
		}

		physicalDeviceId := *frame.physicalDeviceId
		if physicalDeviceId&0x1000 > 0 {
			errorLog("physicalDeviceId exceeds limit")
			return HdlcErrorInvalidValue
		}

		v16 = (physicalDeviceId << 1) | 0x0001

		p[0] = byte((v16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		p[0] = byte(v16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
	} else {
		panic("wrong expected server address length value")
	}

	return nil
}

func (htran *HdlcTransport) lengthServerAddress(frame *HdlcFrame) (n int) {

	n = 0

	if !((HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength) || (HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength)) {
		panic("wrong expected server address length value")
	}

	if HDLC_ADDRESS_LENGTH_1 == htran.expectedServerAddrLength {
		n = 1
	} else if HDLC_ADDRESS_LENGTH_2 == htran.expectedServerAddrLength {
		n += 2
	} else if HDLC_ADDRESS_LENGTH_4 == htran.expectedServerAddrLength {
		n += 4
	} else {
		panic("wrong expected server address length value")
	}

	return n
}

func (htran *HdlcTransport) decodeClientAddress(frame *HdlcFrame) (err error, n int) {
	var r io.Reader = htran.rw
	n = 0
	var b0 byte
	p := make([]byte, 1)

	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n += 1
	frame.fcs16 = pppfcs16(frame.fcs16, p)
	b0 = p[0]

	if b0&0x01 > 0 {
		frame.clientId = (uint8(b0) & 0xFE) >> 1
	} else {
		errorLog("long client address")
		return HdlcErrorMalformedSegment, n
	}

	return nil, n
}

func (htran *HdlcTransport) encodeClientAddress(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw
	var b0 byte
	p := make([]byte, 1)

	clientId := frame.clientId
	if clientId&0x80 > 0 {
		errorLog("clientId exceeds limit")
		return HdlcErrorInvalidValue
	}

	b0 = (clientId << 1) | 0x01

	p[0] = b0
	_, err = w.Write(p)
	if nil != err {
		errorLog("r.Write() failed: %v", err)
		return err
	}
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	return nil
}

func (htran *HdlcTransport) lengthClientAddress(frame *HdlcFrame) int {
	return 1
}

func (htran *HdlcTransport) decodeFrameInfo(frame *HdlcFrame, l int) (err error, n int) {
	var r io.Reader = htran.rw
	p := make([]byte, 1)

	// HCS - header control sum

	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n += 1
	l += 1
	frame.fcs16 = pppfcs16(frame.fcs16, p)
	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n += 1
	l += 1
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	if PPPGOODFCS16 != frame.fcs16 {
		errorLog("wrong HCS")
		return HdlcErrorMalformedSegment, n
	}

	// read information field

	infoFieldLength := frame.length - l

	if infoFieldLength > 0 {

		if (HDLC_FRAME_DIRECTION_CLIENT_INBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_INBOUND == frame.direction) {
			if infoFieldLength > htran.maxInfoFieldLengthReceive {
				errorLog("long info field")
				return HdlcErrorMalformedSegment, n
			}
		} else {
			panic("frame direction is not inbound")
		}

		p = make([]byte, infoFieldLength)
		err = binary.Read(r, binary.BigEndian, p)
		if nil != err {
			errorLog("binary.Read() failed: %v", err)
			return err, n
		}
		n += len(p)
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		frame.infoField = p

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else {
		frame.infoField = make([]byte, 0)
	}

	return nil, n
}

func (htran *HdlcTransport) encodeFrameInfo(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw
	p := make([]byte, 1)

	infoFieldLength := len(frame.infoField)

	if (HDLC_FRAME_DIRECTION_CLIENT_INBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_INBOUND == frame.direction) {
		if infoFieldLength > htran.maxInfoFieldLengthReceive {
			errorLog("long info field")
			return HdlcErrorMalformedSegment
		}
	} else if (HDLC_FRAME_DIRECTION_CLIENT_OUTBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_OUTBOUND == frame.direction) {
		if infoFieldLength > htran.maxInfoFieldLengthTransmit {
			errorLog("long info field")
			return HdlcErrorMalformedSegment
		}
	} else {
		panic("assertion failed")
	}

	// HCS - header control sum

	fcs16 := frame.fcs16
	p[0] = byte(^fcs16 & 0x00FF)
	_, err = w.Write(p)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}
	frame.fcs16 = pppfcs16(frame.fcs16, p)
	p[0] = byte((^fcs16 & 0xFF00) >> 8)
	_, err = w.Write(p)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	// write information field

	if (nil != frame.infoField) && len(frame.infoField) > 0 {

		if (HDLC_FRAME_DIRECTION_CLIENT_INBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_INBOUND == frame.direction) {
			if infoFieldLength > htran.maxInfoFieldLengthReceive {
				errorLog("long info field")
				return HdlcErrorMalformedSegment
			}
		} else if (HDLC_FRAME_DIRECTION_CLIENT_OUTBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_OUTBOUND == frame.direction) {
			if infoFieldLength > htran.maxInfoFieldLengthTransmit {
				errorLog("long info field")
				return HdlcErrorMalformedSegment
			}
		} else {
			panic("assertion failed")
		}

		p = frame.infoField
		err = binary.Write(w, binary.BigEndian, p)
		if nil != err {
			errorLog("binary.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
	}

	return nil
}

// decode frame address, control and information field

func (htran *HdlcTransport) decodeFrameACI(frame *HdlcFrame, l int) (err error, n int) {
	var r io.Reader = htran.rw
	n = 0
	var b0 byte
	var nn int

	p := make([]byte, 1)

	// dst and src address

	if HDLC_FRAME_DIRECTION_SERVER_INBOUND == frame.direction {
		err, nn = htran.decodeServerAddress(frame)
		if nil != err {
			return err, n
		}
		n += nn
		err, nn = htran.decodeClientAddress(frame)
		if nil != err {
			return err, n
		}
		n += nn
	} else if HDLC_FRAME_DIRECTION_CLIENT_INBOUND == frame.direction {
		err, nn = htran.decodeClientAddress(frame)
		if nil != err {
			return err, n
		}
		n += nn
		err, nn = htran.decodeServerAddress(frame)
		if nil != err {
			return err, n
		}
		n += nn
	} else {
		panic("frame direction is not inbound")
	}

	// control

	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n += 1
	frame.fcs16 = pppfcs16(frame.fcs16, p)
	b0 = p[0]

	// P/F bit
	frame.pf = b0&0x10 > 0

	if b0&0x01 == 0 {
		frame.control = HDLC_CONTROL_I

		frame.nr = b0 & 0xE0 >> 5
		frame.ns = b0 & 0x0E >> 1

		err, nn := htran.decodeFrameInfo(frame, l+n)
		if nil != err {
			return err, n
		}
		n += nn

	} else if (b0&0x08 == 0) && (b0&0x04 == 0) && (b0&0x02 == 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_RR

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x08 == 0) && (b0&0x04 > 0) && (b0&0x02 == 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_RNR

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x80 > 0) && (b0&0x40 == 0) && (b0&0x20 == 0) && (b0&0x08 == 0) && (b0&0x04 == 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {

		frame.control = HDLC_CONTROL_SNRM

		err, nn := htran.decodeFrameInfo(frame, l+n)
		if nil != err {
			return err, n
		}
		n += nn

	} else if (b0&0x80 == 0) && (b0&0x40 > 0) && (b0&0x20 == 0) && (b0&0x08 == 0) && (b0&0x04 == 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_DISC

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x80 == 0) && (b0&0x40 > 0) && (b0&0x20 > 0) && (b0&0x08 == 0) && (b0&0x04 == 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_UA

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x80 == 0) && (b0&0x40 == 0) && (b0&0x20 == 0) && (b0&0x08 > 0) && (b0&0x04 > 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_DM

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x80 > 0) && (b0&0x40 == 0) && (b0&0x20 == 0) && (b0&0x08 == 0) && (b0&0x04 > 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_FRMR

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else if (b0&0x80 == 0) && (b0&0x40 == 0) && (b0&0x20 == 0) && (b0&0x08 == 0) && (b0&0x04 == 0) && (b0&0x02 > 0) && (b0&0x01 > 0) {
		frame.control = HDLC_CONTROL_UI

		// FCS - frame control sum

		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n += 1
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		if PPPGOODFCS16 != frame.fcs16 {
			errorLog("wrong FCS")
			return HdlcErrorMalformedSegment, n
		}
	} else {
		errorLog("malformed control field")
		return HdlcErrorMalformedSegment, n
	}

	return nil, n
}

// encode frame address, control and information field

func (htran *HdlcTransport) encodeFrameACI(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw
	var b0 byte

	p := make([]byte, 1)

	// dst and src address

	if (HDLC_FRAME_DIRECTION_SERVER_OUTBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_CLIENT_INBOUND == frame.direction) {
		err = htran.encodeServerAddress(frame)
		if nil != err {
			return err
		}
		err = htran.encodeClientAddress(frame)
		if nil != err {
			return err
		}
	} else if (HDLC_FRAME_DIRECTION_CLIENT_OUTBOUND == frame.direction) || (HDLC_FRAME_DIRECTION_SERVER_INBOUND == frame.direction) {
		err = htran.encodeClientAddress(frame)
		if nil != err {
			return err
		}
		err = htran.encodeServerAddress(frame)
		if nil != err {
			return err
		}
	} else {
		panic("wrond frame direction")
	}

	// control

	// P/F bit
	b0 = 0
	if frame.pf {
		b0 |= 0x10
	}

	if HDLC_CONTROL_I == frame.control {

		if frame.nr > 0x07 {
			panic("NR exceeds limit")
		}
		b0 |= frame.nr << 5

		if frame.ns > 0x07 {
			panic("NS exceeds limit")
		}
		b0 |= frame.ns << 1

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		err = htran.encodeFrameInfo(frame)
		if nil != err {
			return err
		}

	} else if HDLC_CONTROL_RR == frame.control {
		b0 |= 0x01

		if frame.nr > 0x07 {
			panic("NR exceeds limit")
		}
		b0 |= frame.nr << 5

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_RNR == frame.control {
		b0 |= 0x01
		b0 |= 0x04

		if frame.nr > 0x07 {
			panic("NR exceeds limit")
		}
		b0 |= frame.nr << 5

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_SNRM == frame.control {
		b0 |= 0x01
		b0 |= 0x02
		b0 |= 0x80

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		err = htran.encodeFrameInfo(frame)
		if nil != err {
			return err
		}

	} else if HDLC_CONTROL_DISC == frame.control {
		b0 |= 0x01
		b0 |= 0x02
		b0 |= 0x40

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_UA == frame.control {
		b0 |= 0x01
		b0 |= 0x02
		b0 |= 0x20
		b0 |= 0x40

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_DM == frame.control {
		b0 |= 0x01
		b0 |= 0x02
		b0 |= 0x04
		b0 |= 0x08

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_FRMR == frame.control {
		b0 |= 0x01
		b0 |= 0x02
		b0 |= 0x04
		b0 |= 0x80

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else if HDLC_CONTROL_UI == frame.control {
		b0 |= 0x01
		b0 |= 0x02

		p[0] = b0
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

		// FCS - frame control sum

		fcs16 := frame.fcs16
		p[0] = byte(^fcs16 & 0x00FF)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		p[0] = byte((^fcs16 & 0xFF00) >> 8)
		_, err = w.Write(p)
		if nil != err {
			errorLog("w.Write() failed: %v", err)
			return err
		}
		frame.fcs16 = pppfcs16(frame.fcs16, p)

	} else {
		errorLog("invalid control field value")
		return HdlcErrorInvalidValue
	}

	return nil
}

func (htran *HdlcTransport) lengthOfFrame(frame *HdlcFrame) (n int) {
	n = 0

	// format type
	n += 2

	// src, dst address

	n += htran.lengthServerAddress(frame)
	n += htran.lengthClientAddress(frame)

	// control
	n += 1

	// HCS - header control sum
	n += 2

	if (nil != frame.infoField) && len(frame.infoField) > 0 {
		// FCS - frame control sum
		n += 2
	}

	return n
}

// decode frame format, address, control and information field

func (htran *HdlcTransport) decodeFrameFACI(frame *HdlcFrame, l int) (err error, n int) {
	var r io.Reader = htran.rw
	n = 0

	p := make([]byte, 1)
	var b0, b1 byte

	// expect first byte of format field
	_, err = r.Read(p)
	if nil != err {
		errorLog("r.Read() failed: %v", err)
		return err, n
	}
	n++
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	// format field
	if 0xA0 == p[0]&0xF0 {
		b0 = p[0]

		// expect last second byte of format field
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, n
		}
		n++
		frame.fcs16 = pppfcs16(frame.fcs16, p)
		b1 = p[0]

		frame.formatType = 0xA0

		// test segmentation bit
		if b0&0x08 > 0 {
			frame.segmentation = true
		} else {
			frame.segmentation = false
		}

		frame.length = int((uint16(b0&0x07) << 8) + uint16(b1))

		err, nn := htran.decodeFrameACI(frame, l+n)
		n += nn
		return err, n
	} else {
		return HdlcErrorMalformedSegment, n
	}
}

// encode frame format, address, control and information field

func (htran *HdlcTransport) encodeFrameFACI(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw

	p := make([]byte, 1)
	var b0, b1 byte

	// frame format
	b0 |= 0xA0

	// segmentation
	if frame.segmentation {
		b0 |= 0x08
	}

	length := uint16(htran.lengthOfFrame(frame))
	if length > 0x07FF {
		errorLog("frame length exceeds limt")
		return HdlcErrorInvalidValue
	}
	b0 |= byte((0xFF00 & length) >> 8)
	b1 = byte(0x00FF & length)

	p[0] = b0
	_, err = w.Write(p)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	p[0] = b1
	_, err = w.Write(p)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}
	frame.fcs16 = pppfcs16(frame.fcs16, p)

	return htran.encodeFrameACI(frame)

}

func (htran *HdlcTransport) readFrame(direction int) (err error, frame *HdlcFrame) {
	var r io.Reader = htran.rw
	p := make([]byte, 1)
	for {
		// expect opening flag
		_, err = r.Read(p)
		if nil != err {
			errorLog("r.Read() failed: %v", err)
			return err, nil
		}
		if 0x7E == p[0] { // flag
			frame := new(HdlcFrame)
			frame.direction = direction
			frame.fcs16 = PPPINITFCS16

			err, _ = htran.decodeFrameFACI(frame, 0)
			if nil != err {
				if HdlcErrorMalformedSegment == err {
					// ignore malformed segment and try read next segment
					continue
				} else {
					return nil, frame
				}
			}

		} else {
			// ignore everything until leading flag arrives
			continue
		}
	}
}

func (htran *HdlcTransport) readFrameAsync(direction int) <-chan map[string]interface{} {
	ch := make(chan map[string]interface{})
	func() {
		err, frame := htran.readFrame(direction)
		ch <- map[string]interface{}{"err": err, "frame": frame}
	}()
	return ch
}

func (htran *HdlcTransport) writeFrame(frame *HdlcFrame) (err error) {
	var w io.Writer = htran.rw

	if 0 == frame.direction {
		errorLog("frame direction not specified")
		return HdlcErrorInvalidValue
	}
	if 0 == frame.control {
		errorLog("frame controltype not specified")
		return HdlcErrorInvalidValue
	}

	p := make([]byte, 1)

	// opening flag
	p[0] = 0x7E
	_, err = w.Write(p)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}

	err = htran.encodeFrameFACI(frame)
	if nil != err {
		errorLog("w.Write() failed: %v", err)
		return err
	}

	return nil

}

func (htran *HdlcTransport) writeFrameAsync(frame *HdlcFrame) <-chan map[string]interface{} {
	ch := make(chan map[string]interface{})
	go func() {
		err := htran.writeFrame(frame)
		ch <- map[string]interface{}{"err": err}
	}()
	return ch
}

func (htran *HdlcTransport) sendSNRM() (err error, conn *HdlcClientConnection) {
	if htran.connected {
		return HdlcErrorConnected, nil
	}

	frame := new(HdlcFrame)
	frame.direction = HDLC_FRAME_DIRECTION_SERVER_OUTBOUND
	frame.control = HDLC_CONTROL_SNRM
	htran.writeFrame(frame)

	ch := make(chan map[string]interface{})

	go func(ch chan map[string]interface{}) {
		select {
		case _ = <-time.After(time.Duration(htran.responseTimeout) * time.Millisecond):
			errorLog("SNRM response timeout")
			ch <- map[string]interface{}{"err": HdlcErrorTimeout}
		case msg := <-htran.readFrameAsync(frame.direction):
			ch <- msg
		}
	}(ch)

	msg := <-ch
	if nil != msg["err"] {
		return (msg["err"]).(error), nil
	}
	frame = (msg["frame"]).(*HdlcFrame)
	if HDLC_CONTROL_UA == frame.control {
		conn := new(HdlcClientConnection)
		htran.connected = true
		conn.htran = htran
		conn.vs = 0
		conn.vr = 0
		return nil, conn
	} else if HDLC_CONTROL_DM == frame.control {
		warnLog("received DM")
		// per standard we must ignore DM mesaage which should result in timeout
		return HdlcErrorTimeout, nil
	} else {
		errorLog("unexpected frame: %d", frame.control)
		return HdlcErrorProtocolError, nil
	}

}

func (htran *HdlcTransport) listen() (err error, conn *HdlcServerConnection) {
	if htran.connected {
		return HdlcErrorConnected, nil
	}

	err, frame := htran.readFrame(HDLC_FRAME_DIRECTION_SERVER_INBOUND)
	if HDLC_CONTROL_SNRM == frame.control {
		conn = new(HdlcServerConnection)
		conn.htran = htran
		conn.vs = 0
		conn.vr = 0

		frame = new(HdlcFrame)
		frame.direction = HDLC_FRAME_DIRECTION_SERVER_OUTBOUND
		frame.control = HDLC_CONTROL_UA
		err = htran.writeFrame(frame)
		if nil != err {
			return err, nil
		}

		return nil, conn
	} else {
		warnLog("unexpected frame: %s", frame.control)
		return HdlcErrorProtocolError, nil
	}
}

func (conn *HdlcClientConnection) Write(p []byte) (n int, err error) {
	var segment []byte
	var i, l int

	maxSegemntLen := conn.htran.maxInfoFieldLengthTransmit

	// split to segments

	ns = len(p) / maxSegemntLen
	if len(p)%maxSegemntLen > 0 {
		ns += 1
	}
	segments = make([]byte, ns)
	_p := p
	for i = 0; i < ns; i++ {
		if len(_p) >= maxSegemntLen {
			segments[i] = _p[0:maxSegemntLen]
			_p = _p[maxSegemntLen:]
		} else {
			segments[i] = _p[0:len(_p)]
		}
	}

	for len(segment) > 0 {

		if conn.vs+1 > 7 {
			// we ran out all available sequence numbers therefore we must wait for acknowledgement of last received sequence number

			ch := make(chan map[string]interface{})
			go func(ch chan map[string]interface{}) {
				select {
				case _ = <-time.After(time.Duration(htran.responseTimeout) * time.Millisecond):
					errorLog("SNRM response timeout")
					ch <- map[string]interface{}{"err": HdlcErrorTimeout}
				case msg := <-htran.readFrameAsync(frame.direction):
					ch <- msg
				}
			}(ch)

		} else {
			conn.vs += 1

			frame := new(HdlcFrame)
			frame.direction = HDLC_FRAME_DIRECTION_CLIENT_OUTBOUND
			if len(p) > 0 {
				frame.segmentation = true
			} else {
				frame.segmentation = false
				frame.poll = true
			}
			frame.ns = conn.vs
			frame.nr = conn.vr
			err := writeFrame(frame)
			if nil != err {
				//TODO: disconnect
				return err
			}
		}

		if len(p) > maxSegemntLen {
			l = maxSegemntLen
		} else {
			l = len(p)
		}
		segment = p[0:l]
	}
	return 0, nil
}
