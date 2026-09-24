import { useCallback, useEffect, useRef, useState } from 'react';

export function useApiResource(fetcher, deps = [], { enabled = true } = {}) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [version, setVersion] = useState(0);

  const fetcherRef = useRef(fetcher);
  useEffect(() => {
    fetcherRef.current = fetcher;
  });

  const reload = useCallback(() => setVersion((current) => current + 1), []);

  const depsKey = JSON.stringify(deps);
  const fetchKey = `${depsKey}|${version}|${enabled}`;
  const [prevFetchKey, setPrevFetchKey] = useState(fetchKey);

  if (prevFetchKey !== fetchKey) {
    setPrevFetchKey(fetchKey);
    if (enabled) {
      setLoading(true);
      setError(null);
    } else {
      setData(null);
      setLoading(false);
      setError(null);
    }
  }

  useEffect(() => {
    if (!enabled) return;

    const controller = new AbortController();
    let alive = true;

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
  }, [depsKey, version, enabled]);

  return {
    data: enabled ? data : null,
    loading: enabled ? loading : false,
    error: enabled ? error : null,
    reload,
  };
}
