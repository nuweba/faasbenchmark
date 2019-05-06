const DEFAULT_MEMORY_AMOUNT_IN_MB = 127;
const MEGABYTE = 1024 * 1024;
const ITERATION_DEVIATOR = 4;
// More or less relative overhead in that environment + way of code works
var MEMORY_OVERHEAD_ADJUSTER = getMemoryOverhead(DEFAULT_MEMORY_AMOUNT_IN_MB);

function getRandomInt(max) {
	return Math.floor(Math.random() * Math.floor(max));
}

function getMemoryOverhead(amount) {
	var deviator = Math.min(amount, 256);
	var base = 0.1;
	if (amount < 192) {
		base = 0.0018 * (256 - amount);
	}
	return 1.0 + (base * (1 / (amount / deviator)));
}

function memIntensiveCalculation(amountInMB) {
	var iterationCount = amountInMB * ITERATION_DEVIATOR;
	var iterationAllocSize = MEGABYTE / ITERATION_DEVIATOR / getMemoryOverhead(amountInMB);
	console.log("Will allocate " + amountInMB + "MB in " + iterationCount + " iterations (overhead ratio: " + MEMORY_OVERHEAD_ADJUSTER + ")");
	var result = 0;
	var allocatedArrays = [];
	for (var i = iterationCount; i >= 0; i--) {
		var allocatedArray = new Uint8Array(iterationAllocSize);
		for (var j = 0; j < iterationAllocSize / ITERATION_DEVIATOR / 2; j++) {
			allocatedArray[j * ITERATION_DEVIATOR * 2] = getRandomInt(255);
		}
		allocatedArrays.push(allocatedArray);
		result += allocatedArray.length;
	}
	console.log("Total allocated arrays: " + allocatedArrays.length);
	return result;
}

function getMemoryAmount(event) {
	var memoryAmount = DEFAULT_MEMORY_AMOUNT_IN_MB;
	if (event["amount"] && event["amount"] !== "0") {
		memoryAmount = parseInt(event["amount"])
	}
	return memoryAmount;
}

exports.handler = async (event) => {
	var memoryAmount = getMemoryAmount(event);
	var startTime = process.hrtime();
	memIntensiveCalculation(memoryAmount);
	let is_warm = process.env.warm ? true : false;
	let end = process.hrtime(startTime);
	return {
		"reused" : is_warm,
		"duration" : end[1] + (end[0] * 1e9)
	};
};


