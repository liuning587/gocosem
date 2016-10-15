package gocosem

import (
	"testing"
)

func TestActions_hdlcMeter_StateOfDisconnector(t *testing.T) {
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0xFE, 0x1D, 0x00, 0xEF, 0x00, 0x07}

	dconn, err := HdlcConnect(hdlcTestMeterIp, 4059, 1, 1, nil, hdlcTestResponseTimeout, &hdlcTestCosemWaitTime, hdlcTestSnrmTimeout, hdlcTestDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x36, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0xFF, 0xFF, 0xFF, 0x02, 0x00}
	aconn, err := dconn.AppConnectRaw(01, 01, 4, aarq, aare)
	if nil != err {
		t.Fatal(err)
	}
	defer aconn.Close()

	instanceId := &DlmsOid{0x00, 0x00, 0x60, 0x03, 0x0A, 0xFF}
	classId := DlmsClassId(70)
	attributeIdControlState := DlmsAttributeId(3)
	attributeIdControlMode := DlmsAttributeId(4)

	// Read control mode

	val := new(DlmsRequest)
	val.ClassId = classId
	val.InstanceId = instanceId
	val.AttributeId = attributeIdControlMode
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
	if DATA_TYPE_ENUM != data.GetType() {
		t.Fatalf("not integer")
	}
	controlMode := data.GetEnum()
	t.Logf("control mode: %d", controlMode)

	// Check connected state.

	val = new(DlmsRequest)
	val.ClassId = classId
	val.InstanceId = instanceId
	val.AttributeId = attributeIdControlState
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
	if DATA_TYPE_ENUM != data.GetType() {
		t.Fatalf("not integer")
	}

	controlState := data.GetEnum()
	t.Logf("control state: %d", controlState)

}

