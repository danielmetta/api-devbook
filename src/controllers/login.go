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
)

var (
	// ErrCredenciaisInvalidas é retornado quando as credenciais do usuário estão incorretas
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas")
)

// Login é responsável por autenticar um usuário na API
// @Summary Autenticar usuário
// @Description Autentica um usuário na API e retorna um token JWT
// @Tags autenticacao
// @Accept  json
// @Produce  json
// @Param   credentials body modelos.Usuario true "Credenciais do usuário"
// @Success 200 {object} modelos.DadosAutenticacao
// @Failure 400 {object} respostas.Erro
// @Failure 401 {object} respostas.Erro
// @Failure 422 {object} respostas.Erro
// @Failure 500 {object} respostas.Erro
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
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

	if erro = validarCredenciais(&usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	repos, erro := utils.ExtrairRepositorios(r)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioSalvoNoBanco, erro := repos.Usuario.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, ErrCredenciaisInvalidas)
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioID := strconv.FormatUint(usuarioSalvoNoBanco.ID, 10)
	respostas.JSON(w, http.StatusOK, modelos.DadosAutenticacao{ID: usuarioID, Token: token})
}

// validarCredenciais verifica se as credenciais do usuário são válidas
func validarCredenciais(usuario *modelos.Usuario) error {
	if usuario.Email == "" {
		return errors.New("o e-mail é obrigatório")
	}
	if usuario.Senha == "" {
		return errors.New("a senha é obrigatória")
	}
	return nil
}
