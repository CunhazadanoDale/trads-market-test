# TRADS — Frontend

Painel React (Vite) consumindo a API do `backend-golang`.

## Rodando

```bash
# 1. backend (na raiz do projeto)
cd backend-golang
go run ./cmd/api            # precisa do .env e do Postgres (docker compose up -d postgres + make migrate-up)

# 2. frontend
cd frontend
npm install
npm run dev                 # http://localhost:5173
```

Em desenvolvimento o Vite faz proxy de `/api` e `/health` para `http://localhost:8081`
(definido em `vite.config.js`), então **não há CORS**. Para apontar para outra base,
crie um `.env` a partir do `.env.example` com `VITE_API_BASE_URL`.

## Endpoints consumidos

| Método | Caminho | Onde é usado |
| --- | --- | --- |
| GET | `/api/v1/states` | Dashboard, select do topo, módulo **Estados**, select do módulo **Cidades**, selects dos módulos **Mercados** e **Público** |
| GET | `/api/v1/states/{ibgeCode}/cities?page&pageSize` | Painel de cidades do Dashboard e módulo **Cidades** (paginação) |
| GET | `/api/v1/cities/{ibgeCode}` | Rota de detalhe `/cidades/:ibgeCode` |
| GET | `/api/v1/dashboard/national` | Painel "Panorama nacional" do Dashboard |
| GET | `/api/v1/dashboard/states?regiao` | Painel "UFs por indicador" do módulo **Mercados** |
| GET | `/api/v1/dashboard/top-cities?limit` | Rankings do Dashboard (top PIB, renda e população) |
| GET | `/api/v1/dashboard/age?regiao&ibge` | Painel "Distribuição por faixa etária" do módulo **Público** |
| GET | `/api/v1/dashboard/ans?regiao&ibge` | Painel "Penetração ANS" do módulo **Mercados** |
| GET | `/health` | Chip de status no TopBar e painel "Status dos Serviços" |
| GET | `/health/db` | Idem, informando se o Postgres responde |

## Estrutura

```
src/
  services/        # camada de API (padrões do backend: snake_case,
                   # envelope {dados,pagina,tamanho,total}, erro em JSON {"error":{code,message}})
    api.js         # fetch base + ApiError + VITE_API_BASE_URL
    states.js      # GET /api/v1/states
    cities.js      # GET /api/v1/states/{ibge}/cities e /cities/{code}
    health.js      # GET /health e /health/db com latência
    dashboard.js   # GET /api/v1/dashboard/{national,states,top-cities,age,ans}
  hooks/
    useApiResource.js  # loading/erro/cancelamento/reload para qualquer serviço
  utils/
    format.js      # formatação pt-BR de indicadores (null => "—")
  app/
    design-system/ # variáveis + estilos globais
    components/    # DataGrid, FilterPanel, Panel, StatusBadge, ApiStatus,
                   # HorizontalBarChart, ChartTooltip
    layout/        # AppShell, Sidebar, TopBar
  modules/
    dashboard/     # cards nacionais + rankings + painéis (cidades, estados, status)
    metodologia/   # premissas, fontes das bases e dicionário de indicadores
    mercados/      # painéis Mercados por Região, UFs por indicador e Penetração ANS
    publico/       # painel Distribuição por faixa etária
    states/        # grid de estados com filtros (pesquisa/ região)
    cities/        # grid de cidades (paginação do backend) + rota de detalhe
```

## Adicionando um módulo novo (ex.: indicators)

1. Criar `src/services/<recurso>.js` chamando `api.get('/api/v1/...')`.
2. Criar `src/modules/<recurso>/<Recurso>.jsx` usando `useApiResource`.
3. Registrar a rota em `src/App.jsx` e o item no `src/app/layout/Sidebar.jsx`.
4. Ajustar o título em `src/app/layout/TopBar.jsx`.

## Scripts

- `npm run dev` — servidor de desenvolvimento
- `npm run build` — build de produção em `dist/`
- `npm run lint` — oxlint
