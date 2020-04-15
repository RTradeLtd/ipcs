package ipcs

import (
	sdkc "github.com/RTradeLtd/go-temporalx-sdk/client"

	"github.com/containerd/containerd/content"
)

type Config struct {
	ipfsCln *sdkc.Client
}

type store struct {
	cln *sdkc.Client
}

func NewContentStore(cfg Config) (content.Store, error) {
	return &store{
		cln: cfg.ipfsCln,
	}, nil
}

func NewContentStoreFromCoreAPI(cln *sdkc.Client) content.Store {
	return &store{cln}
}
