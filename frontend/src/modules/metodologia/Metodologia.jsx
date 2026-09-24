import React from 'react';
import { BookOpen, Calculator, Calendar, Code, Filter } from 'lucide-react';
import './Metodologia.css';
import { Panel } from '../../app/components/Panel';

const SOURCES = [
  {
    indicator: 'População das cidades',
    year: '2022',
    source: 'IBGE/SIDRA — agregado 4709, variável 93 (população residente), Censo Demográfico',
  },
  {
    indicator: 'Faixas etárias',
    year: '2022',
    source: 'IBGE/SIDRA — agregado 9514, variável 93, classificação de idade, Censo Demográfico',
  },
  {
    indicator: 'Renda média',
    year: '2022',
    source: 'IBGE/SIDRA — agregado 10295, variável 13431 (rendimento domiciliar mensal per capita), Censo Demográfico',
  },
  {
    indicator: 'PIB dos municípios',
    year: '2023',
    source: 'IBGE/SIDRA — agregado 5938, variável 37 (PIB a preços correntes, Mil R$)',
  },
  {
    indicator: 'Beneficiários de planos de saúde',
    year: 'conforme importação',
    source: 'ANS — dados abertos, agregado por município',
  },
  {
    indicator: 'Geografia (UF, região, município)',
    year: '—',
    source: 'IBGE — códigos oficiais de UF, região e município',
  },
];

export default function Metodologia() {
  return (
    <div className="metodologia-page">
      <header className="metodologia-intro">
        <h1 className="metodologia-title">
          <BookOpen size={20} />
          Metodologia e glossário
        </h1>
        <p className="metodologia-subtitle">
          De onde vem cada número, com que ano e com que fórmula — a base de confiança das
          páginas Mercados e Público e de tudo que o dashboard mostra.
        </p>
      </header>

      <Panel title="Indicadores, anos e fontes" icon={<Calendar size={16} />}>
        <table className="met-table">
          <thead>
            <tr>
              <th>Indicador</th>
              <th>Ano</th>
              <th>Fonte</th>
            </tr>
          </thead>
          <tbody>
            {SOURCES.map((row) => (
              <tr key={row.indicator}>
                <td className="met-table-indicator">{row.indicator}</td>
                <td className="met-table-year">{row.year}</td>
                <td>{row.source}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <p className="panel-hint">
          Todo dado entra pela importação das fontes oficiais (API do IBGE e dados
          abertos da ANS) e vive no banco próprio da Trads — a interface não consulta
          as fontes direto.
        </p>
      </Panel>

      <Panel title="Como cada número é calculado" icon={<Calculator size={16} />}>
        <dl className="metodologia-guide">
          <div className="metodologia-guide-item">
            <dt>População (cidade → UF → região)</dt>
            <dd>Soma simples: a UF soma suas cidades; a região soma suas UFs.</dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Renda média (país e por UF)</dt>
            <dd>
              O rendimento domiciliar mensal per capita de cada município entra
              multiplicado pela população da própria cidade — metrópole pesa mais que vila.
              <code className="met-formula">Σ(renda_cidade × população_cidade) ÷ Σ(população_cidade)</code>
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Renda média (por região)</dt>
            <dd>
              O painel de regiões pondera a renda de cada UF pela população da UF, no
              navegador.
              <code className="met-formula">Σ(renda_UF × população_UF) ÷ Σ(população_UF)</code>
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Renda média das cidades por faixa etária</dt>
            <dd>
              O entorno econômico do público: a renda da cidade entra multiplicada pela
              população da faixa que mora nela — não é a renda da pessoa.
              <code className="met-formula">Σ(renda_cidade × população_faixa_na_cidade) ÷ Σ(população_faixa)</code>
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>PIB (Mil R$)</dt>
            <dd>Soma do PIB municipal: a UF soma suas cidades; a região soma suas UFs.</dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Penetração ANS (%)</dt>
            <dd>
              Beneficiários de planos médico-hospitalares somados por município (sexo e
              faixa etária somados) divididos pela população do Censo 2022.
              <code className="met-formula">beneficiários_ans ÷ população_cidade × 100</code>
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>% por faixa etária</dt>
            <dd>
              Participação da faixa no recorte escolhido (país, região ou UF).
              <code className="met-formula">população_faixa ÷ total_do_recorte</code>
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Barras e ordenação</dt>
            <dd>
              As barras são proporcionais ao maior valor em tela; os mercados por região
              saem ordenados por população.
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Rankings de cidades</dt>
            <dd>As maiores cidades do país em PIB, renda e população.</dd>
          </div>
        </dl>
      </Panel>

      <Panel title="Escopo, filtros e origem dos dados" icon={<Filter size={16} />}>
        <dl className="metodologia-guide">
          <div className="metodologia-guide-item">
            <dt>Regiões</dt>
            <dd>
              As cinco regiões oficiais do IBGE: Norte, Nordeste, Centro-Oeste, Sudeste e
              Sul.
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Onde cada filtro age</dt>
            <dd>
              Faixa etária (?regiao= e ?ibge=), Penetração ANS (?regiao=) e tabela por
              UF (?regiao=) filtram no servidor; Estados filtra no navegador; Cidades
              busca e ordena no servidor. O painel Mercados por Região não é afetado
              pelo filtro da tabela — ele sempre compara as cinco regiões.
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Caminho do dado</dt>
            <dd>
              API do IBGE → importação (serviço import do docker compose) e dados
              abertos da ANS → importação (serviço import_ans, roda junto no
              docker compose up) → PostgreSQL → API Go → esta interface
              (React + Vite).
            </dd>
          </div>
          <div className="metodologia-guide-item">
            <dt>Disponibilidade</dt>
            <dd>
              GET /health e /health/db; o topo da tela verifica a cada 30s e o painel de
              Status mostra a latência.
            </dd>
          </div>
        </dl>
      </Panel>

      <Panel title="Herança PHP" icon={<Code size={16} />}>
        <p className="metodologia-note">
          O enunciado do Trads descreve uma versão anterior em PHP como herança do
          projeto. Esta entrega é a reconstrução em Go + PostgreSQL + React, com banco e
          API próprios; as fontes, fórmulas e filtros desta página são os vigentes nesta
          versão.
        </p>
      </Panel>
    </div>
  );
}
