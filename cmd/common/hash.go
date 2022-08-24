/*
Copyright © 2022 ZeYanLin <zeyanlin@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
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
