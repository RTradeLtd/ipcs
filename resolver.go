package ipcs

import (
	"context"
	"io"
	"io/ioutil"

	xpb "github.com/RTradeLtd/TxPB/v3/go"

	"github.com/RTradeLtd/ipcs/digestconv"
	"github.com/ipfs/interface-go-ipfs-core/path"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

func (s *store) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {
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
	return ioutil.NopCloser(data), nil
}
