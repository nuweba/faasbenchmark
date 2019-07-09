const fs = require('fs');
const PATH = '/tmp/faastest';
const proc = require( 'child_process' );

function ioIntensive(baseNumber) {
    var amountInMB = 10 ** baseNumber;
    proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    fs.unlinkSync(PATH);
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
    ioIntensive(intensityLevel)
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



