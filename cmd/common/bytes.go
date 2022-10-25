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

import "fmt"

type ByteSize float64

func (b ByteSize) String() string {
	switch {
	case b >= YiB:
		return fmt.Sprintf("%.2fYiB", b/YiB)
	case b >= ZiB:
		return fmt.Sprintf("%.2fZiB", b/ZiB)
	case b >= EiB:
		return fmt.Sprintf("%.2fEiB", b/EiB)
	case b >= PiB:
		return fmt.Sprintf("%.2fPiB", b/PiB)
	case b >= TiB:
		return fmt.Sprintf("%.2fTiB", b/TiB)
	case b >= GiB:
		return fmt.Sprintf("%.2fGiB", b/GiB)
	case b >= MiB:
		return fmt.Sprintf("%.2fMiB", b/MiB)
	case b >= KiB:
		return fmt.Sprintf("%.2fKiB", b/KiB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func ToByte(i any) string {
	switch n := i.(type) {
	case int:
		return ByteSize(n).String()
	case int8:
		return ByteSize(n).String()
	case int16:
		return ByteSize(n).String()
	case int32:
		return ByteSize(n).String()
	case int64:
		return ByteSize(n).String()
	case uint:
		return ByteSize(n).String()
	case uint8:
		return ByteSize(n).String()
	case uint16:
		return ByteSize(n).String()
	case uint32:
		return ByteSize(n).String()
	case uint64:
		return ByteSize(n).String()
	case float32:
		return ByteSize(n).String()
	case float64:
		return ByteSize(n).String()
	default:
		return ""
	}
}
