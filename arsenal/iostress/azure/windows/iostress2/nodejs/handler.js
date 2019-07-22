const PATH = 'D:\\local\\temp\\faastest';
const ddPath = 'D:\\Program Files\\Git\\usr\\bin\\dd.exe'
const proc = require('child_process');

function ioIntensive(baseNumber) {
    var amountInMB = 10 ** baseNumber;
    var out = proc.spawnSync(ddPath, ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    return out
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

function runTest(intensityLevel) {
    return ioIntensive(intensityLevel)
}

module.exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"error": params.error}
    }

    testOut = runTest(params);
    if (testOut.error) {
        return testOut.error
    }

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        "reused": reused,
        "duration": duration,
    };
};


