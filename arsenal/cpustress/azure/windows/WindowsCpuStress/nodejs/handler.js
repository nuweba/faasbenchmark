function cpuIntensiveCalculation(baseNumber) {
    var iterationCount = 50000 * Math.pow(baseNumber, 3);
    var result = 0;
    for (var i = iterationCount; i >= 0; i--) {
        result += Math.atan(i) * Math.tan(i);
    }
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getLevel(event) {
    let level_input = event.req.query["level"];
    let intensityLevel = level_input ? parseInt(level_input) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel){
    cpuIntensiveCalculation(intensityLevel);
}

module.exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"body": `{"error": ${params.error}}`}
    }

    runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        body: `{"reused": ${reused}, "duration": ${duration}}`
    };
};

