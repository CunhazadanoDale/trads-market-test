import { api } from './api';

export function getStates({ signal } = {}) {
  return api.get('/api/v1/states', { signal });
}
