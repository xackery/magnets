package library

import "io"

type assetReader interface {
	io.ReadSeeker
}

type assetReadCloser struct {
	io.ReadCloser
	io.ReadSeeker
	r assetReader
}

func newAssetReadCloser(r assetReader) *assetReadCloser {
	rc := &assetReadCloser{
		r: r,
	}
	return rc
}

func (rc *assetReadCloser) Close() error {
	rc.r = nil
	return nil
}

func (rc *assetReadCloser) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}
