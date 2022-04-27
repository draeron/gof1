package button

//go:generate go-enum -f=$GOFILE --noprefix

/*
BtnType x ENUM(
	Push
	Absolute
	Relative
) */
type BtnType int

func (t BtnType) Buttons() (s []Button) {
	for _, btn := range Values() {
		if btn.Type() == t {
			s = append(s, btn)
		}
	}
	return
}
