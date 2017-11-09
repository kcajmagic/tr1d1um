/**
 * Copyright 2017 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package retryUtilities

import (
	"time"
	"github.com/go-kit/kit/log"
	"errors"
)

var (
	invalidMaxRetriesValErr = errors.New("retry: MaxRetries should be >= 1")
	undefinedShouldRetryErr = errors.New("retry: ShouldRetry function is required")
	undefinedLogger = errors.New("retry: logger is undefined")
	undefinedRetryOp = errors.New("retry: operation to retry is undefined")
)
type RetryStrategy interface {
	Execute(func(... interface{})(interface{}, error), ... interface{}) (interface{}, error)
}

type Retry struct {
	log.Logger
	Interval time.Duration // time we wait between retries
	MaxRetries int //maximum number of retries
	ShouldRetry func(interface{}, error) bool // provided function to determine whether or not to retry
	OnInternalFail func() interface{} // provided function to define some result in the case of failure
}

func (r *Retry) Execute(op func(... interface{}) (interface{}, error), arguments ... interface{}) (result interface{}, err error) {
	var (
		//todo: add logs
		//errorLogger = logging.Error(r.Logger)
		//debugLogger = logging.Debug(r.Logger)
	)

	if err = r.checkDependencies(); err != nil {
		result = r.OnInternalFail()
		return
	}

	if op == nil {
		err = undefinedRetryOp
		return
	}

	for attempt := 0; attempt < r.MaxRetries; attempt++ {
		result, err = op(arguments...)
		if !r.ShouldRetry(result, err) {
			break
		}
		time.Sleep(r.Interval)
	}
	return
}

func (r *Retry) checkDependencies() (err error) {
	if r.ShouldRetry == nil{
		err = undefinedShouldRetryErr
		return
	}

	if r.MaxRetries < 1 {
		err = invalidMaxRetriesValErr
		return
	}

	if r.Logger == nil {
		err = undefinedLogger
		return
	}

	return
}
