import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { Info, PieChart, Users } from 'lucide-react';
import './Publico.css';
import { HorizontalBarChart } from '../../app/components/HorizontalBarChart';
import { Panel } from '../../app/components/Panel';
import { useApiResource } from '../../hooks/useApiResource';
import { getStates } from '../../services/states';
import { getAgeDistribution } from '../../services/dashboard';
import { formatIncome, formatInteger, formatPercent } from '../../utils/format';

export default function Publico() {
  const statesRes = useApiResource(getStates, []);
  const [ageRegion, setAgeRegion] = useState('');
  const [ageIbge, setAgeIbge] = useState('');
  const [ageFaixa, setAgeFaixa] = useState('');
  const [faixaOptions, setFaixaOptions] = useState([]);

  const states = useMemo(() => statesRes.data ?? [], [statesRes.data]);

  const ageRes = useApiResource(
    ({ signal }) => getAgeDistribution({
      regiao: ageRegion,
      ibge: ageIbge,
      faixa: ageFaixa,
      signal,
    }),
    [ageRegion, ageIbge, ageFaixa],
  );

  const regionOptions = useMemo(
    () => [...new Set(states.map((state) => state.region))].sort((a, b) => a.localeCompare(b)),
    [states],
  );

  const ageStates = useMemo(
    () => states.filter((state) => !ageRegion || state.region === ageRegion),
    [states, ageRegion],
  );

  const selectedAgeState = useMemo(
    () => states.find((state) => String(state.ibge_code) === ageIbge) ?? null,
    [states, ageIbge],
  );

  const handleAgeRegionChange = (region) => {
    setAgeRegion(region);
    setAgeIbge((current) => {
      if (!current) return current;
      const state = states.find((item) => String(item.ibge_code) === current);
      return state && (!region || state.region === region) ? current : '';
    });
  };

  const age = ageRes.data;
  const ageGroups = useMemo(() => age?.grupos ?? [], [age]);

  useEffect(() => {
    if (!ageFaixa && ageGroups.length > 0) {
      setFaixaOptions(ageGroups.map((group) => group.faixa));
    }
  }, [ageFaixa, ageGroups]);

  const ageTooltipContent = (group) => {
    const share = age && age.total > 0 ? (group.populacao / age.total) * 100 : 0;
    return {
      title: group.faixa,
      rows: [
        { label: 'População', value: formatInteger(group.populacao) },
        { label: '% do total', value: formatPercent(share) },
        { label: 'Renda média', value: formatIncome({ value: group.renda_media_cidades }) },
      ],
    };
  };

  const ageScope = [
    ageRegion ? `Região ${ageRegion}` : '',
    selectedAgeState ? `UF ${selectedAgeState.uf}` : '',
    ageFaixa || '',
  ]
    .filter(Boolean)
    .join(' · ');

  return (
    <div className="publico-page">
      <header className="publico-intro">
        <h1 className="publico-title">
          <Users size={20} />
          Para qual público vender?
        </h1>
        <p className="publico-subtitle">
          Distribuição etária do Brasil por região e UF, com o entorno econômico de cada
          faixa — a segunda metade da pergunta que guia a Trads: em quais regiões estão os
          melhores mercados, e para qual público.
        </p>
      </header>

      <Panel title="Distribuição por faixa etária" icon={<PieChart size={16} />}>
        <div className="panel-controls">
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="age-region-filter">Região</label>
            <select
              id="age-region-filter"
              className="filter-select"
              value={ageRegion}
              onChange={(event) => handleAgeRegionChange(event.target.value)}
            >
              <option value="">Todas</option>
              {regionOptions.map((region) => (
                <option key={region} value={region}>{region}</option>
              ))}
            </select>
          </div>
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="age-uf-filter">UF</label>
            <select
              id="age-uf-filter"
              className="filter-select"
              value={ageIbge}
              onChange={(event) => setAgeIbge(event.target.value)}
            >
              <option value="">Todas</option>
              {ageStates.map((state) => (
                <option key={state.ibge_code} value={String(state.ibge_code)}>
                  {`${state.uf} — ${state.name}`}
                </option>
              ))}
            </select>
          </div>
          <div className="toolbar-field">
            <label className="toolbar-label" htmlFor="age-faixa-filter">Faixa</label>
            <select
              id="age-faixa-filter"
              className="filter-select"
              value={ageFaixa}
              onChange={(event) => setAgeFaixa(event.target.value)}
            >
              <option value="">Todas</option>
              {faixaOptions.map((faixa) => (
                <option key={faixa} value={faixa}>{faixa}</option>
              ))}
            </select>
          </div>
        </div>
        {ageRes.loading && !age ? (
          <div className="loading-box">Carregando faixa etária…</div>
        ) : ageRes.error ? (
          <div className="error-box">
            {ageRes.error.message}
            <button type="button" className="error-retry" onClick={ageRes.reload}>
              Tentar novamente
            </button>
          </div>
        ) : ageGroups.length === 0 ? (
          <div className="loading-box">Sem dado de faixa etária.</div>
        ) : (
          <>
            <HorizontalBarChart
              data={ageGroups}
              labelKey="faixa"
              barKey="populacao"
              rowHeight={30}
              barSize={14}
              labelWidth={130}
              valueAxis={{
                dataKey: 'renda_media_cidades',
                width: 130,
                format: (value) => formatIncome({ value }),
              }}
              tooltip={ageTooltipContent}
            />
            <p className="panel-hint">
              Renda média das cidades na coluna à direita: renda das cidades onde vivem as pessoas de cada faixa, ponderada pela população da faixa · passe o mouse sobre a barra para ver população da faixa e % do total.
            </p>
            <p className="panel-hint">
              Total {formatInteger(age.total)} pessoas · {ageScope || 'todo o país'} · Censo {age.ano} · IBGE/SIDRA 9514
            </p>
          </>
        )}
      </Panel>

      <Panel
        title="Como ler estes números"
        icon={<Info size={16} />}
        actions={<Link className="panel-link" to="/metodologia">Ver metodologia</Link>}
      >
        <dl className="publico-guide">
          <div className="publico-guide-item">
            <dt>População · %</dt>
            <dd>
              Quantas pessoas da faixa vivem no recorte escolhido (Região/UF) e que
              porcentagem do total elas representam ali.
            </dd>
          </div>
          <div className="publico-guide-item">
            <dt>Renda média das cidades</dt>
            <dd>
              Não é a renda da pessoa — uma criança de 0 a 4 anos não trabalha. É a renda
              média das cidades onde esse público mora, ponderada por onde ele se
              concentra. Por isso a faixa dos 0 a 4 anos fica abaixo da média nacional:
              as crianças se concentram em cidades de menor renda.
            </dd>
          </div>
          <div className="publico-guide-item">
            <dt>Ano e fonte</dt>
            <dd>
              População por idade: IBGE/SIDRA 9514 · Renda média por município: Censo 2022
              (IBGE). Os dados são importados da API do IBGE para o banco da Trads.
            </dd>
          </div>
          <div className="publico-guide-item">
            <dt>Como usar</dt>
            <dd>
              Filtre por Região e UF, compare as faixas entre si e entre recortes, e cruze
              com os indicadores das regiões para decidir onde e para quem vender.
            </dd>
          </div>
        </dl>
      </Panel>
    </div>
  );
}
