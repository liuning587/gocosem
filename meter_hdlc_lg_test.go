package gocosem

import (
	"fmt"
	"testing"
	"time"
)

func init_TestMeterLgHdlc() {
	testMeterIp = "172.16.123.206"
	testHdlcResponseTimeout = time.Duration(1) * time.Hour
	testHdlcCosemWaitTime = time.Duration(5000) * time.Millisecond
	testHdlcSnrmTimeout = time.Duration(45) * time.Second
	testHdlcDiscTimeout = time.Duration(45) * time.Second
}

func TestMeterLgHdlc_TcpConnect(t *testing.T) {
	init_TestMeterLgHdlc()
	dconn, err := TcpConnect(testMeterIp, 4059)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()
}

func TestMeterLgHdlc_HdlcConnect(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()
}

func TestMeterLgHdlc_AppConnect(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()
}

func TestMeterLgHdlc_GetTime(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	val := new(DlmsRequest)
	val.ClassId = 8
	val.InstanceId = &DlmsOid{0x00, 0x00, 0x01, 0x00, 0x00, 0xFF}
	val.AttributeId = 0x02
	vals := make([]*DlmsRequest, 1)
	vals[0] = val
	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf(fmt.Sprintf("%s\n", err))
	}
	t.Logf("response delivered: in %v", rep.DeliveredIn())
	if 0 != rep.DataAccessResultAt(0) {
		t.Fatalf("dataAccessResult: %d\n", rep.DataAccessResultAt(0))
	}
	data := rep.DataAt(0)
	t.Logf("value read %#v", data.Val)
	t.Logf("datetime: %s", DlmsDateTimeFromBytes(data.GetOctetString()).PrintDateTime())
}

func TestMeterLgHdlc_SetTime(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// read time

	val := new(DlmsRequest)
	val.ClassId = 8
	val.InstanceId = &DlmsOid{0x00, 0x00, 0x01, 0x00, 0x00, 0xFF}
	val.AttributeId = 0x02
	vals := make([]*DlmsRequest, 1)
	vals[0] = val
	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("%s\n", err)
	}
	t.Logf("response delivered: in %v", rep.DeliveredIn())
	if 0 != rep.DataAccessResultAt(0) {
		t.Fatalf("dataAccessResult: %d\n", rep.DataAccessResultAt(0))
	}
	data := rep.DataAt(0)
	t.Logf("value read %#v", data.Val)
	t.Logf("datetime: %s", DlmsDateTimeFromBytes(data.GetOctetString()).PrintDateTime())

	// set time

	val = new(DlmsRequest)
	val.ClassId = 8
	val.InstanceId = &DlmsOid{0x00, 0x00, 0x01, 0x00, 0x00, 0xFF}
	val.AttributeId = 0x02
	val.Data = data
	vals = make([]*DlmsRequest, 1)
	vals[0] = val
	rep, err = aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("%s\n", err)
	}
	t.Logf("response delivered: in %v", rep.DeliveredIn())
	if 0 != rep.DataAccessResultAt(0) {
		t.Fatalf("dataAccessResult: %d\n", rep.DataAccessResultAt(0))
	}
	t.Logf("time set successfully")

	// read time again

	val = new(DlmsRequest)
	val.ClassId = 8
	val.InstanceId = &DlmsOid{0x00, 0x00, 0x01, 0x00, 0x00, 0xFF}
	val.AttributeId = 0x02
	vals = make([]*DlmsRequest, 1)
	vals[0] = val
	rep, err = aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("%s\n", err)
	}
	t.Logf("response delivered: in %v", rep.DeliveredIn())
	if 0 != rep.DataAccessResultAt(0) {
		t.Fatalf("dataAccessResult: %d\n", rep.DataAccessResultAt(0))
	}
	data = rep.DataAt(0)
	t.Logf("value read %#v", data.Val)
	t.Logf("datetime: %s", DlmsDateTimeFromBytes(data.GetOctetString()).PrintDateTime())

}

