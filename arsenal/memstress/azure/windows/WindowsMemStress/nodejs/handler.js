const MEGABYTE = 1024 * 1024;

function memIntensive(level) {
    var available_memory = process.env.AWS_LAMBDA_FUNCTION_MEMORY_SIZE;
    let amountInMB = available_memory - (available_memory / 10) * (4 - level);
    console.log(amountInMB);
    Buffer.alloc(amountInMB * MEGABYTE, 'a');
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

function runTest(intensityLevel){
    memIntensive(intensityLevel)
}

module.exports.handler = async (event) => {
    var startTime = process.hrtime();
    let params = getParameters(event);
    if (params.error) {
        return {"body": `{"error": ${params.error}}`}
    }

    runTest(params);

    var reused = isWarm();
    var duration = getDuration(startTime);

    return {
        body: `{"reused": ${reused}, "duration": ${duration}}`
    };
};

