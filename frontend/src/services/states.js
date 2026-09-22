import { api } from './api';

/**
 * GET /api/v1/states
 *
 * @returns {Promise<Array<{ id: number, ibge_code: number, name: string, uf: string, region: string }>>}
 */
export function getStates({ signal } = {}) {
  return api.get('/api/v1/states', { signal });
}
