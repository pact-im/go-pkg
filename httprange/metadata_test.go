package httprange

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"go.pact.im/x/httpclient/mockhttpclient"
)

func TestHTTPResponseMetadataExtractor(t *testing.T) {
	tests := []struct {
		name      string
		extractor HTTPResponseMetadataExtractor
		expect    func(*mockhttpclient.MockClientMockRecorder)
		want      HTTPMetadata
		wantErr   bool
	}{
		{
			name: "successful extraction with all headers",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderETag:         {`"512aedef7775096f9e152526a30a0ce7"`},
						httpHeaderLastModified: {testTime.UTC().Format(http.TimeFormat)},
						httpHeaderDate:         {testTime.Add(time.Second).UTC().Format(http.TimeFormat)},
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: 1024,
				ETag:          `"512aedef7775096f9e152526a30a0ce7"`,
				LastModified:  testTime,
				Date:          testTime.Add(time.Second),
				AcceptRanges:  []string{"bytes"},
			},
			wantErr: false,
		},
		{
			name: "successful extraction with weak ETag",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 2048,
					Header: http.Header{
						httpHeaderETag:         {`W/"weak-etag"`},
						httpHeaderLastModified: {testTime.UTC().Format(http.TimeFormat)},
						httpHeaderDate:         {testTime.Add(time.Second).UTC().Format(http.TimeFormat)},
						httpHeaderAcceptRanges: {"bytes", "seconds"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: 2048,
				ETag:          `W/"weak-etag"`,
				LastModified:  testTime,
				Date:          testTime.Add(time.Second),
				AcceptRanges:  []string{"bytes", "seconds"},
			},
			wantErr: false,
		},
		{
			name: "successful extraction with empty headers",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 512,
					Header:        http.Header{},
					Body:          http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: 512,
				ETag:          "",
				LastModified:  time.Time{},
				Date:          time.Time{},
				AcceptRanges:  nil,
			},
			wantErr: false,
		},
		{
			name: "client error",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(nil, errors.New("oops"))
			},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "request builder error",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return nil, errors.New("oops")
				}),
			},
			expect:  func(*mockhttpclient.MockClientMockRecorder) {},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "non-2xx status code",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Status:     "404 Not Found",
					Header:     http.Header{},
					Body:       http.NoBody,
				}, nil)
			},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "invalid ETag format",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderETag: {"invalid-etag"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "invalid Last-Modified format",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderLastModified: {"not a valid time"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "invalid Date format",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderDate: {"not a valid time"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want:    HTTPMetadata{},
			wantErr: true,
		},
		{
			name: "multiple Accept-Ranges values",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderAcceptRanges: {"bytes", "seconds", "none"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: 1024,
				ETag:          "",
				LastModified:  time.Time{},
				Date:          time.Time{},
				AcceptRanges:  []string{"bytes", "seconds", "none"},
			},
			wantErr: false,
		},
		{
			name: "unknown content length",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: -1,
					Header: http.Header{
						httpHeaderAcceptRanges: {"bytes"},
					},
					Body: http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: -1,
				ETag:          "",
				LastModified:  time.Time{},
				Date:          time.Time{},
				AcceptRanges:  []string{"bytes"},
			},
			wantErr: false,
		},
		{
			name: "empty Accept-Ranges header",
			extractor: HTTPResponseMetadataExtractor{
				Request: HTTPRequestBuilderFunc(func(_ context.Context) (*http.Request, error) {
					return &http.Request{}, nil
				}),
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					StatusCode:    http.StatusOK,
					ContentLength: 1024,
					Header: http.Header{
						httpHeaderAcceptRanges: {""},
					},
					Body: http.NoBody,
				}, nil)
			},
			want: HTTPMetadata{
				ContentLength: 1024,
				ETag:          "",
				LastModified:  time.Time{},
				Date:          time.Time{},
				AcceptRanges:  []string{""},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := mockhttpclient.NewMockClient(ctrl)

			tt.expect(mockClient.EXPECT())

			extractor := tt.extractor
			extractor.Client = mockClient

			provider, err := extractor.Extract(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(provider, &tt.want) {
				t.Fatalf("Extract() mismatch:\n got: %#v\nwant: %#v", provider, &tt.want)
			}
		})
	}
}

func TestIsValidETag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid strong ETag simple",
			input: `"etag"`,
			want:  true,
		},
		{
			name:  "valid strong ETag with special characters",
			input: `"!etag-123_a.b/c~"`,
			want:  true,
		},
		{
			name:  "valid strong ETag empty",
			input: `""`,
			want:  true,
		},
		{
			name:  "valid strong ETag single character",
			input: `"a"`,
			want:  true,
		},
		{
			name:  "valid weak ETag simple",
			input: `W/"etag"`,
			want:  true,
		},
		{
			name:  "valid weak ETag with special characters",
			input: `W/"etag-123_a.b/c!"`,
			want:  true,
		},
		{
			name:  "valid weak ETag empty",
			input: `W/""`,
			want:  true,
		},
		{
			name:  "valid weak ETag single character",
			input: `W/"a"`,
			want:  true,
		},
		{
			name:  "missing opening quote",
			input: `etag"`,
			want:  false,
		},
		{
			name:  "missing closing quote",
			input: `"etag`,
			want:  false,
		},
		{
			name:  "no quotes",
			input: "etag",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "only quote",
			input: `"`,
			want:  false,
		},
		{
			name:  "space inside ETag",
			input: `"etag with space"`,
			want:  false,
		},
		{
			name:  "tab and newline inside ETag",
			input: "\"etag\twith\ntab\"",
			want:  false,
		},
		{
			name:  "delete character inside ETag",
			input: "\"etag\x7F\"",
			want:  false,
		},
		{
			name:  "quote inside ETag",
			input: `"etag"inside"`,
			want:  false,
		},
		{
			name:  "weak ETag lowercase w",
			input: `w/"etag"`,
			want:  false,
		},
		{
			name:  "weak ETag no slash",
			input: `W"etag"`,
			want:  false,
		},
		{
			name:  "weak ETag missing quotes",
			input: `W/etag`,
			want:  false,
		},
		{
			name:  "weak ETag space before W",
			input: ` W/"etag"`,
			want:  false,
		},
		{
			name:  "weak ETag space after W",
			input: `W /"etag"`,
			want:  false,
		},
		{
			name:  "weak ETag space after slash",
			input: `W/ "etag"`,
			want:  false,
		},
		{
			name:  "weak ETag extra text after ETag",
			input: `W/"etag"extra`,
			want:  false,
		},
		{
			name:  "unicode character",
			input: `"etag©"`,
			want:  true,
		},
		{
			name:  "non-ASCII byte",
			input: "\"etag\x80\"",
			want:  true,
		},
		{
			name:  "only weak prefix",
			input: `W/`,
			want:  false,
		},
		{
			name:  "weak prefix with empty quotes",
			input: `W/""`,
			want:  true,
		},
		{
			name:  "multiple weak prefixes",
			input: `W/W/"etag"`,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidETag(tt.input)
			if got != tt.want {
				t.Errorf("isValidETag(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
