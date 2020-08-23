import sqlConnector
import mysql.connector
from datetime import datetime
import json
import requests
import logging
import re

import pytz


logging.basicConfig(level=logging.INFO, filename='logs/app.log', format='%(asctime)s -%(levelname)s - %(message)s')
# Scheduled Job to run hourly
def price_job():
    singapore=pytz.timezone('Asia/Singapore')
    logging.info("-----------Start time of job: " + singapore.localize(datetime.now()).strftime('%Y-%m-%d %H:%M:%S') + "-----------")
    get_flash_deal()
    update_items()
    logging.info("-----------End time of job: " + singapore.localize(datetime.now()).strftime('%Y-%m-%d %H:%M:%S') + "-----------")

# Shopee API 
api = 'https://shopee.sg/api/v2'

def get_flash_deal():
    """
        Obtain all items in flash deal and store them in DB if not exist
    """
    try:
        url = api + '/flash_sale/get_items'
        req = requests.get(url)
    except requests.exceptions.HTTPError as err:
        logging.exception("Failed to fetch Shopee API :{}".format(err))
    
    sql = 'INSERT INTO item(item_id, shop_id, item_name, price) VALUES(%s, %s, %s, %s) ON DUPLICATE KEY UPDATE item_id = item_id' 
    
    flash_items_data = req.json()['data']['items']
    my_data = []
    logging.info("-----------START: Updating item data with new Flash Deals!-----------")
    for item in flash_items_data:
        item_id = item['itemid']
        shop_id = item['shopid']
        # try del special chars to allow insertion
        item_name = re.sub(r'([^\s\w]|_)+', '', item['name'].strip().encode("utf-8"))
        item_price = item['price']
        temp = [item_id, shop_id, item_name, item_price]
        my_data.append(tuple(temp))
    try:
        mydb = sqlConnector.connection_pool.get_connection() 
        mycursor = mydb.cursor(prepared=True)
        mycursor.executemany(sql, my_data)
    except mysql.connector.Error as err:
        logging.exception("Failed FlashSales insert :{}".format(err)) 
    finally:
        mydb.commit()
        logging.info('%s new flash deals items inserted into item table' %mycursor.rowcount)
        mycursor.close()
        del mycursor
        mydb.close()
        del mydb
        logging.info("-----------END: Updating item data with new Flash Deals!-----------")


def update_items():
    """
        Update price changelog of each item in the list
    """
    url = api + '/item/get'
    try:
        mydb = sqlConnector.connection_pool.get_connection()  
        mycursor = mydb.cursor(prepared=True)
        mycursor.execute("SELECT item_id, shop_id, price FROM item")
        myresult = mycursor.fetchall()
        logging.info("-----------START: Updating item_price data!-----------")
        updateCount = 0
        insertCount = 0
        for x in myresult:
            item_id = x[0]
            shop_id = x[1]
            price = x[2]
            req = requests.get(url, params={'itemid': item_id, 'shopid': shop_id})
            req.raise_for_status()
            json_data = json.loads(req.text)
            item = json_data['item']
            if item is None:
                logging.error("Item not found! ItemID = %s, ShopID = %s", item_id, shop_id)
                continue
            fetched_item_price = item['price']
            fetched_item_name = item['name']
            fetched_flash_sale = False if item['flash_sale'] == None else True
            singapore=pytz.timezone('Asia/Singapore')
            fetched_datetime = singapore.localize(datetime.now()).strftime('%Y-%m-%d %H:%M:%S')
            # update price on item table and item_price table if price changed
            # if price check if there is an entry already in the item_price table if not still insert
            if price != fetched_item_price:
                # update price in item table
                sql = "UPDATE item SET price = %s WHERE item_id = %s"
                adr = (fetched_item_price, item_id)
                mycursor.execute(sql, adr)
                # add new entry to item_price table
                sql = 'INSERT INTO item_price(item_id, price_datetime, price, flash_sale) VALUES(%s, %s, %s, %s)'
                adr = (item_id, fetched_datetime, fetched_item_price, fetched_flash_sale)
                mycursor.execute(sql, adr)
                logging.debug('Updated price of %s' %fetched_item_name)
                updateCount+=1
            else:
                mycursor2 = mydb.cursor(prepared=True)
                sql = 'SELECT item_id FROM item_price WHERE item_id = %s'
                adr = (item_id, )
                mycursor2.execute(sql, adr)
                records = mycursor2.fetchall()
                if mycursor2.rowcount == 0:
                    sql = 'INSERT INTO item_price(item_id, price_datetime, price, flash_sale) VALUES(%s, %s, %s, %s)'
                    adr = (item_id, fetched_datetime, fetched_item_price, fetched_flash_sale)
                    mycursor2.execute(sql, adr)
                    logging.debug('Inserted new price of %s' %fetched_item_name)   
                    insertCount+=1
    except requests.exceptions.HTTPError as err:
        logging.exception("Failed to fetch Shopee API :{}".format(err))
    except mysql.connector.Error as err:
        logging.exception("Failed updating item_price table :{}".format(err))
    finally:
        mydb.commit()
        mycursor.close()
        del mycursor
        mycursor2.close()
        del mycursor2
        mydb.close()
        del mydb
        logging.info('%s items added into item_price table with their latest price' %updateCount)
        logging.info('%s new items inserted into item_price with their first price changelog'  %insertCount)
        logging.info("-----------END: Updating item_price data!-----------")