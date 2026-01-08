package store

import "context"

const fbpixelkey = "capi:fbpixelkey"
const fbpixelToken = "capi:fbpixelToken"

func SetFBPixel(ctx context.Context, data string) {
	s := Get()
	s.SetDefault(ctx, fbpixelkey, data)
}

func GetFBPixel(ctx context.Context) string {
	s := Get()
	v, ok := s.Get(ctx, fbpixelkey)
	if !ok {
		return ""
	}
	u, _ := v.(string)
	return u
}

func SetFBPixelToken(ctx context.Context, data string) {
	s := Get()
	s.SetDefault(ctx, fbpixelToken, data)
}

func GetFBPixelToken(ctx context.Context) string {
	s := Get()
	v, ok := s.Get(ctx, fbpixelToken)
	if !ok {
		return ""
	}
	u, _ := v.(string)
	return u
}
