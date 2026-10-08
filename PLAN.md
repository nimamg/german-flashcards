# Flashcards — German vocabulary & articles

A personal flashcard app for learning German nouns (with der/die/das), built to learn Go and React.

## Goals
- Learn idiomatic Go (stdlib-first) and enough React + TypeScript to build a real UI.
- Build something I actually use daily for German vocab.
- Proceed in small, reviewable steps. Each milestone ends with working code.

## Stack
| Layer     | Choice                                          | Why                                        |
|-----------|-------------------------------------------------|--------------------------------------------|
| Backend   | Go, stdlib `net/http` (1.22+ router)            | Learn the language, not a framework        |
| Storage   | SQLite (`modernc.org/sqlite`, pure Go, no cgo)  | Single file, trivial to deploy             |
| Frontend  | React + TypeScript + Vite                       |                                            |
| SRS       | Leitner boxes first, maybe SM-2 later           | Simple to reason about, easy to test       |
| Deploy    | Single Go binary embedding the React build      | One artifact, run anywhere (Fly.io / VPS)  |
| Users     | Single user for now; accounts maybe later       |                                            |

## Data model (initial sketch)
- **Word**: id, lemma (`Tisch`), gender (`der`/`die`/`das`, nullable for non-nouns), plural (`Tische`), part of speech, translation (`table`), example sentence, tags.
- **Card**: one word can produce several cards, one per *mode*:
  - `article`: show `Tisch` → answer `der`
  - `meaning`: `der Tisch` → `table` (and reverse)
  - `plural`: `der Tisch` → `die Tische`
- **Review state** per card: Leitner box, due date, history of answers.

## Word source
1. **Seed**: a small hand-curated CSV (~50 common nouns) checked into the repo. Enough for milestones 1–5.
2. **Import**: CSV import so I can paste my own lists.
3. **Enrichment (later)**: build a lookup table from Wiktionary data
   ([kaikki.org](https://kaikki.org/dictionary/German/) JSONL extract, CC BY-SA)
   so typing `Tisch` auto-fills gender, plural and translation. This also covers
   streaming a large JSON file in Go.

## Milestones
1. **Go basics: CLI quiz.** Load words from CSV, quiz articles in the terminal.
   *Learn:* modules, packages, structs, slices/maps, errors, `encoding/csv`, `bufio`.
2. **Spaced repetition.** Leitner scheduling as a pure package with table-driven tests.
   *Learn:* `testing`, interfaces, `time`, injecting a clock.
3. **Persistence.** SQLite schema + migrations, repository layer.
   *Learn:* `database/sql`, `context`, embedding SQL files with `embed`.
4. **HTTP API.** JSON REST API (`GET /api/cards/due`, `POST /api/reviews`, `POST /api/words`, …).
   *Learn:* `net/http`, handlers, middleware, JSON encoding, `httptest`.
5. **React frontend.** Vite + TS: review screen (card flip, der/die/das buttons), gender color-coding.
   *Learn:* components, state, hooks, fetching data, dev proxy to the Go API.
6. **Polish.** Word management UI, CSV import, stats (accuracy per gender, streaks), Wiktionary enrichment.
7. **Deploy.** `embed` the frontend build into the Go binary, Dockerfile, deploy, back up the SQLite DB.

Later/maybe: accounts, audio (TTS), other parts of speech (verbs with irregular forms), mobile-friendly PWA.

## Proposed layout (grows over time)
```
flashcards/
  cmd/
    quiz/        # milestone 1 CLI
    server/      # milestone 4 HTTP server
  internal/
    vocab/       # Word, Card types; CSV loading
    srs/         # scheduling logic
    store/       # SQLite repository
    api/         # HTTP handlers
  web/           # React app (milestone 5)
  data/
    seed.csv
```
