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

var testTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func TestHTTPStrongValidatorBuilder(t *testing.T) {
	tests := []struct {
		name     string
		builder  HTTPStrongValidatorBuilder
		metadata HTTPMetadata
		want     HTTPStrongValidator
		wantErr  bool
	}{
		{
			name:     "no applicable validators",
			builder:  HTTPStrongValidatorBuilder{},
			metadata: HTTPMetadata{},
			want:     HTTPStrongValidator{},
			wantErr:  true,
		},
		{
			name:    "weak ETag",
			builder: HTTPStrongValidatorBuilder{},
			metadata: HTTPMetadata{
				ETag: `W/"weak-etag"`,
			},
			want:    HTTPStrongValidator{},
			wantErr: true,
		},
		{
			name:    "strong ETag",
			builder: HTTPStrongValidatorBuilder{},
			metadata: HTTPMetadata{
				ETag: `"strong-etag"`,
			},
			want: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfMatch: {`"strong-etag"`},
				},
				ETag: `"strong-etag"`,
			},
			wantErr: false,
		},
		{
			name:    "disabled Last-Modified",
			builder: HTTPStrongValidatorBuilder{},
			metadata: HTTPMetadata{
				LastModified: testTime,
				Date:         testTime.Add(time.Second),
			},
			want:    HTTPStrongValidator{},
			wantErr: true,
		},
		{
			name: "strong Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				LastModified: testTime,
				Date:         testTime.Add(time.Second),
			},
			want: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfUnmodifiedSince: {
						testTime.UTC().Format(http.TimeFormat),
					},
				},
				LastModified: testTime,
			},
			wantErr: false,
		},
		{
			name: "weak Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				LastModified: testTime,
				Date:         testTime,
			},
			want:    HTTPStrongValidator{},
			wantErr: true,
		},
		{
			name: "unset Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				Date: testTime,
			},
			want:    HTTPStrongValidator{},
			wantErr: true,
		},
		{
			name: "unset Date",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				LastModified: testTime,
			},
			want:    HTTPStrongValidator{},
			wantErr: true,
		},
		{
			name: "both ETag and Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				ETag:         `"strong-etag"`,
				LastModified: testTime,
				Date:         testTime.Add(time.Second),
			},
			want: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfMatch: {`"strong-etag"`},
				},
				ETag: `"strong-etag"`,
			},
			wantErr: false,
		},
		{
			name: "weak ETag and strong Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				ETag:         `W/"weak-etag"`,
				LastModified: testTime,
				Date:         testTime.Add(time.Second),
			},
			want: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfUnmodifiedSince: {
						testTime.UTC().Format(http.TimeFormat),
					},
				},
				LastModified: testTime,
			},
			wantErr: false,
		},
		{
			name: "strong ETag and weak Last-Modified",
			builder: HTTPStrongValidatorBuilder{
				UseLastModified: true,
			},
			metadata: HTTPMetadata{
				ETag:         `"strong-etag"`,
				LastModified: testTime,
				Date:         testTime,
			},
			want: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfMatch: {`"strong-etag"`},
				},
				ETag: `"strong-etag"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.builder.Build(context.Background(), &tt.metadata, nil)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(client, &tt.want) {
				t.Fatalf("Build() mismatch:\n got: %#v\nwant: %#v", client, &tt.want)
			}
		})
	}
}

func TestHTTPStrongValidator(t *testing.T) {
	tests := []struct {
		name      string
		validator HTTPStrongValidator
		request   http.Request
		expect    func(*mockhttpclient.MockClientMockRecorder)
		wantErr   bool
	}{
		{
			name:      "client error",
			validator: HTTPStrongValidator{},
			request:   http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(nil, errors.New("oops"))
			},
			wantErr: true,
		},
		{
			name:      "without validators",
			validator: HTTPStrongValidator{},
			request:   http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Body: http.NoBody,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "precondition without validators",
			validator: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfMatch: {`"strong-etag"`},
				},
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderIfMatch: {`"strong-etag"`},
					},
				}).Return(&http.Response{
					Body: http.NoBody,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "precondition overrides existing headers",
			validator: HTTPStrongValidator{
				Precondition: http.Header{
					httpHeaderIfMatch: {`"strong-etag"`},
					httpHeaderIfUnmodifiedSince: {
						testTime.UTC().Format(http.TimeFormat),
					},
				},
			},
			request: http.Request{
				Header: http.Header{
					httpHeaderIfMatch: {`"different-etag"`},
					httpHeaderIfUnmodifiedSince: {
						testTime.Add(-time.Second).UTC().Format(http.TimeFormat),
					},
				},
			},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{
					Header: http.Header{
						httpHeaderIfMatch: {`"strong-etag"`},
						httpHeaderIfUnmodifiedSince: {
							testTime.UTC().Format(http.TimeFormat),
						},
					},
				}).Return(&http.Response{
					Body: http.NoBody,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "ETag",
			validator: HTTPStrongValidator{
				ETag: `"strong-etag"`,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Header: http.Header{
						httpHeaderETag: {`"strong-etag"`},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "changed ETag",
			validator: HTTPStrongValidator{
				ETag: `"strong-etag"`,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Header: http.Header{
						httpHeaderETag: {`"changed-etag"`},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid ETag",
			validator: HTTPStrongValidator{
				ETag: `"strong-etag"`,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Header: http.Header{
						httpHeaderETag: {`not a valid ETag`},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unset ETag",
			validator: HTTPStrongValidator{
				ETag: `"strong-etag"`,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "Last-Modified",
			validator: HTTPStrongValidator{
				LastModified: testTime,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Header: http.Header{
						httpHeaderLastModified: {
							testTime.UTC().Format(http.TimeFormat),
						},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "changed Last-Modified",
			validator: HTTPStrongValidator{
				LastModified: testTime,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Header: http.Header{
						httpHeaderLastModified: {
							testTime.Add(time.Second).UTC().Format(http.TimeFormat),
						},
					},
					Body: http.NoBody,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "unset Last-Modified",
			validator: HTTPStrongValidator{
				LastModified: testTime,
			},
			request: http.Request{},
			expect: func(r *mockhttpclient.MockClientMockRecorder) {
				r.Do(&http.Request{}).Return(&http.Response{
					Body: http.NoBody,
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

			validator := tt.validator
			validator.Client = mockClient

			_, err := validator.Do(&tt.request)
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
