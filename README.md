# its-a-key-value-database
This is an in-memory key value database similar to Redis (but a _bit_ more basic).
The program works by creating a tcp server on localhost:3332 that can then be connected to by a client. 
 - Commands are entered in the format **verb key value**
 - Error responses are preceded by a "-" Ex: "-ERR key already exists\r\n" 
 - Non-error responses are preceded by a "+" Ex: "+OK\r\n"
 - All responses are terminated with "\r\n" Ex: "+OK\r\n"

## Commands:
#### **SET** key value
Adds a new entry to the database.
```
> SET user1 Aden
> +OK
> SET user1 Eilers
> -ERR key already exists
```
#### **GET** key
Retrives the value associeated with the given key.
```
> GET user1
> +Aden
> GET user2
> -ERR key does not exist
```
#### **GETALL**
Retrieves all keys and values in the database.
```
> GETALL
> +user1
> Aden
> user2
> Mike
> GETALL
> -ERR no entries in database
```
#### **UPDATE** key value
Changes the value associated with the given key.
```
> UPDATE user1
> +OK
> UPDATE user2
> -ERR key does not exist
```
#### **DELETE** key
Deletes the entry associated with the given key.
```
> DELETE user1
> +OK
> DELETE user2
> -ERR key does not exist
```
#### **LEN**
Retrieves the total number of entries in the database.
```
> LEN
> +5
> LEN
> +0
```
#### **GETVALUES**
Retrieves all values in the database.
```
> GETVALUES
> +Aden
> Eilers
> Mike
> GETVALUES
> -ERR no entries in database
```
#### **GETKEYS**
Retrieves all keys in the database.
```
> GETKEYS
> +user1
> user2
> user3
> GETKEYS
> -ERR no entries in database
```
#### **EXISTS** key
Checks if the given key exists in the database. 1 means 'yes', 0 means 'no'.
```
> EXISTS user1
> +1
> EXISTS user2
> +0
```
#### **PING**
Pings the database to check if the client is still connected.
```
> PING
> +PONG
```
