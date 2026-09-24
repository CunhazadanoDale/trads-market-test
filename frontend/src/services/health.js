import { api } from './api';

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

export async function getHealth({ signal } = {}) {
  const [app, database] = await Promise.all([
    probe('/health', signal),
    probe('/health/db', signal),
  ]);

  return { app, database };
}
