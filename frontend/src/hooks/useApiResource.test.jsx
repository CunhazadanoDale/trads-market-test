import { act, renderHook, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useApiResource } from './useApiResource';

describe('useApiResource', () => {
  it('carrega os dados com sucesso', async () => {
    const fetcher = vi.fn().mockResolvedValue({ id: 7 });
    const { result } = renderHook(() => useApiResource(fetcher, []));

    expect(result.current.loading).toBe(true);
    expect(result.current.data).toBeNull();

    await waitFor(() => expect(result.current.data).toEqual({ id: 7 }));

    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(fetcher).toHaveBeenCalledTimes(1);
    expect(fetcher.mock.calls[0][0]).toHaveProperty('signal');
  });

  it('expõe o erro quando a busca falha', async () => {
    const fetcher = vi.fn().mockRejectedValue(new Error('Falhou na rede'));
    const { result } = renderHook(() => useApiResource(fetcher, []));

    await waitFor(() => expect(result.current.error).not.toBeNull());

    expect(result.current.error.message).toBe('Falhou na rede');
    expect(result.current.data).toBeNull();
    expect(result.current.loading).toBe(false);
  });

  it('aborta a busca ao desmontar', async () => {
    let captured;
    const fetcher = vi.fn(({ signal }) => {
      captured = signal;
      return new Promise(() => {});
    });

    const { unmount } = renderHook(() => useApiResource(fetcher, []));

    await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));
    expect(captured.aborted).toBe(false);

    unmount();

    expect(captured.aborted).toBe(true);
  });

  it('recarrega quando reload é chamado', async () => {
    const fetcher = vi.fn().mockResolvedValueOnce('primeira').mockResolvedValueOnce('segunda');
    const { result } = renderHook(() => useApiResource(fetcher, []));

    await waitFor(() => expect(result.current.data).toBe('primeira'));

    act(() => {
      result.current.reload();
    });

    await waitFor(() => expect(result.current.data).toBe('segunda'));
    expect(fetcher).toHaveBeenCalledTimes(2);
  });
});
