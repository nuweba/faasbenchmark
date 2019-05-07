const DEFAULT_INTENSITY_LEVEL = 1;

function cpuIntensiveCalculation(baseNumber) {
    var startTime = process.hrtime();
    var iterationCount = 50000 * Math.pow(baseNumber, 3);
    var result = 0;
    for (var i = iterationCount; i >= 0; i--) {
        result += Math.atan(i) * Math.tan(i);
    }
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9)
}

exports.handler = async (event) => {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    let custom_level = event["level"] && event["level"] !== "0";
    let intensityLevel = custom_level ? parseInt(event["level"]) : DEFAULT_INTENSITY_LEVEL;

    cpuIntensiveCalculation(intensityLevel);
    return {
        "reused": is_warm,
        "duration": cpuIntensiveCalculation(intensityLevel)
    };
};
