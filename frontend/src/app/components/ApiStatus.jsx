import { useEffect } from 'react';
import { useApiResource } from '../../hooks/useApiResource';
import { getHealth } from '../../services/health';
import { API_BASE_URL } from '../../services/api';
import './ApiStatus.css';

const REFRESH_INTERVAL_MS = 30000;

export function ApiStatus() {
  const { data, loading, reload } = useApiResource(getHealth, []);

  useEffect(() => {
    const id = setInterval(reload, REFRESH_INTERVAL_MS);
    return () => clearInterval(id);
  }, [reload]);

  const app = data?.app;
  const database = data?.database;
  const hasData = Boolean(data);

  let label = 'Verificando…';
  let tone = '';

  if (hasData) {
    if (app.ok && database.ok) {
      label = 'API OK';
      tone = 'ok';
    } else if (!app.ok) {
      label = 'API Offline';
      tone = 'down';
    } else {
      label = 'Banco Offline';
      tone = 'down';
    }
  } else if (loading) {
    tone = '';
  }

  const base = API_BASE_URL || window.location.origin;

  return (
    <button
      type="button"
      className={`api-status ${tone}`.trim()}
      onClick={reload}
      title={`API: ${base} · /health e /health/db · atualiza a cada 30s (clique para verificar agora)`}
    >
      <span className="api-status-dot" />
      {label}
    </button>
  );
}
