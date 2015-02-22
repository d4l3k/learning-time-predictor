require 'time'
require 'pry'
require 'shellwords'

begin
  Dir.mkdir('images')
rescue
end

files = []
Dir['../output/*/*.*'].each do |file|
  escapedPath = Shellwords.shellescape(file)
  mime = `file --mime-type #{escapedPath}`.split(' ').last
  if !mime.start_with? 'image'
    next
  end
  newPath = "images/#{file.split('/').last.gsub(' ', '_')}"
  if !File.exists?(newPath)
    system("cp #{escapedPath} #{newPath}")
  end
  timeStr = file.split('/').last.split('.jpg').first.split()
  time = Time.parse(timeStr[1])
  hours = time.hour#*60 + time.min#/60.0 + time.sec/60.0/60.0

  files.push([newPath, hours])
end
files = files.map do |file|
  file.join(' ')
end

test_index = files.length - 500
train = files[0...test_index].join("\n")
test = files[test_index..(files.length)].join("\n")
File.write('images_train.txt', train)
File.write('images_test.txt', test)
