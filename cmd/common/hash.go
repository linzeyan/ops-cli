package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	HashMd5        = "md5"
	HashSha1       = "sha1"
	HashSha224     = "sha224"
	HashSha256     = "sha256"
	HashSha384     = "sha384"
	HashSha512     = "sha512"
	HashSha512_224 = "sha512_224"
	HashSha512_256 = "sha512_256"
)

func HashAlgorithm(alg string) hash.Hash {
	m := map[string]hash.Hash{
		HashMd5:        md5.New(),
		HashSha1:       sha1.New(),
		HashSha224:     sha256.New224(),
		HashSha256:     sha256.New(),
		HashSha384:     sha512.New384(),
		HashSha512:     sha512.New(),
		HashSha512_224: sha512.New512_224(),
		HashSha512_256: sha512.New512_256(),
	}
	if h, ok := m[alg]; ok {
		return h
	}
	return nil
}
