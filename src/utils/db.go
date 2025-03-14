package utils

import (
	"api/src/middlewares"
	"api/src/repositorios"
	"errors"
	"net/http"
)

var (
	// ErrRepositoriosNaoEncontrados é retornado quando os repositórios não são encontrados no contexto
	ErrRepositoriosNaoEncontrados = errors.New("repositórios não encontrados no contexto")
)

// ExtrairRepositorios retorna os repositórios do contexto da requisição
func ExtrairRepositorios(r *http.Request) (*repositorios.Repositories, error) {
	repos, ok := r.Context().Value(middlewares.ChaveRepositorios).(*repositorios.Repositories)
	if !ok || repos == nil {
		return nil, ErrRepositoriosNaoEncontrados
	}
	return repos, nil
} 