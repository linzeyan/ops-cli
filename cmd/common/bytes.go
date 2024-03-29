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
	"fmt"
)

/* Convert value to byte size. */
type byteSize float64

func (b byteSize) Convert() string {
	switch {
	case b >= YiB:
		return fmt.Sprintf("%.2f%s", b/YiB, YiB)
	case b >= ZiB:
		return fmt.Sprintf("%.2f%s", b/ZiB, ZiB)
	case b >= EiB:
		return fmt.Sprintf("%.2f%s", b/EiB, EiB)
	case b >= PiB:
		return fmt.Sprintf("%.2f%s", b/PiB, PiB)
	case b >= TiB:
		return fmt.Sprintf("%.2f%s", b/TiB, TiB)
	case b >= GiB:
		return fmt.Sprintf("%.2f%s", b/GiB, GiB)
	case b >= MiB:
		return fmt.Sprintf("%.2f%s", b/MiB, MiB)
	case b >= KiB:
		return fmt.Sprintf("%.2f%s", b/KiB, KiB)
	case b < 0:
		return ""
	default:
		return fmt.Sprintf("%.2fB", b)
	}
}

func (b byteSize) String() string {
	switch b {
	case YiB:
		return "YiB"
	case ZiB:
		return "ZiB"
	case EiB:
		return "EiB"
	case PiB:
		return "PiB"
	case TiB:
		return "TiB"
	case GiB:
		return "GiB"
	case MiB:
		return "MiB"
	case KiB:
		return "KiB"
	default:
		return fmt.Sprintf("byteSize(%f)", b)
	}
}

/* ByteSize return byte size string from giving value. */
func ByteSize(i any) string {
	switch n := i.(type) {
	case int:
		return byteSize(n).Convert()
	case int8:
		return byteSize(n).Convert()
	case int16:
		return byteSize(n).Convert()
	case int32:
		return byteSize(n).Convert()
	case int64:
		return byteSize(n).Convert()
	case uint:
		return byteSize(n).Convert()
	case uint8:
		return byteSize(n).Convert()
	case uint16:
		return byteSize(n).Convert()
	case uint32:
		return byteSize(n).Convert()
	case uint64:
		return byteSize(n).Convert()
	case float32:
		return byteSize(n).Convert()
	case float64:
		return byteSize(n).Convert()
	default:
		stdLogger.Log.Debug(ErrInvalidArg.Error(), DefaultField(i))
		return ""
	}
}
