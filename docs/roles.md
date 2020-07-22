# Roles

This document describes `Roles` within Cape. From what they are, how you use them, and how they will affect the
user experience.

In code, `Role` can be found [here](../models/role.go)

## Overview

A `Role` is used to alter the permissions that a `User` in cape has. There are two different sets of roles depending
on your [`Operating Context`](./operating_context.md).

### Global Roles
#### `Admin`

An admin is much like a superuser in a classic linux environment. They can create accounts for team members,
revoke accounts, etc. See the `Permission Chart` at below for a comprehensive list.

#### `User`

A user is someone that uses and consumes information from Cape.

### Project Roles

For more information on projects, see [`Projects`](./projects.md). Your role within a project determines which
`project_actions` you can take

### `Owner`

Owners have permission to do all project actions.

### `Reviewer`

Reviewers can do anything a member or contributor can do, as well as accept [`Policy`](./policy.md) changes.

### `Contributor`

Contributors can do anything a member can do, as well as suggest [`Policy`](./policy.md) changes.

### `Member`

Members can read all project data, as well as view and use the [`Policy`](./policy.md) that is active within this project.


### Permission Matrix

-- TODO .. would be built around https://gist.github.com/kitschysynq/a1135369fc571a8c619cb1715509b3ec

## Examples

You interact with roles through the Cape CLI.

### View your role
```
$ cape roles me
```

### See all roles
```
$ cape roles list
Global Roles
============
admin
user

Project Roles
=============
project-owner
project-reviewer
project-contributor
project-member

```

### View members of a role (admin only)
To see who belongs to a role, use `$ cape roles members <role>` E.g.

```
$ cape roles members project-owner
```