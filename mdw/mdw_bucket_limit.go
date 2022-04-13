package mdw

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

//BucketLimitMiddleware add bucket in fillInterval time
func BucketLimitMiddleware(fillInterval time.Duration, cap int64, limited func(c *gin.Context)) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) == 0 {
			limited(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
