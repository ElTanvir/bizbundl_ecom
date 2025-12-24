package store

const fbpixelkey = "capi:fbpixelkey"
const fbpixelToken = "capi:fbpixelToken"

func SetFBPixel(data string) {
	s := Get()
	s.Set(fbpixelkey, data)
}

func GetFBPixel() string {
	s := Get()
	v, ok := s.Get(fbpixelkey)
	if !ok {
		return ""
	}
	u, _ := v.(string)
	return u
}

func SetFBPixelToken(data string) {
	s := Get()
	s.Set(fbpixelToken, data)
}

func GetFBPixelToken() string {
	s := Get()
	v, ok := s.Get(fbpixelToken)
	if !ok {
		return ""
	}
	u, _ := v.(string)
	return u
}

