package ipcs

import (
	"context"
	"io"
	"io/ioutil"

	xpb "github.com/RTradeLtd/TxPB/v3/go"
	"github.com/RTradeLtd/ipcs/digestconv"
	"github.com/containerd/containerd/content"
	"github.com/ipfs/interface-go-ipfs-core/path"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// ReaderAt only requires desc.Digest to be set.
// Other fields in the descriptor may be used internally for resolving
// the location of the actual data.
func (s *store) ReaderAt(ctx context.Context, desc ocispec.Descriptor) (content.ReaderAt, error) {
	c, err := digestconv.DigestToCid(desc.Digest)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert digest '%s' to cid", desc.Digest)
	}
	data, err := s.cln.DownloadFile(ctx, &xpb.DownloadRequest{
		Hash: path.IpfsPath(c).String(),
	}, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get unixfs node %q", c)
	}

	return &sizeReaderAt{
		size:   desc.Size,
		reader: data,
	}, nil
}

type sizeReaderAt struct {
	size   int64
	reader io.Reader
	n      int64
}

func (ra *sizeReaderAt) ReadAt(p []byte, offset int64) (n int, err error) {
	if offset < ra.n {
		return 0, errors.New("invalid offset")
	}
	diff := offset - ra.n
	written, err := io.CopyN(ioutil.Discard, ra.reader, diff)
	ra.n += written
	if err != nil {
		return int(written), err
	}

	n, err = ra.reader.Read(p)
	ra.n += int64(n)
	return
}

func (ra *sizeReaderAt) Size() int64 {
	return ra.size
}

func (ra *sizeReaderAt) Close() error {
	return nil
}
