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

========  =================================  ===============  ===============  ===============
Language  Plugin                             MySQL            PostgreSQL       SQLite
========  =================================  ===============  ===============  ===============
C#        `DaredevilOSS/sqlc-gen-csharp`_    Stable           Stable           Stable
F#        `kaashyapan/sqlc-gen-fsharp`_      N/A              Beta             Beta
Java      `tandemdude/sqlc-gen-java`_        Beta             Beta             N/A 
PHP       `lcarilla/sqlc-plugin-php-dbal`_   Beta             N/A              N/A    
Ruby      `DaredevilOSS/sqlc-gen-ruby`_      Beta             Beta             Beta           
Zig       `tinyzimmer/sqlc-gen-zig`_         N/A              Beta             Beta            
[Any]     `fdietze/sqlc-gen-from-template`_  Stable           Stable           Stable
========  =================================  ===============  ===============  ===============

.. _sqlc-gen-go: https://github.com/sqlc-dev/sqlc-gen-go
.. _kaashyapan/sqlc-gen-fsharp: https://github.com/kaashyapan/sqlc-gen-fsharp
.. _sqlc-gen-kotlin: https://github.com/sqlc-dev/sqlc-gen-kotlin
.. _sqlc-gen-python: https://github.com/sqlc-dev/sqlc-gen-python
.. _sqlc-gen-typescript: https://github.com/sqlc-dev/sqlc-gen-typescript
.. _DaredevilOSS/sqlc-gen-csharp: https://github.com/DaredevilOSS/sqlc-gen-csharp
.. _DaredevilOSS/sqlc-gen-ruby: https://github.com/DaredevilOSS/sqlc-gen-ruby
.. _fdietze/sqlc-gen-from-template: https://github.com/fdietze/sqlc-gen-from-template
.. _lcarilla/sqlc-plugin-php-dbal: https://github.com/lcarilla/sqlc-plugin-php-dbal
.. _tandemdude/sqlc-gen-java: https://github.com/tandemdude/sqlc-gen-java
.. _tinyzimmer/sqlc-gen-zig: https://github.com/tinyzimmer/sqlc-gen-zig

Future language support
************************

