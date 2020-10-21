"# restfullapi" 

NOTE:- I have used the inbuilt salt feature for generating hash. I am not maintaing salt manually.

1 - Develop a Restful service for user creation. The service provides the following APIs 
1.1 Create user -
 Performs payload validation,
  Password should contain one Capital letter, one special character, and a number
  User name should not be more than 10 characters,
  the email should be a valid email Id
Creates the user by inserting user details into the user table and password details in the password table.


1.3 Login API -
  Takes input email and password
  Compares hash from the password table for the user and the hash generated for the password from the payload, If matches send JWT token, else 400 Bad Request
2. The database contain two tables
2.1 User Table which contains [ user Id, Name, email, address, phone ]
2.2 Password Table which contains [User Id (primary key of user table), password_hash, status ( active, inactive) ]
