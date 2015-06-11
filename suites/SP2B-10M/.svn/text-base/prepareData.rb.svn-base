require 'fileutils' 

basePath = '../../'
files = Dir.glob("setup/*.prop")

# files are like "setup/fill-SP2B-a-config.prop"
Dir.chdir(basePath)
files.each do |item| 
	system("call runEval.bat suites/SP2B-10M/" + item)
	# copy results
	name = item[11...-5]
	FileUtils.mv("result/loadTimes.csv", "suites/SP2B-10M/result/load-" + name + ".csv");
	FileUtils.mv("result/result.nt", "suites/SP2B-10M/result/result-" + name + ".nt");
	FileUtils.mv("suites/SP2B-10M/" + item, "suites/SP2B/" + item + ".done");
end