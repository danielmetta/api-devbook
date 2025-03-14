package repositorios

import "api/src/modelos"

// IUsuarioRepository define as operações disponíveis para o repositório de usuários
type IUsuarioRepository interface {
	Criar(usuario modelos.Usuario) (uint64, error)
	Buscar(nomeOuNick string) ([]modelos.Usuario, error)
	BuscarPorID(ID uint64) (modelos.Usuario, error)
	Atualizar(ID uint64, usuario modelos.Usuario) error
	Deletar(ID uint64) error
	BuscarPorEmail(email string) (modelos.Usuario, error)
	Seguir(usuarioID, seguidorID uint64) error
	PararDeSeguir(usuarioID, seguidorID uint64) error
	BuscarSeguidores(usuarioID uint64) ([]modelos.Usuario, error)
	BuscarSeguindo(usuarioID uint64) ([]modelos.Usuario, error)
	BuscarSenha(usuarioID uint64) (string, error)
	AtualizarSenha(usuarioID uint64, senha string) error
}

// IPublicacaoRepository define as operações disponíveis para o repositório de publicações
type IPublicacaoRepository interface {
	Criar(publicacao modelos.Publicacao) (uint64, error)
	BuscarPorID(publicacaoID uint64) (modelos.Publicacao, error)
	Buscar(usuarioID uint64) ([]modelos.Publicacao, error)
	Atualizar(publicacaoID uint64, publicacao modelos.Publicacao) error
	Deletar(publicacaoID uint64) error
	BuscarPorUsuario(usuarioID uint64) ([]modelos.Publicacao, error)
	Curtir(publicacaoID uint64) error
	Descurtir(publicacaoID uint64) error
} 