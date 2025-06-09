CREATE SCHEMA negociacoes;

SET search_path = negociacoes;


CREATE TABLE negociacoes (
    id                             BIGSERIAL PRIMARY KEY,
    data_referencia                DATE                NOT NULL,
    codigo_instrumento             VARCHAR(20)         NOT NULL,
    acao_atualizacao               SMALLINT            NOT NULL,
    preco_negocio                  NUMERIC(15,3)       NOT NULL,
    quantidade_negociada           BIGINT              NOT NULL,
    hora_fechamento                TIME                NOT NULL,
    codigo_identificador_negocio   BIGINT              NOT NULL,
    tipo_sessao_pregao             SMALLINT            NOT NULL,
    data_negocio                   DATE                NOT NULL,
    codigo_participante_comprador  INTEGER                     ,
    codigo_participante_vendedor   INTEGER                     ,
    CONSTRAINT uq_neg_ticker_data  UNIQUE (codigo_instrumento, data_referencia, hora_fechamento, codigo_identificador_negocio)
);

-- Índice para otimizar busca do máximo de preço (com e sem filtro de data)
CREATE INDEX idx_price_stats
    ON negociacoes.negociacoes (
                                codigo_instrumento,
                                preco_negocio DESC,
                                data_negocio DESC
        );

-- Índice para otimizar agregação do volume diário (com e sem filtro de data)
CREATE INDEX idx_volume_stats
    ON negociacoes.negociacoes (
                                codigo_instrumento,
                                data_negocio
        )
    INCLUDE (quantidade_negociada);
