import React, { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { BarChart2, Globe, Info, Map as MapIcon, Percent } from 'lucide-react';
import './Mercados.css';
import { DataGrid } from '../../app/components/DataGrid';
import { HorizontalBarChart } from '../../app/components/HorizontalBarChart';
import { Panel } from '../../app/components/Panel';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { getANSMetrics, getStateMetrics } from '../../services/dashboard';
import { ANS_SORT_FIELDS, buildANSRanking } from '../../utils/ansRanking';
import { formatGDP, formatIncome, formatInteger, formatPenetration, formatPopulation } from '../../utils/format';

function regionTooltipContent(market) {
  return {
    title: market.region,
    rows: [
      { label: 'População', value: formatPopulation({ value: market.population }) },
      { label: 'PIB (Mil R$)', value: formatGDP({ value: market.gdp }) },
      { label: 'Renda média', value: formatIncome({ value: market.income }) },
      { label: 'UFs', value: formatInteger(market.ufs) },
      { label: 'Municípios', value: formatInteger(market.municipios) },
    ],
  };
}

export default function Mercados() {
  const statesRes = useApiResource(getStates, []);
  const [metricsRegion, setMetricsRegion] = useState('');
  const [ansRegion, setAnsRegion] = useState('');
  const [ansSort, setAnsSort] = useState('penetracao');

  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);

  const regionMetricsRes = useApiResource(({ signal }) => getStateMetrics({ signal }), []);
  const stateMetricsRes = useApiResource(
    ({ signal }) => getStateMetrics({ regiao: metricsRegion, signal }),
    [metricsRegion],
    { enabled: metricsRegion !== '' },
  );

  const ansRes = useApiResource(
    ({ signal }) => getANSMetrics({ regiao: ansRegion, ordenar: ansSort, limit: 10, signal }),
    [ansRegion, ansSort],
  );
  const ansMetrics = ansRes.data;

  const allMetrics = regionMetricsRes.data;
  const activeMetricsRes = metricsRegion ? stateMetricsRes : regionMetricsRes;
  const stateMetrics = activeMetricsRes.data;

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

  const ansRanking = useMemo(
    () => buildANSRanking(ansMetrics?.municipios),
    [ansMetrics],
  );

  const ansColumns = [
    { label: '#', field: 'posicao', width: '6%' },
    { label: 'Município', field: 'name', width: '30%' },
    { label: 'UF', field: 'uf', width: '8%' },
    { label: 'Região', field: 'region', width: '14%' },
    {
      label: 'Beneficiários',
      field: 'beneficiarios',
      width: '16%',
      render: (value) => formatInteger(value),
    },
    {
      label: 'População',
      field: 'populacao',
      width: '14%',
      render: (value) => formatInteger(value),
    },
    {
      label: 'Penetração',
      field: 'penetracao',
      width: '12%',
      render: (value) => formatPenetration(value),
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
            <HorizontalBarChart
              data={regionMarkets}
              labelKey="region"
              barKey="population"
              rowHeight={44}
              barSize={18}
              labelWidth={110}
              valueWidth={110}
              formatValue={(value) => formatPopulation({ value })}
              tooltip={regionTooltipContent}
            />
            <p className="panel-hint">
              Ordenado por população · barra proporcional à maior região · valores de população ao fim de cada barra ·
              passe o mouse sobre a barra para ver PIB, renda, UFs e municípios · renda média ponderada pela população de cada UF.
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
        {activeMetricsRes.loading && !stateMetrics ? (
          <div className="loading-box">Carregando indicadores por UF…</div>
        ) : activeMetricsRes.error ? (
          <div className="error-box">
            {activeMetricsRes.error.message}
            <button type="button" className="error-retry" onClick={activeMetricsRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : (
          <>
            <DataGrid columns={stateMetricColumns} data={stateMetrics ?? []} />
            <p className="panel-hint">
              {metricsCount} UF{metricsCount === 1 ? '' : 's'} · {metricsRegion ? `Região ${metricsRegion}` : 'todo o país'}
            </p>
          </>
        )}
      </Panel>

      <Panel title="Penetração ANS" icon={<Percent size={16} />}>
        <div className="panel-controls">
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="ans-region-filter">Região</label>
            <select
              id="ans-region-filter"
              className="filter-select"
              value={ansRegion}
              onChange={(event) => setAnsRegion(event.target.value)}
            >
              <option value="">Todas</option>
              {regionOptions.map((region) => (
                <option key={region} value={region}>{region}</option>
              ))}
            </select>
          </div>
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="ans-sort-filter">Ordenar por</label>
            <select
              id="ans-sort-filter"
              className="filter-select"
              value={ansSort}
              onChange={(event) => setAnsSort(event.target.value)}
            >
              {ANS_SORT_FIELDS.map((option) => (
                <option key={option.value} value={option.value}>{option.label}</option>
              ))}
            </select>
          </div>
        </div>
        {ansRes.loading && !ansMetrics ? (
          <div className="loading-box">Carregando penetração ANS…</div>
        ) : ansRes.error ? (
          <div className="error-box">
            {ansRes.error.message}
            <button type="button" className="error-retry" onClick={ansRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : ansMetrics?.municipios?.length === 0 ? (
          <div className="loading-box">Sem dado de beneficiários ANS para este recorte.</div>
        ) : (
          <>
            <DataGrid columns={ansColumns} data={ansRanking} />
            <p className="panel-hint">
              Top 10 de {formatInteger(ansMetrics?.total_municipios ?? 0)} municípios
              {' · '}{ansRegion ? `Região ${ansRegion}` : 'todo o país'}
              {' · '}ordenado por {ANS_SORT_FIELDS.find((option) => option.value === ansSort)?.label.toLowerCase()}
              {ansMetrics?.fonte ? (
                <>
                  {' · '}fonte {ansMetrics.fonte}
                </>
              ) : null}
              {' · '}população: Censo 2022.
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
              IBGE — população e renda: Censo 2022 · PIB: 2023.
              {ansMetrics?.fonte
                ? ` ANS — beneficiários de planos: ${ansMetrics.fonte} (${ansMetrics.ano}).`
                : ''}{' '}
              Os dados são importados das fontes oficiais para o
              banco da Trads.
            </dd>
          </div>
          <div className="mercados-guide-item">
            <dt>Penetração ANS</dt>
            <dd>
              Beneficiários de planos do município divididos pela população do Censo 2022 —
              quanto da população está coberta por plano. Top 10 com filtro de região.
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
