import { describe, expect, it } from 'vitest';
import { ANS_SORT_FIELDS, buildANSRanking } from './ansRanking';

const municipios = [
  { name: 'Alpha', penetracao: 10, beneficiarios: 500, populacao: 5000 },
  { name: 'Beta', penetracao: 30, beneficiarios: 300, populacao: 1000 },
  { name: 'Gama', penetracao: 20, beneficiarios: 800, populacao: 4000 },
];

describe('buildANSRanking', () => {
  it('ordena por penetração decrescente por padrão', () => {
    expect(buildANSRanking(municipios).map((m) => m.name)).toEqual(['Beta', 'Gama', 'Alpha']);
  });

  it('ordena por beneficiários quando pedido', () => {
    const ranking = buildANSRanking(municipios, { sortBy: 'beneficiarios' });
    expect(ranking.map((m) => m.name)).toEqual(['Gama', 'Alpha', 'Beta']);
  });

  it('ordena por população quando pedido', () => {
    const ranking = buildANSRanking(municipios, { sortBy: 'populacao' });
    expect(ranking.map((m) => m.name)).toEqual(['Alpha', 'Gama', 'Beta']);
  });

  it('limita o ranking ao top 10 por padrão', () => {
    const muitos = Array.from({ length: 25 }, (_, i) => ({ name: `Cidade ${i}`, penetracao: i }));
    const ranking = buildANSRanking(muitos);
    expect(ranking).toHaveLength(10);
    expect(ranking[0].name).toBe('Cidade 24');
  });

  it('respeita o limite informado', () => {
    expect(buildANSRanking(municipios, { limit: 2 })).toHaveLength(2);
  });

  it('desempata por nome', () => {
    const empatados = [
      { name: 'Zebra', penetracao: 15 },
      { name: 'Avião', penetracao: 15 },
    ];
    expect(buildANSRanking(empatados).map((m) => m.name)).toEqual(['Avião', 'Zebra']);
  });

  it('usa penetração quando o campo de ordenação é inválido', () => {
    const ranking = buildANSRanking(municipios, { sortBy: 'inexistente' });
    expect(ranking[0].name).toBe('Beta');
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
