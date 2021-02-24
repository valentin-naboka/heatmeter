package mbus

/*

//#cgo LDFLAGS: -L/Users/valentin/dev/3rd_party/libmbus/mbus/.libs/ -lmbus1
#cgo LDFLAGS:  /usr/local/lib/libmbus.a

#include "/Users/valentin/dev/3rd_party/libmbus/mbus/mbus-protocol-aux.h"
#include "/Users/valentin/dev/3rd_party/libmbus/mbus/mbus-serial.h"

// Workaround to explicitly convert void* to mbus_frame*,
// since Go can't convert unsafe.Pointer to concrete type C.mbus_frame
mbus_frame* to_frame(void* p){
	return p;
}
*/
import "C"

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/net/html/charset"
)

type Reader struct {
	handle  *C.mbus_handle
	address C.int
}

func (r *Reader) Open(device string, primaryAddress uint8, baudrate uint16) error {
	dev := C.CString(device)
	defer C.free(unsafe.Pointer(dev))

	r.handle = C.mbus_context_serial(dev)

	if C.mbus_connect(r.handle) != 0 {
		return fmt.Errorf("failed to setup connection to M-bus gateway: %s", device)
	}

	if C.mbus_serial_set_baudrate(r.handle, C.long(baudrate)) != 0 {
		return fmt.Errorf("failed to set baud rate: %d", baudrate)
	}

	r.address = C.int(primaryAddress)
	if C.mbus_send_ping_frame(r.handle, r.address, 1) != 0 {
		return fmt.Errorf("failed to deselect address: %d", primaryAddress)
	}
	return nil
}

func (r *Reader) Close() error {
	defer C.mbus_context_free(r.handle)
	if C.mbus_disconnect(r.handle) != 0 {
		return errors.New("failed to disconnect")
	}
	return nil
}

func (r *Reader) ReadData() (*Measurement, error) {
	// Application reset subcode to request "instant values" instead of a default telegram.
	// According to https://www.diehl.com/cms/files/254320-FR-EN-SHARKY_774_CommunicationDescription_v1_0_EN.pdf?download=1
	subcode := 0x50
	if C.mbus_send_application_reset_frame(r.handle, r.address, C.int(subcode)) == -1 {
		return nil, fmt.Errorf("failed to send reset frame: %s", C.GoString(C.mbus_error_str()))
	}

	var reply C.mbus_frame
	ret := C.mbus_recv_frame(r.handle, &reply)

	if ret == C.MBUS_RECV_RESULT_TIMEOUT {
		return nil, fmt.Errorf("failed to get a reply from device: timeout expired")
	}

	if C.mbus_frame_type(&reply) != C.MBUS_FRAME_TYPE_ACK {
		return nil, fmt.Errorf("unexpected frame type, receiving ACK telegram is failed")
	}

	const maxFrames C.int = 16
	if C.mbus_sendrecv_request(r.handle, r.address, &reply, maxFrames) != 0 {
		C.mbus_frame_free(C.to_frame(reply.next))
		return nil, fmt.Errorf("failed to send/receive M-Bus request: %s", C.GoString(C.mbus_error_str()))
	}

	var frameData C.mbus_frame_data
	if C.mbus_frame_data_parse(&reply, &frameData) != 0 {
		return nil, fmt.Errorf("M-bus data parse error: %s", C.GoString(C.mbus_error_str()))
	}

	xmlOutput := C.mbus_frame_data_xml_normalized(&frameData)
	defer C.free(unsafe.Pointer(xmlOutput))
	if frameData.data_var.record != nil {
		defer C.mbus_data_record_free(frameData.data_var.record)
	}

	if xmlOutput == nil {
		return nil, fmt.Errorf("failed to generate XML output of the frame: %s", C.GoString(C.mbus_error_str()))
	}

	reader := bytes.NewReader([]byte(C.GoString(xmlOutput)))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var measurement Measurement
	err := decoder.Decode(&measurement)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML output: %w", err)
	}

	return &measurement, nil
}
