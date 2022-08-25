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
	"context"
	"errors"
	"time"
)

var (
	ErrConfigContent = errors.New("config content is incorrect")
	ErrConfigTable   = errors.New("table not found in the config")
	ErrIllegalPath   = errors.New("illegal file path")
	ErrInvalidArg    = errors.New("invalid argument")
	ErrInvalidIP     = errors.New("invalid IP")
	ErrInvalidFile   = errors.New("invalid file format")
	ErrInvalidURL    = errors.New("invalid URL")
	ErrResponse      = errors.New("response error")
	ErrStatusCode    = errors.New("status code is not 200")
)

var (
	Context = context.Background()
	TimeNow = time.Now().Local()
)
