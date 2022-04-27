module github.com/draeron/gof1

go 1.18

require (
	github.com/TheCodeTeam/goodbye v0.0.0-20170927022442-a83968bda2d3
	github.com/bearsh/hid v1.4.1-0.20220201202009-adffbb8a8616
	github.com/draeron/gopkgs v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
)

require (
	gitlab.com/gomidi/midi/v2 v2.0.7 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)

// replace gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi => gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi v0.0.0-20220419143954-65c25ed8cc67

replace github.com/draeron/gopkgs => ../gopkgs

//replace github.com/bearsh/hid => ../../external/hid
