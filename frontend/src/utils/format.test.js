import { describe, expect, it } from 'vitest';
import { formatIncome, formatInteger, formatPercent, formatPopulation } from './format';

const NO_BREAK_SPACE = '\u00A0';

describe('formatPopulation', () => {
  it('formata um valor de população', () => {
    expect(formatPopulation({ value: 89589000 })).toBe('89.589.000');
  });

  it('devolve traço quando o indicador é nulo', () => {
    expect(formatPopulation(null)).toBe('—');
  });

  it('formata zero', () => {
    expect(formatPopulation({ value: 0 })).toBe('0');
  });
});

describe('formatIncome', () => {
  it('formata um valor de renda em reais', () => {
    expect(formatIncome({ value: 3480.5 })).toBe(`R$${NO_BREAK_SPACE}3.480,50`);
  });

  it('devolve traço quando o indicador é nulo', () => {
    expect(formatIncome(null)).toBe('—');
  });

  it('formata zero', () => {
    expect(formatIncome({ value: 0 })).toBe(`R$${NO_BREAK_SPACE}0,00`);
  });
});

describe('formatInteger', () => {
  it('formata um inteiro com separador de milhar', () => {
    expect(formatInteger(1668)).toBe('1.668');
  });

  it('devolve traço quando o valor é nulo', () => {
    expect(formatInteger(null)).toBe('—');
  });

  it('formata zero', () => {
    expect(formatInteger(0)).toBe('0');
  });
});

describe('formatPercent', () => {
  it('formata uma porcentagem com uma casa decimal', () => {
    expect(formatPercent(12.34)).toBe('12,3%');
  });

  it('devolve traço quando o valor é nulo', () => {
    expect(formatPercent(null)).toBe('—');
  });

  it('formata zero', () => {
    expect(formatPercent(0)).toBe('0%');
  });
});
