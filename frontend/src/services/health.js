import { api } from './api';

/**
 * Sonda um endpoint de health medindo a latência.
 * Falha individual não derruba a chamada: vira { ok: false, message }.
 */
async function probe(path, signal) {
  const startedAt = performance.now();

  try {
    await api.get(path, { signal });
    return {
      path,
      ok: true,
      latencyMs: Math.round(performance.now() - startedAt),
      message: null,
    };
  } catch (error) {
    if (error.name === 'AbortError') {
      throw error;
    }
    return { path, ok: false, latencyMs: null, message: error.message };
  }
}

/**
 * GET /health  -> { status: "OK" }
 * GET /health/db -> { status: "OK" } | 500 (texto puro)
 */
export async function getHealth({ signal } = {}) {
  const [app, database] = await Promise.all([
    probe('/health', signal),
    probe('/health/db', signal),
  ]);

  return { app, database };
}
