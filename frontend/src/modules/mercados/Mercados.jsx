import React, { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { BarChart2, Globe, Info, Map as MapIcon } from 'lucide-react';
import './Mercados.css';
import { DataGrid } from '../../app/components/DataGrid';
import { Panel } from '../../app/components/Panel';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { getStateMetrics } from '../../services/dashboard';
import { formatGDP, formatIncome, formatInteger, formatPopulation } from '../../utils/format';

export default function Mercados() {
  const statesRes = useApiResource(getStates, []);
  const [metricsRegion, setMetricsRegion] = useState('');

  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);

  const regionMetricsRes = useApiResource(({ signal }) => getStateMetrics({ signal }), []);
  const stateMetricsRes = useApiResource(
    ({ signal }) => getStateMetrics({ regiao: metricsRegion, signal }),
    [metricsRegion],
  );

  const allMetrics = regionMetricsRes.data;
  const stateMetrics = stateMetricsRes.data;

  const regionMarkets = useMemo(() => {
    const markets = new Map();

    for (const item of allMetrics ?? []) {
      const market = markets.get(item.region) ?? {
        region: item.region,
        ufs: 0,
        municipios: 0,
        population: 0,
        gdp: 0,
        incomeWeighted: 0,
        incomeBase: 0,
      };

      const population = item.indicators?.population?.value ?? 0;
      const gdp = item.indicators?.gdp?.value ?? 0;
      const income = item.indicators?.income?.value ?? 0;

      market.ufs += 1;
      market.municipios += item.municipios ?? 0;
      market.population += population;
      market.gdp += gdp;

      if (population > 0 && income > 0) {
        market.incomeWeighted += income * population;
        market.incomeBase += population;
      }

      markets.set(item.region, market);
    }

    return [...markets.values()]
      .map((market) => ({
        ...market,
        income: market.incomeBase > 0 ? market.incomeWeighted / market.incomeBase : 0,
      }))
      .sort((a, b) => b.population - a.population);
  }, [allMetrics]);

  const maxRegionPopulation = regionMarkets.length > 0 ? regionMarkets[0].population : 1;

  const regionOptions = useMemo(
    () => [...new Set(states.map((state) => state.region))].sort((a, b) => a.localeCompare(b)),
    [states],
  );

  const metricsCount = stateMetrics?.length ?? 0;

  const stateMetricColumns = [
    { label: 'UF', field: 'uf', width: '8%' },
    { label: 'Estado', field: 'name', width: '24%' },
    {
      label: 'Municípios',
      field: 'municipios',
      width: '14%',
      render: (value) => formatInteger(value),
    },
    {
      label: 'População',
      field: 'indicators',
      width: '18%',
      render: (value) => formatPopulation(value?.population),
    },
    {
      label: 'Renda média',
      field: 'indicators',
      width: '18%',
      render: (value) => formatIncome(value?.income),
    },
    {
      label: 'PIB (Mil R$)',
      field: 'indicators',
      width: '18%',
      render: (value) => formatGDP(value?.gdp),
    },
  ];

  return (
    <div className="mercados-page">
      <header className="mercados-intro">
        <h1 className="mercados-title">
          <Globe size={20} />
          Em quais regiões estão os melhores mercados?
        </h1>
        <p className="mercados-subtitle">
          População, PIB e renda das cidades, agregados por região e detalhados por UF —
          a primeira metade da pergunta que guia a Trads.
        </p>
      </header>

      <Panel title="Mercados por Região" icon={<BarChart2 size={16} />}>
        {regionMetricsRes.loading && regionMarkets.length === 0 ? (
          <div className="loading-box">Carregando mercados…</div>
        ) : regionMetricsRes.error ? (
          <div className="error-box">
            {regionMetricsRes.error.message}
            <button type="button" className="error-retry" onClick={regionMetricsRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : (
          <>
            <ul className="region-list">
              {regionMarkets.map((market) => (
                <li key={market.region} className="region-item">
                  <span className="region-name" title={market.region}>{market.region}</span>
                  <span className="region-bar">
                    <span
                      className="region-bar-fill"
                      style={{ width: `${(market.population / maxRegionPopulation) * 100}%` }}
                    />
                  </span>
                  <span className="region-count">{formatPopulation({ value: market.population })}</span>
                  <span className="region-meta">
                    {market.ufs} UFs · {formatInteger(market.municipios)} municípios · PIB (Mil R$) {formatGDP({ value: market.gdp })} · renda média {formatIncome({ value: market.income })}
                  </span>
                </li>
              ))}
            </ul>
            <p className="panel-hint">
              Ordenado por população · barra proporcional à maior região · renda média ponderada pela população de cada UF.
            </p>
          </>
        )}
      </Panel>

      <Panel title="UFs por indicador" icon={<MapIcon size={16} />}>
        <div className="panel-controls">
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="metrics-region-filter">Região</label>
            <select
              id="metrics-region-filter"
              className="filter-select"
              value={metricsRegion}
              onChange={(event) => setMetricsRegion(event.target.value)}
            >
              <option value="">Todas</option>
              {regionOptions.map((region) => (
                <option key={region} value={region}>{region}</option>
              ))}
            </select>
          </div>
        </div>
        {stateMetricsRes.loading && !stateMetrics ? (
          <div className="loading-box">Carregando indicadores por UF…</div>
        ) : stateMetricsRes.error ? (
          <div className="error-box">
            {stateMetricsRes.error.message}
            <button type="button" className="error-retry" onClick={stateMetricsRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : (
          <>
            <DataGrid columns={stateMetricColumns} data={stateMetrics ?? []} selectable={false} />
            <p className="panel-hint">
              {metricsCount} UF{metricsCount === 1 ? '' : 's'} · {metricsRegion ? `Região ${metricsRegion}` : 'todo o país'}
            </p>
          </>
        )}
      </Panel>

      <Panel
        title="Como ler estes números"
        icon={<Info size={16} />}
        actions={<Link className="panel-link" to="/metodologia">Ver metodologia</Link>}
      >
        <dl className="mercados-guide">
          <div className="mercados-guide-item">
            <dt>População · Municípios · PIB</dt>
            <dd>
              Somas das cidades de cada UF e, no painel de regiões, das UFs de cada
              região. PIB em Mil R$.
            </dd>
          </div>
          <div className="mercados-guide-item">
            <dt>Renda média</dt>
            <dd>
              Média da renda das cidades ponderada pela população: dentro do estado,
              cidade grande pesa mais; dentro da região, UF grande pesa mais.
            </dd>
          </div>
          <div className="mercados-guide-item">
            <dt>Ano e fonte</dt>
            <dd>
              IBGE — população e renda: Censo 2022 · PIB: 2023. Os dados são importados
              da API do IBGE para o banco da Trads.
            </dd>
          </div>
          <div className="mercados-guide-item">
            <dt>Filtro Região</dt>
            <dd>
              Recarrega só a tabela por UF, com ida ao servidor. O painel de regiões acima
              não muda — ele sempre compara as cinco regiões do país.
            </dd>
          </div>
        </dl>
      </Panel>
    </div>
  );
}
