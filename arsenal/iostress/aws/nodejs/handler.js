const fs = require('fs');

const DEFAULT_INTENSITY_LEVEL = 5;
const OUTPUT_FILE_PATH = "/tmp/current_value";

function ioIntensiveCalculation(baseNumber) {
	console.time('ioIntensiveCalculation_'+baseNumber);
	var iterationCount = 30 * Math.pow(baseNumber, 3) + 1000;
	console.log("Will perform IO-read-write function for " + iterationCount + " iterations");
	var finalBufferLength = 0;
	var totalBytesRead = 0;
	for (var i = iterationCount; i >= 0; i--) {
		var bufToWrite = "Iteration Number #" + i + "wrote this message to this file, ";
		fs.appendFileSync(OUTPUT_FILE_PATH, bufToWrite);
		var newerFileBuf = fs.readFileSync(OUTPUT_FILE_PATH);
		finalBufferLength = newerFileBuf.length;
		totalBytesRead += newerFileBuf.length;
		console.log('Output file was written ' + bufToWrite.length + ' bytes and has been re-read from with '+ newerFileBuf.length + ' length value!');
	}
	console.timeEnd('ioIntensiveCalculation_'+baseNumber);
	return [finalBufferLength, totalBytesRead];
}

function getIntensityLevel(event) {
	var intensityLevel = DEFAULT_INTENSITY_LEVEL;
	if (event["level"] && event["level"]!== "0") {
		intensityLevel = parseInt(event["level"])
	}
	return intensityLevel;
}

exports.handler = async (event) => {
	var intensityLevel = getIntensityLevel(event);
	var is_warm = process.env.warm ? true : false;
	var startTime = process.hrtime();
	ioIntensiveCalculation(intensityLevel);
	var end = process.hrtime(startTime);
	return {
		"reused" : is_warm,
		"duration" : end[1] + (end[0] * 1e9)
	};
};

