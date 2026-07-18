package mockhttpclient

//go:generate go tool mockgen -package=mockhttpclient -destination=mockhttpclient.go -write_package_comment=false go.pact.im/x/httpclient Client
