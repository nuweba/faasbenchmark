var buf = Buffer.allocUnsafe(10 * 1024 * 1024);
const AWS = require('aws-sdk');
const s3 = new AWS.S3();

async function upload() {
    await s3.putObject({
        Bucket: process.env.TEST_BUCKET,
        Key: "faastest.dat",
        Body: buf,
    }).promise();
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
    upload(intensityLevel)
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

