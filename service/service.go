package service

type Services interface {
	GetDTO() interface{}
	GetMultipartDTO() interface{}
}
