import React, { useMemo, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ChevronLeft, ChevronRight, RefreshCw } from 'lucide-react';
import '../shared.css';
import { DataGrid } from '../../app/components/DataGrid';
import { FilterGroup, FilterPanel } from '../../app/components/FilterPanel';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { DEFAULT_PAGE_SIZE, PAGE_SIZE_OPTIONS, getCities } from '../../services/cities';
import { formatGDP, formatIncome, formatPopulation } from '../../utils/format';

const EMPTY_CITY_FILTERS = { nome: '', ordenar: '', ordem: 'asc' };

const ORDER_OPTIONS = [
  { value: '', label: 'Nome' },
  { value: 'populacao', label: 'População' },
  { value: 'renda', label: 'Renda' },
  { value: 'pib', label: 'PIB' },
];

const DIRECTION_OPTIONS = [
  { value: 'asc', label: 'Crescente' },
  { value: 'desc', label: 'Decrescente' },
];

export default function Cities() {
  const location = useLocation();
  const navigate = useNavigate();

  const statesRes = useApiResource(getStates, []);
  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);

  // Pode chegar pré-selecionado do Dashboard ("Ver todas");
  // sem escolha explícita, vale o primeiro estado da lista (derivado)
  const [rawSelected, setRawSelected] = useState(() => String(location.state?.ibge ?? ''));
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [draftFilters, setDraftFilters] = useState(EMPTY_CITY_FILTERS);
  const [filters, setFilters] = useState(EMPTY_CITY_FILTERS);

  const selected = rawSelected || (states.length > 0 ? String(states[0].ibge_code) : '');

  const selectedState = useMemo(
    () => states.find((state) => String(state.ibge_code) === selected) ?? null,
    [states, selected],
  );

  const handleStateChange = (ibge) => {
    setRawSelected(ibge);
    setPage(1);
  };

  const handlePageSizeChange = (size) => {
    setPageSize(size);
    setPage(1);
  };

  const citiesRes = useApiResource(
    ({ signal }) => getCities(selected, { page, pageSize, signal, ...filters }),
    [selected, page, pageSize, filters],
  );

  const cities = citiesRes.data;
  const total = cities?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const currentPage = cities?.pagina ?? page;

  const rows = useMemo(
    () => (cities?.dados ?? []).map((city) => ({ ...city, uf: selectedState?.uf ?? '—' })),
    [cities, selectedState],
  );

  const columns = [
    {
      label: 'Código IBGE',
      field: 'ibge_code',
      width: '14%',
      render: (value) => String(value),
    },
    { label: 'Nome', field: 'name', width: '26%' },
    { label: 'Estado', field: 'uf', width: '6%' },
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

  const handleApply = () => {
    setPage(1);
    setFilters(draftFilters);
    citiesRes.reload();
  };

  const handleClear = () => {
    setDraftFilters(EMPTY_CITY_FILTERS);
    setFilters(EMPTY_CITY_FILTERS);
    handleStateChange(states.length > 0 ? String(states[0].ibge_code) : '');
    handlePageSizeChange(DEFAULT_PAGE_SIZE);
  };

  const renderContent = () => {
    if (states.length === 0) {
      if (statesRes.error) {
        return (
          <div className="error-box">
            {statesRes.error.message}
            <button type="button" className="error-retry" onClick={statesRes.reload}>
              Tentar novamente
            </button>
          </div>
        );
      }
      return <div className="loading-box">Carregando estados…</div>;
    }

    if (!selected) {
      return <div className="loading-box">Selecione um estado para ver as cidades.</div>;
    }

    if (citiesRes.loading && !cities) {
      return <div className="loading-box">Carregando cidades…</div>;
    }

    if (citiesRes.error) {
      return (
        <div className="error-box">
          {citiesRes.error.message}
          <button type="button" className="error-retry" onClick={citiesRes.reload}>
            Tentar novamente
          </button>
        </div>
      );
    }

    return (
      <DataGrid
        columns={columns}
        data={rows}
        onRowClick={(row) => navigate(`/cidades/${row.ibge_code}`, { state: { ibge: selected } })}
      />
    );
  };

  return (
    <div className="module-container">
      <div className="module-main">
        <div className="module-toolbar">
          <div className="module-toolbar-title">
            Cidades
            {selectedState ? ` — ${selectedState.uf}` : ''}
            {cities ? ` (${total})` : ''}
          </div>
          <div className="module-toolbar-actions">
            <button
              type="button"
              className="btn btn--secondary"
              onClick={citiesRes.reload}
              disabled={citiesRes.loading}
              title="Recarregar da API"
            >
              <RefreshCw size={14} className={citiesRes.loading ? 'spin' : undefined} />
              Atualizar
            </button>
          </div>
        </div>

        <div className="module-content">{renderContent()}</div>

        {cities && !citiesRes.error && (
          <div className="pagination-bar">
            <span className="pagination-info">
              Página {currentPage} de {totalPages} · {total} registros
            </span>
            <div className="pagination-actions">
              <button
                type="button"
                className="btn btn--secondary"
                onClick={() => setPage((current) => Math.max(1, current - 1))}
                disabled={currentPage <= 1 || citiesRes.loading}
              >
                <ChevronLeft size={14} />
                Anterior
              </button>
              <button
                type="button"
                className="btn btn--secondary"
                onClick={() => setPage((current) => current + 1)}
                disabled={currentPage >= totalPages || citiesRes.loading}
              >
                Próxima
                <ChevronRight size={14} />
              </button>
            </div>
          </div>
        )}
      </div>

      <FilterPanel onApply={handleApply} onClear={handleClear}>
        <FilterGroup label="Estado">
          <select
            className="filter-select"
            value={selected}
            onChange={(event) => handleStateChange(event.target.value)}
          >
            {states.length === 0 && <option value="">Carregando…</option>}
            {states.map((state) => (
              <option key={state.ibge_code} value={String(state.ibge_code)}>
                {`${state.uf} — ${state.name}`}
              </option>
            ))}
          </select>
        </FilterGroup>
        <FilterGroup label="Nome da cidade">
          <input
            type="text"
            className="filter-input"
            placeholder="Ex.: são"
            value={draftFilters.nome}
            onChange={(event) => setDraftFilters({ ...draftFilters, nome: event.target.value })}
          />
        </FilterGroup>
        <FilterGroup label="Ordenar por">
          <select
            className="filter-select"
            value={draftFilters.ordenar}
            onChange={(event) => setDraftFilters({ ...draftFilters, ordenar: event.target.value })}
          >
            {ORDER_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
        </FilterGroup>
        <FilterGroup label="Ordem">
          <select
            className="filter-select"
            value={draftFilters.ordem}
            onChange={(event) => setDraftFilters({ ...draftFilters, ordem: event.target.value })}
          >
            {DIRECTION_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>{option.label}</option>
            ))}
          </select>
        </FilterGroup>
        <FilterGroup label="Itens por página">
          <select
            className="filter-select"
            value={pageSize}
            onChange={(event) => handlePageSizeChange(Number(event.target.value))}
          >
            {PAGE_SIZE_OPTIONS.map((size) => (
              <option key={size} value={size}>{size}</option>
            ))}
          </select>
        </FilterGroup>
      </FilterPanel>
    </div>
  );
}
