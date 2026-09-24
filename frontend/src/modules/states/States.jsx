import React, { useMemo, useState } from 'react';
import { RefreshCw } from 'lucide-react';
import '../shared.css';
import { DataGrid } from '../../app/components/DataGrid';
import { FilterGroup, FilterPanel } from '../../app/components/FilterPanel';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';

const ALL_REGIONS = 'Todas';
const EMPTY_FILTERS = { search: '', region: ALL_REGIONS };

export default function States() {
  const { data: states, loading, error, reload } = useApiResource(getStates, []);

  // Rascunho editado nos filtros + filtros efetivamente aplicados
  const [draft, setDraft] = useState(EMPTY_FILTERS);
  const [filters, setFilters] = useState(EMPTY_FILTERS);

  const regionOptions = useMemo(
    () => [ALL_REGIONS, ...new Set((states ?? []).map((state) => state.region))],
    [states],
  );

  const rows = useMemo(() => {
    const search = filters.search.trim().toLowerCase();

    return (states ?? []).filter((state) => {
      const matchesRegion = filters.region === ALL_REGIONS || state.region === filters.region;
      const matchesSearch = !search
        || state.name.toLowerCase().includes(search)
        || state.uf.toLowerCase().includes(search)
        || String(state.ibge_code).includes(search);

      return matchesRegion && matchesSearch;
    });
  }, [states, filters]);

  const columns = [
    { label: 'Nome', field: 'name', width: '35%' },
    { label: 'UF', field: 'uf', width: '10%' },
    { label: 'Região', field: 'region', width: '25%' },
    {
      label: 'Código IBGE',
      field: 'ibge_code',
      width: '30%',
      render: (value) => String(value),
    },
  ];

  const handleApply = () => setFilters(draft);

  const handleClear = () => {
    setDraft(EMPTY_FILTERS);
    setFilters(EMPTY_FILTERS);
  };

  return (
    <div className="module-container">
      <div className="module-main">
        <div className="module-toolbar">
          <div className="module-toolbar-title">
            Estados {loading ? '' : `(${rows.length}${rows.length !== (states ?? []).length ? ` de ${states.length}` : ''})`}
          </div>
          <div className="module-toolbar-actions">
            <button
              type="button"
              className="btn btn--secondary"
              onClick={reload}
              disabled={loading}
              title="Recarregar da API"
            >
              <RefreshCw size={14} className={loading ? 'spin' : undefined} />
              Atualizar
            </button>
          </div>
        </div>

        <div className="module-content">
          {loading && !states ? (
            <div className="loading-box">Carregando estados…</div>
          ) : error ? (
            <div className="error-box">
              {error.message}
              <button type="button" className="error-retry" onClick={reload}>
                Tentar novamente
              </button>
            </div>
          ) : (
            <DataGrid columns={columns} data={rows} />
          )}
        </div>
      </div>

      <FilterPanel onApply={handleApply} onClear={handleClear}>
        <FilterGroup label="Pesquisa">
          <input
            type="text"
            className="filter-input"
            placeholder="Nome, UF ou código IBGE"
            value={draft.search}
            onChange={(event) => setDraft({ ...draft, search: event.target.value })}
          />
        </FilterGroup>
        <FilterGroup label="Região">
          <select
            className="filter-select"
            value={draft.region}
            onChange={(event) => setDraft({ ...draft, region: event.target.value })}
          >
            {regionOptions.map((region) => (
              <option key={region} value={region}>{region}</option>
            ))}
          </select>
        </FilterGroup>
      </FilterPanel>
    </div>
  );
}
