import urllib 
import urllib2 
import json
import time
import geopy

from geopy.distance import vincenty


url = 'http://127.0.0.1:8000/search' 
base_lat = 23.0
base_lng = 45.0

values = {'lat' : 23.0, 'lng' :24.0, 'id' : 'liu' } 

for i in range(100):

	values['lat'] = base_lat + i*0.01
	values['lng'] = base_lng + i* 0.01
	values['id'] = str(i)
	
	begin = time.time()
	data = urllib.urlencode(values) 
	req = urllib2.Request(url, data) 
	response = urllib2.urlopen(req) 
	the_page = response.read()

	cost1 = time.time() - begin
	result = json.loads(the_page)
	#

	max_ = 0
	for p in result["Data"]["points"]:

		dis = vincenty((values[u'lat'],values[u'lng']),(p[u'lat'],p[u'lng']))
		if dis > max_:
			max_ = dis

	print len(result["Data"]["points"]),max_, cost1, time.time()-begin