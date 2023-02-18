package rest

// MediaType - Custom type to hold value for find and replace on media type.
type MediaType int

// Declare related constants for each MediaType starting with index 1.
const (
	UnknownMedia MediaType = iota
	ContentType
	ApplicationJSON
	FormURLEncoded
	MultipartForm
	TextPlain
	XContentTypOptions
)

// String - Creating common behavior - give the type a String function.
func (m MediaType) String() string {
	return [...]string{
		"text/unknown",
		"Content-Type",
		"application/json",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
		"text/plain",
		"X-Content-Type-Options",
	}[m]
}

// Index - Return index of the Constant.
func (m MediaType) Index() int {
	return int(m)
}
