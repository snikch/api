# Go API Framework

Note: This is not even Alpha software.

API View Controller Toolkit.

## Ethos

Composable. Pick and choose what you want to use, nothing is absolute. You're just writing Go. HTTP handlers everywhere. No unknown magic.

# Packages


* [changes](https://github.com/snikch/api/tree/master/changes) Generate diffs between type instances for audit logs, and update management.

* [ctx](https://github.com/snikch/api/tree/master/ctx) Tightly coupled, lockable contexts used in most packages.

* [fail](https://github.com/snikch/api/tree/master/fail) Return intelligent, api friendly, and log friendly errors.

* [lifecycle](https://github.com/snikch/api/tree/master/lifecycle) Manage the lifecycle of your application, e.g. shutdown callbacks.

* [log](https://github.com/snikch/api/tree/master/log) Sane defaults for logging via Logrus.

* [sideload](https://github.com/snikch/api/tree/master/sideload) Automatically load related entities

* [lynx](https://github.com/snikch/api/tree/master/lynx) Encrypt and decrypt your data as required.

* [vc](https://github.com/snikch/api/tree/master/vc) Handle request and response lifecycle, including encryption, sideloading and rendering.


# TODO

- [x] Sideloading of related entities
- [ ] Types for actors, users and parameters
- [x] View Helpers for rendering responses and errors
