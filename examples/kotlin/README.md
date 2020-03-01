# Kotlin examples

This is a Kotlin gradle project configured to compile and test all examples. Currently tests have only been written for the `authors` example.

To run tests:

```shell script
docker run --name dinosql-postgres -d -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=postgres -p 5432:5432 postgres:11
./gradlew clean test
```

The project can be easily imported into Intellij.

1. Install Java if you don't already have it
1. Download Intellij IDEA Community Edition
1. In the "Welcome" modal, click "Import Project"
1. Open the `build.gradle` file adjacent to this README file
1. Wait for Intellij to sync the gradle modules and complete indexing
