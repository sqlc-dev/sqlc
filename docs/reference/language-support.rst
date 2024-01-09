Database and language support
#############################

==========  =======================  ============  ============  ===============
Language    Plugin                   MySQL         PostgreSQL    SQLite
==========  =======================  ============  ============  ===============
Go          (built-in)               Stable        Stable        Beta
Go          `sqlc-gen-go`_           Stable        Stable        Beta
Kotlin      `sqlc-gen-kotlin`_       Beta          Beta          Not implemented
Python      `sqlc-gen-python`_       Beta          Beta          Not implemented
TypeScript  `sqlc-gen-typescript`_   Beta          Beta          Not implemented
==========  =======================  ============  ============  ===============

Community language support
**************************

New languages can be added via :doc:`plugins <../guides/plugins>`.

========  ==============================  ===============  ============  ===============
Language  Plugin                          MySQL            PostgreSQL    SQLite
========  ==============================  ===============  ============  ===============
F#        `kaashyapan/sqlc-gen-fsharp`_   Not implemented  Beta          Beta
========  ==============================  ===============  ============  ===============

.. _sqlc-gen-go: https://github.com/sqlc-dev/sqlc-gen-go
.. _kaashyapan/sqlc-gen-fsharp: https://github.com/kaashyapan/sqlc-gen-fsharp
.. _sqlc-gen-kotlin: https://github.com/sqlc-dev/sqlc-gen-kotlin
.. _sqlc-gen-python: https://github.com/sqlc-dev/sqlc-gen-python
.. _sqlc-gen-typescript: https://github.com/sqlc-dev/sqlc-gen-typescript

Future language support
************************

- `C# <https://github.com/sqlc-dev/sqlc/issues/373>`_
