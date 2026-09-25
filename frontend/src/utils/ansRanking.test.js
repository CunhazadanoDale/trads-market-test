import { describe, expect, it } from 'vitest';
import { ANS_SORT_FIELDS, buildANSRanking } from './ansRanking';

const municipios = [
  { name: 'Beta', penetracao: 30, beneficiarios: 300, populacao: 1000 },
  { name: 'Gama', penetracao: 20, beneficiarios: 800, populacao: 4000 },
  { name: 'Alpha', penetracao: 10, beneficiarios: 500, populacao: 5000 },
];

describe('buildANSRanking', () => {
  it('numera a ordem que vem do servidor', () => {
    expect(buildANSRanking(municipios).map((m) => m.posicao)).toEqual([1, 2, 3]);
    expect(buildANSRanking(municipios).map((m) => m.name)).toEqual(['Beta', 'Gama', 'Alpha']);
  });

  it('mantém os campos originais da linha', () => {
    const [primeira] = buildANSRanking(municipios);
    expect(primeira).toEqual({
      name: 'Beta',
      penetracao: 30,
      beneficiarios: 300,
      populacao: 1000,
      posicao: 1,
    });
  });

  it('não reordena nem corta o que o servidor devolveu', () => {
    const muitos = Array.from({ length: 100 }, (_, i) => ({ name: `Cidade ${i}`, penetracao: i }));
    const ranking = buildANSRanking(muitos);
    expect(ranking).toHaveLength(100);
    expect(ranking[0].name).toBe('Cidade 0');
    expect(ranking[99].posicao).toBe(100);
  });

  it('devolve lista vazia para entrada ausente', () => {
    expect(buildANSRanking(undefined)).toEqual([]);
    expect(buildANSRanking(null)).toEqual([]);
  });
});

describe('ANS_SORT_FIELDS', () => {
  it('expõe os três campos de ordenação', () => {
    expect(ANS_SORT_FIELDS.map((option) => option.value)).toEqual([
      'penetracao',
      'beneficiarios',
      'populacao',
    ]);
  });
});
