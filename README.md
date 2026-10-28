# Scholarly

<!--toc:start-->
- [Scholarly](#scholarly)
  - [Visão geral](#visão-geral)
  - [Diferenciais do projeto](#diferenciais-do-projeto)
  - [Stack de tecnologias](#stack-de-tecnologias)
  - [Arquitetura](#arquitetura)
    - [Estrutura do projeto](#estrutura-do-projeto)
  - [Banco de dados](#banco-de-dados)
    - [SGBD](#sgbd)
    - [Principais tabelas](#principais-tabelas)
    - [View, Function e Procedure](#view-function-e-procedure)
  - [Endpoints da API](#endpoints-da-api)
    - [Infraestrutura & Autenticação (Público)](#infraestrutura--autenticação-público)
    - [Administração (Basic Auth)](#administração-basic-auth)
    - [Recursos Protegidos (Bearer Token)](#recursos-protegidos-bearer-token)
  - [Como executar](#como-executar)
    - [Requisitos](#requisitos)
    - [Passos](#passos)
  - [Decisões de projeto](#decisões-de-projeto)
    - [Arquitetura hexagonal](#arquitetura-hexagonal)
    - [Autenticação e Segurança](#autenticação-e-segurança)
    - [Frontend e Experiência do Usuário](#frontend-e-experiência-do-usuário)
    - [Uso dos recursos do PostgreSQL](#uso-dos-recursos-do-postgresql)
  - [Evolução em relação ao semestre anterior](#evolução-em-relação-ao-semestre-anterior)
  - [Performance](#performance)
  - [Licença](#licença)
<!--toc:end-->

Sistema de gestão acadêmica Fullstack desenvolvido com Go, PostgreSQL e SvelteKit.

---

## Visão geral

O Scholarly é uma aplicação web completa para gestão acadêmica que lida com alunos, professores, disciplinas, provas, matrículas e notas. O projeto foi desenvolvido como parte de um trabalho da disciplina de Banco de Dados, com o objetivo de demonstrar, na prática, a integração entre uma aplicação Fullstack moderna e recursos avançados do PostgreSQL.

O sistema atende à necessidade de uma plataforma centralizada de gestão acadêmica que vá além do CRUD tradicional. No backend, utiliza recursos nativos do banco de dados para processar e consolidar informações (como cálculo de média ponderada, boletins consolidados e matrícula atômica com validação de vagas). No frontend, entrega uma interface de usuário (UI) fluida, responsiva e moderna para a manipulação desses dados.

> **Nota de transparência:** O frontend foi construído com o auxílio de IA, refletindo a adoção de ferramentas modernas de produtividade para acelerar a entrega de interfaces, enquanto toda a lógica de negócio, arquitetura e infraestrutura do backend foram desenvolvidas do zero.

---

## Diferenciais do projeto

Além dos requisitos mínimos, o projeto foi construído seguindo práticas adotadas em ambientes de produção:

- **Arquitetura Fullstack desacoplada:** Backend em Go (API REST) e Frontend em SvelteKit (SPA), comunicando-se de forma segura via Bearer Token.
- **Arquitetura hexagonal (ports and adapters):** Separação estrita no backend entre domínio, casos de uso, contratos e adaptadores de infraestrutura.
- **Camada de domínio pura:** Independente de frameworks, drivers de banco de dados ou bibliotecas HTTP, permitindo que as regras de negócio sejam testadas sem dependências externas.
- **Autenticação stateless com PASETO v4:** Utilizando criptografia simétrica (XChaCha20-Poly1305) e proteção contra ataques de análise de tempo (*timing attacks*).
- **Hash de senha com Argon2id:** Utilizando parâmetros acima da recomendação mínima do OWASP.
- **Rate Limiting e proteção contra DDoS/Brute Force:** Limitando requisições por IP para garantir a estabilidade da API.
- **Observabilidade nativa:** Endpoints de profiling (`pprof`) e métricas (`expvar`) protegidos por autenticação.
- **Frontend moderno:** Interface construída com SvelteKit, TypeScript e Tailwind CSS, com validação de formulários em tempo real e feedback visual de erros.
- **Tradução de erros do PostgreSQL:** Mapeando códigos SQLSTATE (`23505`, `23503` e códigos customizados) para erros de domínio na API.
- **Documentação interativa com Swagger UI:** Permitindo que todos os endpoints sejam explorados e testados.
- **Orquestração completa com Docker Compose:** Banco de dados, migrações, backend e frontend sobem com um único comando.
- **Integração contínua com GitHub Actions:** Executando `go vet`, `staticcheck` e testes automatizados a cada push.

---

## Stack de tecnologias

| Camada | Tecnologia |
| --- | --- |
| **Backend** | |
| Linguagem | Go 1.27 |
| Banco de dados | PostgreSQL 17 |
| Cache | Valkey 9 |
| Roteador HTTP | chi |
| Driver PostgreSQL | pgx/v5 |
| Autenticação | PASETO v4 (aidantwoods) |
| Rate Limiting | httprate |
| Hash de senha | Argon2id |
| Documentação | Swagger (swaggo) |
| **Frontend** | |
| Framework | SvelteKit (Svelte 5) |
| Linguagem | TypeScript |
| Estilização | Tailwind CSS |
| **Infraestrutura** | |
| Containerização | Docker e Docker Compose |
| Testes | Testify e Testcontainers |
| CI | GitHub Actions |

---

## Arquitetura

O projeto adota uma separação clara entre as responsabilidades de backend e frontend.

O **Backend** segue os princípios da **arquitetura hexagonal**, com dependências fluindo exclusivamente de fora para dentro:

```
┌─────────────────────────────────────────┐
│  HTTP (adapter de entrada)              │
│  ────────────────────────────────────   │
│  Service (casos de uso)                 │
│  ────────────────────────────────────   │
│  Ports (interfaces)                     │
│  ────────────────────────────────────   │
│  Domain (entidades e regras)            │
└─────────────────────────────────────────┘
         ▲
         │  Adapters de saída (Postgres, Argon2, Valkey, PASETO)
```

O **Frontend** atua como um consumidor puro da API REST. Ele gerencia o estado de autenticação (armazenando o token PASETO) e consome os endpoints protegidos, garantindo que a lógica de negócio permaneça centralizada no backend.

### Estrutura do projeto

```
Scholarly/
├── cmd/
│   └── api/                 # Entrypoint da aplicação (Backend)
├── internal/
│   ├── domain/              # Entidades e regras de negócio
│   ├── ports/               # Interfaces (contratos)
│   ├── service/             # Casos de uso
│   ├── env/                 # Suporte a variáveis de ambiente
│   ├── config/              # Configurações básicas
│   └── adapters/            # Implementações de infraestrutura (HTTP, DB, Cache, Auth)
├── database/
│   ├── tables/              # Scripts de criação de tabelas
│   ├── views/               # Scripts de criação de Views
│   ├── functions/           # Scripts de criação de Functions
│   ├── procedures/          # Scripts de criação de Procedures
│   └── inserts/             # Dados de teste (seed)
├── frontend/                # Aplicação SvelteKit (Frontend)
│   ├── src/
│   │   ├── routes/          # Páginas (login, dashboard, CRUDs)
│   │   └── lib/             # Componentes e utilitários
│   ├── package.json
│   └── Dockerfile
├── docs/                    # Documentação Swagger gerada
├── docker-compose.yml       # Orquestração de toda a stack
├── Makefile
└── README.md
```

---

## Banco de dados

### SGBD
**PostgreSQL 17**, executado via Docker. O projeto também possui suporte para conexão com o **Supabase** em nuvem.

### Principais tabelas
O modelo utiliza **herança por tabela** para separar a identidade (`users`) das especializações (`students`, `teachers`), evitando duplicação de dados e permitindo múltiplos papéis para um mesmo usuário.

| Tabela | Descrição |
| --- | --- |
| `users` | Identidade central (autenticação) |
| `students` | Dados específicos de alunos |
| `teachers` | Dados específicos de professores |
| `courses` | Disciplinas ofertadas |
| `enrollments` | Matrículas de alunos em disciplinas |
| `tests` | Provas por disciplina, com peso |
| `grades` | Nota de cada aluno por prova |

### View, Function e Procedure
A inteligência de dados foi delegada ao PostgreSQL para garantir performance e atomicidade:

- **View (`vw_student_report_card`):** Consolida múltiplos JOINs para gerar o boletim acadêmico (média e situação do aluno).
- **Function (`fn_weighted_average`):** Calcula a média ponderada de um aluno em uma disciplina.
- **Procedure (`sp_enroll_student`):** Realiza a matrícula de forma atômica, validando vagas disponíveis e duplicidade em uma única transação.

---

## Endpoints da API

### Infraestrutura & Autenticação (Público)

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/health` | Health check |
| POST | `/authentication/register` | Cria um novo usuário (Registro) |
| POST | `/authentication/token` | Autentica e gera o token PASETO (Login) |

### Administração (Basic Auth)

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/swagger/*` | Documentação interativa Swagger UI |
| GET | `/metrics` | Métricas da aplicação (expvar) |
| GET | `/debug/*` | Profiling da aplicação (pprof) |

### Recursos Protegidos (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| GET/PUT/DELETE | `/users` | Listagem, atualização e remoção de usuários |
| GET/POST/PUT/DELETE | `/students` | Gerenciamento de alunos |
| GET/POST/PUT/DELETE | `/teachers` | Gerenciamento de professores |
| GET/POST/PUT/DELETE | `/courses` | Gerenciamento de disciplinas |
| POST | `/courses/enroll/` | Matrícula de aluno (consome a Procedure do banco) |
| GET/PUT/DELETE | `/enrollments` | Gerenciamento de matrículas |
| GET/POST/PUT/DELETE | `/tests` | Gerenciamento de provas |
| GET/POST/PUT/DELETE | `/grades` | Gerenciamento de notas |
| GET | `/grades/average/list` | Lista as médias consolidadas (View) |
| GET | `/grades/average/{s_id}/{c_id}` | Consulta a média ponderada (Function) |

A documentação completa, com payloads e códigos de resposta, está disponível via Swagger UI em `http://localhost:8080/swagger/index.html`.

---

## Como executar

### Requisitos
- Docker e Docker Compose
- Make (Opcional, para atalhos de desenvolvimento)

### Passos

**1. Clone o repositório e configure o ambiente**
```bash
git clone https://github.com/alvarolucio2007/Scholarly.git
cd Scholarly
cp .env.example .env
cp /front/.env.example /front/.env
```
*(Edite o arquivo `.env` com as suas credenciais e defina um `BASIC_AUTH_USER` e `BASIC_AUTH_PASS` para acessar o Swagger e Métricas).*

**2. Suba a aplicação completa (Backend + Frontend + Banco + Cache)**
```bash
docker compose up --build
```
O Docker Compose já monta toda a aplicação, incluindo as migrações do banco de dados e o servidor do frontend.

**3. Acesse a aplicação**
*   **Frontend:** `http://localhost/login`
*   **Swagger UI:** `http://127.0.0.1:8080/swagger/index.html`

**(Opcional) Desenvolvimento com Hot Reload**
Para rodar o backend via Air e o frontend via Vite localmente:
```bash
# Backend
docker compose up db cache -d
make migrateup
air

# Frontend
cd frontend
npm install
npm run dev
```

---

## Decisões de projeto

### Arquitetura hexagonal
A arquitetura hexagonal foi escolhida no backend para isolar o domínio de dependências externas. Isso permite que as regras de negócio sejam testadas sem banco de dados, que adaptadores sejam substituídos sem impacto no núcleo, e que o projeto evolua sem acumular acoplamento.

### Autenticação e Segurança
- **PASETO v4 em vez de JWT:** Elimina vulnerabilidades comuns de JWT, forçando versões estritas e seguras por padrão.
- **Rate Limiting (`httprate`):** Aplicado em duas camadas (restrição severa no login para evitar *brute force* e limite global nas rotas da API).
- **Basic Auth:** Utilizado para proteger rotas administrativas e de infraestrutura (Swagger, Pprof e Métricas).

### Frontend e Experiência do Usuário
- **SvelteKit e Tailwind CSS:** Fornecem uma interface de usuário extremamente rápida, com bundle reduzido e estilização moderna sem dependência de bibliotecas pesadas.
- **Validação e Máscaras:** O frontend implementa validação de formulários (ex: máscara de CPF) e gerencia os estados de *loading* e erros amigáveis, garantindo uma boa experiência mesmo em falhas de rede.

### Uso dos recursos do PostgreSQL
A View, a Function e a Procedure foram escolhidas para resolver problemas que fazem sentido no nível do banco:
- **View:** Consultas consolidadas com múltiplos JOINs.
- **Function:** Cálculo de média ponderada.
- **Procedure:** Matrícula de aluno, envolvendo múltiplas validações e uma inserção atômica.

---

## Evolução em relação ao semestre anterior

Este projeto é uma reescrita de um sistema acadêmico desenvolvido no semestre anterior, agora com foco em qualidade arquitetural, segurança, performance e entrega Fullstack.

| Aspecto | Versão anterior | Versão atual |
| --- | --- | --- |
| Arquitetura | Arquitetura em camadas | Hexagonal (backend) + SPA (frontend) |
| Persistência | Sem FKs, sem índices | FKs, índices e constraints |
| Recursos SQL | Sem View, Function ou Procedure | View, Function e Procedure integrados |
| Autenticação | Inexistente | PASETO v4 (XChaCha20-Poly1305) |
| Hash de senha | Inexistente | Argon2id |
| Segurança da API | Nenhuma | Rate Limiting e Basic Auth (admin) |
| Observabilidade | Inexistente | pprof e expvar (métricas) |
| Documentação | Apenas README | Swagger |
| Interface | Streamlit | SvelteKit + Tailwind CSS |
| Testes | Testes unitários básicos | Testes unitários e de integração |
| CI/CD | Inexistente | GitHub Actions |
| Containerização | Docker e Docker Compose (Fullstack) | Docker e Docker Compose (Fullstack) |

---

## Performance

Os testes de carga foram realizados utilizando a ferramenta **Autocannon** com 1.000 conexões simultâneas no backend.

### Cache Hit vs Cache Miss
Otimização da rota `GET /users/1` com implementação do padrão *Cache-Aside* no Valkey:

| Cenário | Requisições/Seg | Latência (p50) | Latência (p99) |
| --- | --- | --- | --- |
| Sem cache | 4.900 | 202ms | 290ms |
| Com cache | 47.000 | 21ms | 40ms |

### Resiliência do Rate Limiter
Teste de estresse forçando o bloqueio de IP para validar a proteção contra DDoS e *Brute Force*:

| Métrica | Valor |
| --- | --- |
| Requisições/Seg (Bloqueadas) | ~49.000 |
| Latência Média de Bloqueio | 20 ms |
| Total de Bloqueios (429) | 979.702 |

*Mesmo sob um ataque massivo de 1 milhão de requisições em 20 segundos, a API manteve a estabilidade e respondeu instantaneamente com os bloqueios.*

---
## Vídeos

- [Parte 1: Demonstração do sistema](https://youtu.be/W15-8AW7WaY)
- [Parte 2: Recursos do banco de dados](https://youtu.be/InDuAxvD600)
- [Parte 3: Integração com Redis e Supabase](https://youtu.be/SAQKUrK_t3c)

---

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
