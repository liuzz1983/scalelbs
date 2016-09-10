import urllib 
import urllib2 
import json
url = 'http://127.0.0.1:8000/search' 
base_lat = 23.0
base_lng = 45.0

values = {'lat' : 23.0, 'lng' :24.0, 'id' : 'liu' } 

for i in range(100):

	values['lat'] = base_lat + i*0.01
	values['lng'] = base_lng + i* 0.01
	values['id'] = str(i)
	
	data = urllib.urlencode(values) 
	req = urllib2.Request(url, data) 
	response = urllib2.urlopen(req) 
	the_page = response.read()
	values = json.loads(the_page)

	print len(values["Data"]["Points"])
