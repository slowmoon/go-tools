package xtime

import (
    "context"
    "time"
)

type Duration time.Duration

//parse duration
func (d *Duration)UnmarshalText(text []byte) error {
    s, err := time.ParseDuration(string(text))
    if err != nil {
        return err
    }
    *d = Duration(s)
    return nil
}

//shrink the context
func (d Duration)Shrink(ctx context.Context) (Duration, context.Context, context.CancelFunc) {
    if deadline , ok := ctx.Deadline(); ok {
        if delta := time.Until(deadline); delta < time.Duration(d) {
            return Duration(delta), ctx, func() {}
        }
    }
    ctx , cancel := context.WithTimeout(ctx, time.Duration(d))
    return d, ctx, cancel
}
