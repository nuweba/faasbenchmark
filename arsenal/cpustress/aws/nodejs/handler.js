const DEFAULT_INTENSITY_LEVEL = 1;

function cpuIntensiveCalculation(baseNumber) {
	console.time('cpuIntensiveCalculation_'+baseNumber);
	var iterationCount = 50000 * Math.pow(baseNumber, 3);
	console.log("Will perform tan(x) functions for " + iterationCount + " iterations");
	var result = 0;
	for (var i = iterationCount; i >= 0; i--) {
		result += Math.atan(i) * Math.tan(i);
	}
	console.timeEnd('cpuIntensiveCalculation_'+baseNumber);
	return result;
}

function getIntensityLevel(event) {
	var intensityLevel = DEFAULT_INTENSITY_LEVEL;
	if (event["level"] && event["level"] !== "0") {
		intensityLevel = parseInt(event["level"])
	}
	return intensityLevel;
}

exports.handler = async (event) => {
	var intensityLevel = getIntensityLevel(event);
	var is_warm = process.env.warm ? true : false;
	var startTime = process.hrtime();
	cpuIntensiveCalculation(intensityLevel);
	var end = process.hrtime(startTime);
	return {
		"reused" : is_warm,
		"duration" : end[1] + (end[0] * 1e9)
	};
};
