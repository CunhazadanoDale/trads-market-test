import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import States from './States';
import { getStates } from '../../services/states';

vi.mock('../../services/states', () => ({
  getStates: vi.fn(),
}));

const states = [
  { name: 'São Paulo', uf: 'SP', region: 'Sudeste', ibge_code: 35 },
  { name: 'Rio de Janeiro', uf: 'RJ', region: 'Sudeste', ibge_code: 33 },
  { name: 'Bahia', uf: 'BA', region: 'Nordeste', ibge_code: 29 },
  { name: 'Ceará', uf: 'CE', region: 'Nordeste', ibge_code: 23 },
];

beforeEach(() => {
  getStates.mockReset();
  getStates.mockResolvedValue(states);
});

describe('States', () => {
  it('carrega os estados pelo serviço', async () => {
    render(<States />);

    await waitFor(() => expect(screen.getByText('São Paulo')).toBeInTheDocument());

    expect(getStates).toHaveBeenCalledTimes(1);
    expect(screen.getByText('Bahia')).toBeInTheDocument();
    expect(screen.getByText('Ceará')).toBeInTheDocument();
  });

  it('filtra por região com serviço mockado', async () => {
    render(<States />);

    await waitFor(() => expect(screen.getByText('São Paulo')).toBeInTheDocument());

    fireEvent.change(screen.getByRole('combobox'), { target: { value: 'Nordeste' } });
    fireEvent.click(screen.getByRole('button', { name: 'Aplicar' }));

    await waitFor(() => expect(screen.queryByText('São Paulo')).not.toBeInTheDocument());

    expect(screen.getByText('Bahia')).toBeInTheDocument();
    expect(screen.getByText('Ceará')).toBeInTheDocument();
    expect(screen.queryByText('Rio de Janeiro')).not.toBeInTheDocument();
    expect(getStates).toHaveBeenCalledTimes(1);
  });

  it('volta a listar todos os estados ao limpar os filtros', async () => {
    render(<States />);

    await waitFor(() => expect(screen.getByText('São Paulo')).toBeInTheDocument());

    fireEvent.change(screen.getByRole('combobox'), { target: { value: 'Nordeste' } });
    fireEvent.click(screen.getByRole('button', { name: 'Aplicar' }));
    await waitFor(() => expect(screen.queryByText('São Paulo')).not.toBeInTheDocument());

    fireEvent.click(screen.getByRole('button', { name: 'Limpar' }));

    await waitFor(() => expect(screen.getByText('São Paulo')).toBeInTheDocument());
    expect(screen.getByText('Bahia')).toBeInTheDocument();
  });
});
