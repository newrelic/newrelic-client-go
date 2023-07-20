//go:build unit
// +build unit

package http

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var testCompressionCases = []struct {
	compress bool
	data     []byte
	gzip     []byte
}{
	{
		compress: false,
		data:     []byte("test"),
	},
	{
		compress: false,
		data:     []byte(`abcdefghijklmnopqrstuvwxyz1234567890`),
	},
	{
		compress: true,
		data:     []byte(`{"data": "json", "maybe":"didn't check", "handcrafted":true, "example": { "sub": "object", "arry": [ "why", "not" ] }, "todo": [ "make", "it", "larger" }`),
		gzip: []byte{
			0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x24, 0x8c, 0x4d, 0xae, 0xc2, 0x30, 0xc, 0x6, 0xaf, 0xf2,
			0xe9, 0xdb, 0xbc, 0x4d, 0x4f, 0x90, 0xab, 0x3c, 0xb1, 0x70, 0x63, 0x43, 0x7f, 0x13, 0xe4, 0xba, 0x82, 0xaa, 0xea, 0xdd, 0x51, 0x60, 0x3b, 0xa3,
			0x99, 0x93, 0x2a, 0x21, 0x4c, 0xe0, 0xb4, 0xd5, 0xc2, 0xe, 0x5c, 0xe5, 0xe8, 0x8d, 0x89, 0x3a, 0x6a, 0xf9, 0xb, 0xe4, 0xc1, 0xf2, 0xdc, 0xf8,
			0x20, 0x45, 0xb3, 0xcb, 0x3d, 0x4c, 0x99, 0xc2, 0x77, 0xeb, 0x40, 0x7b, 0xcb, 0xfa, 0x5c, 0x8c, 0x9, 0x27, 0xb8, 0xed, 0x7d, 0xfb, 0xd4, 0x7e,
			0xb2, 0x1c, 0xad, 0x10, 0xf7, 0x83, 0x9, 0xff, 0xe0, 0x6b, 0x38, 0x1a, 0x28, 0x35, 0x88, 0x1b, 0xae, 0xe, 0x8c, 0xaa, 0xf5, 0xe7, 0x56, 0x99,
			0xad, 0xc9, 0xf1, 0xdb, 0x2c, 0xe2, 0xf, 0x73, 0xe2, 0xfa, 0x4, 0x0, 0x0, 0xff, 0xff, 0x60, 0x64, 0xbb, 0x7c, 0x99, 0x0, 0x0, 0x0,
		},
	},
}

func TestNoneCompressor(t *testing.T) {
	t.Parallel()

	tc := mock.NewTestConfig(t, nil)
	c := NewClient(tc)

	req, err := c.NewRequest("POST", c.config.Region().RestURL("path"), nil, nil, nil)
	assert.NoError(t, err)

	compress := NoneCompressor{}
	var bodyReader io.Reader

	for _, d := range testCompressionCases {
		req.SetHeader("content-encoding", "<invalid>")

		bodyReader, err = compress.Compress(req, d.data)
		assert.NoError(t, err)
		assert.Equal(t, "", req.GetHeader("content-encoding"))

		res, err := io.ReadAll(bodyReader)
		assert.NoError(t, err)
		assert.Equal(t, d.data, res)
	}
}

func TestGzipCompressor(t *testing.T) {
	t.Parallel()

	tc := mock.NewTestConfig(t, nil)
	c := NewClient(tc)

	req, err := c.NewRequest("POST", c.config.Region().RestURL("path"), nil, nil, nil)
	assert.NoError(t, err)

	compress := GzipCompressor{}
	var bodyReader io.Reader

	for _, d := range testCompressionCases {
		req.SetHeader("content-encoding", "<invalid>")

		bodyReader, err = compress.Compress(req, d.data)
		assert.NoError(t, err)

		res, err := io.ReadAll(bodyReader)
		assert.NoError(t, err)

		if d.compress {
			assert.Equal(t, "gzip", req.GetHeader("content-encoding"))
			assert.Equal(t, d.gzip, res)
		} else {
			assert.Equal(t, "", req.GetHeader("content-encoding"))
			assert.Equal(t, d.data, res)
		}
	}
}
