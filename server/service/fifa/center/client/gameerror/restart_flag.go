package gameerror

type RestartFlag struct {
	err error
}

func (r *RestartFlag) Error() string {
	return "账号重启：" + r.err.Error()
}

func NewRestartErr(err error) error {
	return &RestartFlag{err: err}
}
