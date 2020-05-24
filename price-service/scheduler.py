from datetime import datetime
from datetime import timedelta
import os
import time

from pricetracker import price_job
from apscheduler.schedulers.blocking import BlockingScheduler

import logging

logging.basicConfig(filename='logs/app.log', format='%(asctime)s -%(levelname)s - %(message)s')
logging.getLogger('apscheduler').setLevel(logging.INFO)



if __name__ == '__main__':
    scheduler = BlockingScheduler()
    # # Cron scheduling
    # scheduler.add_job(price_job, 'cron', hour='0-23', minute= '10')
    # Interval scheduling
    start_time = datetime.now() + timedelta(seconds=60)
    scheduler.add_job(price_job, 'interval', hours=1, start_date=start_time)
    print('Press Ctrl+{0} to exit'.format('Break' if os.name == 'nt' else 'C'))
    try:
        scheduler.start()
    except (KeyboardInterrupt, SystemExit):
        pass