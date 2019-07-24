var wait = ms => new Promise((r, j) => setTimeout(r, ms));

function getDuration(startTime) {
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function getSleep(event) {
    let sleep_input = event.req.query["sleep"];
    let sleep_time = sleep_input ? parseInt(sleep_input) : null;
    if (!sleep_time && sleep_time !== 0) {
        return {"error": "invalid sleep parameter"};
    }
    return sleep_time;
}

function getParameters(event) {
    return getSleep(event);
}

async function runTest(sleep_time){
    await wait(sleep_time);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

module.exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"error": params.error}
    }

    await runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        "reused": reused,
        "duration": duration
    };
};



