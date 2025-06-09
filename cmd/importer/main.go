package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// transformRecord converte um slice de strings em uma fatia de interfaces pronta para COPY.
func transformRecord(rec []string) ([]interface{}, error) {
	// Parse das datas
	dataRef, err := time.Parse("2006-01-02", rec[0])
	if err != nil {
		return nil, fmt.Errorf("data_referencia inválida: %w", err)
	}
	dataNeg, err := time.Parse("2006-01-02", rec[8])
	if err != nil {
		return nil, fmt.Errorf("data_negocio inválida: %w", err)
	}
	// Parse do preço substituindo "." e ","
	precoStr := strings.ReplaceAll(rec[3], ".", "")
	precoStr = strings.ReplaceAll(precoStr, ",", ".")
	preco, err := strconv.ParseFloat(precoStr, 64)
	if err != nil {
		return nil, fmt.Errorf("preco_negocio inválido: %w", err)
	}
	// Parse da hora no formato "hhmmssSSS"
	h := rec[5]
	horaFmt := fmt.Sprintf("%s:%s:%s.%s", h[0:2], h[2:4], h[4:6], h[6:])
	hora, err := time.Parse("15:04:05.000", horaFmt)
	if err != nil {
		return nil, fmt.Errorf("hora_fechamento inválida: %w", err)
	}
	// Field converter: vazio -> nil
	toNull := func(s string) interface{} {
		if s == "" {
			return nil
		}
		v, _ := strconv.ParseInt(s, 10, 64)
		return v
	}
	return []interface{}{dataRef, rec[1], rec[2], preco, rec[4], hora, rec[6], rec[7], dataNeg, toNull(rec[9]), toNull(rec[10])}, nil
}

func importCSV(ctx context.Context, db *sql.DB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("abrindo %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	// Pular header
	if _, err := r.Read(); err != nil {
		return fmt.Errorf("lendo header em %s: %w", path, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("iniciando transação: %w", err)
	}

	// Preparar COPY IN
	stmt, err := tx.Prepare(pq.CopyIn("negociacoes",
		"data_referencia",
		"codigo_instrumento",
		"acao_atualizacao",
		"preco_negocio",
		"quantidade_negociada",
		"hora_fechamento",
		"codigo_identificador_negocio",
		"tipo_sessao_pregao",
		"data_negocio",
		"codigo_participante_comprador",
		"codigo_participante_vendedor",
	))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("preparando COPY: %w", err)
	}

	count := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("erro lendo %s: %v", path, err)
			continue
		}

		args, err := transformRecord(rec)
		if err != nil {
			log.Printf("transformando linha %v: %v", rec, err)
			continue
		}

		if _, err := stmt.Exec(args...); err != nil {
			log.Printf("COPY exec erro: %v", err)
			continue
		}

		count++
		if count%100000 == 0 {
			log.Printf("%d registros processados em %s", count, filepath.Base(path))
		}
	}

	// Finalizar COPY
	if _, err := stmt.Exec(); err != nil {
		stmt.Close()
		tx.Rollback()
		return fmt.Errorf("finalizando COPY: %w", err)
	}
	stmt.Close()

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transação: %w", err)
	}

	log.Printf("importados %d registros de %s", count, filepath.Base(path))
	return nil
}

func main() {
	// 1) flags
	dir := flag.String("dir", "../data-csv", "diretório com arquivos CSV")
	workers := flag.Int("workers", runtime.NumCPU(), "número de goroutines concorrentes")
	flag.Parse() // <— não esqueça

	// 2) string de conexão fixa com search_path
	connStr := "host=localhost " +
		"port=5432 " +
		"user=user " +
		"password=pass " +
		"dbname=trading " +
		"sslmode=disable " +
		"search_path=negociacoes,public"

	log.Printf("▶️  Usando conexão: %s", connStr)

	// 3) abre e testa a conexão
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ sql.Open deu erro: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ db.Ping deu erro: %v", err)
	}
	log.Println("✅ Conectado com sucesso ao Postgres!")

	// 4) lista CSVs
	entries, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatalf("lendo diretório %q: %v", *dir, err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".csv") {
			continue
		}
		files = append(files, filepath.Join(*dir, e.Name()))
	}

	// 5) pool de workers
	sem := make(chan struct{}, *workers)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer wg.Done()
			if err := importCSV(ctx, db, path); err != nil {
				log.Printf("erro importando %s: %v", filepath.Base(path), err)
			}
			<-sem
		}(file)
	}
	wg.Wait()
}
