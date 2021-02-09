
# Installation

Every release is published on [github](https://github.com/fernandezvara/certsfor/releases), where you can download the appropiate file for your OS/Architecture and add it to a directory on your PATH.

To allow a better lifecycle for the software software can be installed using several package managers.

## Homebrew (macOs/linux)

**Homebrew** is a package manager, initially created for macOs, that allows the installation of free and open-source software directly from the terminal. The linux version can be used on the Windows Subsystem for Linux (WSL) too.

Please refer to the [instructions](https://brew.sh/) for help on installation and usage.

Installation process (example on macOs):

```bash
> # ensure you have brew installed on your machine
> brew tap back-io/tap
> brew install cfd
Updating Homebrew...
==> Installing cfd from backd-io/tap
==> Downloading https://github.com/fernandezvara/certsfor/releases/download/v0.2.2/cfd_0.2.2_Darwin_x86_64.tar.gz

:beer:  /usr/local/Cellar/cfd/0.2.2: 5 files, 21.8MB, built in 3 seconds

```

## Scoop (Windows)

**Scoop** is a package manager for the Windows command line that allows to install familiar Unix tools. Scoop installs programs on the home directory by default, to ensure a limited set of permissions.

Please refer to the [instructions](https://scoop.sh/) for help on installation and usage.

> [!NOTE|label:Scoop dependencies]
> Installing new buckets in Scoop require some dependencies, like `git` that need to be present on the system. If not found it will be requested.

```cmd
c:\Users\demo>scoop bucket add backd https://github.com/backd-io/scoop-bucket
c:\Users\demo>scoop install cfd

```

## RPM / DEB / APK (linux)

Packages for RedHat, Debian and Alpine (and others using these package managers) are created for easy install. But there are not a repository to import on the system and maintain the software lifecycle still.

Please download the appropiate file from the software repository. *Publication is in roadmap.*
