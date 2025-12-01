package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*************** MODELOS ****************/

type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:text;not null;uniqueIndex" json:"email"`
	Senha     string    `gorm:"type:text;not null" json:"-"` // não expor senha no JSON
	Cargo     string    `gorm:"type:text;not null" json:"cargo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Almoxarifado struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Nome          string     `gorm:"type:text;not null" json:"nome"`
	Categoria     string     `gorm:"type:text;not null;index:idx_almox_categoria" json:"categoria"`
	DataValidade  *time.Time `gorm:"type:date" json:"data_validade,omitempty"`
	DataEntrada   time.Time  `gorm:"not null;default:now();index:idx_almox_data_entrada" json:"data_entrada"`
	DataSaida     *time.Time `json:"data_saida,omitempty"`
	CriadoPor     uint       `gorm:"not null;index" json:"criado_por"`
	CriadoPorUser Usuario    `gorm:"foreignKey:CriadoPor;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE" json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Patrimonio struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	Nome                string     `gorm:"type:text;not null" json:"nome"`
	IdentificacaoFisica string     `gorm:"type:text;uniqueIndex:uid_patrimonio_ident_fisica;not null" json:"identificacao_fisica"`
	Localizacao         string     `gorm:"type:text" json:"localizacao"`
	Status              string     `gorm:"type:text;not null;default:'ativo';index:idx_patr_status" json:"status"`
	DataEntrada         time.Time  `gorm:"not null;default:now();index:idx_patr_data_entrada" json:"data_entrada"`
	DataSaida           *time.Time `json:"data_saida,omitempty"`
	CriadoPor           uint       `gorm:"not null;index" json:"criado_por"`
	CriadoPorUser       Usuario    `gorm:"foreignKey:CriadoPor;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE" json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

/*************** DB ****************/

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("variável de ambiente %s não definida", key)
	}
	return v
}

func openDB() (*gorm.DB, string, error) {
	host := mustEnv("DB_URL") // ex.: db (service do compose)
	user := mustEnv("DB_USER")
	pass := mustEnv("DB_PASSWORD")
	name := mustEnv("DB_NAME")
	port := mustEnv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				// SingularTable deixa nome de tabela no singular (Usuario -> usuario)
				SingularTable: true,
			},
		},
	)
	return db, dsn, err
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Usuario{},
		&Almoxarifado{},
		&Patrimonio{},
	)
}

/*************** UTIL ***************/

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// parseID extrai o último segmento numérico do path (/api/{recurso}/{id})
func parseID(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, errors.New("id ausente")
	}
	last := parts[len(parts)-1]
	i, err := strconv.ParseUint(last, 10, 64)
	return uint(i), err
}

func hashPassword(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", errors.New("senha vazia")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(b), err
}

/*************** REQS ****************/

// Usuário
type usuarioCreateReq struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
	Cargo string `json:"cargo"`
}
type usuarioUpdateReq struct {
	Email *string `json:"email,omitempty"`
	Senha *string `json:"senha,omitempty"`
	Cargo *string `json:"cargo,omitempty"`
}

// Almoxarifado
type almoxCreateReq struct {
	Nome         string     `json:"nome"`
	Categoria    string     `json:"categoria"`
	DataValidade *time.Time `json:"data_validade,omitempty"`
	CriadoPor    uint       `json:"criado_por"`
}
type almoxUpdateReq struct {
	Nome         *string    `json:"nome,omitempty"`
	Categoria    *string    `json:"categoria,omitempty"`
	DataValidade *time.Time `json:"data_validade,omitempty"`
	DataSaida    *time.Time `json:"data_saida,omitempty"`
}

// Patrimônio
type patrCreateReq struct {
	Nome                string     `json:"nome"`
	IdentificacaoFisica string     `json:"identificacao_fisica"`
	Localizacao         string     `json:"localizacao"`
	Status              string     `json:"status"`
	CriadoPor           uint       `json:"criado_por"`
	DataSaida           *time.Time `json:"data_saida,omitempty"`
}
type patrUpdateReq struct {
	Nome                *string    `json:"nome,omitempty"`
	IdentificacaoFisica *string    `json:"identificacao_fisica,omitempty"`
	Localizacao         *string    `json:"localizacao,omitempty"`
	Status              *string    `json:"status,omitempty"`
	DataSaida           *time.Time `json:"data_saida,omitempty"`
}

/*************** HANDLERS: USUÁRIOS ****************/

func handleListUsuarios(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []Usuario
		if err := db.Order("id ASC").Find(&users).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, users)
	}
}

func handleGetUsuario(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var u Usuario
		if err := db.First(&u, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, u)
	}
}

func handleCreateUsuario(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req usuarioCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Cargo) == "" || strings.TrimSpace(req.Senha) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email, cargo e senha são obrigatórios"})
			return
		}
		hashed, err := hashPassword(req.Senha)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		u := Usuario{Email: req.Email, Senha: hashed, Cargo: req.Cargo}
		if err := db.Create(&u).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, u)
	}
}

func handleUpdateUsuario(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var req usuarioUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		var u Usuario
		if err := db.First(&u, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}

		updates := map[string]any{}
		if req.Email != nil {
			if strings.TrimSpace(*req.Email) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email inválido"})
				return
			}
			updates["email"] = *req.Email
		}
		if req.Cargo != nil {
			if strings.TrimSpace(*req.Cargo) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cargo inválido"})
				return
			}
			updates["cargo"] = *req.Cargo
		}
		if req.Senha != nil {
			if strings.TrimSpace(*req.Senha) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "senha inválida"})
				return
			}
			h, err := hashPassword(*req.Senha)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			updates["senha"] = h
		}

		if len(updates) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nenhum campo para atualizar"})
			return
		}
		if err := db.Model(&u).Updates(updates).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if err := db.First(&u, id).Error; err == nil {
			writeJSON(w, http.StatusOK, u)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "updated": true})
	}
}

func handleDeleteUsuario(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		if err := db.Delete(&Usuario{}, id).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted_id": id})
	}
}

/*************** HANDLERS: ALMOXARIFADO ****************/

func handleListAlmox(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var itens []Almoxarifado
		if err := db.Order("id ASC").Find(&itens).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, itens)
	}
}

func handleGetAlmox(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var a Almoxarifado
		if err := db.First(&a, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, a)
	}
}

func handleCreateAlmox(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req almoxCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(req.Nome) == "" || strings.TrimSpace(req.Categoria) == "" || req.CriadoPor == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nome, categoria e criado_por são obrigatórios"})
			return
		}
		a := Almoxarifado{
			Nome:         req.Nome,
			Categoria:    req.Categoria,
			DataValidade: req.DataValidade,
			CriadoPor:    req.CriadoPor,
		}
		if err := db.Create(&a).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, a)
	}
}

func handleUpdateAlmox(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var req almoxUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		var a Almoxarifado
		if err := db.First(&a, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}

		updates := map[string]any{}
		if req.Nome != nil {
			if strings.TrimSpace(*req.Nome) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nome inválido"})
				return
			}
			updates["nome"] = *req.Nome
		}
		if req.Categoria != nil {
			if strings.TrimSpace(*req.Categoria) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "categoria inválida"})
				return
			}
			updates["categoria"] = *req.Categoria
		}
		if req.DataValidade != nil {
			updates["data_validade"] = *req.DataValidade
		}
		if req.DataSaida != nil {
			updates["data_saida"] = *req.DataSaida
		}

		if len(updates) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nenhum campo para atualizar"})
			return
		}
		if err := db.Model(&a).Updates(updates).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if err := db.First(&a, id).Error; err == nil {
			writeJSON(w, http.StatusOK, a)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "updated": true})
	}
}

func handleDeleteAlmox(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		if err := db.Delete(&Almoxarifado{}, id).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted_id": id})
	}
}

/*************** HANDLERS: PATRIMÔNIO ****************/

func handleListPatr(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var itens []Patrimonio
		if err := db.Order("id ASC").Find(&itens).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, itens)
	}
}

func handleGetPatr(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var p Patrimonio
		if err := db.First(&p, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}
		writeJSON(w, http.StatusOK, p)
	}
}

func handleCreatePatr(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req patrCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if strings.TrimSpace(req.Nome) == "" ||
			strings.TrimSpace(req.IdentificacaoFisica) == "" ||
			strings.TrimSpace(req.Status) == "" ||
			req.CriadoPor == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nome, identificacao_fisica, status e criado_por são obrigatórios"})
			return
		}
		p := Patrimonio{
			Nome:                req.Nome,
			IdentificacaoFisica: req.IdentificacaoFisica,
			Localizacao:         req.Localizacao,
			Status:              req.Status,
			CriadoPor:           req.CriadoPor,
			DataSaida:           req.DataSaida,
		}
		if err := db.Create(&p).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, p)
	}
}

func handleUpdatePatr(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		var req patrUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		var p Patrimonio
		if err := db.First(&p, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			}
			return
		}

		updates := map[string]any{}
		if req.Nome != nil {
			if strings.TrimSpace(*req.Nome) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nome inválido"})
				return
			}
			updates["nome"] = *req.Nome
		}
		if req.IdentificacaoFisica != nil {
			if strings.TrimSpace(*req.IdentificacaoFisica) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "identificacao_fisica inválida"})
				return
			}
			updates["identificacao_fisica"] = *req.IdentificacaoFisica
		}
		if req.Localizacao != nil {
			updates["localizacao"] = *req.Localizacao
		}
		if req.Status != nil {
			if strings.TrimSpace(*req.Status) == "" {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "status inválido"})
				return
			}
			updates["status"] = *req.Status
		}
		if req.DataSaida != nil {
			updates["data_saida"] = *req.DataSaida
		}

		if len(updates) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nenhum campo para atualizar"})
			return
		}
		if err := db.Model(&p).Updates(updates).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if err := db.First(&p, id).Error; err == nil {
			writeJSON(w, http.StatusOK, p)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "updated": true})
	}
}

func handleDeletePatr(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.URL.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
			return
		}
		if err := db.Delete(&Patrimonio{}, id).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted_id": id})
	}
}

func handleLogin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string `json:"email"`
			Senha string `json:"senha"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		// Buscar usuário pelo email
		var u Usuario
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "email ou senha incorretos"})
			return
		}

		// Comparar senha
		if bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(req.Senha)) != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "email ou senha incorretos"})
			return
		}

		// OK → devolve dados do usuário (sem senha)
		resp := map[string]any{
			"id":         u.ID,
			"email":      u.Email,
			"cargo":      u.Cargo,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

/*************** CORS + ROUTER ****************/

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ajuste o domínio do seu front quando necessário
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func routes(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// ---------- USUÁRIOS ----------
	// /api/usuarios
	mux.Handle("/api/usuarios", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListUsuarios(db)(w, r)
		case http.MethodPost:
			handleCreateUsuario(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
	// /api/usuarios/{id}
	mux.Handle("/api/usuarios/", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsuario(db)(w, r)
		case http.MethodPut:
			handleUpdateUsuario(db)(w, r)
		case http.MethodDelete:
			handleDeleteUsuario(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	// ---------- ALMOXARIFADO ----------
	// /api/almoxarifado
	mux.Handle("/api/almoxarifado", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListAlmox(db)(w, r)
		case http.MethodPost:
			handleCreateAlmox(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
	// /api/almoxarifado/{id}
	mux.Handle("/api/almoxarifado/", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAlmox(db)(w, r)
		case http.MethodPut:
			handleUpdateAlmox(db)(w, r)
		case http.MethodDelete:
			handleDeleteAlmox(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	// ---------- PATRIMÔNIOS ----------
	// /api/patrimonios
	mux.Handle("/api/patrimonios", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleListPatr(db)(w, r)
		case http.MethodPost:
			handleCreatePatr(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
	// /api/patrimonios/{id}
	mux.Handle("/api/patrimonios/", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPatr(db)(w, r)
		case http.MethodPut:
			handleUpdatePatr(db)(w, r)
		case http.MethodDelete:
			handleDeletePatr(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	// ---------- LOGIN ----------
	mux.Handle("/api/login", cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleLogin(db)(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	return mux
}

/*************** MAIN ****************/

func main() {
	db, dsn, err := openDB()
	if err != nil {
		log.Fatalf("erro ao abrir conexão: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("erro ao obter sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("erro ao conectar no banco (DSN: %s): %v", dsn, err)
	}

	if err := migrate(db); err != nil {
		log.Fatalf("erro no AutoMigrate: %v", err)
	}
	log.Println("Migração concluída.")

	// HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("API ouvindo em %s", addr)
	if err := http.ListenAndServe(addr, routes(db)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("erro no servidor http: %v", err)
	}
}
