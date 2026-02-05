package httprange

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/rogpeppe/ioseq"

	"go.pact.im/x/httpclient/mockhttpclient"
)

type multipartPart struct {
	Header textproto.MIMEHeader
	Reader io.Reader
}

func newMultipartReader(boundary string, parts []multipartPart) io.ReadCloser {
	return ioseq.ReaderWithContent(func(w io.Writer) error {
		writer := multipart.NewWriter(w)
		if err := writer.SetBoundary(boundary); err != nil {
			return err
		}
		for _, p := range parts {
			part, err := writer.CreatePart(p.Header)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, p.Reader); err != nil {
				return err
			}
		}
		return writer.Close()
	})
}

func TestHTTPRanger(t *testing.T) {
	type rangeData struct {
		Resp string
		Data []byte
	}
	tests := []struct {
		name    string
		spec    Specifier
		ranger  HTTPRanger
		expect  func(*mockhttpclient.MockClientMockRecorder)
		want    []rangeData
		wantErr bool
	}{
		{
			name: "successful single range request",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentRange: {"bytes 0-99/100"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: io.NopCloser(bytes.NewReader(bytes.Repeat([]byte("a"), 100))),
				}, nil)
			},
			want: []rangeData{
				{
					Resp: "0-99/100",
					Data: bytes.Repeat([]byte("a"), 100),
				},
			},
			wantErr: false,
		},
		{
			name: "request builder error",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return nil, errors.New("oops")
				}),
			},
			expect:  func(*mockhttpclient.MockClientMockRecorder) {},
			wantErr: true,
		},
		{
			name: "client error",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-99"},
					},
				}).Return(nil, errors.New("oops"))
			},
			wantErr: true,
		},
		{
			name: "unsatisfied range with content-range",
			spec: "bytes=1000-1999",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=1000-1999"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusRequestedRangeNotSatisfiable,
					Header: http.Header{
						httpHeaderContentRange: {"bytes */100"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unsatisfied range without content-range",
			spec: "bytes=1000-1999",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=1000-1999"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusRequestedRangeNotSatisfiable,
					Header: http.Header{
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unsatisfied range with invalid content-range",
			spec: "bytes=1000-1999",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=1000-1999"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusRequestedRangeNotSatisfiable,
					Header: http.Header{
						httpHeaderContentRange: {"invalid"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unexpected status code",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Header: http.Header{
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "ranges implicitly not supported",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "ranges explicitly not supported",
			spec: "bytes=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						httpHeaderAcceptRanges: {"none"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unaccepted unit",
			spec: "custom=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"custom=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unexpected unit in content-range",
			spec: "custom=0-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"custom=0-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentRange: {"bytes 0-99/100"},
					},
					Body: io.NopCloser(bytes.NewReader(bytes.Repeat([]byte("a"), 100))),
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "successful multipart range request",
			spec: "bytes=0-49,100-149",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-49,100-149"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType:  {"multipart/byteranges; boundary=test"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: newMultipartReader("test", []multipartPart{
						{
							Header: textproto.MIMEHeader{
								httpHeaderContentRange: {"bytes 0-49/200"},
							},
							Reader: bytes.NewReader(bytes.Repeat([]byte("a"), 50)),
						},
						{
							Header: textproto.MIMEHeader{
								httpHeaderContentRange: {"bytes 100-149/200"},
							},
							Reader: bytes.NewReader(bytes.Repeat([]byte("b"), 50)),
						},
					}),
				}, nil)
			},
			want: []rangeData{
				{
					Resp: "0-49/200",
					Data: bytes.Repeat([]byte("a"), 50),
				},
				{
					Resp: "100-149/200",
					Data: bytes.Repeat([]byte("b"), 50),
				},
			},
			wantErr: false,
		},
		{
			name: "multipart coalesced ranges",
			spec: "bytes=0-49,90-99",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-49,90-99"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType:  {"multipart/byteranges; boundary=test"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: newMultipartReader("test", []multipartPart{
						{
							Header: textproto.MIMEHeader{
								httpHeaderContentRange: {"bytes 0-99/100"},
							},
							Reader: bytes.NewReader(bytes.Repeat([]byte("a"), 100)),
						},
					}),
				}, nil)
			},
			want: []rangeData{
				{
					Resp: "0-99/100",
					Data: bytes.Repeat([]byte("a"), 100),
				},
			},
			wantErr: false,
		},
		{
			name: "multipart with empty boundary",
			spec: "bytes=0-49,100-149",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-49,100-149"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType:  {"multipart/byteranges"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "multipart with invalid parameters",
			spec: "bytes=0-49,100-149",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-49,100-149"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType:  {"multipart/byteranges; invalid"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "multipart without ranges",
			spec: "bytes=0-49,100-149",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bytes=0-49,100-149"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType:  {"multipart/byteranges; boundary=test"},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: newMultipartReader("test", nil),
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "multipart with unexpected unit in content-range",
			spec: "bits=0-399",
			ranger: HTTPRanger{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderRange: {"bits=0-399"},
					},
				}).Return(&http.Response{
					StatusCode: http.StatusPartialContent,
					Header: http.Header{
						httpHeaderContentType: {"multipart/byteranges; boundary=test"},
					},
					Body: newMultipartReader("test", []multipartPart{
						{
							Header: textproto.MIMEHeader{
								httpHeaderContentRange: {"bytes 0-49/50"},
							},
							Reader: bytes.NewReader(bytes.Repeat([]byte("a"), 50)),
						},
					}),
				}, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mockhttpclient.NewMockClient(ctrl)

			tt.expect(mockClient.EXPECT())

			ranger := tt.ranger
			ranger.Client = mockClient

			seq := ranger.Range(context.Background(), tt.spec)

			var ranges []rangeData
			var err error
			for r, e := range seq {
				if err != nil {
					t.Fatal("expected no ranges after error")
				}
				if e != nil {
					err = e
					continue
				}
				data, err := io.ReadAll(r.Reader)
				if err != nil {
					t.Fatal(err)
				}
				ranges = append(ranges, rangeData{
					Resp: r.Resp,
					Data: data,
				})
			}
			if len(ranges) != len(tt.want) {
				t.Errorf("got %d ranges, want %d", len(ranges), len(tt.want))
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

			for i := range min(len(ranges), len(tt.want)) {
				got, want := ranges[i], tt.want[i]
				if got.Resp != want.Resp {
					t.Errorf("range %d: got Resp %q, want %q", i, got.Resp, want.Resp)
				}
				if !bytes.Equal(got.Data, want.Data) {
					t.Errorf("range %d: got Data %q, want %q", i, got.Data, want.Data)
				}
			}
		})
	}
}

