package controllers

import (
	"api/src/autenticacao"
	"api/src/modelos"
	"api/src/respostas"
	"api/src/seguranca"
	"api/src/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario cria um novo usuário no sistema
// @Summary Criar um novo usuário
// @Description Cria um novo usuário no sistema
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuario body modelos.Usuario true "Dados do usuário"
// @Success 201 {object} modelos.Usuario
// @Failure 400 {object} respostas.Erro
// @Failure 422 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Router /usuarios [post]
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuario.ID, erro = repos.Usuario.Criar(usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}

// BuscarUsuarios busca todos os usuários que atendam um filtro de nome ou nick
// @Summary Buscar usuários
// @Description Busca usuários por nome ou nick
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   nome query string false "Nome ou nick do usuário"
// @Success 200 {array} modelos.Usuario
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios [get]
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	
	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarios, erro := repos.Usuario.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuario busca os dados detalhados de um usuário específico
// @Summary Buscar um usuário específico
// @Description Retorna os dados de um usuário específico
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Success 200 {object} modelos.Usuario
// @Failure 400 {object} respostas.Erro
// @Failure 404 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId} [get]
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuario, erro := repos.Usuario.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario altera as informações de um usuário no banco de dados
// @Summary Atualizar um usuário
// @Description Atualiza os dados de um usuário específico
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Param   usuario body modelos.Usuario true "Novos dados do usuário"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 403 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId} [put]
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Usuario.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario exclui as informações de um usuário no banco de dados
// @Summary Deletar um usuário
// @Description Remove um usuário do sistema
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 403 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId} [delete]
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível deletar um usuário que não seja o seu"))
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Usuario.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// SeguirUsuario permite que um usuário siga outro
// @Summary Seguir um usuário
// @Description Faz o usuário autenticado seguir outro usuário
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário a ser seguido"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/seguir [post]
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Usuario.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// PararDeSeguirUsuario permite que um usuário deixe de seguir outro
// @Summary Parar de seguir um usuário
// @Description Faz o usuário autenticado parar de seguir outro usuário
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário a deixar de seguir"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/parar-de-seguir [post]
func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível parar de seguir você mesmo"))
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Usuario.PararDeSeguir(usuarioID, seguidorID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarSeguidores retorna todos os seguidores de um usuário
// @Summary Buscar seguidores
// @Description Retorna todos os seguidores de um usuário
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Success 200 {array} modelos.Usuario
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/seguidores [get]
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	seguidores, erro := repos.Usuario.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, seguidores)
}

// BuscarSeguindo retorna todos os usuários que um usuário específico está seguindo
// @Summary Buscar usuários seguidos
// @Description Retorna todos os usuários que um usuário específico está seguindo
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Success 200 {array} modelos.Usuario
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/seguindo [get]
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarios, erro := repos.Usuario.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarSenha permite alterar a senha de um usuário
// @Summary Atualizar senha
// @Description Atualiza a senha do usuário autenticado
// @Tags usuarios
// @Accept  json
// @Produce  json
// @Param   senhas body map[string]string true "Nova senha e senha atual"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/atualizar-senha [post]
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	senhaSalvaNoBanco, erro := repos.Usuario.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não condiz com a que está salva no banco"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repos.Usuario.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
