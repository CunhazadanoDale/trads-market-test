import React from 'react';
import './ChartTooltip.css';

export function ChartTooltip({ title, rows }) {
  if (!rows || rows.length === 0) return null;

  return (
    <div className="chart-tooltip">
      {title ? <p className="chart-tooltip-title">{title}</p> : null}
      <ul className="chart-tooltip-list">
        {rows.map((row) => (
          <li key={row.label} className="chart-tooltip-row">
            <span className="chart-tooltip-label">{row.label}</span>
            <span className="chart-tooltip-value">{row.value}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
