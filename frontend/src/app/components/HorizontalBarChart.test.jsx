import React from 'react';
import { fireEvent, render, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { HorizontalBarChart } from './HorizontalBarChart';

const CHART_WIDTH = 720;
const ROW_HEIGHT = 36;

const groups = [
  { faixa: '0 a 4 anos', populacao: 8100000 },
  { faixa: '5 a 9 anos', populacao: 9200000 },
  { faixa: '10 a 14 anos', populacao: 9600000 },
  { faixa: '15 a 19 anos', populacao: 9400000 },
];

class ResizeObserverStub {
  constructor(callback) {
    this.callback = callback;
  }

  observe() {
    this.callback([
      {
        contentRect: {
          width: CHART_WIDTH,
          height: groups.length * ROW_HEIGHT,
        },
      },
    ]);
  }

  unobserve() {}

  disconnect() {}
}

const tooltipContent = (group) => ({
  title: group.faixa,
  rows: [{ label: 'População', value: String(group.populacao) }],
});

const readTooltip = () => ({
  title: document.querySelector('.chart-tooltip-title')?.textContent ?? null,
  value: document.querySelector('.chart-tooltip-value')?.textContent ?? null,
});

describe('HorizontalBarChart', () => {
  beforeEach(() => {
    vi.stubGlobal('ResizeObserver', ResizeObserverStub);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('mostra no tooltip os dados da barra sob o cursor', async () => {
    const { container } = render(
      <HorizontalBarChart
        data={groups}
        labelKey="faixa"
        barKey="populacao"
        rowHeight={ROW_HEIGHT}
        tooltip={tooltipContent}
      />,
    );

    await waitFor(() =>
      expect(container.querySelectorAll('.recharts-bar-rectangle')).toHaveLength(groups.length),
    );

    const wrapper = container.querySelector('.recharts-wrapper');

    for (let index = 0; index < groups.length; index += 1) {
      fireEvent.mouseMove(wrapper, {
        clientX: CHART_WIDTH / 2,
        clientY: (index + 0.5) * ROW_HEIGHT,
      });

      await waitFor(() =>
        expect(readTooltip()).toEqual({
          title: groups[index].faixa,
          value: String(groups[index].populacao),
        }),
      );
    }
  });
});
