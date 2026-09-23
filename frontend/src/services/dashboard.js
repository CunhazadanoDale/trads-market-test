import { api } from './api';

export function getNationalMetrics({ signal } = {}) {
  return api.get('/api/v1/dashboard/national', { signal });
}

export function getStateMetrics({ regiao = '', signal } = {}) {
  const params = new URLSearchParams();

  if (regiao) params.set('regiao', regiao);

  const query = params.toString();

  return api.get(`/api/v1/dashboard/states${query ? `?${query}` : ''}`, { signal });
}

export function getTopCities({ signal } = {}) {
  return api.get('/api/v1/dashboard/top-cities', { signal });
}

export function getAgeDistribution({ regiao = '', ibge = '', signal } = {}) {
  const params = new URLSearchParams();

  if (regiao) params.set('regiao', regiao);
  if (ibge) params.set('ibge', ibge);

  const query = params.toString();

  return api.get(`/api/v1/dashboard/age${query ? `?${query}` : ''}`, { signal });
}
