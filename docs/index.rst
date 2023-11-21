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
   :caption: Overview
   :hidden:

   overview/install.md

.. toctree::
   :maxdepth: 2
   :caption: Tutorials
   :hidden:

   tutorials/getting-started-mysql.md
   tutorials/getting-started-postgresql.md
   tutorials/getting-started-sqlite.md

.. toctree::
   :maxdepth: 2
   :caption: Commands
   :hidden:

   howto/generate.md
   howto/push.md
   howto/verify.md
   howto/vet.md

.. toctree::
   :maxdepth: 2
   :caption: How-to Guides
   :hidden:

   howto/select.md
   howto/query_count.md
   howto/insert.md
   howto/update.md
   howto/delete.md

   howto/prepared_query.md
   howto/transactions.md
   howto/named_parameters.md

   howto/ddl.md
   howto/structs.md
   howto/embedding.md
   howto/overrides.md
   howto/rename.md

.. toctree::
   :maxdepth: 3
   :caption: sqlc Cloud
   :hidden:

   howto/managed-databases.md

.. toctree::
   :maxdepth: 3
   :caption: Reference
   :hidden:

   reference/changelog.md
   reference/cli.md
   reference/config.md
   reference/datatypes.md
   reference/environment-variables.md
   reference/language-support.rst
   reference/macros.md
   reference/query-annotations.md

.. toctree::
   :maxdepth: 2
   :caption: Conceptual Guides
   :hidden:

   howto/ci-cd.md
   guides/using-go-and-pgx.rst
   guides/plugins.md
   guides/development.md
   guides/privacy.md
