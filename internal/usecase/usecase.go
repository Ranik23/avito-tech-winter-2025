package usecase


import (
	"avito/internal/storage/postgres"
)

type UserCase interface {

}


type UserOperator struct {
	strg postgres.Storage
}