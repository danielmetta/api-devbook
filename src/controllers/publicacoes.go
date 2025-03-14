package controllers

import (
	"api/src/autenticacao"
	"api/src/modelos"
	"api/src/respostas"
	"api/src/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	// mensagens de erro comuns
	msgErroPublicacaoNaoAutorizada = "Não é possível realizar operações em uma publicação que não seja sua"
)

// CriarPublicacao cria uma nova publicação no sistema
// @Summary Criar uma nova publicação
// @Description Cria uma nova publicação para o usuário autenticado
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacao body modelos.Publicacao true "Dados da publicação"
// @Success 201 {object} modelos.Publicacao
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 422 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes [post]
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publicacao.AutorID = usuarioID

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacao.ID, erro = repos.Publicacao.Criar(publicacao)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, publicacao)
}

// BuscarPublicacoes retorna as publicações que devem aparecer no feed do usuário
// @Summary Buscar publicações
// @Description Retorna as publicações que aparecem no feed do usuário autenticado
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Success 200 {array} modelos.Publicacao
// @Failure 401 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes [get]
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacoes, erro := repos.Publicacao.Buscar(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacoes)
}

// BuscarPublicacao retorna uma única publicação
// @Summary Buscar uma publicação específica
// @Description Retorna os dados de uma publicação específica
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacaoId path int true "ID da Publicação"
// @Success 200 {object} modelos.Publicacao
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes/{publicacaoId} [get]
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	publicacaoID, erro := strconv.ParseUint(mux.Vars(r)["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacao, erro := repos.Publicacao.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacao)
}

// AtualizarPublicacao altera os dados de uma publicação
// @Summary Atualizar uma publicação
// @Description Atualiza os dados de uma publicação específica
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacaoId path int true "ID da Publicação"
// @Param   publicacao body modelos.Publicacao true "Novos dados da publicação"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 403 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes/{publicacaoId} [put]
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	publicacaoID, erro := strconv.ParseUint(mux.Vars(r)["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacaoSalvaNoBanco, erro := repos.Publicacao.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if publicacaoSalvaNoBanco.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New(msgErroPublicacaoNaoAutorizada))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = publicacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repos.Publicacao.Atualizar(publicacaoID, publicacao); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarPublicacao exclui os dados de uma publicação
// @Summary Deletar uma publicação
// @Description Remove uma publicação específica
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacaoId path int true "ID da Publicação"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 403 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes/{publicacaoId} [delete]
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	publicacaoID, erro := strconv.ParseUint(mux.Vars(r)["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacaoSalvaNoBanco, erro := repos.Publicacao.BuscarPorID(publicacaoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if publicacaoSalvaNoBanco.AutorID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New(msgErroPublicacaoNaoAutorizada))
		return
	}

	if erro = repos.Publicacao.Deletar(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// BuscarPublicacoesPorUsuario retorna todas as publicações de um usuário específico
// @Summary Buscar publicações de um usuário
// @Description Retorna todas as publicações de um usuário específico
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   usuarioId path int true "ID do Usuário"
// @Success 200 {array} modelos.Publicacao
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /usuarios/{usuarioId}/publicacoes [get]
func BuscarPublicacoesPorUsuario(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := strconv.ParseUint(mux.Vars(r)["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	publicacoes, erro := repos.Publicacao.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacoes)
}

// CurtirPublicacao adiciona uma curtida na publicação
// @Summary Curtir uma publicação
// @Description Adiciona uma curtida em uma publicação específica
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacaoId path int true "ID da Publicação"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes/{publicacaoId}/curtir [post]
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	publicacaoID, erro := strconv.ParseUint(mux.Vars(r)["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Publicacao.Curtir(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DescurtirPublicacao remove uma curtida da publicação
// @Summary Descurtir uma publicação
// @Description Remove uma curtida de uma publicação específica
// @Tags publicacoes
// @Accept  json
// @Produce  json
// @Param   publicacaoId path int true "ID da Publicação"
// @Success 204 "No Content"
// @Failure 400 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Security ApiKeyAuth
// @Router /publicacoes/{publicacaoId}/descurtir [post]
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	publicacaoID, erro := strconv.ParseUint(mux.Vars(r)["publicacaoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repos.Publicacao.Descurtir(publicacaoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
