import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { DataGrid } from './DataGrid';

const columns = [
  { label: 'Nome', field: 'name', width: '60%' },
  { label: 'UF', field: 'uf', width: '40%' },
];

const rows = [
  { name: 'São Paulo', uf: 'SP' },
  { name: 'Rio de Janeiro', uf: 'RJ' },
];

describe('DataGrid', () => {
  it('renderiza as colunas informadas', () => {
    render(<DataGrid columns={columns} data={rows} />);

    expect(screen.getByText('Nome')).toBeInTheDocument();
    expect(screen.getByText('UF')).toBeInTheDocument();
  });

  it('renderiza uma linha para cada registro', () => {
    render(<DataGrid columns={columns} data={rows} />);

    expect(screen.getAllByRole('row')).toHaveLength(3);
    expect(screen.getByText('São Paulo')).toBeInTheDocument();
    expect(screen.getByText('Rio de Janeiro')).toBeInTheDocument();
    expect(screen.getByText('RJ')).toBeInTheDocument();
  });

  it('dispara onRowClick com a linha clicada', () => {
    const onRowClick = vi.fn();
    render(<DataGrid columns={columns} data={rows} onRowClick={onRowClick} />);

    fireEvent.click(screen.getByText('Rio de Janeiro'));

    expect(onRowClick).toHaveBeenCalledTimes(1);
    expect(onRowClick).toHaveBeenCalledWith(rows[1]);
  });

  it('mostra mensagem quando não há registros', () => {
    render(<DataGrid columns={columns} data={[]} />);

    expect(screen.getByText('Nenhum registro encontrado.')).toBeInTheDocument();
  });
});
