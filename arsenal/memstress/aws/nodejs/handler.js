const DEFAULT_INTENSITY_LEVEL = 2;
const MEGABYTE = 1024 * 1024;
const ITERATION_DEVIATOR = 4;

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

function memIntensiveCalculation(level) {
    var startTime = process.hrtime();
    let amountInMB = level * 64 - 1;
    var iterationCount = amountInMB * ITERATION_DEVIATOR;
    var iterationAllocSize = MEGABYTE / ITERATION_DEVIATOR / getMemoryOverhead(amountInMB);
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
    let end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    let got_custom_level = event["level"] && event["level"] !== "0";
    let intensityLevel = got_custom_level ? parseInt(event["level"]) : DEFAULT_INTENSITY_LEVEL;

    return {
        "reused": isWarm(),
        "duration": memIntensiveCalculation(intensityLevel)
    };
};


