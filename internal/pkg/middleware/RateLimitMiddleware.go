// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"gotribe-admin/pkg/api/response"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.ResponseWithCode(c, 429, nil, "访问限流")
			c.Abort()
			return
		}
		c.Next()
	}
}
