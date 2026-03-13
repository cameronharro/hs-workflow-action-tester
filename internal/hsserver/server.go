package hsserver

type HSServer struct {
	clientSecret string
}

func NewHSServer(clientSecret string) HSServer {
	return HSServer{
		clientSecret: clientSecret,
	}
}
