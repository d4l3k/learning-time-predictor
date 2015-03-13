import caffe
import urllib
import cStringIO as StringIO
import os
from dateutil import parser
from visual import * # must import visual or vis first
from visual.graph import *	# import graphing features


print("########### CLASS")
net = caffe.Classifier('./traffic_deploy.prototxt', './snapshot_iter_11000.caffemodel')
print("########### CLASSDONE")

def classify(imageurl):
    #imageurl = 'images/2015-02-04_20:36:01.346726_+0000_UTC.jpg'
    string_buffer = StringIO.StringIO(
    urllib.urlopen(imageurl).read())
    image = caffe.io.load_image(string_buffer)
    scores = net.predict([image]).flatten()
    return scores

total = 0.0
totalDrift = 0
count = 0.0

timeI = [0.0]*24
timeC = [0.0]*24

maxa = 0
mina = 0

for filename in os.listdir('images')[-1000:]:
    prediction = classify('images/'+filename)
    parts = filename.split('_')
    time = parser.parse(parts[1]).time().hour
    if prediction[0] > maxa:
        maxa = prediction[0]
    if prediction[0] < mina:
        mina = prediction[0]
    #Adjusted = (prediction[0] + 0.50477093)/3.787515
    Adjusted = (prediction[0] + 0.6879)/(3.3925991+0.6879)
    pTime = Adjusted * 24.0
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
print(maxa, mina)

