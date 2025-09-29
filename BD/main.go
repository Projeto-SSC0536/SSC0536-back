package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/*************** MODELOS ****************/

type Usuario struct {
	ID        uint   `gorm:"primaryKey"`
	Nome      string `gorm:"type:text;not null"`
	Senha     string `gorm:"type:text;not null"`
	Cargo     string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Almoxarifado struct {
	ID            uint       `gorm:"primaryKey"`
	Nome          string     `gorm:"type:text;not null"`
	Categoria     string     `gorm:"type:text;not null;index:idx_almox_categoria"`
	DataValidade  *time.Time `gorm:"type:date"`
	DataEntrada   time.Time  `gorm:"not null;default:now();index:idx_almox_data_entrada"`
	DataSaida     *time.Time
	CriadoPor     uint    `gorm:"not null;index"`
	CriadoPorUser Usuario `gorm:"foreignKey:CriadoPor;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Patrimonio struct {
	ID                  uint      `gorm:"primaryKey"`
	Nome                string    `gorm:"type:text;not null"`
	IdentificacaoFisica string    `gorm:"type:text;uniqueIndex:uid_patrimonio_ident_fisica;not null"`
	Localizacao         string    `gorm:"type:text"`
	Status              string    `gorm:"type:text;not null;default:'ativo';index:idx_patr_status"`
	DataEntrada         time.Time `gorm:"not null;default:now();index:idx_patr_data_entrada"`
	DataSaida           *time.Time
	CriadoPor           uint    `gorm:"not null;index"`
	CriadoPorUser       Usuario `gorm:"foreignKey:CriadoPor;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
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
				// Deixa SingularTable como quiser; TableName() já fixa os nomes.
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

func verify(db *gorm.DB) {
	// Banco conectado
	var dbname string
	_ = db.Raw("select current_database()").Scan(&dbname)
	log.Printf("Conectado ao banco: %s", dbname)

	// Tabelas existem?
	type pair struct {
		name  string
		model interface{}
	}
	checks := []pair{
		{"usuarios", &Usuario{}},
		{"almoxarifado", &Almoxarifado{}},
		{"patrimonios", &Patrimonio{}},
	}
	for _, c := range checks {
		ok := db.Migrator().HasTable(c.model)
		log.Printf("Tabela '%s' existe? %v", c.name, ok)
	}
}

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

	log.Println("Migração concluída: tabelas criadas/atualizadas (usuarios, almoxarifado, patrimonios).")
	verify(db)
}
