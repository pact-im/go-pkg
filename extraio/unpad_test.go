package extraio

import (
	"bytes"
	"io"
	"strconv"
	"testing"
)

func TestUnpadReader(t *testing.T) {
	testCases := []struct {
		Data []byte
		Size uint8
	}{
		{
			Data: nil,
			Size: 0,
		},
		{
			Data: []byte{0},
			Size: 1,
		},
		{
			Data: []byte{1, 1, 1},
			Size: 1,
		},
		{
			Data: []byte{1, 2, 3},
			Size: 1,
		},
		{
			Data: []byte{3, 2, 1},
			Size: 1,
		},
		{
			Data: []byte{1, 2, 3},
			Size: 3,
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			runTestUnpadReader(t, tc.Data, tc.Size)
		})
	}
}

func runTestUnpadReader(t testing.TB, data []byte, blockSize uint8) {
	var success bool
	var expect []byte
	if n := len(data); n > 0 && blockSize > 0 && len(data)%int(blockSize) == 0 {
		fillByte := data[len(data)-1]
		if fillByte > 0 && fillByte <= blockSize {
			expect, success = unpadPayload(data, fillByte)
		}
	} else if blockSize == 0 {
		success = true
		expect = data
	}
	if !success {
		n := len(data)
		bs := int(blockSize)
		if n > bs {
			n -= bs
		} else {
			n = 0
		}
		expect = data[:n]
	}

	unpadded, err := io.ReadAll(NewUnpadReader(bytes.NewReader(data), blockSize))
	if success != (err == nil) {
		t.Errorf("expected success=%v, got err=%v", success, err)
	}
	if !bytes.Equal(expect, unpadded) {
		t.Errorf("expected %v, got %v", expect, unpadded)
	}
}

func TestUnpadPayload(t *testing.T) {
	testCases := []struct {
		Data    []byte
		Fill    byte
		Payload []byte
		Fail    bool
	}{
		{
			Data: []byte{0},
			Fill: 1,
			Fail: true,
		},
		{
			Data: []byte{7, 7, 7},
			Fill: 7,
			Fail: true,
		},
		{
			Data:    bytes.Repeat([]byte{7}, 7),
			Fill:    7,
			Payload: []byte{},
		},
		{
			Data: []byte{1, 2, 3},
			Fill: 1,
			Fail: true,
		},
		{
			Data:    []byte{3, 2, 1},
			Fill:    1,
			Payload: []byte{3, 2},
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload, ok := unpadPayload(tc.Data, tc.Fill)
			if ok != !tc.Fail {
				t.Errorf("expected ok=%v, got %v", !tc.Fail, ok)
			}
			if !bytes.Equal(tc.Payload, payload) {
				t.Errorf("expected %v, got %v", tc.Payload, payload)
			}
		})
	}
}
