package repository

type Repo interface {
	Get()
	Set()
	Delete()
	Init()
}
