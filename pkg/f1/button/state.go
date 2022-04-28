package button

//go:generate go-enum -f=$GOFILE --noprefix

/*
PushState x ENUM(
	Released
	Pushed
) */
type PushState int

