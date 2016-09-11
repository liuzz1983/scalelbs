import urllib 
import urllib2 

import random

url = 'http://127.0.0.1:8000/add' 
base_lat = 23.0
base_lng = 45.0

values = {'lat' : 23.0, 'lng' :24.0, 'id' : 'liu' } 

for i in range(10000):

	values['lat'] = base_lat + random.randint(0,1000)*0.001
	values['lng'] = base_lng + random.randint(0,1000)* 0.001
	values['id'] = str(i)
	data = urllib.urlencode(values) 
	req = urllib2.Request(url, data) 
	response = urllib2.urlopen(req) 
	the_page = response.read()
	if i%1000 == 0 :
		print the_page
