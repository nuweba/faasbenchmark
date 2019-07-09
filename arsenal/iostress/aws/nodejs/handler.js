const fs = require('fs');
const PATH = '/tmp/faastest';
const proc = require( 'child_process' );

function ioIntensiveCalculation(baseNumber) {
    var amountInMB = 10 ** baseNumber;
    proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    fs.unlinkSync(PATH);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    var startTime = process.hrtime();
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if(!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"}
    }

    ioIntensiveCalculation(intensityLevel);

    let retval = {
        "reused": isWarm(),
    };

    var end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};

