# Data Stores

## badger

[badger](https://github.com/dgraph-io/badger) is the default database driver for `cfd`. Badger is an embeddable and fast key-value database written in Go.

Use this database if using a local development machine for simplicity (you don't need to modify the configuration file to start using it). If you can live with just one server for the API, this setup will be reliable enough for a moderated size of clients.

**Backup**: Copy full content of the selected database directory.

**Restore**: Just copy the back up back and start the service.

>The service needs to be stopped to ensure back up and restore, `badger` *locks data files* making the process to (normally) fail.

#### Configuration

All configuration resides on the `config.yaml` file (or whatever other name you selected).

>[!INFO]
>`type` must be set to `badger`.
>
>`connection` is the directory where the data will be stored.

```yaml
db:
  connection: /your_home_directory/.cfd/db
  type: badger
```

## firebase (firestore)

[firebase](https://firebase.google.com/) is a backend service by Google that helps building applications.

`firebase` has generous limits on its free plan so unless you are creating a big amount of certificates it will remain free *(you need to store more than 1Gb of data, write more than 20.000 docs/day or read more than 50.000 docs/day to pass the limits)*. Review firebase [pricing](https://firebase.google.com/pricing) to get full information.

?>This database allows High Availability so you can run several API servers that will work with the same database storage (usually behind a load balancer).

#### Database creation

- First, you need to register in firebase, if you have a Google account it will be a straigh forward step.

- Go to the Console, and you will be prompted with a 'Create a new Project', click it. A wizard will guide you in the process.

- **Select a project name.**
![Step 1. Select a project name](./images/firestore-1.png)

- **Add Analytics if you think you will need them**
If you select to use analytics you will be asked for the project to use or set up a new one. *This is not covered here.*
![Step 2. Analytics needed?](./images/firestore-2.png)

- **Wait a few moments to project creation**
![Step 3. Project Creation](./images/firestore-3.png)

- **Go 'Cloud Firestore', Then press 'Create database'**
![Step 4. Create database](./images/firestore-4.png)

- **Select 'Start in production mode' to secure the data store from guests**
![Step 5. Production mode](./images/firestore-5.png)

- **Select a location near to you.** *This cannot be changed later*
![Step 6. Location](./images/firestore-6.png)

- **Database is created. Go to Project Settings**
![Step 7. Go to project settings](./images/firestore-7.png)

- **Get the Firebase Admin SDK file**
  - Go to 'Service Accounts'
  - Press 'Generate new private key'
  - Save the file downloaded which will be used to identify your cfd client. Called something like `project-name-firebase-admin-sdk.json`

![Step 8. Select a project name](./images/firestore-8.png)

#### Configuration

All configuration resides on the `config.yaml` file (or whatever other name you selected).

>[!INFO]
>`type` must be set to `firestore`.
>
>`connection` is the file with firebase admin SDK configured (the one you downloaded from the service account configuration on Firebase Console).
