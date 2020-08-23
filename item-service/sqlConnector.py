import mysql.connector
from mysql.connector import Error
from mysql.connector.pooling import MySQLConnectionPool
import logging
import time
import os


logging.basicConfig(level=logging.INFO, filename='logs/app.log', format='%(asctime)s -%(levelname)s - %(message)s')

dbPass = os.environ['DBPASS']
dbHost = os.environ['DBHOST']
dbPort = os.environ['DBPORT']

# def initSQL():
#     try:
#         hostname = "items-db"
#         return mysql.connector.connect(
#             host=hostname,
#             port="3306",
#             user="root",
#             passwd="root",
#             database="itemdb"
#             )
#     except Error as e:
#         print(e)

try:
    logging.info("Trying to connect to SQL_Connection_Pool!")
    time.sleep(60)
    connection_pool = MySQLConnectionPool(pool_name="SQL_connection_pool",
                                                                  pool_size=30,
                                                                  pool_reset_session=True,
                                                                #   connection_timeout=0,
                                                                  host=dbHost, 
                                                                  port=dbPort,
                                                                  user="root",
                                                                  passwd=dbPass,
                                                                  database="itemdb")
    logging.info("Connection Pool Name - %s | Connection Pool size - %s", connection_pool.pool_name, connection_pool.pool_size)
except Error as e :
    logging.error("Error while connecting to MySQL using Connection pool: %s", e)