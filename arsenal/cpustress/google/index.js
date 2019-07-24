function cpuIntensiveCalculation(baseNumber) {
    var iterationCount = 50000 * Math.pow(baseNumber, 3);
    var result = 0;
    for (var i = iterationCount; i >= 0; i--) {
        result += Math.atan(i) * Math.tan(i);
    }
}

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    let sleep_time = event.query.level ? parseInt(event.query.level) : null;
    if (!sleep_time || sleep_time < 1) {
        return {"error": "invalid level parameter"};
    }
    return sleep_time;
}

function getParameters(event) {
    return getSleep(event);
}

function runTest(intensityLevel) {
    cpuIntensiveCalculation(intensityLevel);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

async function main(req, res) {
    let startTime = process.hrtime();
    let params = getParameters(req);
    if (params.error) {
        return {"error": params.error}
    }

    runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    res.status(200).send({
        reused: reused,
        duration: duration
    });

}

exports.test128 = main;
exports.test256 = main;
exports.test512 = main;
exports.test1024 = main;
exports.test2048 = main;