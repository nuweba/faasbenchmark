var buf = Buffer.allocUnsafe(10 * 1024 * 1024);

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
    return {}
}

function getParameters(event) {
    return getLevel(event);
}

function runTest(intensityLevel){
}

module.exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"body": `{"error": ${params.error}}`}
    }

    runTest(params);
    event.bindings.blob = buf.toString('binary');

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        body: `{"reused": ${reused}, "duration": ${duration}}`
    };
};

