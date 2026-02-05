package httprange

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"slices"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestBytesSpecifier(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []int64
		expected Specifier
	}{
		{
			name:     "Empty positions",
			inputs:   []int64{},
			expected: "bytes=-0",
		},
		{
			name:     "Single suffix range",
			inputs:   []int64{-50},
			expected: "bytes=-50",
		},
		{
			name:     "Single open-ended range",
			inputs:   []int64{0},
			expected: "bytes=0-",
		},
		{
			name:     "Single closed range",
			inputs:   []int64{0, 99},
			expected: "bytes=0-99",
		},
		{
			name:     "Multiple closed ranges",
			inputs:   []int64{0, 99, 100, 199},
			expected: "bytes=0-99,100-199",
		},
		{
			name:     "Open then suffix",
			inputs:   []int64{0, -50},
			expected: "bytes=0-,-50",
		},
		{
			name:     "Suffix then open",
			inputs:   []int64{-50, 0},
			expected: "bytes=-50,0-",
		},
		{
			name:     "Closed then open",
			inputs:   []int64{0, 99, 200},
			expected: "bytes=0-99,200-",
		},
		{
			name:     "Single-element range",
			inputs:   []int64{0, 0},
			expected: "bytes=0-0",
		},
		{
			name:     "Multiple suffix ranges",
			inputs:   []int64{-10, -20},
			expected: "bytes=-10,-20",
		},
		{
			name:     "Descending positions",
			inputs:   []int64{100, 50},
			expected: "bytes=100-,50-",
		},
		{
			name:     "Triple with middle negative",
			inputs:   []int64{0, -50, 100},
			expected: "bytes=0-,-50,100-",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesSpecifier(tt.inputs...)
			if got != tt.expected {
				t.Errorf("BytesSpecifier() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestBytesReader(t *testing.T) {
	tests := []struct {
		name    string
		expect  func(*MockRangerMockRecorder)
		offset  int64
		buf     []byte
		want    []byte
		eof     bool
		wantErr bool
	}{
		{
			name: "successful read",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "0-9/100",
								Reader: bytes.NewReader([]byte("0123456789")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte("0123456789"),
			eof:     false,
			wantErr: false,
		},
		{
			name: "read at EOF returns EOF",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=90-99")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "90-99/100",
								Reader: bytes.NewReader([]byte("0123456789")),
							},
						},
					})))
			},
			offset:  90,
			buf:     make([]byte, 10),
			want:    []byte("0123456789"),
			eof:     true,
			wantErr: false,
		},
		{
			name:    "empty buffer returns zero",
			expect:  func(*MockRangerMockRecorder) {},
			offset:  0,
			buf:     []byte{},
			want:    []byte{},
			eof:     false,
			wantErr: false,
		},
		{
			name:    "negative offset returns error",
			expect:  func(*MockRangerMockRecorder) {},
			offset:  -1,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name:    "offset overflow returns error",
			expect:  func(*MockRangerMockRecorder) {},
			offset:  (1<<63 - 1) - 4, // MaxInt64 - 4
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "unsatisfied range with unknown length returns EOF",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=100-109")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							Second: &UnsatisfiedRangeError{},
						},
					})))
			},
			offset:  100,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     true,
			wantErr: false,
		},
		{
			name: "unsatisfied range with known length at EOF returns EOF",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=100-109")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							Second: &UnsatisfiedRangeError{Resp: "*/100"},
						},
					})))
			},
			offset:  100,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     true,
			wantErr: false,
		},
		{
			name: "unsatisfied range with known length before EOF returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=90-99")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							Second: &UnsatisfiedRangeError{Resp: "*/100"},
						},
					})))
			},
			offset:  90,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "invalid unsatisfied range returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=100-109")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							Second: &UnsatisfiedRangeError{Resp: "invalid"},
						},
					})))
			},
			offset:  100,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "range error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							Second: errors.New("oops"),
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "invalid bytes range returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "invalid",
								Reader: bytes.NewReader([]byte("data")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "unexpected first byte position returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "5-14/100",
								Reader: bytes.NewReader([]byte("56789")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "unexpected last byte position returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "0-5/100",
								Reader: bytes.NewReader([]byte("012345")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "short read from range returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "0-9/100",
								Reader: bytes.NewReader([]byte("short")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "unexpected multiple ranges returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(toSeq2(slices.Values([]seqPair[Range, error]{
						{
							First: Range{
								Resp:   "0-9/100",
								Reader: bytes.NewReader([]byte("0123456789")),
							},
						},
						{
							First: Range{
								Resp:   "10-19/100",
								Reader: bytes.NewReader([]byte("0123456789")),
							},
						},
					})))
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
		{
			name: "empty sequence returns error",
			expect: func(r *MockRangerMockRecorder) {
				r.Range(gomock.Any(), Specifier("bytes=0-9")).
					Return(func(func(Range, error) bool) {})
			},
			offset:  0,
			buf:     make([]byte, 10),
			want:    []byte{},
			eof:     false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRanger := NewMockRanger(ctrl)

			tt.expect(mockRanger.EXPECT())

			reader := BytesReader{
				Context: context.Background(),
				Ranger:  mockRanger,
			}

			n, err := reader.ReadAt(tt.buf, tt.offset)
			if got := tt.buf[:n]; !bytes.Equal(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if err == io.EOF {
				if tt.eof {
					return
				}
				t.Error("unexpected EOF")
			}
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type seqPair[T, U any] struct {
	First  T
	Second U
}

func toSeq2[T, U any](seq iter.Seq[seqPair[T, U]]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		for p := range seq {
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}
