require 'fileutils' 

basePath = '../../'
files = Dir.glob("fillSingle/*.prop")
filesMixed = Dir.glob("fillSingleMixed/*.prop")

# files are like "fillSingle/fill-%name%-config.prop"
Dir.chdir(basePath)
files.each do |item| 
	system("call runEval.bat suites/setup/" + item)
	# copy results
	name = item[11...-5]
	FileUtils.mv("result/loadTimes.csv", "suites/setup/result/load-" + name + ".csv");
	FileUtils.mv("result/result.nt", "suites/setup/result/result-" + name + ".nt");
	FileUtils.mv("suites/setup/" + item, "suites/setup/" + item + ".done");
end

# files are like "fillSingleMixed/fill-%name%-config.prop"
filesMixed.each do |item| 
	system("call runEval.bat suites/setup/" + item)
	# copy results
	name = item[16...-5]
	FileUtils.mv("result/loadTimes.csv", "suites/setup/result/load-" + name + ".csv");
	FileUtils.mv("result/result.nt", "suites/setup/result/result-" + name + ".nt");
	FileUtils.mv("suites/setup/" + item, "suites/setup/" + item + ".done");
end