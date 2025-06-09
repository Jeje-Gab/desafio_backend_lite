package entity

import (
	"time"
)

// Negociacao representa uma linha de negociação extraída do CSV.
type Negociacao struct {
	DataReferencia              time.Time `json:"data_referencia"`
	CodigoInstrumento           string    `json:"codigo_instrumento"`
	AcaoAtualizacao             int       `json:"acao_atualizacao"`
	PrecoNegocio                float64   `json:"preco_negocio"`
	QuantidadeNegociada         int64     `json:"quantidade_negociada"`
	HoraFechamento              time.Time `json:"hora_fechamento"`
	CodigoIdentificadorNegocio  int64     `json:"codigo_identificador_negocio"`
	TipoSessaoPregao            int       `json:"tipo_sessao_pregao"`
	DataNegocio                 time.Time `json:"data_negocio"`
	CodigoParticipanteComprador *int64    `json:"codigo_participante_comprador,omitempty"`
	CodigoParticipanteVendedor  *int64    `json:"codigo_participante_vendedor,omitempty"`
}

// Stats agrupa os resultados de estatísticas de negociações.
type Stats struct {
	MaxPrice       float64 `json:"max_price"`
	MaxDailyVolume int64   `json:"max_daily_volume"`
}
