import time
import urllib.request, json
import requests

start= time.time()
failed=0
for i in range(200):
   innernow=time.time()
   r = requests.get('http://35.245.85.172/users/')
   innerthen=time.time()
   stat=r.status_code
   if (stat == 200 and len(r.text) > 0):
       print((innerthen-innernow))
   else:
       failed= failed+1
end = time.time()
print('200 requests in', end - start, 'seconds')
print('Requests per second:', 200 / (end - start))
print('failed :',failed)
