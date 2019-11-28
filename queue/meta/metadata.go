package meta

import (
    "context"
    "strconv"
)

type  MD map[string]interface{}

type mdKey struct {}


//存储key value
func Put(ctx context.Context, key string, value interface{})  {
    if md , ok := FromContext(ctx); ok {
        md[key] = value
    } else {
        //is not ok
        md := MD{}
        md[key] = value
        context.WithValue(ctx, mdKey{}, md)
    }
}

func Value(ctx context.Context, key string) interface{} {
    if res, ok :=  ctx.Value(mdKey{}).(MD); !ok {
        return  nil
    } else {
        return res[key]
    }
}

func Bool(ctx context.Context, key string) bool {
    if res, ok := ctx.Value(mdKey{}).(MD); !ok {
        return  false
    } else {
        switch t := res[key].(type) {
        case bool:
            return t
        case string:
            ok,  _ = strconv.ParseBool(t)
            return  ok
        default:
            return false
        }
    }
}

func Int64(ctx context.Context, key string) int64 {
    if res, ok := ctx.Value(mdKey{}).(MD); !ok {
        return 0
    } else {
        switch t := res[key].(type) {
        case int:
            return int64(t)
        case int32:
            return int64(t)
        case int64:
            return int64(t)
        case string:
            if r, err := strconv.ParseInt(t, 10, 64) ;err == nil {
                return r
            }
            return 0
        default:
            return 0
        }
    }
}

func String(ctx context.Context, key string) string {
    if res, ok := ctx.Value(mdKey{}).(MD); ok {
        return res[key].(string)
    }
	return ""
}

func FromContext(ctx context.Context) (md MD, ok bool) {
    md, ok =   ctx.Value(mdKey{}).(MD)
    return
}

func NewContext(ctx context.Context, key MD) context.Context {
	return context.WithValue(ctx, mdKey{}, key)
}

func WithContext(ctx context.Context) context.Context {
    if md, ok := FromContext(ctx); ok {
        nmd := md.Copy()
        return NewContext(context.Background(), nmd)
    }
    return context.Background()
}

func New(params map[string]interface{})MD  {
    md := MD{}
    for key, value := range params {
        md[key]  = value
    }
    return md
}

func (md MD)Copy() MD {
    nmd := make(MD)
    for key, value := range md {
        nmd[key] = value
    }
    return nmd
}

func (md MD)Len()int  {
    return len(md)
}

