import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Publico from './Publico';
import { getStates } from '../../services/states';
import { getAgeDistribution } from '../../services/dashboard';

vi.mock('../../services/states', () => ({
  getStates: vi.fn(),
}));

vi.mock('../../services/dashboard', () => ({
  getAgeDistribution: vi.fn(),
}));

const states = [
  { name: 'Rio Grande do Sul', uf: 'RS', region: 'Sul', ibge_code: 43 },
  { name: 'Santa Catarina', uf: 'SC', region: 'Sul', ibge_code: 42 },
];

const faixas = ['0 a 4 anos', '5 a 9 anos', '10 a 14 anos'];

const ageDistribution = {
  ano: 2022,
  total: 1000,
  grupos: faixas.map((faixa, index) => ({
    faixa,
    populacao: 100 + index,
    renda_media_cidades: 2000 + index,
  })),
};

const renderPublico = () =>
  render(
    <MemoryRouter>
      <Publico />
    </MemoryRouter>,
  );

const faixaSelect = () => screen.getByLabelText('Faixa');

const faixaValues = (select) => Array.from(select.options).map((option) => option.value);

beforeEach(() => {
  getStates.mockReset();
  getStates.mockResolvedValue(states);
  getAgeDistribution.mockReset();
  getAgeDistribution.mockResolvedValue(ageDistribution);
});

describe('Publico', () => {
  it('carrega a distribuição etária sem faixa selecionada', async () => {
    renderPublico();

    await waitFor(() => expect(getAgeDistribution).toHaveBeenCalledTimes(1));

    expect(getAgeDistribution).toHaveBeenCalledWith(
      expect.objectContaining({ regiao: '', ibge: '', faixa: '' }),
    );

    await waitFor(() => expect(faixaValues(faixaSelect())).toEqual(['', ...faixas]));
  });

  it('refaz a requisição ao trocar a faixa etária', async () => {
    renderPublico();

    await waitFor(() => expect(faixaValues(faixaSelect())).toEqual(['', ...faixas]));

    fireEvent.change(faixaSelect(), { target: { value: '5 a 9 anos' } });

    await waitFor(() => expect(getAgeDistribution).toHaveBeenCalledTimes(2));

    expect(getAgeDistribution).toHaveBeenLastCalledWith(
      expect.objectContaining({ regiao: '', ibge: '', faixa: '5 a 9 anos' }),
    );
  });

  it('mantém todas as faixas disponíveis depois de filtrar', async () => {
    renderPublico();

    await waitFor(() => expect(faixaValues(faixaSelect())).toEqual(['', ...faixas]));

    fireEvent.change(faixaSelect(), { target: { value: '10 a 14 anos' } });

    await waitFor(() => expect(getAgeDistribution).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(faixaSelect().value).toBe('10 a 14 anos'));

    expect(faixaValues(faixaSelect())).toEqual(['', ...faixas]);
  });

  it('mantém a faixa ao trocar a região e soma o recorte no rodapé', async () => {
    renderPublico();

    await waitFor(() => expect(faixaValues(faixaSelect())).toEqual(['', ...faixas]));

    fireEvent.change(screen.getByLabelText('Região'), { target: { value: 'Sul' } });
    fireEvent.change(faixaSelect(), { target: { value: '0 a 4 anos' } });

    await waitFor(() =>
      expect(getAgeDistribution).toHaveBeenLastCalledWith(
        expect.objectContaining({ regiao: 'Sul', ibge: '', faixa: '0 a 4 anos' }),
      ),
    );
  });
});
