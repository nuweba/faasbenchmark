var wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    let sleep_time = event.query.sleep ? parseInt(event.query.sleep) : null;
    if (!sleep_time && sleep_time !== 0) {
        return {"error": "invalid sleep parameter"};
    }
    return sleep_time;
}

function getParameters(event) {
    return getSleep(event);
}

async function runTest(sleep_time) {
    await wait(sleep_time);
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

    await runTest(params);

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