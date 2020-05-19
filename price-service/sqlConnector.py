import mysql.connector

def initSQL():
    try:
        hostname = "items-db"
        return mysql.connector.connect(
            host=hostname,
            port="3306",
            user="root",
            passwd="root",
            database="itemdb"
            )
    except mysql.connector.Error as e:
        print(e)