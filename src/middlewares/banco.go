package middlewares

import (
	"api/src/banco"
	"api/src/repositorios"
	"context"
	"net/http"
)

// ChaveRepositorios é a chave que será usada para armazenar os repositórios no contexto
const ChaveRepositorios = "repositories"

// InjetarDependencias é um middleware que injeta as dependências necessárias no contexto da requisição
func InjetarDependencias(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, erro := banco.Conectar()
		if erro != nil {
			http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Cria os repositórios
		repos := repositorios.NovoRepositories(db)

		// Adiciona os repositórios no contexto da requisição
		ctx := context.WithValue(r.Context(), ChaveRepositorios, repos)
		r = r.WithContext(ctx)

		next(w, r)
	}
} 