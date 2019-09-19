var buf = Buffer.allocUnsafe(10 * 1024 * 1024);
const AWS = require('aws-sdk');
const s3 = new AWS.S3();

async function upload(requestID) {
    await s3.putObject({
        Bucket: process.env.TEST_BUCKET,
        Key: requestID,
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

function getID(context) {
    let requestID = context.awsRequestId;
    if (!requestID) {
        return {"error": "invalid request ID header"}
    }
    return requestID
}

function getParameters(context) {
    return getID(context);
}

async function runTest(requestID) {
    await upload(requestID)
}

exports.handler = async (event, context) => {
    let startTime = process.hrtime();
    let params = getParameters(context);
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

