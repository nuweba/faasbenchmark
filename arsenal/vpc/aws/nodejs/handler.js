var wait = ms => new Promise((r, j) => setTimeout(r, ms));

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
    let sleep_time = event.sleep ? parseInt(event.sleep) : null;
    if (!sleep_time || sleep_time === 0) {
        return {"error": "invalid sleep parameter"};
    }
    return sleep_time;
}

async function runTest(sleep_time){
    await wait(sleep_time);
}

exports.handler = async (event) => {
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
        "duration": duration,
        "response": data
    };
};





