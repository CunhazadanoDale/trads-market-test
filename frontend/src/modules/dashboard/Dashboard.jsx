import React, { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Activity, BarChart2, Map as MapIcon, MapPin, PieChart, RefreshCw, Users, Wallet } from 'lucide-react';
import './Dashboard.css';
import { DataGrid } from '../../app/components/DataGrid';
import { Panel } from '../../app/components/Panel';
import { StatusBadge } from '../../app/components/StatusBadge';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { getCities } from '../../services/cities';
import { getHealth } from '../../services/health';
import { getAgeDistribution, getNationalMetrics, getStateMetrics, getTopCities } from '../../services/dashboard';
import { API_BASE_URL } from '../../services/api';
import { formatGDP, formatIncome, formatInteger, formatPercent, formatPopulation } from '../../utils/format';

const CITIES_PANEL_SIZE = 8;

function ServiceRow({ label, path, probe }) {
  let status = <StatusBadge status="neutral" label="..." />;

  if (probe?.ok) {
    status = (
      <span className="service-meta">
        {probe.latencyMs} ms
        <StatusBadge status="success" label="OK" />
      </span>
    );
  } else if (probe) {
    status = <StatusBadge status="danger" label="Offline" />;
  }

  return (
    <li className="service-item">
      <span className="service-name">
        {label} <code>{path}</code>
      </span>
      {status}
    </li>
  );
}

