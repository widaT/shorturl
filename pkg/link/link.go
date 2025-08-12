package link

import (
	"context"
	"errors"
	"shorturl/pkg/radix"
	"shorturl/pkg/shint"
)

var (
	charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	ranges  = [][]uint64{
		{6, 56800235584},
		{7, 3521614606208},
		{8, 218340105584896},
		{9, 13537086546263552},
		{10, 839299365868340224},
	}

	keyPrefix = "link:"
)

type Link struct {
	idg *shint.Shint
	rdx *radix.Radix
}

func New(idg *shint.Shint) *Link {
	return &Link{
		idg: idg,
		rdx: radix.New(radix.Charset(charset)),
	}
}

func (l *Link) NewID(ctx context.Context) (string, error) {
	i, err := l.idg.Incr(ctx)
	if err != nil {
		return "", err
	}
	for _, r := range ranges {
		if i < r[1] {
			return l.rdx.Itoa(int(i), int(r[0])), nil
		}
	}
	return "", errors.New("link: out of range")
}

func (l *Link) AddLink(ctx context.Context, meta string, scene string) (string, string, error) {
	if scene == "" {
		scene = "default"
	}

	//TODO: check scene

	id, err := l.NewID(ctx)
	if err != nil {
		return "", "", err
	}

	//set key to redis
	// key := keyPrefix + id

	// rdb.SetNX(ctx, key, meta, ttl)

	return id, scene, nil

}
