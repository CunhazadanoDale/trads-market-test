import React from 'react';
import './DataGrid.css';

export function DataGrid({ columns, data, onRowClick }) {
  if (!data || data.length === 0) {
    return (
      <div className="datagrid-container">
        <div className="datagrid-empty">Nenhum registro encontrado.</div>
      </div>
    );
  }

  return (
    <div className="datagrid-container">
      <table className="datagrid">
        <thead>
          <tr>
            {columns.map((col, idx) => (
              <th key={idx} style={{ width: col.width || 'auto' }}>
                {col.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr
              key={rowIndex}
              style={onRowClick ? { cursor: 'pointer' } : undefined}
              onClick={() => onRowClick?.(row)}
            >
              {columns.map((col, colIndex) => (
                <td key={colIndex}>
                  {col.render ? col.render(row[col.field], row) : row[col.field]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
