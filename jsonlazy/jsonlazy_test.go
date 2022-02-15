package jsonlazy

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

type mockMarshaler struct {
	ctrl *gomock.Controller
}

func newMockMarshaler(ctrl *gomock.Controller) *mockMarshaler {
	return &mockMarshaler{ctrl: ctrl}
}

func (m *mockMarshaler) MarshalJSON() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarshalJSON")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *mockMarshaler) expectMarshalJSON() *gomock.Call {
	m.ctrl.T.Helper()
	return m.ctrl.RecordCallWithMethodType(m, "MarshalJSON", reflect.TypeOf((*mockMarshaler)(nil).MarshalJSON))
}

func TestMarshaler(t *testing.T) {
	t.Run("RawMessage", func(t *testing.T) {
		expect := []byte(`[1,2,3]`)
		m := NewMarshaler(json.RawMessage(expect))
		got, err := m.MarshalJSON()
		assert.NilError(t, err)
		assert.Equal(t, string(expect), string(got))
	})

	t.Run("Once", func(t *testing.T) {
		expect := []byte(`[4,2]`)

		ctrl := gomock.NewController(t)
		mock := newMockMarshaler(ctrl)
		mock.expectMarshalJSON().Return(expect, nil)

		m := Once(NewMarshaler(mock))

		got1, err := m.MarshalJSON()
		assert.NilError(t, err)
		assert.Equal(t, string(expect), string(got1))

		got2, err := m.MarshalJSON()
		assert.NilError(t, err)
		assert.Equal(t, string(expect), string(got2))
	})
}
