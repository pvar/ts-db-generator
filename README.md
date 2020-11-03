ts-db-generator
===============

#### Description

ts-db-generator generates a sqlite database that contains all known timezones and the corresponding zone transitions.
If the database exists, the program attempts to update stored data. That is, to create new tables for the updated
zone transitions and to add/remove original timezones and replicas (links), as appropriate.

#### Conditions for successful update
1. Ammount of new timezones should not supersede 5% of the ammount of stored ones.
2. Ammount of replicas should not supersede 5% of the ammount of stored ones.
3. Version of parsed TZdata should be newer that the version of stored data.
4. If parsed and stored data are of the same version, the ammount of new zones <br>
   and the ammount of stored zones for each timezone should be the same.
5. The ammount of new zones should not supersede 5% of the stored zone for a given timezone.

If any of the conditions 1, 2 and 3 is not met, the update proceedure aborts.
If any of the conditions 4 an5 is not met, the update of a given timezone is skipped, but the overall process will continue.

#### Invocation

ts-db-genmerator has an optional parameter that specifies the name of the database file to work with.

`./ts-db-generator {db_filename}`

If the parameter is omitted, the program used the default database name (`tsdb.sqlite`).
