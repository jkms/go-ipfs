package nilrouting

import (
	"context"
	"errors"

	repo "github.com/ipfs/go-ipfs/repo"

	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	p2phost "gx/ipfs/QmcQJspoQP914EDCTemcXz8yQkL4xvc49pTWFhzmifvUZY/go-libp2p-host"
	routing "gx/ipfs/QmeCCYJLNMQnEZyeJGeytEytx6J2hWXGSXHnJvzD5rACCz/go-libp2p-routing"
	cid "gx/ipfs/QmetUj3ZqWMDVeFMRq7S9PdMauXCwBZuggfHqoS4MPt1Vy/go-cid"
)

type nilclient struct {
}

func (c *nilclient) PutValue(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (c *nilclient) GetValue(_ context.Context, _ string) ([]byte, error) {
	return nil, errors.New("Tried GetValue from nil routing.")
}

func (c *nilclient) GetValues(_ context.Context, _ string, _ int) ([]routing.RecvdVal, error) {
	return nil, errors.New("Tried GetValues from nil routing.")
}

func (c *nilclient) FindPeer(_ context.Context, _ peer.ID) (pstore.PeerInfo, error) {
	return pstore.PeerInfo{}, nil
}

func (c *nilclient) FindProvidersAsync(_ context.Context, _ *cid.Cid, _ int) <-chan pstore.PeerInfo {
	out := make(chan pstore.PeerInfo)
	defer close(out)
	return out
}

func (c *nilclient) Provide(_ context.Context, _ *cid.Cid, _ bool) error {
	return nil
}

func (c *nilclient) Bootstrap(_ context.Context) error {
	return nil
}

func ConstructNilRouting(_ context.Context, _ p2phost.Host, _ repo.Datastore) (routing.IpfsRouting, error) {
	return &nilclient{}, nil
}

//  ensure nilclient satisfies interface
var _ routing.IpfsRouting = &nilclient{}
