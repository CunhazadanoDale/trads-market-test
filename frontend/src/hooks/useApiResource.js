import { useCallback, useEffect, useRef, useState } from 'react';

/**
 * Executa uma chamada de serviço (fetcher) com loading, erro, cancelamento
 * e recarga manual.
 *
 * @param {(ctx: { signal: AbortSignal }) => Promise<any>} fetcher
 * @param {any[]} deps dependências estáveis que disparam nova chamada
 */
export function useApiResource(fetcher, deps = []) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [version, setVersion] = useState(0);

  const fetcherRef = useRef(fetcher);
  fetcherRef.current = fetcher;

  const reload = useCallback(() => setVersion((current) => current + 1), []);

  useEffect(() => {
    const controller = new AbortController();
    let alive = true;

    setLoading(true);
    setError(null);

    Promise.resolve(fetcherRef.current({ signal: controller.signal }))
      .then((result) => {
        if (!alive) return;
        setData(result);
        setLoading(false);
      })
      .catch((err) => {
        if (!alive || err?.name === 'AbortError') return;
        setError(err);
        setLoading(false);
      });

    return () => {
      alive = false;
      controller.abort();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [...deps, version]);

  return { data, loading, error, reload };
}
