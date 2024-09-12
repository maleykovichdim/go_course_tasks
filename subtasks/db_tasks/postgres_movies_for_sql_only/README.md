# Task: Develop a Database for an Information Site About Movies

## Overview
This document outlines the tasks required to develop a database for an information site about movies.

---

## Task #1: Develop a Database Schema

Develop a database schema for storing the following entities:

- **Movies**
  - title
  - year of release
  - actors
  - directors
  - box office receipts
  - rating

- **Actor**
  - name
  - date of birth

- **Directors**
  - name
  - date of birth

- **Studios**
  - title

### Requirements
- Each film can have many actors and directors, but one studio.
- Rating is a value of the set PG-10, PG-13, PG-18.
- The year of release cannot be less than 1800.
- There cannot be two movies with the same title in one year.

---

## Task #2: Develop a Query, Supplying the Database with Sample Data

---

## Task #3: Write Queries
Here are the required queries:

- A selection of movies with the name of the studio.
- A selection of movies for leading actors.
- A count of movies for a minor director.
- Select movies for several directors from the list (subquery).
- Count the number of movies for an actor.
- Select actors and directors who participated in more than 2 movies.
- Count the number of movies with box office receipts of more than 1000.
- Count the number of directors who have collected more than 1000 movies.
- Select actors from different families.
- Count the number of movies with duplicates by title.

---

## Conclusion
This project will culminate in a well-structured database that effectively manages movie information and allows for extensive data queries related to actors, directors, studios, and ratings.
