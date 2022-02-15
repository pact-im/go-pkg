package mockclock

//go:generate mockgen -package=mockclock -destination=mockclock.go -write_package_comment=false -source=clock.go
//go:generate mockgen -package=mockclock -destination=timer.go -write_package_comment=false go.pact.im/x/clock Timer
//go:generate mockgen -package=mockclock -destination=ticker.go -write_package_comment=false go.pact.im/x/clock Ticker
