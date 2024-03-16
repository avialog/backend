package controller

type Controllers interface {
}

type controllers struct {
}

func NewControllers() Controllers {
	return &controllers{}
}
