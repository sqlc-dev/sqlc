.. sqlc documentation master file, created by
   sphinx-quickstart on Mon Feb  1 23:18:36 2021.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

sqlc Documentation
==================

  And lo, the Great One looked down upon the people and proclaimed:
    "SQL is actually pretty great"

sqlc generates **fully type-safe idiomatic Go code** from SQL. Here's how it
works:

1. You write SQL queries
2. You run sqlc to generate Go code that presents type-safe interfaces to those
   queries
3. You write application code that calls the methods sqlc generated

Seriously, it's that easy. You don't have to write any boilerplate SQL querying
code ever again.

.. toctree::
   :maxdepth: 2
   :caption: Overview:
   :hidden:

   overview/install.md

.. toctree::
   :maxdepth: 2
   :caption: Tutorials:
   :hidden:

   tutorials/getting-started.md

.. toctree::
   :maxdepth: 2
   :caption: How-to Guides:
   :hidden:

   howto/alter_table.md
   howto/annotations.md
   howto/any.md
   howto/arrays.md
   howto/delete.md
   howto/enums.md
   howto/insert.md
   howto/json_tags.md
   howto/migrations.md
   howto/named_parameters.md
   howto/null.md
   howto/prepared_query.md
   howto/privacy.md
   howto/query_count.md
   howto/query_one.md
   howto/returning.md
   howto/table.md
   howto/time.md
   howto/transactions.md
   howto/update.md
   howto/uuid.md

.. toctree::
   :maxdepth: 2
   :caption: Reference:
   :hidden:

   reference/cli.md
   reference/config.md

.. toctree::
   :maxdepth: 2
   :caption: Conceptual Guides:
   :hidden:

   guides/development.md
   guides/privacy.md