func TestMeterLgHdlc_ProfileCaptureObjects(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// capture objects definitions

	t.Logf("read objects captured by profile...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 3
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	if DATA_TYPE_ARRAY != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("profile captures follwing objects:")
	for i, st := range data.Arr {
		if DATA_TYPE_STRUCTURE != st.GetType() {
			t.Fatalf("wrong data type")
		}
		t.Logf("capture object [%d]:", i)
		t.Logf("\tclass id: %d", st.Arr[0].GetLongUnsigned())
		t.Logf("\tlogical name: % 02X", st.Arr[1].GetOctetString())
		t.Logf("\tattribute index: %d", st.Arr[2].GetInteger())
		t.Logf("\tdata index: % 02X", st.Arr[3].GetLongUnsigned())
	}
}

//@@@@@@@@@@@@@@@@@@@@@@@@
func TestMeterLgHdlc_ProfileEntriesInUse(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// profile entries in use

	t.Logf("read profile entries in use...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 7
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	t.Logf("profile entries in use: %d", data.GetDoubleLongUnsigned())
}

func TestMeterLgHdlc_ProfileEntries(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// profile entries

	t.Logf("read profile entries in use...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 8
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	t.Logf("maximum profile entries: %d", data.GetDoubleLongUnsigned())
}

func TestMeterLgHdlc_ProfileSortMethod(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// sort method

	t.Logf("read sort method ...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 5
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	t.Logf("sort method: %d", data.GetEnum())
}

func TestMeterLgHdlc_ProfileSortObject(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// sort object

	t.Logf("read sort object ...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 6
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)

	if DATA_TYPE_STRUCTURE != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("sort object:")
	t.Logf("\tclass id: %d", data.Arr[0].GetLongUnsigned())
	t.Logf("\tlogical name: % 02X", data.Arr[1].GetOctetString())
	t.Logf("\tattribute index: %d", data.Arr[2].GetInteger())
	t.Logf("\tdata index: % 02X", data.Arr[3].GetLongUnsigned())
}

func TestMeterLgHdlc_ProfileCapturePeriod(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// capture period

	t.Logf("read capture period ...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 4
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	t.Logf("capture period: %d seconds", data.GetDoubleLongUnsigned())
}

//@@@@@@@@@@@@@@@@@@@@@@@@

func TestMeterLgHdlc_ProfileFirstEntries(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// request first 10 entries

	vals := make([]*DlmsRequest, 1)

	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 2
	val.AccessSelector = 2
	val.AccessParameter = new(DlmsData)
	val.AccessParameter.SetStructure(4)
	val.AccessParameter.Arr[0].SetDoubleLongUnsigned(1)  // from_entry
	val.AccessParameter.Arr[1].SetDoubleLongUnsigned(10) // to_entry
	val.AccessParameter.Arr[2].SetLongUnsigned(1)        // from_selected_value
	val.AccessParameter.Arr[3].SetLongUnsigned(0)        // to_selected_value

	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}

	data := rep.DataAt(0) // first request
	if DATA_TYPE_ARRAY != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("profile entries read:\n")
	for i := 0; i < len(data.Arr); i++ {
		d := data.Arr[i]
		d4 := d.Arr[4]
		t.Logf("\t%d: %s %s: ", i, DlmsDateTimeFromBytes(d4.GetOctetString()).PrintDateTime(), d.Print())
	}
}

func TestMeterLgHdlc_ProfileLastEntries(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// profile entries in use

	t.Logf("read profile entries ...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 7
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	entriesInUse := data.GetDoubleLongUnsigned()
	t.Logf("profile entries in use: %d", entriesInUse)

	vals = make([]*DlmsRequest, 1)

	// read last 10 entries

	val = new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 2
	val.AccessSelector = 2
	val.AccessParameter = new(DlmsData)
	val.AccessParameter.SetStructure(4)
	if entriesInUse > 10 {
		val.AccessParameter.Arr[0].SetDoubleLongUnsigned(entriesInUse - 10 + 1) // from_entry
	} else {
		val.AccessParameter.Arr[0].SetDoubleLongUnsigned(1) // from_entry
	}
	val.AccessParameter.Arr[1].SetDoubleLongUnsigned(entriesInUse) // to_entry
	val.AccessParameter.Arr[2].SetLongUnsigned(1)                  // from_selected_value
	val.AccessParameter.Arr[3].SetLongUnsigned(0)                  // to_selected_value

	vals[0] = val

	rep, err = aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult = rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}

	data = rep.DataAt(0) // first request
	if DATA_TYPE_ARRAY != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("profile entries read:\n")
	for i := 0; i < len(data.Arr); i++ {
		d := data.Arr[i]
		d4 := d.Arr[4]
		t.Logf("\t%d: %s %s: ", i, DlmsDateTimeFromBytes(d4.GetOctetString()).PrintDateTime(), d.Print())
	}
}

func failing_TestMeterLgHdlc_ProfileTimeRange(t *testing.T) {
	init_TestMeterLgHdlc()

	applicationClient := uint16(0x11)
	logicalDevice := uint16(0x01)
	physicalDevice := uint16(0x7F)

	dconn, err := HdlcConnect(testMeterIp, 4059, applicationClient, logicalDevice, &physicalDevice, nil, testHdlcResponseTimeout, &testHdlcCosemWaitTime, testHdlcSnrmTimeout, testHdlcDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x3A, 0x80, 0x02, 0x02, 0x84, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x41, 0x32, 0x42, 0x33, 0x43, 0x34, 0x44, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x04, 0x00}
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0x18, 0x1D, 0x01, 0xFA, 0x00, 0x07}
	aconn, err := dconn.AppConnectRaw(applicationClient, logicalDevice, 8, aarq, aare)

	if nil != err {
		t.Fatal(err)
	}
	t.Logf("application connected")
	defer aconn.Close()

	// profile entries in use

	t.Logf("read profile entries ...")
	vals := make([]*DlmsRequest, 1)
	val := new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 7
	vals[0] = val

	rep, err := aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult := rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}
	data := rep.DataAt(0)
	entriesInUse := data.GetDoubleLongUnsigned()
	t.Logf("profile entries in use: %d", entriesInUse)

	vals = make([]*DlmsRequest, 1)

	// read last 10 entries

	val = new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 2
	val.AccessSelector = 2
	val.AccessParameter = new(DlmsData)
	val.AccessParameter.SetStructure(4)
	if entriesInUse > 10 {
		val.AccessParameter.Arr[0].SetDoubleLongUnsigned(entriesInUse - 10 + 1) // from_entry
	} else {
		val.AccessParameter.Arr[0].SetDoubleLongUnsigned(1) // from_entry
	}
	val.AccessParameter.Arr[1].SetDoubleLongUnsigned(entriesInUse) // to_entry
	val.AccessParameter.Arr[2].SetLongUnsigned(1)                  // from_selected_value
	val.AccessParameter.Arr[3].SetLongUnsigned(0)                  // to_selected_value

	vals[0] = val

	rep, err = aconn.SendRequest(vals)
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult = rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}

	data = rep.DataAt(0) // first request
	if DATA_TYPE_ARRAY != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("profile entries read:\n")
	for i := 0; i < len(data.Arr); i++ {
		d := data.Arr[i]
		d4 := d.Arr[4]
		t.Logf("\t%d: %s %s: ", i, DlmsDateTimeFromBytes(d4.GetOctetString()).PrintDateTime(), d.Print())
	}

	d1 := data.Arr[0]
	if nil != d1.Err {
		t.Fatalf("data error: %v", d1.Err)
	}
	d2 := data.Arr[len(data.Arr)-1]
	if nil != d2.Err {
		t.Fatalf("data error: %v", d2.Err)
	}

	fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ cp 100: last 10 entries read ok\n")

	// read last 10 entries using time interval selection

	vals = make([]*DlmsRequest, 1)

	val = new(DlmsRequest)
	val.ClassId = 7
	val.InstanceId = &DlmsOid{1, 0, 99, 1, 0, 255}
	val.AttributeId = 2
	val.AccessSelector = 1
	val.AccessParameter = new(DlmsData)
	val.AccessParameter.SetStructure(4)

	// selecting according fourh column which is the time

	restrictingObject := new(DlmsData)
	restrictingObject.SetStructure(4)
	restrictingObject.Arr[0].SetLongUnsigned(8)                                         // class_id
	restrictingObject.Arr[1].SetOctetString([]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0xFF}) // logical_name
	restrictingObject.Arr[2].SetInteger(2)                                              // attribute_index
	restrictingObject.Arr[3].SetLongUnsigned(0)                                         // data_index

	tim := DlmsDateTimeFromBytes(d1.Arr[4].GetOctetString())
	/*
		tim := new(DlmsDateTime)

		tim.Year = 2016
		tim.Month = 2
		tim.DayOfMonth = 22
		tim.DayOfWeek = 1
		tim.Hour = 4
		tim.Minute = 16
		tim.Second = 39
		tim.Hundredths = 0
		tim.Deviation = 0
		tim.ClockStatus = 0
	*/

	t.Logf("time from: %s", tim.PrintDateTime())

	// for some reason deviation and status must be zeroed or else this meter reports error
	//tim.Deviation = 0
	//tim.ClockStatus = 0

	fromValue := new(DlmsData)
	fromValue.SetOctetString(tim.ToBytes())

	tim = DlmsDateTimeFromBytes(d2.Arr[4].GetOctetString())
	/*
		tim = new(DlmsDateTime)

		tim.Year = 2016
		tim.Month = 2
		tim.DayOfMonth = 22
		tim.DayOfWeek = 1
		tim.Hour = 5
		tim.Minute = 16
		tim.Second = 39
		tim.Hundredths = 0
		tim.Deviation = 0
		tim.ClockStatus = 0
	*/

	// for some reason deviation and status must be zeroed or else this meter reports error
	//tim.Deviation = 0
	//tim.ClockStatus = 0

	t.Logf("time to: %s", tim.PrintDateTime())

	toValue := new(DlmsData)
	toValue.SetOctetString(tim.ToBytes())

	selectedValues := new(DlmsData)
	selectedValues.SetArray(0)

	val.AccessParameter.Arr[0] = restrictingObject
	val.AccessParameter.Arr[1] = fromValue
	val.AccessParameter.Arr[2] = toValue
	val.AccessParameter.Arr[3] = selectedValues

	vals[0] = val

	fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ cp 200\n")
	rep, err = aconn.SendRequest(vals)
	fmt.Printf("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ cp 201\n")
	if nil != err {
		t.Fatalf("read failed: %s", err)
		return
	}
	dataAccessResult = rep.DataAccessResultAt(0)
	if 0 != dataAccessResult {
		t.Fatalf("data access result: %d", dataAccessResult)
	}

	data = rep.DataAt(0) // first request
	if DATA_TYPE_ARRAY != data.GetType() {
		t.Fatalf("wrong data type")
	}
	t.Logf("profile entries read:\n")
	for i := 0; i < len(data.Arr); i++ {
		d := data.Arr[i]
		d4 := d.Arr[4]
		t.Logf("\t%d: %s %s: ", i, DlmsDateTimeFromBytes(d4.GetOctetString()).PrintDateTime(), d.Print())
	}
}
