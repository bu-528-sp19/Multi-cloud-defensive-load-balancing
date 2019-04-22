import time
import urllib.request, json
import requests

#data = {"UserID": 1, "Model": "Model-T"}
NUM_REQUESTS = 200

start = time.time()
failed = 0
for i in range(NUM_REQUESTS):
   innernow=time.time()
   data = {"UserID": 1, "Model": str(i)}
   json_data = json.dumps(data)
   r = requests.post('http://35.245.85.172/cars/', json_data)
   #r = requests.post('http://35.245.114.96:8888/cars/', json_data)
   stat=r.status_code
   innerthen=time.time()
   if (stat == 200):
       print((innerthen-innernow))
   else:
       failed= failed+1
end = time.time()
print(NUM_REQUESTS,'requests in', end - start, 'seconds')
print('Requests per second:', NUM_REQUESTS / (end - start))
print('failed :',failed)
