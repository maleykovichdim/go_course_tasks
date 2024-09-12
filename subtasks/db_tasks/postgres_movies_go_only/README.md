# Task: Develop a Package for Working with the Database of an Information Site About Films

## Overview
Develop a package for working with the database of an information site about films, based on the database schema from the previous task.

---

## Task #1: Develop an Interface for Interacting with the Database

### Requirements
The interface should allow you to:
- Add an array of films.
- Delete a film.
- Update a film.
- Get an array of films.

When getting an array of films, it should be possible to specify the studio ID for filtering by studio. If the ID is not specified (== 0), then return all films.

---

## Task #2: Develop a Package for the PostgreSQL Database

Develop a package for the Postgres database that implements the interface contract from Task #1.

---

## Task #3: Write Tests for the Developed Functions

Write tests for the developed functions.
