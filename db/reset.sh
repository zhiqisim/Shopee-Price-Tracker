## RUN this bash script when you want to reset the database.
## DELETES volume of storage for the databases.

# remove storage folder in items-db
rm -rf items-db/storage
# remove storage folder in user-db
rm -rf user-db/storage