func TestHTTPRangerStopMultipart(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockhttpclient.NewMockClient(ctrl)

	mockClient.EXPECT().Do(&http.Request{
		Header: http.Header{
			httpHeaderRange: {"bytes=0-9,10-19"},
		},
	}).Return(&http.Response{
		StatusCode: http.StatusPartialContent,
		Header: http.Header{
			httpHeaderContentType: {"multipart/byteranges; boundary=test1"},
		},
		Body: newMultipartReader("test1", []multipartPart{
			{
				Header: textproto.MIMEHeader{
					httpHeaderContentRange: {"bytes 0-9/20"},
				},
				Reader: bytes.NewReader(bytes.Repeat([]byte("a"), 10)),
			},
			{
				Header: textproto.MIMEHeader{
					httpHeaderContentRange: {"bytes 10-19/20"},
				},
				Reader: bytes.NewReader(bytes.Repeat([]byte("b"), 10)),
			},
		}),
	}, nil)

	ranger := HTTPRanger{
		Client: mockClient,
		Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
			return &http.Request{}, nil
		}),
	}

	seq := ranger.Range(context.Background(), "bytes=0-9,10-19")

	next, stop := iter.Pull2(seq)
	_, _, _ = next()
	stop()
}

func TestEqualFoldASCII(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{
			name: "identical lowercase",
			a:    "bytes",
			b:    "bytes",
			want: true,
		},
		{
			name: "identical uppercase",
			a:    "BYTES",
			b:    "BYTES",
			want: true,
		},
		{
			name: "case-insensitive match",
			a:    "ByTeS",
			b:    "bYtEs",
			want: true,
		},
		{
			name: "different strings same length",
			a:    "bytes",
			b:    "other",
			want: false,
		},
		{
			name: "different lengths",
			a:    "bytes",
			b:    "byte",
			want: false,
		},
		{
			name: "empty strings",
			a:    "",
			b:    "",
			want: true,
		},
		{
			name: "non-ASCII characters",
			a:    "байты",
			b:    "БАЙТЫ",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := equalFoldASCII(tt.a, tt.b)
			if got != tt.want {
				t.Errorf(
					"equalFoldASCII(%q, %q) = %t, want %t",
					tt.a, tt.b, got, tt.want,
				)
			}
		})
	}
}
