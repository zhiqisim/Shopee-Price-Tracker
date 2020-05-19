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
    scheduler.add_job(price_job, 'cron', hour='0-23', minute= '35')
    # scheduler.add_job(price_job, 'interval', hours=1)
    print('Press Ctrl+{0} to exit'.format('Break' if os.name == 'nt' else 'C'))

    try:
        scheduler.start()
    except (KeyboardInterrupt, SystemExit):
        pass