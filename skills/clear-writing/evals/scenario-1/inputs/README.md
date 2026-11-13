# Ledger

Ledger records every payment your service makes. It stores each
payment once and keeps it for seven years.

## Install

Run the installer with the version you want. It needs Go 1.23 or
newer.

```bash
go install example.com/ledger/cmd/ledger@v2.4.0 # installs the CLI into your Go bin folder, which must be on your PATH for the commands below to work
```

## How it works

Ledger is designed to seamlessly leverage the write-ahead log of
your database in order to capture every payment the moment it is
committed, which means that no payment is lost even when the
service crashes, because the log is replayed on start, and the
replay skips payments that are already stored.

Payments are grouped into daily files. Each file is named after its
date, such as `2026-03-14.ledger`.

## Retention

When a file is older than seven years, which is the period the tax
office requires, and no audit has put a hold on it, Ledger deletes
it during the nightly cleanup, so you do not need to delete old
files yourself.
