# User

type

- webauthn user
- roles number corresponding of bytes associte to Permissions
- Email
- password backup to webauthn
- Incredentials webauthn credentials  in database and credentials are charge in code

# Role

type

- name of the role
- byte number

# api request

**POST** _register/start/:username_ begin registration and send Incredentials to user create a new session form user

**POST** _register/end/:username_ finish the registration and return User Incredentials and save in database

**POST** _register/password/:username_ register with password and username

**POST** _login/start/:username_ begin the login and send Incredentials to user create a new session form the user if not exist

**POST** _login/end/:username_ finish the login and return User Incredentials and update in database

**POST** _login/password/:username_ login with password and username

**user** path prefix for user

- **GET** return the user
- **GET** _logout_ force close api session
- **Patch** take a JSON body and update a user
- **DELETE** delete a user
- **DELETE** _cred_ remove all login Incredentials

# Run the project

go version go1.19.2
docker version 1.1.4
docker compose version v2.12.0

```sh
docker compose up -d
go mod tidy
go run .
```

New terminal or same

```sh
cd front
npm i 
npm run serve
```