export default function Dashboard() {
  const statesRes = useApiResource(getStates, []);
  const healthRes = useApiResource(getHealth, []);
  // Estado escolhido; quando vazio, vale o primeiro da lista (derivado)
  const [rawIbge, setRawIbge] = useState('');
  const [ageRegion, setAgeRegion] = useState('');
  const [ageIbge, setAgeIbge] = useState('');
  const [metricsRegion, setMetricsRegion] = useState('');

  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);
  const selectedIbge = rawIbge || (states.length > 0 ? String(states[0].ibge_code) : '');

  const citiesRes = useApiResource(
    ({ signal }) => getCities(selectedIbge, { page: 1, pageSize: CITIES_PANEL_SIZE, signal }),
    [selectedIbge],
  );

  const nationalRes = useApiResource(getNationalMetrics, []);
  const stateMetricsRes = useApiResource(
    ({ signal }) => getStateMetrics({ regiao: metricsRegion, signal }),
    [metricsRegion],
  );
  const topCitiesRes = useApiResource(getTopCities, []);
  const ageRes = useApiResource(
    ({ signal }) => getAgeDistribution({ regiao: ageRegion, ibge: ageIbge, signal }),
    [ageRegion, ageIbge],
  );

  const selectedState = useMemo(
    () => states.find((state) => String(state.ibge_code) === selectedIbge) ?? null,
    [states, selectedIbge],
  );

  const stateMetrics = stateMetricsRes.data;

  const regionMarkets = useMemo(() => {
    const markets = new Map();

    for (const item of stateMetrics ?? []) {
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
  }, [stateMetrics]);

  const maxRegionPopulation = regionMarkets.length > 0 ? regionMarkets[0].population : 1;

  const regionOptions = useMemo(
    () => [...new Set(states.map((state) => state.region))].sort((a, b) => a.localeCompare(b)),
    [states],
  );

  const ageStates = useMemo(
    () => states.filter((state) => !ageRegion || state.region === ageRegion),
    [states, ageRegion],
  );

  const selectedAgeState = useMemo(
    () => states.find((state) => String(state.ibge_code) === ageIbge) ?? null,
    [states, ageIbge],
  );

  const handleAgeRegionChange = (region) => {
    setAgeRegion(region);
    setAgeIbge((current) => {
      if (!current) return current;
      const state = states.find((item) => String(item.ibge_code) === current);
      return state && (!region || state.region === region) ? current : '';
    });
  };

  const health = healthRes.data;
  const apiUp = Boolean(health?.app?.ok);
  const dbUp = Boolean(health?.database?.ok);
  const allUp = apiUp && dbUp;

  const summaryCards = [
    {
      title: 'Estados',
      value: statesRes.loading && states.length === 0 ? '…' : String(states.length),
      hint: 'GET /api/v1/states',
    },
    {
      title: selectedState ? `Cidades em ${selectedState.uf}` : 'Cidades',
      value: !selectedState
        ? '—'
        : citiesRes.loading && !citiesRes.data
          ? '…'
          : String(citiesRes.data?.total ?? 0),
      hint: 'total do estado selecionado',
    },
    {
      title: 'Regiões',
      value: stateMetricsRes.loading && regionMarkets.length === 0 ? '…' : String(regionMarkets.length),
      hint: 'agregados por região',
    },
    {
      title: 'Status da API',
      value: health ? (allUp ? 'OK' : 'Offline') : '…',
      hint: apiUp && !dbUp ? 'banco indisponível' : '/health + /health/db',
    },
  ];

  const national = nationalRes.data;
  const nationalPending = nationalRes.loading && !national;

  const yearHint = (indicator) => {
    if (nationalPending) return '…';
    return indicator ? `Ano ${indicator.year}` : 'sem dado';
  };

  const nationalCards = [
    {
      title: 'Municípios',
      value: nationalPending ? '…' : national ? formatInteger(national.municipios) : '—',
      hint: 'todo o país',
    },
    {
      title: 'População',
      value: nationalPending ? '…' : formatPopulation(national?.indicators?.population),
      hint: yearHint(national?.indicators?.population),
    },
    {
      title: 'Renda média',
      value: nationalPending ? '…' : formatIncome(national?.indicators?.income),
      hint: yearHint(national?.indicators?.income),
    },
    {
      title: 'PIB (Mil R$)',
      value: nationalPending ? '…' : formatGDP(national?.indicators?.gdp),
      hint: yearHint(national?.indicators?.gdp),
    },
  ];

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

  const rankings = [
    {
      title: 'Top PIB',
      icon: <BarChart2 size={16} />,
      cities: topCitiesRes.data?.top_pib,
      pick: (city) => city.indicators?.gdp,
      format: formatGDP,
    },
    {
      title: 'Top renda',
      icon: <Wallet size={16} />,
      cities: topCitiesRes.data?.top_renda,
      pick: (city) => city.indicators?.income,
      format: formatIncome,
    },
    {
      title: 'Top população',
      icon: <Users size={16} />,
      cities: topCitiesRes.data?.top_populacao,
      pick: (city) => city.indicators?.population,
      format: formatPopulation,
    },
  ];

  const age = ageRes.data;
  const ageGroups = useMemo(() => age?.grupos ?? [], [age]);
  const maxAgePopulation =
    ageGroups.length > 0 ? Math.max(...ageGroups.map((group) => group.populacao)) : 1;

  const ageScope = [
    ageRegion ? `Região ${ageRegion}` : '',
    selectedAgeState ? `UF ${selectedAgeState.uf}` : '',
  ]
    .filter(Boolean)
    .join(' · ');

  const handleRefresh = () => {
    statesRes.reload();
    healthRes.reload();
    citiesRes.reload();
    nationalRes.reload();
    stateMetricsRes.reload();
    topCitiesRes.reload();
    ageRes.reload();
  };

  return (
    <div className="dashboard-container">
      <div className="dashboard-toolbar">
        <div className="toolbar-field">
          <span className="toolbar-label">Estado:</span>
          <select
            className="filter-select"
            value={selectedIbge}
            onChange={(event) => setRawIbge(event.target.value)}
          >
            {states.length === 0 && <option value="">Carregando…</option>}
            {states.map((state) => (
              <option key={state.ibge_code} value={String(state.ibge_code)}>
                {`${state.uf} — ${state.name}`}
              </option>
            ))}
          </select>
        </div>

        <button
          type="button"
          className="btn btn--secondary"
          onClick={handleRefresh}
          disabled={statesRes.loading || healthRes.loading}
        >
          <RefreshCw size={14} className={statesRes.loading ? 'spin' : undefined} />
          Atualizar
        </button>
      </div>

      {statesRes.error && (
        <div className="error-box">
          {statesRes.error.message}
          <button type="button" className="error-retry" onClick={statesRes.reload}>
            Tentar novamente
          </button>
        </div>
      )}

      <div className="summary-cards">
        {[...summaryCards, ...nationalCards].map((card) => (
          <div key={card.title} className="summary-card">
            <span className="summary-card-title">{card.title}</span>
            <span className="summary-card-value">{card.value}</span>
            <span className="summary-card-hint">{card.hint}</span>
          </div>
        ))}
      </div>

      <div className="dashboard-grid">
        <Panel title="Mercados por Região" icon={<BarChart2 size={16} />}>
          {stateMetricsRes.loading && regionMarkets.length === 0 ? (
            <div className="loading-box">Carregando mercados…</div>
          ) : stateMetricsRes.error ? (
            <div className="error-box">
              {stateMetricsRes.error.message}
              <button type="button" className="error-retry" onClick={stateMetricsRes.reload}>
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

        <Panel
          title="Status dos Serviços"
          icon={<Activity size={16} />}
          actions={(
            <button
              type="button"
              className="topbar-action-btn"
              onClick={healthRes.reload}
              title="Verificar novamente"
            >
              <RefreshCw size={14} />
            </button>
          )}
        >
          <ul className="service-list">
            <ServiceRow label="API" path="/health" probe={health?.app} />
            <ServiceRow label="Banco" path="/health/db" probe={health?.database} />
          </ul>
          <p className="panel-hint">
            Base: <code>{API_BASE_URL || window.location.origin}</code> · verificação automática a cada 30s no topo da tela.
          </p>
        </Panel>

        <Panel
          title={selectedState ? `Cidades — ${selectedState.uf}` : 'Cidades'}
          icon={<MapPin size={16} />}
          actions={(
            <Link className="panel-link" to="/cidades" state={{ ibge: selectedIbge }}>
              Ver todas
            </Link>
          )}
        >
          {!selectedState ? (
            <div className="loading-box">Selecione um estado.</div>
          ) : citiesRes.loading && !citiesRes.data ? (
            <div className="loading-box">Carregando cidades…</div>
          ) : citiesRes.error ? (
            <div className="error-box">
              {citiesRes.error.message}
              <button type="button" className="error-retry" onClick={citiesRes.reload}>
                Tentar novamente
              </button>
            </div>
          ) : (
            <>
              <ul className="city-list">
                {(citiesRes.data?.dados ?? []).map((city) => (
                  <li key={city.ibge_code} className="city-item">
                    <span className="city-ibge">{city.ibge_code}</span>
                    <span className="city-name">{city.name}</span>
                  </li>
                ))}
              </ul>
              <p className="panel-hint">
                {citiesRes.data?.total ?? 0} cidades no total · exibindo as primeiras {CITIES_PANEL_SIZE}.
              </p>
            </>
          )}
        </Panel>

        <Panel title="Estados (UF)" icon={<MapIcon size={16} />}>
          {states.length === 0 ? (
            <div className="loading-box">Carregando estados…</div>
          ) : (
            <>
              <div className="uf-chips">
                {states.map((state) => (
                  <button
                    key={state.ibge_code}
                    type="button"
                    className={`uf-chip ${String(state.ibge_code) === selectedIbge ? 'active' : ''}`}
                    onClick={() => setRawIbge(String(state.ibge_code))}
                    title={state.name}
                  >
                    {state.uf}
                  </button>
                ))}
              </div>
              <p className="panel-hint">Clique numa UF para trocar o estado dos painéis.</p>
            </>
          )}
        </Panel>

        {rankings.map((ranking) => (
          <Panel key={ranking.title} title={ranking.title} icon={ranking.icon}>
            {topCitiesRes.loading && !topCitiesRes.data ? (
              <div className="loading-box">Carregando ranking…</div>
            ) : topCitiesRes.error ? (
              <div className="error-box">
                {topCitiesRes.error.message}
                <button type="button" className="error-retry" onClick={topCitiesRes.reload}>
                  Tentar novamente
                </button>
              </div>
            ) : (
              <ul className="city-list">
                {(ranking.cities ?? []).map((city, index) => (
                  <li key={city.ibge_code} className="city-item">
                    <span className="city-ibge">#{index + 1}</span>
                    <Link to={`/cidades/${city.ibge_code}`}>
                      {city.name} ({city.state.uf})
                    </Link>
                    <span className="service-meta" style={{ marginLeft: 'auto' }}>
                      {ranking.format(ranking.pick(city))}
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </Panel>
        ))}
      </div>

      <Panel title="Distribuição por faixa etária" icon={<PieChart size={16} />}>
        <div className="panel-controls">
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="age-region-filter">Região</label>
            <select
              id="age-region-filter"
              className="filter-select"
              value={ageRegion}
              onChange={(event) => handleAgeRegionChange(event.target.value)}
            >
              <option value="">Todas</option>
              {regionOptions.map((region) => (
                <option key={region} value={region}>{region}</option>
              ))}
            </select>
          </div>
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="age-uf-filter">UF</label>
            <select
              id="age-uf-filter"
              className="filter-select"
              value={ageIbge}
              onChange={(event) => setAgeIbge(event.target.value)}
            >
              <option value="">Todas</option>
              {ageStates.map((state) => (
                <option key={state.ibge_code} value={String(state.ibge_code)}>
                  {`${state.uf} — ${state.name}`}
                </option>
              ))}
            </select>
          </div>
        </div>
        {ageRes.loading && !age ? (
          <div className="loading-box">Carregando faixa etária…</div>
        ) : ageRes.error ? (
          <div className="error-box">
            {ageRes.error.message}
            <button type="button" className="error-retry" onClick={ageRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : ageGroups.length === 0 ? (
          <div className="loading-box">Sem dado de faixa etária.</div>
        ) : (
          <>
            <ul className="age-list">
              {ageGroups.map((group) => (
                <li key={group.faixa} className="age-item">
                  <span className="age-name" title={group.faixa}>{group.faixa}</span>
                  <span className="age-bar">
                    <span
                      className="age-bar-fill"
                      style={{ width: `${(group.populacao / maxAgePopulation) * 100}%` }}
                    />
                  </span>
                  <span className="age-count">
                    {formatInteger(group.populacao)} · {formatPercent((group.populacao / age.total) * 100)}
                  </span>
                </li>
              ))}
            </ul>
            <p className="panel-hint">
              Total {formatInteger(age.total)} pessoas · {ageScope || 'todo o país'} · Censo {age.ano} · IBGE/SIDRA 9514
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
    </div>
  );
}
