package link

import (
	"context"
	"errors"
	"shorturl/config"
	"shorturl/radix"
	"shorturl/shint"

	"github.com/redis/go-redis/v9"
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
	idg        *shint.Shint
	rdx        *radix.Radix
	redisCient *redis.Client
}

func New(idg *shint.Shint, redisClient *redis.Client) *Link {
	return &Link{
		idg:        idg,
		rdx:        radix.New(radix.Charset(charset)),
		redisCient: redisClient,
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
	conf, ok := config.Get().Scenes[scene]
	if !ok {
		return "", "", errors.New("link: scene not found")
	}
	id, err := l.NewID(ctx)
	if err != nil {
		return "", "", err
	}
	key := keyPrefix + id
	err = l.redisCient.SetNX(ctx, key, meta, conf.TTL).Err()
	if err != nil {
		return "", "", err
	}
	return id, scene, nil
}

func (l *Link) GetLink(ctx context.Context, scene string, id string) (string, error) {
	if scene == "" {
		scene = "default"
	}
	conf, ok := config.Get().Scenes[scene]
	if !ok {
		return "", errors.New("link: scene not found")
	}
	key := keyPrefix + id
	l.redisCient.Get(ctx, key).Result()
}
