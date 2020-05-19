from concurrent import futures
import logging

import grpc

import item_pb2
import item_pb2_grpc
from sqlConnector import initSQL

import mysql.connector

import datetime


class ItemService(item_pb2_grpc.ItemServiceServicer):

    def ListAllItems(self, request, context):
        arr = []
        try:
            mydb = initSQL()  
            mycursor = mydb.cursor(prepared=True)
            mycursor.execute("SELECT item_id, shop_id, item_name, price FROM item")
            myresult = mycursor.fetchall()
            logging.info("Compiling query!")
            for x in myresult:
                item_id = str.encode(str(x[0]))
                shop_id = str.encode(str(x[1]))
                item_name = str.encode(str(x[2]))
                item_price = x[3]
                item = item_pb2.Item(item_id = item_id, shop_id = shop_id, item_name = item_name, item_price = item_price)
                arr.append(item)
        except mysql.connector.Error as err:
            logging.error("Failed ListAllItem Select :{}".format(err))
            return item_pb2.ListAllItemsResponse(message="error")   
        finally:
            mycursor.close()
            del mycursor
            mydb.close()
            del mydb
        logging.info("ListAllItems Success!")
        return item_pb2.ListAllItemsResponse(message="success", items= arr)

    def ItemPrice(self, request, context):
        arr = []
        try:
            mydb = initSQL()  
            item_proto = item_pb2.Item()
            mycursor = mydb.cursor(prepared=True)
            sql = "SELECT price_datetime, price, flash_sale FROM item_price WHERE item_id=%s ORDER BY price_datetime DESC"
            adr = (request.item_id, )
            mycursor.execute(sql, adr)
            myresult = mycursor.fetchall()
            for x in myresult:
                item_date = str.encode(str(x[0]))
                item_price = x[1]
                flash_sale = x[2]
                price = item_pb2.ItemPrice(price_datetime = item_date, price = item_price, flash_sale = flash_sale)
                arr.append(price)
        except mysql.connector.Error as err:
            logging.error("Failed ItemPrice Select :{}".format(err))
            return item_pb2.ItemPriceResponse(message="error")   
        finally:
            mycursor.close()  
            del mycursor
            mydb.close()
            del mydb
        logging.info("ItemPrice Success!")
        return item_pb2.ItemPriceResponse(message="success", itemPrice= arr)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    item_pb2_grpc.add_ItemServiceServicer_to_server(ItemService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    logging.info("Serving gRPC server on port 50051")
    serve()
