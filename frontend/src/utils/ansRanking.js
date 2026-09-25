export const ANS_SORT_FIELDS = [
  { value: 'penetracao', label: 'Penetração' },
  { value: 'beneficiarios', label: 'Beneficiários' },
  { value: 'populacao', label: 'População' },
];

export function buildANSRanking(municipios) {
  const rows = Array.isArray(municipios) ? municipios : [];

  return rows.map((row, index) => ({
    ...row,
    posicao: index + 1,
  }));
}
