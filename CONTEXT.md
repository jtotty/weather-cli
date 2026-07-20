# Domain Glossary

## Credential Store

The single answer to "where does the API key come from".
Owns retrieval, storage, deletion, availability, and precedence between sources.
Not responsible for how a human enters a key - that is the Setup Wizard's job.

## Store Adapter

A concrete implementation of the Credential Store interface.
Production uses the OS keyring adapter; tests use the in-memory fake.

## Env Override

The rule that `WEATHER_API_KEY` beats the stored key when reading.
It affects reads only: setup still writes to the keyring and delete-key still deletes from it while the env var is set.
Precedence is a Credential Store concern, not a config concern.

## Availability

Whether the OS keyring can actually be used at runtime (e.g. present dbus/Secret Service on Linux), probed before prompting so users on unsupported systems are routed to the env var without typing a key.

## Setup Wizard

The interactive flow behind `weather-cli setup`: checks Availability, prompts for a key, stores it via the Credential Store.
