const PATH = '/tmp/faastest';
const proc = require('child_process');

function ioIntensive(baseNumber) {
    var amountInMB = 10 ** (baseNumber - 1);
    var out = proc.spawnSync('dd', ['if=/dev/zero', `of=${PATH}`, `bs=${amountInMB}M`, 'count=1', 'oflag=direct']);
    if (out.status !== 0)
        return out.stderr.toString();
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
    let intensityLevel = event.query.level ? parseInt(event.query.level) : null;
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

async function main(req, res) {
    let startTime = process.hrtime();
    let params = getParameters(req);
    if (params.error) {
        return {"error": params.error}
    }

    var error = runTest(params);
    if (error) {
        return {"error": error}
    }

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