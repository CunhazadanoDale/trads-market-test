import React, { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Activity, BarChart2, ChevronRight, Globe, Map as MapIcon, MapPin, RefreshCw, Users, Wallet } from 'lucide-react';
import './Dashboard.css';
import { Panel } from '../../app/components/Panel';
import { StatusBadge } from '../../app/components/StatusBadge';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { getCities } from '../../services/cities';
import { getHealth } from '../../services/health';
import { getNationalMetrics, getTopCities } from '../../services/dashboard';
import { API_BASE_URL } from '../../services/api';
import { formatGDP, formatIncome, formatInteger, formatPopulation } from '../../utils/format';

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
  const [rawIbge, setRawIbge] = useState('');

  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);
  const selectedIbge = rawIbge || (states.length > 0 ? String(states[0].ibge_code) : '');

  const citiesRes = useApiResource(
    ({ signal }) => getCities(selectedIbge, { page: 1, pageSize: CITIES_PANEL_SIZE, signal }),
    [selectedIbge],
  );

  const nationalRes = useApiResource(getNationalMetrics, []);
  const topCitiesRes = useApiResource(getTopCities, []);

  const selectedState = useMemo(
    () => states.find((state) => String(state.ibge_code) === selectedIbge) ?? null,
    [states, selectedIbge],
  );

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
        : citiesRes.loading
          ? '…'
          : String(citiesRes.data?.total ?? 0),
      hint: 'total do estado selecionado',
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

  const handleRefresh = () => {
    statesRes.reload();
    healthRes.reload();
    citiesRes.reload();
    nationalRes.reload();
    topCitiesRes.reload();
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
          disabled={
            statesRes.loading ||
            healthRes.loading ||
            citiesRes.loading ||
            nationalRes.loading ||
            topCitiesRes.loading
          }
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

      <div className="synthesis-grid">
        <Link to="/mercados" className="synthesis-card">
          <span className="synthesis-card-icon"><Globe size={20} /></span>
          <span className="synthesis-card-body">
            <span className="synthesis-card-title">Onde vender?</span>
            <span className="synthesis-card-desc">
              Mercados por região e indicadores por UF: população, PIB e renda de cada mercado.
            </span>
          </span>
          <ChevronRight size={18} className="synthesis-card-arrow" />
        </Link>
        <Link to="/publico" className="synthesis-card">
          <span className="synthesis-card-icon"><Users size={20} /></span>
          <span className="synthesis-card-body">
            <span className="synthesis-card-title">Para quem vender?</span>
            <span className="synthesis-card-desc">
              Faixa etária por região e UF, com a renda das cidades onde cada público mora.
            </span>
          </span>
          <ChevronRight size={18} className="synthesis-card-arrow" />
        </Link>
      </div>

      <div className="summary-cards">
        {[...nationalCards, ...summaryCards].map((card) => (
          <div key={card.title} className="summary-card">
            <span className="summary-card-title">{card.title}</span>
            <span className="summary-card-value">{card.value}</span>
            <span className="summary-card-hint">{card.hint}</span>
          </div>
        ))}
      </div>

      <div className="dashboard-grid">
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
              <>
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
                <p className="panel-hint">
                  Ranking nacional — não muda com a UF selecionada.
                </p>
              </>
            )}
          </Panel>
        ))}

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
          ) : citiesRes.loading ? (
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
              <p className="panel-hint">Clique numa UF para trocar o painel Cidades.</p>
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
      </div>

    </div>
  );
}
