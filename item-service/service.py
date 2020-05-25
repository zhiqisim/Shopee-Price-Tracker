from concurrent import futures
import logging

import json
import requests
import grpc

import item_pb2
import item_pb2_grpc
# from sqlConnector import initSQL
import sqlConnector

import mysql.connector

import datetime

class ItemService(item_pb2_grpc.ItemServiceServicer):
    def AddNewItem(self, request, context):
        url = 'https://shopee.sg/api/v2/item/get'
        try:
            req = requests.get(url, params={'itemid': request.item.item_id, 'shopid': request.item.shop_id})
            req.raise_for_status()
            json_data = json.loads(req.text)
            item = json_data['item']
            if item is None:
                logging.error("Item not found in Shopee's API!")
                return item_pb2.AddNewItemResponse(message="error")
            fetched_item_price = item['price']
            # Get connection object from a pool
            mydb = sqlConnector.connection_pool.get_connection() 
            cursor = mydb.cursor(prepared=True)
            sql = 'INSERT INTO item(item_id, shop_id, item_name, price) VALUES(%s, %s, %s, %s) ON DUPLICATE KEY UPDATE item_id = item_id'
            adr = (request.item.item_id, request.item.shop_id, request.item.item_name, fetched_item_price,)
            cursor.execute(sql, adr)
            logging.debug("Inserting item "+ request.item.item_name)
        except requests.exceptions.HTTPError as err:
            logging.exception("Failed to fetch Shopee API :{}".format(err))
            return item_pb2.ItemPriceResponse(message="error") 
        except mysql.connector.Error as err:
            logging.exception("Failed AddNewItem Insert :{}".format(err))
            return item_pb2.ItemPriceResponse(message="error")   
        finally:
            mydb.commit()
            cursor.close()  
            del cursor
            mydb.close()
            del mydb
            logging.info('Success! Insertion of item : %s', request.item.item_name)
        return item_pb2.ItemPriceResponse(message="success")

    def ListAllItems(self, request, context):
        arr = []
        logging.info("----------START: Listing all items----------")
        myresult = []
        try:
            # Get connection object from a pool
            mydb = sqlConnector.connection_pool.get_connection()  
            mycursor = mydb.cursor(prepared=True)
            sql = 'SELECT item_id, shop_id, item_name, price FROM item WHERE id > %s LIMIT %s'
            offset = str(request.offset)
            limit = str(request.limit)
            adr = (offset, limit, )
            mycursor.execute(sql, adr)
            myresult = mycursor.fetchall()
        except mysql.connector.Error as err:
            logging.exception("Failed ListAllItem Select :{}".format(err))
            return item_pb2.ListAllItemsResponse(message="error")   
        finally:
            logging.debug('Success! Obtained %s items' %mycursor.rowcount)
            mycursor.close()
            del mycursor
            mydb.close()
            del mydb
        for x in myresult:
            item_id = str.encode(str(x[0]))
            shop_id = str.encode(str(x[1]))
            item_name = x[2].decode('utf-8')
            item_price = x[3]
            item = item_pb2.Item(item_id = item_id, shop_id = shop_id, item_name = item_name, item_price = item_price)
            arr.append(item)
        logging.info("----------END: Listing all items----------")
        return item_pb2.ListAllItemsResponse(message="success", items= arr)

    def ItemPrice(self, request, context):
        arr = []
        logging.info("----------START: Obtaining price of item----------")
        myresult = []
        try:
            # Get connection object from a pool
            mydb = sqlConnector.connection_pool.get_connection()  
            item_proto = item_pb2.Item()
            mycursor = mydb.cursor(prepared=True)
            sql = "SELECT price_datetime, price, flash_sale FROM item_price WHERE item_id=%s ORDER BY price_datetime DESC"
            adr = (request.item_id, )
            mycursor.execute(sql, adr)
            myresult = mycursor.fetchall()
        except mysql.connector.Error as err:
            logging.exception("Failed ItemPrice Select :{}".format(err))
            return item_pb2.ItemPriceResponse(message="error")   
        finally:
            logging.info('Success! Obtained %s prices for item %s' %mycursor.rowcount %item_id)
            mycursor.close()  
            del mycursor
            mydb.close()
            del mydb
            for x in myresult:
                item_date = str.encode(str(x[0]))
                item_price = x[1]
                flash_sale = x[2]
                price = item_pb2.ItemPrice(price_datetime = item_date, price = item_price, flash_sale = flash_sale)
                arr.append(price)
        logging.info("ItemPrice Success!")
        logging.info("----------END: Obtaining price of item----------")
        return item_pb2.ItemPriceResponse(message="success", itemPrice= arr)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    item_pb2_grpc.add_ItemServiceServicer_to_server(ItemService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO, filename='logs/app.log', format='%(asctime)s -%(levelname)s - %(message)s')
    logging.info("Serving gRPC server on port 50051")
    serve()
