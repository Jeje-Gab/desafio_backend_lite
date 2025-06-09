SET search_path = negociacoes;

CREATE TABLE negociacoes_daily_volume (
    codigo_instrumento VARCHAR(20) NOT NULL,
    data_negocio       DATE         NOT NULL,
    daily_sum          BIGINT       NOT NULL,
    PRIMARY KEY (codigo_instrumento, data_negocio)
);
CREATE INDEX idx_ndv_max
ON negociacoes_daily_volume (codigo_instrumento, daily_sum DESC);


INSERT INTO negociacoes_daily_volume
SELECT codigo_instrumento, data_negocio, SUM(quantidade_negociada)
FROM negociacoes.negociacoes
GROUP BY 1,2
    ON CONFLICT (codigo_instrumento, data_negocio)
DO UPDATE SET daily_sum = EXCLUDED.daily_sum;

CREATE INDEX idx_ndv_max
    ON negociacoes_daily_volume (codigo_instrumento, daily_sum DESC);