func TestActions_hdlcMeter_Disconnector(t *testing.T) {
	aare := []byte{0x61, 0x29, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0xA2, 0x03, 0x02, 0x01, 0x00, 0xA3, 0x05, 0xA1, 0x03, 0x02, 0x01, 0x00, 0xBE, 0x10, 0x04, 0x0E, 0x08, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0x00, 0xFE, 0x1D, 0x00, 0xEF, 0x00, 0x07}

	dconn, err := HdlcConnect(hdlcTestMeterIp, 4059, 1, 1, nil, hdlcTestResponseTimeout, &hdlcTestCosemWaitTime, hdlcTestSnrmTimeout, hdlcTestDiscTimeout)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("transport connected")
	defer dconn.Close()

	aarq := []byte{0x60, 0x36, 0xA1, 0x09, 0x06, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x01, 0x01, 0x8A, 0x02, 0x07, 0x80, 0x8B, 0x07, 0x60, 0x85, 0x74, 0x05, 0x08, 0x02, 0x01, 0xAC, 0x0A, 0x80, 0x08, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0xBE, 0x10, 0x04, 0x0E, 0x01, 0x00, 0x00, 0x00, 0x06, 0x5F, 0x1F, 0x04, 0x00, 0xFF, 0xFF, 0xFF, 0x02, 0x00}
	aconn, err := dconn.AppConnectRaw(01, 01, 4, aarq, aare)
	if nil != err {
		t.Fatal(err)
	}
	defer aconn.Close()

	instanceId := &DlmsOid{0x00, 0x00, 0x60, 0x03, 0x0A, 0xFF}
	classId := DlmsClassId(70)
	attributeIdControlState := DlmsAttributeId(3)
	attributeIdControlMode := DlmsAttributeId(4)
	methodIdRemoteDisconnect := DlmsMethodId(1)
	methodIdRemoteConnect := DlmsMethodId(2)
	//stateDisconnected := uint8(0)
	//stateConnected := uint8(1)

	// Read control mode

	val := new(DlmsRequest)
	val.ClassId = classId
	val.InstanceId = instanceId
	val.AttributeId = attributeIdControlMode
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
	if DATA_TYPE_ENUM != data.GetType() {
		t.Fatalf("not integer")
	}
	controlMode := data.GetEnum()
	t.Logf("control mode: %d", data.GetEnum())

	// Check if control mode is acceptable.
	if controlMode != 2 {
		t.Fatalf("unsupported control mode")
	}

	// Check connected state.

	val = new(DlmsRequest)
	val.ClassId = classId
	val.InstanceId = instanceId
	val.AttributeId = attributeIdControlState
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
	if DATA_TYPE_ENUM != data.GetType() {
		t.Fatalf("not integer")
	}

	controlState := data.GetEnum()
	t.Logf("control state: %d", controlState)

	// Based on current control state try to disconnect or connect.
	// At the end of test always return meter to connected state.

	if 1 == controlState { // connected

		// Call remote_disconnect method.

		method := new(DlmsRequest)
		method.ClassId = classId
		method.InstanceId = instanceId
		method.MethodId = methodIdRemoteDisconnect
		methodParameters := new(DlmsData)
		methodParameters.SetInteger(1)
		method.MethodParameters = methodParameters
		methods := make([]*DlmsRequest, 1)
		methods[0] = method
		rep, err = aconn.SendRequest(methods)
		if nil != err {
			t.Fatalf("%s\n", err)
		}
		t.Logf("response delivered: in %v", rep.DeliveredIn())
		if 0 != rep.ActionResultAt(0) {
			t.Fatalf("actionResult: %d\n", rep.ActionResultAt(0))
		}

		// Check connected state.

		val = new(DlmsRequest)
		val.ClassId = classId
		val.InstanceId = instanceId
		val.AttributeId = attributeIdControlState
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
		if DATA_TYPE_ENUM != data.GetType() {
			t.Fatalf("not integer")
		}

		controlState := data.GetEnum()
		t.Logf("control state: %d", controlState)
		if 0 != controlState {
			t.Fatalf("meter did not disconnect, controlState: %d", controlState)
		}

		// Call remote_connect method.

		method = new(DlmsRequest)
		method.ClassId = classId
		method.InstanceId = instanceId
		method.MethodId = methodIdRemoteConnect
		methodParameters = new(DlmsData)
		methodParameters.SetInteger(1)
		method.MethodParameters = methodParameters
		methods = make([]*DlmsRequest, 1)
		methods[0] = method
		rep, err = aconn.SendRequest(methods)
		if nil != err {
			t.Fatalf("%s\n", err)
		}
		t.Logf("response delivered: in %v", rep.DeliveredIn())
		if 0 != rep.ActionResultAt(0) {
			t.Fatalf("actionResult: %d\n", rep.ActionResultAt(0))
		}

		controlState = data.GetEnum()
		t.Logf("control state: %d", controlState)
		if 1 != controlState {
			t.Fatalf("meter did not connect, controlState: %d", controlState)
		}

	} else if 0 == controlState { // disconncted

		// Call remote_connect method.

		method := new(DlmsRequest)
		method.ClassId = classId
		method.InstanceId = instanceId
		method.MethodId = methodIdRemoteConnect
		methodParameters := new(DlmsData)
		methodParameters.SetInteger(1)
		method.MethodParameters = methodParameters
		methods := make([]*DlmsRequest, 1)
		methods[0] = method
		rep, err = aconn.SendRequest(methods)
		if nil != err {
			t.Fatalf("%s\n", err)
		}
		t.Logf("response delivered: in %v", rep.DeliveredIn())
		if 0 != rep.ActionResultAt(0) {
			t.Fatalf("actionResult: %d\n", rep.ActionResultAt(0))
		}

		// Check connected state.

		val = new(DlmsRequest)
		val.ClassId = classId
		val.InstanceId = instanceId
		val.AttributeId = attributeIdControlState
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
		if DATA_TYPE_ENUM != data.GetType() {
			t.Fatalf("not integer")
		}

		controlState := data.GetEnum()
		t.Logf("control state: %d", controlState)
		if 1 != controlState {
			t.Fatalf("meter did not connect, controlState: %d", controlState)
		}

	} else if 3 == controlState { // ready for connection

		// Call remote_connect method.

		method := new(DlmsRequest)
		method.ClassId = classId
		method.InstanceId = instanceId
		method.MethodId = methodIdRemoteConnect
		methodParameters := new(DlmsData)
		methodParameters.SetInteger(1)
		method.MethodParameters = methodParameters
		methods := make([]*DlmsRequest, 1)
		methods[0] = method
		rep, err = aconn.SendRequest(methods)
		if nil != err {
			t.Fatalf("%s\n", err)
		}
		t.Logf("response delivered: in %v", rep.DeliveredIn())
		if 0 != rep.ActionResultAt(0) {
			t.Fatalf("actionResult: %d\n", rep.ActionResultAt(0))
		}

		// Check connected state.

		val = new(DlmsRequest)
		val.ClassId = classId
		val.InstanceId = instanceId
		val.AttributeId = attributeIdControlState
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
		if DATA_TYPE_ENUM != data.GetType() {
			t.Fatalf("not integer")
		}

		controlState := data.GetEnum()
		t.Logf("control state: %d", controlState)
		if 1 != controlState {
			t.Fatalf("meter did not connect, controlState: %d", controlState)
		}

	} else {
		t.Fatalf("unknown controlState: %d", controlState)
	}

}