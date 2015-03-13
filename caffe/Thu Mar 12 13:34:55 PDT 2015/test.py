import caffe
import urllib
import cStringIO as StringIO
import os
from dateutil import parser
from visual import * # must import visual or vis first
from visual.graph import *	# import graphing features


print("########### CLASS")
net = caffe.Classifier('./traffic_deploy.prototxt', './snapshot_iter_500.caffemodel')
print("########### CLASSDONE")

def classify(imageurl):
    #imageurl = 'images/2015-02-04_20:36:01.346726_+0000_UTC.jpg'
    string_buffer = StringIO.StringIO(
    urllib.urlopen(imageurl).read())
    image = caffe.io.load_image(string_buffer)
    scores = net.predict([image]).flatten()
    print(scores)
    return (-scores).argsort()

total = 0.0
count = 0.0

timeI = [0.0]*24
timeC = [0.0]*24

for filename in os.listdir('images')[-500:]:
    prediction = classify('images/'+filename)
    parts = filename.split('_')
    time = parser.parse(parts[1]).time().hour
    pTime = prediction[1]
    diff = abs(time - pTime)
    if diff > 12:
        diff = 24 - diff
    timeI[time] += diff
    timeC[time] += 1
    total += diff
    count += 1
    print("C {}: D {}, T {}, P {}, A {}, B {}".format(count, diff, time, prediction[0:10], total/count, timeI[time]/timeC[time]))

f1 = gcurve(color=color.cyan)	# a graphics curve
for x in range(0, 24):
    if timeC[x] > 0:
        f1.plot(pos=(x, timeI[x]/timeC[x]))
print(timeI)
print(timeC)
print([x/y for x, y in zip(timeI, timeC)])

