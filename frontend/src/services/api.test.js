import { afterEach, describe, expect, it, vi } from 'vitest';
import { api, API_BASE_URL, ApiError } from './api';

function resposta(status, payload) {
  return {
    ok: status >= 200 && status < 300,
    status,
    headers: {
      get: (name) => (name.toLowerCase() === 'content-type' ? 'application/json' : null),
    },
    json: () => Promise.resolve(payload),
    text: () => Promise.resolve(JSON.stringify(payload)),
  };
}

describe('api', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('devolve o corpo quando a API responde 200 com JSON', async () => {
    const corpo = [{ id: 12, nome: 'Acre' }];
    const fetchMock = vi.fn().mockResolvedValue(resposta(200, corpo));
    vi.stubGlobal('fetch', fetchMock);

    const resultado = await api.get('/api/v1/states');

    expect(resultado).toEqual(corpo);
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(fetchMock.mock.calls[0][0]).toBe(`${API_BASE_URL}/api/v1/states`);
    expect(fetchMock.mock.calls[0][1]).toMatchObject({
      method: 'GET',
      headers: { Accept: 'application/json' },
    });
  });

  it('transforma erro 400 com payload em ApiError com status e mensagem', async () => {
    const payload = {
      error: {
        code: 'invalid_request',
        message: 'regiao deve ser uma das 5 regiões do Brasil',
      },
    };
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(resposta(400, payload)));

    const erro = await api.get('/api/v1/states?regiao=XXX').catch((error) => error);

    expect(erro).toBeInstanceOf(ApiError);
    expect(erro.name).toBe('ApiError');
    expect(erro.status).toBe(400);
    expect(erro.message).toBe(payload.error.message);
    expect(erro.payload).toEqual(payload);
  });

  it('usa mensagem padrão quando o erro da API não vem em JSON', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        headers: { get: () => 'text/plain' },
        json: () => Promise.resolve(null),
        text: () => Promise.resolve('   '),
      }),
    );

    const erro = await api.get('/api/v1/dashboard/national').catch((error) => error);

    expect(erro).toBeInstanceOf(ApiError);
    expect(erro.status).toBe(500);
    expect(erro.message).toBe('Erro 500 na API');
  });

  it('avisa que o backend pode estar fora quando o fetch rejeita', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('Failed to fetch')));

    const erro = await api.get('/api/v1/states').catch((error) => error);

    expect(erro).toBeInstanceOf(ApiError);
    expect(erro.status).toBe(0);
    expect(erro.message).toBe('Não foi possível conectar à API. O backend está rodando?');
  });

  it('deixa o AbortError passar sem virar ApiError', async () => {
    const abort = new Error('aborted');
    abort.name = 'AbortError';
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(abort));

    const erro = await api.get('/api/v1/states').catch((error) => error);

    expect(erro).toBe(abort);
    expect(erro).not.toBeInstanceOf(ApiError);
  });
});
