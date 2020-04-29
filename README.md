# Moodle Toolkit
### Description
It's a toolkit for Moodle, mostly written in crossplatform language - Go
### How to build ?
1. Visit [Go "Getting Started"](https://golang.org/doc/install) page
2. Choose suitable for your OS package, download and install it by the guide given on that page
3. On linux - run `chmod 1777 build.sh && ./build.sh`, on Windows - run Powershell script `build.ps`
### How to run ?
Your executable file will be available at `<package_name>.moodler/build/<os>/<arch>`
For example:
```
user@pc $> cd online.moodler/build/linux/86
user@pc $> ./online.moodler.run
```
##### Before you run: Set your Moodle domain and credentials in `settings.json` file
### FAQ
#### Why I get "Authorizing... Failed" ?
Check your settings and Internet connection.
Are you sure that your credentials are correct ?
### Contacts
 - [Telegram](https://t.me/dimankiev)
