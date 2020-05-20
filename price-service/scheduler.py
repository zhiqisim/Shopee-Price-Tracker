from datetime import datetime
import os
import time

from pricetracker import price_job
from apscheduler.schedulers.blocking import BlockingScheduler

import logging

logging.basicConfig()
logging.getLogger('apscheduler').setLevel(logging.DEBUG)



if __name__ == '__main__':
    scheduler = BlockingScheduler()
    # scheduler.add_job(price_job, 'cron', hour='0-23', minute= '10')
    scheduler.add_job(price_job, 'interval', hours=1, start_date='2020-05-20 10:20:00')
    # scheduler.add_job(price_job, 'interval', seconds=3)
    print('Press Ctrl+{0} to exit'.format('Break' if os.name == 'nt' else 'C'))
    try:
        scheduler.start()
    except (KeyboardInterrupt, SystemExit):
        pass