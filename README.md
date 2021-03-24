# its-a-key-value-database
This is an in-memory key value database similar to Redis (but a _bit_ more basic).
The program works by creating a tcp server on localhost:3332 that can then be connected to by a client. 
Commands are entered in the format **verb key value**

The accepted commands:
* GET key
* SET key value
* UPDATE key value
* DELETE key

## Example inputs and outputs:
#### GET
```
> GET user1
> value
> GET user2
> ERR key does not exist
```
#### SET
```
> SET user1 Aden
> OK
> SET user1 Eilers
> ERR key already exists
```
#### UPDATE
```
> UPDATE user1
> OK
> UPDATE user2
> ERR key does not exist
```
#### DELETE
```
> DELETE user1
> OK
> DELETE user2
> ERR key does not exist
```
