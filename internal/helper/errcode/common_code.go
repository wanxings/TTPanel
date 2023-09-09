package errcode

var (
	Success       = NewError(200, "success")
	ServerError   = NewError(500, "service error")
	InvalidParams = NewError(501, "Parameter error")
	//NotFound      = NewError(404, "404 Not found")

	UnauthorizedTokenError   = NewError(402, "Authentication failed, Token is wrong or missing")
	UnauthorizedTokenTimeout = NewError(403, "Authentication failed, Token timed out")
	TooManyRequests          = NewError(503, "Too many requests")
)
