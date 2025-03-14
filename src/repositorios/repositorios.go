package repositorios

import "database/sql"

// Repositories contém todos os repositórios da aplicação
type Repositories struct {
	Usuario    IUsuarioRepository
	Publicacao IPublicacaoRepository
}

// NovoRepositories cria uma nova instância de Repositories
func NovoRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Usuario:    NovoRepositorioDeUsuarios(db),
		Publicacao: NovoRepositorioDePublicacoes(db),
	}
} 