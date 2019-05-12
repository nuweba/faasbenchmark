const fs = require('fs');
const PATH = '/tmp/faastest';
const DEFAULT_INTENSITY_LEVEL = 1;
const proc = require( 'child_process' );

function ioIntensiveCalculation(baseNumber) {
    var startTime = process.hrtime();
    var amountInMB = 10 ** baseNumber;
    proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    fs.unlinkSync(PATH);
    var end = process.hrtime(startTime);
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
        "duration": ioIntensiveCalculation(intensityLevel)
    };
};

