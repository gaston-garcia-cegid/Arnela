// run_sql.go
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Configura tu conexión aquí
	connStr := "postgres://arnela_user:arnela_secure_pass_2024@localhost:5432/arnela_db?sslmode=disable"
	sqlFile := "D:/repos/Arnela/backend/migrations/000003_add_nif_field.up.sql"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}
	defer db.Close()

	content, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("No se pudo leer el archivo SQL: %v", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Error al ejecutar SQL: %v", err)
	}

	fmt.Println("Migración ejecutada correctamente.")
}
