# Authentication Routes

We have two authentication routes to support password-less login. The first creates a
token which is returned to the user which can then be signed. The second checks that the
signed token is signed by the correct public key.

See the flows below for more details.

## CreateLoginSession

CreateLoginSession creates a token and session which can then be used to log in
via the CreateAuthSession route.

To defend against enumeration attacks (e.g. a attacker trying every combination
of email to figure out which accounts exist and which don't) we've implemented
a feature where we'll always return some salt (random data) to the attacker. This
salt will not actually be useful but they won't be able to tell which accounts
exist and which don't.

Here's the flow of the route for reference:

- Find user by email, if not found return garbage and continue, if found continue
- Generate token for this user
- Use token to create a session
- If user actually exists, add session to the database
- Return fake or real session

## CreateAuthSession

CreateAuthSession finishes the log in flow by checking that the login token
has been properly signed by a user's private key.

Here's the flow for reference:

- Once session and user have been obtained, load user credentials
- Verify the token has been properly signed
- Generate a new authentication token
- Create an authentication session using above token
- Add new session to the database and return it
