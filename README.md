[Read in English](README.en.md)

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
    - [View](#view)
    - [Function](#function)
    - [Procedure](#procedure)
  - [Endpoints da API](#endpoints-da-api)
    - [Infraestrutura](#infraestrutura)
    - [Users](#users)
    - [Students](#students)
    - [Teachers](#teachers)
    - [Courses](#courses)
    - [Enrollments](#enrollments)
    - [Tests](#tests)
    - [Grades](#grades)
    - [Relatórios](#relatórios)
  - [Como executar](#como-executar)
    - [Requisitos](#requisitos)
    - [Passos](#passos)
  - [Decisões de projeto](#decisões-de-projeto)
    - [Arquitetura hexagonal](#arquitetura-hexagonal)
    - [DTOs de entrada e saída](#dtos-de-entrada-e-saída)
    - [Erros de domínio](#erros-de-domínio)
    - [Uso dos recursos do PostgreSQL](#uso-dos-recursos-do-postgresql)
  - [Evolução em relação ao semestre anterior](#evolução-em-relação-ao-semestre-anterior)
  - [Performance](#performance)
  - [Vídeos](#vídeos)
  - [Licença](#licença)
<!--toc:end-->

Sistema de gestão acadêmica desenvolvido em Go e PostgreSQL.

---

## Visão geral

O Scholarly é uma API REST para gestão acadêmica que lida com alunos, professores, disciplinas, provas, matrículas e notas. O projeto foi desenvolvido como parte de um trabalho da disciplina de Banco de Dados, com o objetivo de demonstrar, na prática, a integração entre uma aplicação e recursos avançados do PostgreSQL.

O sistema atende à necessidade de uma plataforma centralizada de gestão acadêmica que vá além do CRUD tradicional, utilizando recursos nativos do banco de dados para processar e consolidar informações, como cálculo de média ponderada, boletins consolidados de alunos e matrícula atômica com validação de vagas.

---

## Diferenciais do projeto

Além dos requisitos mínimos, o projeto foi construído seguindo práticas adotadas em ambientes de produção:

- **Arquitetura hexagonal (ports and adapters)**, com separação estrita entre domínio, casos de uso, contratos e adaptadores de infraestrutura.
- **Camada de domínio pura**, independente de frameworks, drivers de banco de dados ou bibliotecas HTTP, permitindo que as regras de negócio sejam testadas sem dependências externas.
- **Hash de senha com Argon2id**, utilizando parâmetros acima da recomendação mínima do OWASP.
- **DTOs específicos por camada**, isolando o contrato HTTP das entidades de domínio e evitando a exposição de campos sensíveis como `password_hash`.
- **Tradução de erros do PostgreSQL**, mapeando códigos SQLSTATE (`23505`, `23503` e códigos customizados da aplicação como `P0S01`, `P0C01`) para erros de domínio.
- **Documentação interativa com Swagger UI**, permitindo que todos os endpoints sejam explorados e testados sem ferramentas adicionais.
- Orquestração completa com Docker Compose, incluindo banco de dados, migrações e aplicação. O projeto pode ser executado com apenas três comandos, sem necessidade de instalar Go ou PostgreSQL localmente.
- **Integração contínua com GitHub Actions**, executando `go vet`, `staticcheck` e testes automatizados a cada push.
- **Hot reload em desenvolvimento** com Air.
- **Migrações de banco de dados com golang-migrate em docker**, versionando a evolução do schema.
- Programa inteiramente containe

---

## Stack de tecnologias

| Camada | Tecnologia |
| --- | --- |
| Linguagem | Go 1.27 |
| Banco de dados | PostgreSQL 17 |
| Cache | Valkey 9 |
| Roteador HTTP | chi |
| Driver PostgreSQL | pgx/v5 |
| Hash de senha | Argon2id |
| Documentação da API | Swagger (swaggo) |
| Containerização | Docker e Docker Compose |
| Testes | Testify e Testcontainers |
| CI | GitHub Actions |
| Ferramenta de desenvolvimento | Air |
| Migrações de banco de dados | golang-migrate |

---

## Arquitetura

O projeto segue os princípios da **arquitetura hexagonal**, com dependências fluindo exclusivamente de fora para dentro:

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
         │  Adapters de saída (Postgres, Argon2, Valkey)
```

A camada de domínio não possui conhecimento de HTTP, bancos de dados ou qualquer biblioteca de infraestrutura. Os serviços dependem apenas de interfaces (`ports`), e os adaptadores implementam essas interfaces — permitindo substituir PostgreSQL por outro banco de dados, ou chi por outro roteador, sem alterações no núcleo da aplicação.

### Estrutura do projeto

```
Scholarly/
├── cmd/
│   └── api/                 # Entrypoint da aplicação
├── internal/
│   ├── domain/              # Entidades e regras de negócio
│   ├── ports/               # Interfaces (contratos)
│   ├── service/             # Casos de uso
│   ├── env/                 # Suporte a variáveis de ambiente (.env)
│   ├── config/              # Configurações básicas
│   └── adapters/
│       ├── argon2/          # Implementação de hash de senha        
│       ├── http/            # Handlers HTTP, DTOs, middlewares
|       ├── logger/          # Logger
|       ├── paseto/          # Autenticação PASETO
│       ├── postgres/        # PostgreSQL (Banco de dados Relacional)
|       ├── valkey/          # Valkey (Cache não relacional)
├── database/
│   ├── tables/              # Scripts de criação de tabelas
│   ├── views/               # Scripts de criação de Views
│   ├── functions/           # Scripts de criação de Functions
│   ├── procedures/          # Scripts de criação de Procedures
│   └── inserts/             # Dados de teste (seed)
├── docs/                    # Documentação Swagger gerada
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Banco de dados

### SGBD

**PostgreSQL 17**, executado via Docker.
**Supabase**, como opção.

### Principais tabelas

| Tabela | Descrição |
| --- | --- |
| `users` | Identidade central de alunos, professores e administradores |
| `students` | Dados específicos de alunos (herança de tabela a partir de `users`) |
| `teachers` | Dados específicos de professores (herança de tabela a partir de `users`) |
| `courses` | Disciplinas ofertadas por semestre |
| `enrollments` | Matrículas de alunos em disciplinas |
| `tests` | Provas por disciplina, com peso |
| `grades` | Nota de cada aluno por prova |

O modelo utiliza **herança por tabela** para separar a identidade (`users`) das especializações (`students`, `teachers`), evitando duplicação de dados de autenticação e permitindo que um único usuário possua múltiplos papéis.

### View

**`vw_student_report_card`**

Consolida, em uma única linha por aluno e disciplina, os dados do boletim acadêmico: nome do aluno, disciplina, professor, média ponderada e situação (aprovado, recuperação ou reprovado). A View encapsula múltiplos JOINs entre `users`, `students`, `enrollments`, `courses` e `teachers`, e integra a função `fn_weighted_average` para o cálculo da média.

**Uso na aplicação:** `GET /students/{student_id}/report`

### Function

**`fn_weighted_average(p_student_id BIGINT, p_course_id BIGINT)`**

Calcula a média ponderada de um aluno em uma disciplina, considerando o peso de cada prova. Retorna um valor `NUMERIC` — `0` quando o aluno não possui notas, ou levanta uma exceção (`P0S01`, `P0C01`) quando o aluno ou a disciplina não existem.

**Uso na aplicação:**

- Diretamente pelo endpoint `GET /students/{student_id}/courses/{course_id}/average`
- Internamente pela View `vw_student_report_card`

### Procedure

**`sp_enroll_student(p_student_id BIGINT, p_course_id BIGINT, OUT p_enrollment_id BIGINT)`**

Realiza a matrícula de um aluno em uma disciplina com validações atômicas:

1. Aluno existe
2. Disciplina existe
3. Aluno não está já matriculado
4. Disciplina possui vaga disponível

Retorna o `id` da matrícula criada via parâmetro `OUT`.

**Uso na aplicação:** `POST /courses/enroll/`

---

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

### Users (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| GET | `/users` | Lista usuários com filtros |
| PUT | `/users` | Atualiza usuário |
| GET | `/users/{id}` | Busca usuário por ID |
| DELETE | `/users/{id}` | Remove usuário |

### Students (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/students` | Cria aluno |
| PUT | `/students` | Atualiza aluno |
| GET | `/students/enrollment/{id}` | Busca aluno por número de matrícula |
| GET | `/students/{id}` | Busca aluno por ID |
| DELETE | `/students/{id}` | Remove aluno |

### Teachers (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/teachers` | Cria professor |
| PUT | `/teachers` | Atualiza professor |
| GET | `/teachers` | Lista professores |
| GET | `/teachers/{id}` | Busca professor por ID |
| DELETE | `/teachers/{id}` | Remove professor |

### Courses (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/courses` | Cria disciplina |
| PUT | `/courses` | Atualiza disciplina |
| GET | `/courses` | Lista disciplinas |
| GET | `/courses/{id}` | Busca disciplina por ID |
| DELETE | `/courses/{id}` | Remove disciplina |
| POST | `/courses/enroll/` | Matricula aluno em uma disciplina (utiliza Procedure) |

### Enrollments (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| PUT | `/enrollments` | Atualiza matrícula |
| GET | `/enrollments` | Lista matrículas |
| GET | `/enrollments/{id}` | Busca matrícula por ID |
| DELETE | `/enrollments/{id}` | Remove matrícula |

### Tests (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/tests` | Cria prova |
| PUT | `/tests` | Atualiza prova |
| GET | `/tests` | Lista provas |
| GET | `/tests/{id}` | Busca prova por ID |
| DELETE | `/tests/{id}` | Remove prova |

### Grades (Bearer Token)

| Método | Rota | Descrição |
| --- | --- | --- |
| POST | `/grades` | Cria nota |
| PUT | `/grades` | Atualiza nota |
| GET | `/grades` | Lista notas |
| GET | `/grades/{id}` | Busca nota por ID |
| DELETE | `/grades/{id}` | Remove nota |
| GET | `/grades/average/list` | Lista as médias consolidadas |
| GET | `/grades/average/{student_id}/{course_id}` | Consulta a média ponderada (Function) |
A documentação completa, com payloads e códigos de resposta, está disponível via Swagger UI em `http://localhost:8080/swagger/index.html`.

---

## Como executar

### Requisitos

- Docker e Docker Compose
- Make

### Passos

**1. Clone o repositório**

```bash
git clone https://github.com/alvarolucio2007/Scholarly.git
cd Scholarly
```

**2. Remova o .example do arquivo .env.example**

O arquivo .env.example possui exemplos válidos de variáveis de ambiente.

Ele possui as seguintes variáveis:

```env
DATABASE_URL = postgres://admin:root@db:5432/scholarly?sslmode=disable
CACHE_URL = cache:6379
SERVER_ADDR = 0.0.0.0:8080
```

**3. Suba a aplicação completa (com Postgres local, para Supabase, instruções abaixo.)**

```bash
docker compose up --build
```

O docker compose já monta toda a aplicação, incluindo migrações de banco de dados.

Em modo de desenvolvimento com hot reload (é necessário ter a ferramenta [Make](https://kb.shells.com/tutorials/Windows_11/GNU_Make/)):

```bash
docker compose up db cache -d
make migrateup
air

```

**(OPCIONAL) Supabase como banco de dados**

Primeiro, mude a variável `DATABASE_URL` para a URL do supabase (é altamente recomendado usar Session Pooler nas opções de conexão do Supabase)

Após isso, basta executar:

```bash
docker compose up api cache --build --no-deps
```

**3. Acesse o Swagger UI**

```
http://127.0.0.1:8080/swagger/index.html
```

---

## Decisões de projeto

### Arquitetura hexagonal

A arquitetura hexagonal foi escolhida para isolar o domínio de dependências externas. Isso permite que as regras de negócio sejam testadas sem banco de dados, que adaptadores sejam substituídos sem impacto no núcleo, e que o projeto evolua sem acumular acoplamento.

### DTOs de entrada e saída

A aplicação utiliza DTOs específicos para HTTP, mantidos separados das entidades de domínio. Isso evita a exposição de campos sensíveis como `password_hash`, desacopla o contrato da API das entidades internas e simplifica a geração da documentação Swagger.

### Erros de domínio

Erros são definidos no pacote `domain` (`ErrNotFound`, `ErrEmailAlreadyExists`, `ErrCourseFull`, entre outros) e traduzidos no adaptador PostgreSQL a partir de códigos SQLSTATE. O handler HTTP é responsável por mapear cada erro de domínio para o status HTTP apropriado, mantendo o núcleo independente de detalhes de infraestrutura.

### Uso dos recursos do PostgreSQL

A View, a Function e a Procedure foram escolhidas para resolver problemas que fazem sentido no nível do banco:

- **View:** consultas consolidadas com múltiplos JOINs e agregações, que seriam duplicadas na aplicação.
- **Function:** cálculo de média ponderada, executado tanto diretamente quanto dentro da View.
- **Procedure:** matrícula de aluno, envolvendo múltiplas validações e uma inserção que precisa ser atômica.

---

## Evolução em relação ao semestre anterior

Este projeto é uma reescrita de um sistema acadêmico desenvolvido no semestre anterior, agora com foco em qualidade arquitetural e uso avançado de banco de dados.

| Aspecto | Versão anterior | Versão atual |
| --- | --- | --- |
| Arquitetura | Arquitetura em camadas | Hexagonal (ports and adapters) |
| Persistência | Sem FKs, sem índices | FKs, índices e constraints |
| Recursos SQL | Sem View, Function ou Procedure | View, Function e Procedure integrados |
| Hash de senha | Inexistente | Argon2id |
| Documentação | Inexistente | Swagger |
| Testes | Testes unitários | Testes unitários |
| CI/CD | Inexistente | GitHub Actions |
| Containerização | Docker e Docker Compose | Docker e Docker Compose |

---

## Performance

Comando

```bash
npx autocannon -c 1000 -d 20 -V http://127.0.0.1:8080/users/1 --renderStatusCodes
```

| Cenário | Req/s | p50 | p99 |
| --- | --- | --- | --- |
| Sem cache | 4,900 | 202ms | 290ms |
| Com cache | 47,000 | 21ms | 40ms |

---

## Vídeos

- [Parte 1: Demonstração do sistema](https://youtu.be/W15-8AW7WaY)
- [Parte 2: Recursos do banco de dados](https://youtu.be/InDuAxvD600)
- [Parte 3: Integração com Redis e Supabase](https://youtu.be/SAQKUrK_t3c)

---

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
