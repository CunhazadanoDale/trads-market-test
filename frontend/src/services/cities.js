import { api } from './api';

export const DEFAULT_PAGE_SIZE = 20;
export const PAGE_SIZE_OPTIONS = [20, 50, 100];

export async function getCities(
  stateIbgeCode,
  { page = 1, pageSize = DEFAULT_PAGE_SIZE, nome = '', ordenar = '', ordem = '', signal } = {},
) {
  if (!stateIbgeCode) {
    return { dados: [], pagina: 1, tamanho: pageSize, total: 0 };
  }

  const params = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize),
  });

  if (nome) params.set('nome', nome);
  if (ordenar) params.set('ordenar', ordenar);
  if (ordem) params.set('ordem', ordem);

  const response = await api.get(
    `/api/v1/states/${stateIbgeCode}/cities?${params.toString()}`,
    { signal },
  );

  return {
    dados: response?.dados ?? [],
    pagina: response?.pagina ?? page,
    tamanho: response?.tamanho ?? pageSize,
    total: response?.total ?? 0,
  };
}

export function getCityDetail(ibgeCode, { signal } = {}) {
  return api.get(`/api/v1/cities/${ibgeCode}`, { signal });
}
