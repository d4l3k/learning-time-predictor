import caffe
import urllib
import cStringIO as StringIO
import os
from dateutil import parser


net = caffe.Classifier('./traffic_deploy.prototxt', './snapshot_iter_10000.caffemodel')

def classify(imageurl):
    #imageurl = 'images/2015-02-04_20:36:01.346726_+0000_UTC.jpg'
    string_buffer = StringIO.StringIO(
    urllib.urlopen(imageurl).read())
    image = caffe.io.load_image(string_buffer)
    scores = net.predict([image]).flatten()
    return (-scores).argsort()

total = 0.0
count = 0.0

timeI = [0.0]*24
timeC = [0.0]*24

for filename in os.listdir('images'):
    prediction = classify('images/'+filename)
    parts = filename.split('_')
    time = parser.parse(parts[1]).time().hour
    pTime = prediction[1]
    diff = abs(time - pTime)
    if diff > 12:
        diff -= 12
    timeI[pTime] += diff
    timeC[pTime] += 1
    total += diff
    count += 1
    print("D {}, T {}, P {}, A {}, B {}".format(diff, time, prediction[0:10], total/count, timeI[pTime]/timeC[pTime]))
