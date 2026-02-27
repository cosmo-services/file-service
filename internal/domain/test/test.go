package test

type TestRepo interface {
	GetOne() (int, error)
}
