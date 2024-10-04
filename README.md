# gtool

[🇷🇺Русская версия](README.ru.md)

This application is designed for interaction with JIRA and Active Directory systems and offers a range of tools for task automation.

## Why?

During my work, I noticed repetitive, routine tasks that were performed with just a few simple steps. This sparked the idea of automating the execution of these tasks.

## Key Features

- **JIRA Task Management**: Update employee status in CMDB Insight, assign, update, and close tasks related to employee onboarding/offboarding, and grant access.
- **CMDB Object Management**: Generate information for equipment write-offs and display descriptions of laptops for subsequent courier service orders.
- **LDAP Management**: Manage user groups and their access rights.

## Installation

1. Clone the repository

```
git clone https://github.com/gptlv/gtool.git
cd gtool
```

2. Create a `.env` file and provide the necessary data (example is in `.env.example`):

```
touch .env
```

3. Fill in the `config.yml` configuration file. You need to provide up-to-date employee surnames for document generation.

4.	Run `go mod tidy` to ensure all dependencies are properly managed.

5. Build the project:

```
go build -o gtool .
```

## Available Commands

### JIRA Ticket Processing

- `./gtool issue process-insight` Process user deactivation requests in Insight.
- `./gtool issue process-ldap` Process user deactivation requests in Active Directory.
- `./gtool issue process-staff --component=[all|hiring|dismissal]` Process employee hiring and/or dismissal requests.
- `./gtool issue grant-access --key=<key>` Process access request (add group in Active Directory).
- `./gtool issue assign --component=[all|hiring|dismissal|insight|ldap]` Assign requests to the current user.
- `./gtool issue update-trainee --key=<key>` Update subtask names related to intern offboarding.
- `./gtool issue show-empty` Show tickets that require specifying a component.

### CMDB Queries

- `./gtool asset generate-records --start=<id>` Generate a `.csv` file with necessary data for use in the [wroffs](https://github.com/gptlv/wroffs) script.
- `./gtool asset print-description --isc=<isc>` Display a laptop description to add to a courier service order.

### Active Directory Queries

- `./gtool ldap add-group -emails <user1@ex.com user2@ex.com> -cns <cn1 cn2 cn3>` Add multiple users to multiple groups in Active Directory.

### Main Libraries

- [go-arg](https://github.com/alexflint/go-arg) for command-line argument parsing.
- [go-jira](https://github.com/andygrunwald/go-jira) or JIRA API interaction. I used my [own fork](https://github.com/gptlv/go-jira), which extends functionality for working with JIRA Insight
- [go-ldap](https://github.com/go-ldap/ldap) for LDAP operations.
