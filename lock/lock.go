package lock

import "time"

type Locker interface {

    TryLock(key string, expireTime time.Duration)bool

    Release(key string) bool

}

