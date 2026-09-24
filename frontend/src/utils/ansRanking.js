export const ANS_SORT_FIELDS = [
  { value: 'penetracao', label: 'Penetração' },
  { value: 'beneficiarios', label: 'Beneficiários' },
  { value: 'populacao', label: 'População' },
];

export function buildANSRanking(municipios, { sortBy = 'penetracao', limit = 10 } = {}) {
  const rows = Array.isArray(municipios) ? municipios : [];
  const field = ANS_SORT_FIELDS.some((option) => option.value === sortBy)
    ? sortBy
    : 'penetracao';

  return [...rows]
    .sort((a, b) => {
      const diff = (b?.[field] ?? 0) - (a?.[field] ?? 0);
      if (diff !== 0) return diff;
      return String(a?.name ?? '').localeCompare(String(b?.name ?? ''), 'pt-BR');
    })
    .slice(0, limit);
}
