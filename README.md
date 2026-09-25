# Painel de Inteligência de Mercado - Trads Corretora

Aplicação fullstack que responde à pergunta que guia o negócio da Trads: **em quais
regiões do Brasil estão os melhores mercados, e para qual público?**

Ela consome a API pública do IBGE (demografia, renda, PIB) e dados abertos da ANS
(beneficiários de planos de saúde), persiste tudo em um banco próprio e apresenta um
dashboard com filtros por região, estado e município.

[![backend-ci](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/backend-ci.yml)
[![frontend-ci](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/frontend-ci.yml)
[![docker-ci](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/docker-ci.yml/badge.svg)](https://github.com/CunhazadanoDale/trads-market-test/actions/workflows/docker-ci.yml)

---

## Como rodar o projeto do zero

**Pré-requisito único: Docker com Docker Compose.** Nenhuma linguagem ou banco precisa
estar instalado na máquina.

```bash
git clone https://github.com/CunhazadanoDale/trads-market-test.git
cd trads-market-test
docker compose up -d --build
```

Na primeira subida o compose orquestra os serviços desta forma:

1. **postgres** (Postgres 16, aguarda healthcheck);
2. **migrate** — aplica as migrações de `backend-golang/db/migrations/`;
3. **import** — baixa os dados do IBGE para o banco (upsert idempotente, pode rodar
   de novo sem quebrar);
4. **import_ans** — baixa o CSV de beneficiários da ANS;
5. **api** — sobe quando o postgres está saudável, mesmo que a importação ainda esteja
   em andamento ou tenha falhado (a aplicação não depende do IBGE para servir);
6. **frontend** — nginx com o build do React, espera a API ficar saudável.

| Serviço | URL |
|---|---|
| Aplicação (dashboard) | http://localhost:8085 |
| API | http://localhost:8082 |
| Health da API | http://localhost:8082/health e `/health/db` |
| Adminer (banco) | http://localhost:8084 |
| Postgres | localhost:5444 (user/senha/db: `trads`) |

### Comandos úteis (Makefile)

```bash
make up             # sobe tudo (docker compose up -d --build)
make down           # para os containers
make reset          # para e apaga o volume do banco (começa do zero)
make import         # roda a importação do IBGE de novo
make import-ans     # roda a importação da ANS de novo
make test           # go vet + go test (backend)
make lint           # oxlint (frontend)
make migrate-status # versão atual das migrações
```

A importação completa do IBGE (estados, cidades, população, renda, PIB e faixas
etárias) roda em **poucos minutos** na primeira subida e é reexecutável quantas
vezes for preciso, pois usa upsert idempotente. A otimização para escrita em lote
transacional derrubou as medições locais de **7min43 → 19s** no bloco de cidades e
indicadores e de **384s → 21s** na faixa etária. Os tempos por etapa aparecem no log
do serviço `import`.

---

## O que a aplicação faz

| Tela | O que responde |
|---|---|
| **Dashboard** (`/`) | Panorama nacional (municípios, população, renda, PIB), rankings das maiores cidades em população/renda/PIB, distribuição por faixa etária, status dos serviços |
| **Mercados** (`/mercados`) | *Onde vender?* — população, PIB e renda agregados por região (barras comparativas) e por UF (tabela com filtro de região), mais o painel **Penetração ANS** (top 10 municípios com mais beneficiários por população) |
| **Público** (`/publico`) | *Para qual público vender?* — distribuição etária do país com filtro de região e UF, junto da renda média das cidades onde cada faixa mora |
| **Estados** (`/estados`) | 27 UFs com busca por nome/UF/código e filtro de região |
| **Cidades** (`/cidades`) | Municípios de um estado com busca por nome, ordenação por população/renda/PIB, paginação e detalhe por cidade |
| **Metodologia** (`/metodologia`) | De onde vem cada número: agregado, variável, ano, fonte e fórmula de cada indicador |

### Fluxo do dado

```
API do IBGE (localidades v1 + agregados SIDRA)   ANS dados abertos (CSV PDA-047)
                \                                    /
                 \── serviço import / import_ans ──/      (upsert em lote, idempotente)
                                 │
                          PostgreSQL 16
                                 │
                    API Go (net/http + pgx, porta 8082)
                                 │
                 nginx (proxy /api) → React + Vite (porta 8085)
```

A interface **nunca** consulta o IBGE direto: ela fala com o banco da Trads através da
própria API. É o requisito de persistência do desafio.

---

## A herança: análise do código legado

A tentativa anterior estava na pasta `legado/` do repositório do desafio: um coletor em
PHP puro, um painel em HTML + jQuery, um CSV exportado e o bilhete do ex-funcionário.

### O que a tentativa anterior fazia (ou dizia fazer)

O `coleta_ibge.php` percorria os estados, chamava o que o comentário chama de
"endpoint v9 mais atual", media um valor por UF, aplicava um filtro de "2 salários
mínimos" e supostamente gravava numa tabela MySQL. O `painel_antigo.html` exibia uma
tabela estática com botão "Atualizar dados". O bilhete afirmava que "já está 90%
pronto", que o CSV estava "certinho e atualizado de 2020" e que "roda de novo que
costuma passar".

### O que está errado ou não é confiável

Conferi cada afirmação contra a documentação oficial do IBGE e contra o próprio código:

1. **"Já salva no banco" — não salvava.** A função `salvar_no_banco()` escreve uma linha
   em `coleta_debug.txt`. Não existe INSERT. O `config.php` com a conexão MySQL é
   incluído, mas a gravação nunca acontece.
2. **Chama de "renda" o que é PIB.** A função `pega_pib()` busca o agregado 5938,
   variável 37 ( **PIB total do estado** ) enquanto o comentário e o nome da função
   alternam entre "renda per capita" e "PIB". São três grandezas diferentes no mesmo
   trecho.
3. **O filtro de "2 salários mínimos" compara unidades erradas.** Ele compara PIB total
   (em Mil R$, soma de toda a UF) com `2824`. PIB de estado não é renda de pessoa, e
   2824 era 2× salário mínimo de 2024 aplicado a dados de 2020. A premissa do bilhete
   ("abaixo disso o IBGE não tem dado confiável") não encontra respaldo na documentação
   do IBGE e, aplicada de verdade à renda per capita municipal, eliminaria a maioria
   dos municípios do país.
4. **"Endpoint v9" não existe na documentação.** A API oficial de localidades é
   `servicodados.ibge.gov.br/api/v1/localidades`. Também não há retry, timeout ou
   checagem de status HTTP: `file_get_contents` falhando devolve warning e o script
   segue processando `null`.
5. **O CSV não é confiável, apesar do bilhete.** João Pessoa aparece **duas vezes**;
   há um "Municipio Fantasma" (código 9999999); Porto Alegre aparece com `-1` numa
   coluna numérica; Manaus tem população `"1.234.567"` (texto formatado, não número);
   Curitiba tem um campo a mais; os nomes estão com encoding quebrado (`JoÃ£o Pessoa`);
   e o **header do arquivo contradiz o bilhete**. O bilhete diz coluna 3 = população,
   coluna 4 = renda, mas o arquivo é `codigo;municipio;renda;populacao_toal` (com typo
   no próprio header).
6. **O painel tem bugs que impedem o uso.** `montatabela()` é chamado no load, mas a
   função se chama `montaTabela` → erro de JavaScript na abertura. O `sort()` retorna
   booleano, então a ordenação não funciona. O botão "Atualizar dados" joga a resposta
   da API de estados (objetos) numa tabela que espera arrays de colunas. O cache
   embutido mistura população e renda nas posições que o próprio comentário do HTML
   documenta como invertidas.
7. **Credencial no código.** `config.php` tem a senha de produção no repositório com
   um `TODO: tirar a senha daqui antes de subir pro git`, usando a API `mysql_*`, que
   foi **removida no PHP 7** - o script nem roda em PHP moderno.
8. **"Roda de novo que costuma passar"** é a única estratégia de erro do projeto:
   nenhuma transação, nenhum estado consistente, nenhum log do que falhou.

### Por que descartamos e como esta solução resolve

Descartamos o PHP porque o problema não era a linguagem e sim que era que nada persistia e os
números não eram o que diziam ser e nenhuma afirmação era verificável. Reconstruímos
do zero com a mesma pergunta em mente:

| Problema do legado | Como esta solução resolve |
|---|---|
| Dados não persistiam (só log de debug) | PostgreSQL com 7 migrações, upsert transacional em lote |
| "Renda" era PIB total de UF | Cada indicador mapeado para agregado/variável/ano documentados e testados (ver *Metodologia* na aplicação) |
| Filtro em unidade errada e premissa não verificável | Regra do legado não reproduzida; decisão registrada nas suposições abaixo |
| Endpoint inventado, sem erro/retry | Cliente HTTP unificado com retry, backoff exponencial, timeout e User-Agent, testado com `httptest` |
| CSV com duplicatas, fantasma e `-1` | Importação refeita a cada execução a partir da fonte oficial; consistência garantida por constraints (`UNIQUE`, `CHECK >= 0`) |
| Painel quebrado e ordenação falsa | Dashboard com filtros que refazem a consulta no servidor e paginação real |
| Senha no repositório | Configuração por variável de ambiente; `.env` fora do git |
| "Roda de novo que costuma passar" | Erros retornam com contexto (`%w`), a API sobe mesmo se a importação falhar, healthcheck `/health/db` |

---

## Decisões técnicas

### Stack

| Camada | Escolha | Por quê |
|---|---|---|
| Backend | **Go 1.26** + `net/http` da biblioteca padrão | Compilação única estática e imagem final pequena (base Alpine), tipagem forte para dados numéricos, `net/http` já resolve o que precisamos, portanto sem motivos para um framework |
| Banco | **PostgreSQL 16** | Os dados são relacionais e agregados (cidade → estado → região); boa parte do trabalho acontece em SQL (somas, médias ponderadas, janelas). Migrações versionadas com golang-migrate (imagem oficial no compose) |
| Acesso ao banco | **pgx** (pool, transações, `ON CONFLICT`) | Sem ORM: as consultas de agregação são o coração do produto e ficam mais legíveis em SQL explícito |
| Frontend | **React 19 + Vite** | SPA com 6 telas e filtros que refazem consultas; Vite dá build rápido e dev server com hot reload |
| Gráficos | **recharts** | Barras horizontais com tooltip acessível sem escrever SVG à mão |
| Servir o front | **nginx** no container do frontend | Proxy de `/api` e `/health` para a API → mesmo host, **zero CORS** |
| Infra | **Docker Compose** | Exigência do desafio: o avaliador sobe tudo sem instalar nada |
| Lint/teste | **oxlint** (front), `gofmt`/`go vet`/`go test`, **Vitest + Testing Library** | Verificação local e nos 3 workflows de CI |

### Arquitetura

O backend está organizado em camadas com dependência apontando para dentro:

```
backend-golang/
  cmd/api/            entrada da API (graceful shutdown, timeouts)
  cmd/import/         entrada da importação (ibge | ans | all)
  db/migrations/      7 migrações up/down
  internal/
    core/
      domain/         entidades e erros (City, State, Indicator, ErrCityNotFound…)
      ports/in/       interfaces dos casos de uso
      ports/out/      interfaces dos repositórios
      usecases/       regras de negócio (import, busca, métricas)
    adapter/
      http/           router, handlers, dtos (padrão {dados, pagina, total} e erro {error:{code,message}})
      ibge/           cliente HTTP do IBGE com retry
      ans/            cliente do CSV da ANS
      postgres/       repositórios concretos (pgx)
```

Os handlers não conhecem o banco; os repositórios não conhecem HTTP. Isso permitiu
testar casos de uso e handlers com fakes e o cliente do IBGE com `httptest`, sem
subir infraestrutura.

**Trade-off assumido:** essa separação é mais pesada do que um CRUD simples pediria.
Mantive porque torna o código testável e substituível (trocar a fonte de dados é
implementar uma porta `out`), e documentei aqui para que a escolha seja consciente.
Para um serviço pequeno, como este, um pacote único com repositórios seria suficiente.

### Escolha dos indicadores

A pergunta de negócio tem duas metades, e cada uma virou uma tela:

- **Onde vender** → população (agregado 4709, var. 93, Censo 2022), PIB municipal
  (agregado 5938, var. 37, 2023) e renda (agregado 10295, var. 13431, Censo 2022),
  agregados por cidade → UF → região.
- **Para qual público vender** → faixas etárias (agregado 9514, var. 93, Censo 2022),
  cruzadas com a renda da cidade onde cada faixa mora.
- **Quanto já se vende** → beneficiários de planos da ANS (PDA-047), transformados em
  **penetração** = beneficiários ÷ população × 100. Município com penetração alta é
  mercado disputado; com penetração baixa e renda alta é oportunidade.

Renda média é sempre **ponderada pela população** (Σ renda×pop / Σ pop), não média
ingênua de municípios(metrópole tem que pesar mais que vila). A fórmula de cada número
está dentro da aplicação, na tela *Metodologia*.

---

## Núcleo e extras

**Núcleo (obrigatório do desafio)**

- Integração real com a API do IBGE (localidades + agregados SIDRA)
- Persistência em banco próprio; a aplicação não bate no IBGE por requisição
- Consultas e filtros: região, estado, município, busca por nome, ordenação,
  paginação
- Dashboard com mais de duas visualizações (barras, rankings, tabelas, cards) e
  filtros que atualizam o que é exibido
- `docker compose up` sobe aplicação + banco + migração + importação
- Análise do legado (acima) e decisões técnicas (acima)
- Repositório público com histórico de mais de 140 commits incrementais

**Extras (opcionais, feitos)**

- Cruzamento com **dados da ANS** (beneficiários × população → penetração)
- **CI no GitHub Actions**: 3 workflows (backend com `gofmt`/`vet`/`test` + job de
  integração com Postgres real, frontend com lint/build/testes, build das imagens)
- **Testes**: suíte Go (casos de uso, handlers, dtos, cliente IBGE, integração com
  Postgres) e 33 testes de frontend (Vitest + Testing Library)
- **Importação otimizada** (upsert em lote em transação): 7min43 → 19s
- **Resiliência**: retry/backoff no cliente do IBGE, API que sobe mesmo com a
  importação falhando, healthchecks, graceful shutdown
- Página de **Metodologia** dentro do produto (fonte, ano e fórmula de cada número)

**Extras opcionais não feitos** (ver *Limitações*): deploy, atualização agendada e
mapa do Brasil.

---

## Requisitos do desafio → onde está

| Requisito | Onde |
|---|---|
| Análise do legado no README | Seção *A herança* |
| Integração com a API do IBGE | `internal/adapter/ibge`, serviço `import` |
| Banco de dados persistido | `db/migrations`, serviço `postgres` |
| Consultas e filtros | `?regiao=`, `?ibge=`, busca/ordenação/paginação em `/api/v1/states/*` |
| Frontend com dashboard + filtros | Telas *Dashboard*, *Mercados*, *Público*, *Estados*, *Cidades* |
| Docker | `docker-compose.yml` (compose config validado no CI) |
| README com decisões | Este arquivo |
| Repo público com histórico | Mais de 140 commits convencionais (`feat`, `fix`, `test`, `ci`, `chore`) |

---

## Testes e CI

```bash
make test                    # go vet + go test (backend)
make lint                    # oxlint (frontend)
cd frontend && npm test      # Vitest (33 testes)
```

- **Backend**: testes unitários dos casos de uso, handlers, dtos e parsing do IBGE;
  testes de **integração** do repositório Postgres que rodam quando `DATABASE_URL`
  está definida (e reportam `SKIP` quando não está).
- **Frontend**: formatação, hook de dados com fetch mockado, componentes e filtro de
  região da tela de estados.
- **CI** (GitHub Actions, roda a cada push e PR):
  - `backend-ci`: `gofmt`, `go vet`, `go test` + job `integration` com Postgres 16
    como serviço;
  - `frontend-ci`: `npm ci`, lint, build e testes;
  - `docker-ci`: `docker compose config` e build das duas imagens.

---

## Como usei IA no desenvolvimento

Usei IA (OpenCode) como ferramenta dirigida em quatro momentos:

1. **Revisão antes de codar.** Peço uma auditoria do projeto contra o enunciado do
   desafio: comparação requisito a requisito, apontamento de bugs (uma mensagem de
   erro com formato nunca aplicado, typo numa migração, portas inconsistentes entre
   `.env`, Vite e compose), código morto e o que faltava. Essa revisão virou o plano de
   fases executado depois.
2. **Plano antes de execução.** Cada etapa (fundação, CI, resiliência do backend,
   frontend, diferenciais) foi descrita **antes** de qualquer mudança: arquivos no
   escopo, o que fazer, o que não tocar, comandos de verificação e critério de aceite.
3. **Uma sessão de IA por fase, com regras fixas.** Toda sessão recebeu as mesmas
   regras: nenhum comentário no código, escopo estrito de arquivos, mensagens em
   português, um commit convencional por tarefa, rodar `gofmt`/`vet`/`test`/`lint`/
   `build` e reportar a saída. A ordem das fases foi definida por dependência (formatação
   antes do CI, CI antes das fases paralelas) e cada sessão só podia adicionar os
   próprios arquivos ao commit, para fases rodando em paralelo não se pisarem.
4. **Supervisão e verificação.** Revisão do diff e dos testes ao fim de cada fase,
   mesma bateria de verificação que o CI roda, e decisão humana em cada encruzilhada
   (o que entraria, o que seria descartado. inclusive, o deploy foi planejado e cortado por
   decisão minha, não por falha).

O que a IA **não** fez: decidir o escopo sozinha, incluindo a arquitetura adotada, alterar este README sem revisão
humana, ou burlar as verificações e regras. O histórico de commits na branch `main` mostra a evolução
fase a fase.

---

## Limitações conhecidas

- **Sem deploy.** A aplicação roda apenas localmente via Docker Compose. O enunciado
  trata deploy como opcional; o preparo para produção (imagens multi-stage, healthcheck,
  config por variável de ambiente) já existe. Para projetos de pequeno porte eu costumo colocá-los no RAILWAY porque ele nunca adormece como o RENDER quando está sem uso por mais de um minuto. Como minha conta no railway é vinculada ao GITHUB e eu uso gratuitamente a ferramenta, seria necessário criar uma nova conta ou utilizar conta de terceiros para subir o projeto e isso vai contra minha intenção de ser **eu** a pessoa quem montou tudo. 
- **Sem atualização agendada.** A reimportação é manual (`make import` ou rodando novamente o `docker compose up`). Os upserts são
  idempotentes, então rodar de novo é seguro e falta apenas agendar.
- **Anos fixos.** Os indicadores são dos anos disponíveis na fonte (população, renda e
  idade: Censo 2022; PIB: 2023; ANS: ano do arquivo corrente). Não há seletor de ano
  nem série histórica.
- **Sem autenticação.** A API é somente leitura e aberta, como um painel interno de
  demonstração.
- **Bundle acima de 500 kB.** O build avisa; resolveria com code-splitting por rota.
- **Testes de UI não cobrem todas as telas** e não há teste E2E (Playwright/Cypress).
- **Sem mapa do Brasil** (visualização opcional do desafio).

## O que eu faria com mais tempo

1. **Deploy + atualização noturna** dos dados (o passo que ficou de fora desta entrega).
2. **Mapa coroplético** das UFs por penetração/renda -> a visualização mais natural
   para "onde vender".
3. **Seletor de ano** e comparação entre períodos (o IBGE publica séries históricas).
4. Mais datasets ANS (titulares, dependentes, operadoras) para aprofundar o cruzamento.
5. Testes E2E e cache de respostas agregadas (a consulta de penetração varre todos os
   municípios a cada chamada).

## Suposições documentadas

- A "regra dos 2 salários mínimos" do legado **não foi reproduzida**: ela comparava PIB
  com salário mínimo e não tem respaldo na documentação do IBGE. Preferi não inventar
  um corte de negócio.
- "Renda média" = rendimento domiciliar mensal per capita (agregado 10295, var. 13431),
  ponderada por população: é a leitura que a documentação do SIDRA dá ao dado.
- Penetração ANS = beneficiários ÷ população × 100, com os dois dados do mesmo recorte
  municipal; municípios sem população importada ficam fora do ranking.
- Números de ano diferente (PIB 2023 com população 2022) são somados como estão:
  limitação aceita e exibida na tela *Metodologia*.
- Portas de desenvolvimento: backend `8081`, Vite `5173` (proxy), compose `8082/8084/8085`.

## Estrutura do repositório

```
.
├── .github/workflows/      # CI: backend-ci, frontend-ci, docker-ci
├── backend-golang/
│   ├── cmd/api/            # servidor HTTP
│   ├── cmd/import/         # importador (ibge | ans | all)
│   ├── db/migrations/      # migrações do Postgres
│   └── internal/           # core (domain, ports, usecases) + adapters (http, ibge, ans, postgres)
├── frontend/
│   └── src/
│       ├── app/            # design system, componentes, layout
│       ├── modules/        # dashboard, mercados, publico, states, cities, metodologia
│       ├── services/       # cliente HTTP + serviços da API
│       └── hooks/          # useApiResource (loading, erro, cancelamento, reload)
├── docker-compose.yml      # postgres, adminer, migrate, import, import_ans, api, frontend
├── Makefile                # up, down, reset, import, test, lint, migrações
└── README.md
```
