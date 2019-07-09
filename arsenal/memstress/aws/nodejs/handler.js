const MEGABYTE = 1024 * 1024;

function memIntensive(level) {
    let amountInMB = level * 64 - 1;
    Buffer.alloc(amountInMB * MEGABYTE);
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

function getParameters(event) {
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if (!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"};
    }
    return intensityLevel;
}

function runTest(intensityLevel){
    memIntensive(intensityLevel)
}

exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"error": params.error}
    }

    runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        "reused": reused,
        "duration": duration
    };
};


